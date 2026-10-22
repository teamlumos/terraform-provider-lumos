package provider

import (
	"fmt"
	"strings"
	"testing"

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
	return fmt.Sprintf(`
removed {
  from = %s.test
  lifecycle {
    destroy = false
  }
}
`, resourceType)
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
// replacing the bindings with removed blocks, since Lumos can't empty these
// slots. No bound template is updated.
func TestAccAppWorkflowBindings(t *testing.T) {
	prefix := "tf-acc-" + acctest.RandString(8)
	app := "lumos_app.app.id"
	workflows := testAccApprovalTemplate("approval_a", app) + testAccApprovalTemplate("approval_b", app) +
		testAccWaitWorkflow("provision_a", "account-provisioning", app) +
		testAccWaitWorkflow("provision_b", "account-provisioning", app) +
		testAccWaitWorkflow("deprovision_a", "account-deprovisioning", app) +
		testAccWaitWorkflow("deprovision_b", "account-deprovisioning", app)
	base := providerConfig() + testAccCustomApp("app", prefix, "Workflow binding test") + workflows
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
			Check:            resource.ComposeAggregateTestCheckFunc(rebindChecks...),
		},
	)
	for _, s := range slots {
		steps = append(steps, importByAttr(s.resourceType+".test", "app_id"))
	}
	steps = append(steps, resource.TestStep{Config: base + removed.String()})

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
// the override off, and replaces the other two with removed blocks. No
// bound template is updated.
func TestAccEntitlementWorkflowBindings(t *testing.T) {
	prefix := "tf-acc-" + acctest.RandString(8)
	app := "lumos_app.app.id"
	ent := "lumos_entitlement.ent.id"
	workflows := testAccApprovalTemplate("approval_a", app) + testAccApprovalTemplate("approval_b", app) +
		testAccWaitWorkflow("provision_a", "entitlement-provisioning", app) +
		testAccWaitWorkflow("provision_b", "entitlement-provisioning", app) +
		testAccWaitWorkflow("deprovision_a", "entitlement-deprovisioning", app) +
		testAccWaitWorkflow("deprovision_b", "entitlement-deprovisioning", app)
	base := providerConfig() + testAccCustomApp("app", prefix, "Workflow binding test") + fmt.Sprintf(`
resource "lumos_entitlement" "ent" {
  app_id = lumos_app.app.id
  label  = "%s-admin"
}
`, prefix) + workflows
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
			Check:            resource.ComposeAggregateTestCheckFunc(rebindChecks...),
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
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}
