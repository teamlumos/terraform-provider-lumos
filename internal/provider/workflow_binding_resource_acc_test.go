package provider

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// testAccApprovalTemplate is an approval template named name whose only
// chain, created alongside it, has one stage assigned to the app's admins.
func testAccApprovalTemplate(name, appRef string) string {
	return fmt.Sprintf(`
resource "lumos_workflow" "%[1]s_chain" {
  app_id = %[2]s
  type   = "approval-chain"
  name   = "%[1]s chain"
  steps = [{
    step   = "approval_task"
    inputs = jsonencode({ stage_number = 1, assignees = [{ uuid = "app_admins", type = "dynamic_assignee" }] })
  }]
}

resource "lumos_workflow" "%[1]s" {
  app_id = %[2]s
  type   = "approval"
  name   = "%[1]s"
  steps = [{
    step   = "branch"
    inputs = jsonencode({ branches = [], fallback_workflow_id = lumos_workflow.%[1]s_chain.id })
  }]
}
`, name, appRef)
}

// testAccWaitWorkflow is a one-step wait workflow of the given type.
func testAccWaitWorkflow(name, workflowType, appRef string) string {
	return fmt.Sprintf(`
resource "lumos_workflow" "%[1]s" {
  app_id = %[3]s
  type   = %[2]q
  name   = "%[1]s"
  steps = [{
    step   = "wait"
    inputs = jsonencode({ duration = "P1D" })
  }]
}
`, name, workflowType, appRef)
}

// testAccBinding binds lumos_workflow.<workflow> with a resource of the given
// type named "test".
func testAccBinding(resourceType, ownerAttr, ownerRef, workflow string) string {
	return fmt.Sprintf(`
resource %q "test" {
  %s = %s
  workflow_id = lumos_workflow.%s.id
}
`, resourceType, ownerAttr, ownerRef, workflow)
}

func testAccRemovedBinding(resourceType string) string {
	return testAccRemoved(resourceType + ".test")
}

// testAccRemoved stops managing addr without changing Lumos.
func testAccRemoved(addr string) string {
	return fmt.Sprintf(`
removed {
  from = %s
  lifecycle {
    destroy = false
  }
}
`, addr)
}

// testAccBindingWorkflows names the lumos_workflow resources the binding
// tests create, templates before the chains they use, so they can be
// deleted in this order.
var testAccBindingWorkflows = []string{
	"approval_a", "approval_b", "approval_a_chain", "approval_b_chain",
	"provision_a", "provision_b", "deprovision_a", "deprovision_b",
}

// testAccRemovedBindingWorkflows stops managing every workflow in
// testAccBindingWorkflows. Lumos refuses to delete a workflow while its app
// or entitlement still binds it, so a test releases them and lets its owner
// be destroyed first, as a config would.
func testAccRemovedBindingWorkflows() string {
	var b strings.Builder
	for _, name := range testAccBindingWorkflows {
		b.WriteString(testAccRemoved("lumos_workflow." + name))
	}
	return b.String()
}

// testAccCaptureWorkflowIDs records the ID of each workflow in
// testAccBindingWorkflows that state holds.
func testAccCaptureWorkflowIDs(ids map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, name := range testAccBindingWorkflows {
			if rs, ok := s.RootModule().Resources["lumos_workflow."+name]; ok {
				ids[name] = rs.Primary.Attributes["id"]
			}
		}
		return nil
	}
}

// testAccDeleteWorkflows deletes the workflows in ids through the API, in
// testAccBindingWorkflows order. It runs once the test's app and
// entitlement are destroyed. A delete the API refuses is logged rather than
// failed, since whether Lumos lets a deleted owner's workflows go is up to
// the API (GOV-4488); each one left behind is named by its ID.
func testAccDeleteWorkflows(t *testing.T, ids map[string]string) {
	for _, name := range testAccBindingWorkflows {
		id := ids[name]
		if id == "" {
			continue
		}
		if err := testAccDeleteWorkflow(id); err != nil {
			t.Logf("left workflow %s (%s) in Lumos: %v", name, id, err)
		}
	}
}

// testAccDeleteWorkflow deletes one workflow through the API, waiting out
// rate limiting as the provider's client does. A workflow that's already
// gone counts as deleted.
func testAccDeleteWorkflow(id string) error {
	url := strings.TrimSuffix(os.Getenv("LUMOS_API_URL"), "/") + "/v1/workflows/" + id
	for attempt := 1; ; attempt++ {
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+os.Getenv("LUMOS_API_TOKEN"))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		switch {
		case resp.StatusCode < 300 || resp.StatusCode == http.StatusNotFound:
			return nil
		case resp.StatusCode == http.StatusTooManyRequests && attempt < 5:
			wait, err := strconv.Atoi(resp.Header.Get("Retry-After"))
			if err != nil || wait < 1 {
				wait = 1
			}
			time.Sleep(time.Duration(wait) * time.Second)
		default:
			return fmt.Errorf("DELETE returned %d: %s", resp.StatusCode, body)
		}
	}
}

// importByAttr imports addr by the value its state holds for attr.
func importByAttr(addr, attr string) resource.TestStep {
	return resource.TestStep{
		ResourceName: addr,
		ImportState:  true,
		ImportStateIdFunc: func(s *terraform.State) (string, error) {
			rs, ok := s.RootModule().Resources[addr]
			if !ok {
				return "", fmt.Errorf("%s not in state", addr)
			}
			return rs.Primary.Attributes[attr], nil
		},
		ImportStateVerify:                    true,
		ImportStateVerifyIdentifierAttribute: attr,
	}
}

// TestAccAppWorkflowBindings binds an approval template and account
// provisioning and deprovisioning workflows to a custom app the test
// creates, rebinds each slot to a second workflow (planning only
// workflow_id), and imports each binding by the app's ID. It ends by
// replacing the bindings and workflows with removed blocks, since Lumos
// can't empty these slots or delete a workflow its app binds, so the
// post-test destroy deletes only the app; the workflows are then deleted
// through the API. No bound template is updated.
func TestAccAppWorkflowBindings(t *testing.T) {
	prefix := "tf-acc-" + acctest.RandString(8)
	ids := map[string]string{}
	t.Cleanup(func() { testAccDeleteWorkflows(t, ids) })
	app := "lumos_app.app.id"
	workflows := testAccApprovalTemplate("approval_a", app) + testAccApprovalTemplate("approval_b", app) +
		testAccWaitWorkflow("provision_a", "account-provisioning", app) +
		testAccWaitWorkflow("provision_b", "account-provisioning", app) +
		testAccWaitWorkflow("deprovision_a", "account-deprovisioning", app) +
		testAccWaitWorkflow("deprovision_b", "account-deprovisioning", app)
	owners := providerConfig() + testAccCustomApp("app", prefix, "Workflow binding test")
	base := owners + workflows
	slots := []struct{ resourceType, workflow string }{
		{"lumos_app_approval_workflow", "approval"},
		{"lumos_app_provisioning_workflow", "provision"},
		{"lumos_app_deprovisioning_workflow", "deprovision"},
	}
	bindings := func(suffix string) string {
		var b strings.Builder
		for _, s := range slots {
			b.WriteString(testAccBinding(s.resourceType, "app_id", app, s.workflow+"_"+suffix))
		}
		return b.String()
	}

	var bindChecks, rebindChecks []resource.TestCheckFunc
	var rebindPlan []plancheck.PlanCheck
	var removed strings.Builder
	steps := []resource.TestStep{}
	for _, s := range slots {
		addr := s.resourceType + ".test"
		bindChecks = append(bindChecks, resource.TestCheckResourceAttrPair(addr, "workflow_id", "lumos_workflow."+s.workflow+"_a", "id"))
		rebindChecks = append(rebindChecks, resource.TestCheckResourceAttrPair(addr, "workflow_id", "lumos_workflow."+s.workflow+"_b", "id"))
		rebindPlan = append(rebindPlan,
			plancheck.ExpectResourceAction(addr, plancheck.ResourceActionUpdate),
			expectOnlyAttributesChange(addr, "workflow_id"))
		removed.WriteString(testAccRemovedBinding(s.resourceType))
	}
	steps = append(steps,
		resource.TestStep{Config: base + bindings("a"), Check: resource.ComposeAggregateTestCheckFunc(bindChecks...)},
		resource.TestStep{
			Config:           base + bindings("b"),
			ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: rebindPlan},
			Check:            resource.ComposeAggregateTestCheckFunc(append(rebindChecks, testAccCaptureWorkflowIDs(ids))...),
		},
	)
	for _, s := range slots {
		steps = append(steps, importByAttr(s.resourceType+".test", "app_id"))
	}
	steps = append(steps, resource.TestStep{Config: owners + removed.String() + testAccRemovedBindingWorkflows()})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

// TestAccEntitlementWorkflowBindings binds an approval template override
// and entitlement provisioning and deprovisioning workflows to a manual
// entitlement on a custom app the test creates, rebinds each slot to a
// second workflow (planning only workflow_id), and imports each binding by
// the entitlement's ID. It then destroys the override binding, which turns
// the override off, and replaces the other two with removed blocks. It ends
// by replacing the workflows with removed blocks too, so the post-test
// destroy deletes only the entitlement and the app; the workflows are then
// deleted through the API. No bound template is updated.
func TestAccEntitlementWorkflowBindings(t *testing.T) {
	prefix := "tf-acc-" + acctest.RandString(8)
	ids := map[string]string{}
	t.Cleanup(func() { testAccDeleteWorkflows(t, ids) })
	app := "lumos_app.app.id"
	ent := "lumos_entitlement.ent.id"
	workflows := testAccApprovalTemplate("approval_a", app) + testAccApprovalTemplate("approval_b", app) +
		testAccWaitWorkflow("provision_a", "entitlement-provisioning", app) +
		testAccWaitWorkflow("provision_b", "entitlement-provisioning", app) +
		testAccWaitWorkflow("deprovision_a", "entitlement-deprovisioning", app) +
		testAccWaitWorkflow("deprovision_b", "entitlement-deprovisioning", app)
	owners := providerConfig() + testAccCustomApp("app", prefix, "Workflow binding test") + fmt.Sprintf(`
resource "lumos_entitlement" "ent" {
  app_id = lumos_app.app.id
  label  = "%s-admin"
}
`, prefix)
	base := owners + workflows
	slots := []struct{ resourceType, workflow string }{
		{"lumos_entitlement_approval_workflow", "approval"},
		{"lumos_entitlement_provisioning_workflow", "provision"},
		{"lumos_entitlement_deprovisioning_workflow", "deprovision"},
	}
	bindings := func(suffix string) string {
		var b strings.Builder
		for _, s := range slots {
			b.WriteString(testAccBinding(s.resourceType, "entitlement_id", ent, s.workflow+"_"+suffix))
		}
		return b.String()
	}

	var bindChecks, rebindChecks []resource.TestCheckFunc
	var rebindPlan []plancheck.PlanCheck
	steps := []resource.TestStep{}
	for _, s := range slots {
		addr := s.resourceType + ".test"
		bindChecks = append(bindChecks, resource.TestCheckResourceAttrPair(addr, "workflow_id", "lumos_workflow."+s.workflow+"_a", "id"))
		rebindChecks = append(rebindChecks, resource.TestCheckResourceAttrPair(addr, "workflow_id", "lumos_workflow."+s.workflow+"_b", "id"))
		rebindPlan = append(rebindPlan,
			plancheck.ExpectResourceAction(addr, plancheck.ResourceActionUpdate),
			expectOnlyAttributesChange(addr, "workflow_id"))
	}
	steps = append(steps,
		resource.TestStep{Config: base + bindings("a"), Check: resource.ComposeAggregateTestCheckFunc(bindChecks...)},
		resource.TestStep{
			Config:           base + bindings("b"),
			ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: rebindPlan},
			Check:            resource.ComposeAggregateTestCheckFunc(append(rebindChecks, testAccCaptureWorkflowIDs(ids))...),
		},
	)
	for _, s := range slots {
		steps = append(steps, importByAttr(s.resourceType+".test", "entitlement_id"))
	}
	cleanup := base + testAccRemovedBinding("lumos_entitlement_provisioning_workflow") +
		testAccRemovedBinding("lumos_entitlement_deprovisioning_workflow")
	steps = append(steps,
		resource.TestStep{
			Config: cleanup,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("lumos_entitlement_approval_workflow.test", plancheck.ResourceActionDestroy),
				},
			},
		},
		// The entitlement reads its override back as unset once it's off.
		resource.TestStep{
			RefreshState: true,
			Check:        resource.TestCheckNoResourceAttr("lumos_entitlement.ent", "approval_workflow_override.id"),
		},
		resource.TestStep{Config: owners + testAccRemovedBindingWorkflows()},
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}
