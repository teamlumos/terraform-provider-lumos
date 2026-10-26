package provider

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// TestAccWorkflowAssignees creates an approval chain with one stage per kind
// of assignee: the user in USER_ID, the entitlement owner, app and org
// admins, and, when LUMOS_TEST_CUSTOM_ATTRIBUTE_LABEL names a custom
// attribute definition in the test domain, that custom attribute. The plan
// after apply must be empty, so each assignee reads back as written, and the
// import must match. The chain is never bound.
func TestAccWorkflowAssignees(t *testing.T) {
	userID := testAccOptionalEnv(t, "USER_ID")
	attributeLabel := os.Getenv("LUMOS_TEST_CUSTOM_ATTRIBUTE_LABEL")
	prefix := "tf-acc-" + acctest.RandString(8)

	assignees := []string{
		fmt.Sprintf(`[{ uuid = %q, type = "user" }]`, userID),
		`[{ uuid = "entitlement_owner", type = "dynamic_assignee" }]`,
		`[{ uuid = "app_admins", type = "dynamic_assignee" }, { uuid = "org_admins", type = "dynamic_assignee" }]`,
	}
	if attributeLabel != "" {
		assignees = append(assignees, fmt.Sprintf(`[{ label = %q, type = "custom_attribute_definition" }]`, attributeLabel))
	}
	stages := make([]string, len(assignees))
	for i, a := range assignees {
		stages[i] = fmt.Sprintf(`
    {
      step = "approval_task"
      inputs = jsonencode({
        stage_number = %d
        assignees    = %s
      })
    },`, i+1, a)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + testAccCustomApp("app", prefix, "Workflow test") + fmt.Sprintf(`
resource "lumos_workflow" "chain" {
  app_id = lumos_app.app.id
  type   = "approval-chain"
  name   = "%s-assignees"
  steps = [%s
  ]
}
`, prefix, strings.Join(stages, "")),
				Check: resource.TestCheckResourceAttr("lumos_workflow.chain", "steps.#", fmt.Sprint(len(assignees))),
			},
			{
				ResourceName:      "lumos_workflow.chain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccWorkflowProvisioningTypes creates a one-step wait workflow of each
// provisioning and deprovisioning type on a custom app the test creates,
// changes the wait, which must plan only the steps, then imports it. The
// duration is written in the ISO 8601 form the API reads back. Each type
// gets its own app, so the types run in parallel.
func TestAccWorkflowProvisioningTypes(t *testing.T) {
	for _, workflowType := range []string{
		"account-provisioning",
		"account-deprovisioning",
		"entitlement-provisioning",
		"entitlement-deprovisioning",
	} {
		t.Run(workflowType, func(t *testing.T) {
			t.Parallel()
			prefix := "tf-acc-" + acctest.RandString(8)
			config := func(duration string) string {
				return providerConfig() + testAccCustomApp("app", prefix, "Workflow test") + fmt.Sprintf(`
resource "lumos_workflow" "test" {
  app_id = lumos_app.app.id
  type   = %q
  name   = "%s-wait"
  steps = [{
    step   = "wait"
    inputs = jsonencode({ duration = %q })
  }]
}
`, workflowType, prefix, duration)
			}

			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: config("P1D"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("lumos_workflow.test", "type", workflowType),
							resource.TestCheckResourceAttr("lumos_workflow.test", "steps.0.step", "wait"),
						),
					},
					{
						Config: config("PT12H"),
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectResourceAction("lumos_workflow.test", plancheck.ResourceActionUpdate),
								expectOnlyAttributesChange("lumos_workflow.test", "steps"),
							},
						},
					},
					{
						ResourceName:      "lumos_workflow.test",
						ImportState:       true,
						ImportStateVerify: true,
					},
				},
			})
		})
	}
}

// TestAccWorkflowApproval creates an approval chain with one stage assigned
// to GROUP_ID and an approval template whose branch step sends the Eng and
// IT teams to that chain, with the same chain as the fallback, all on a
// custom app the test creates. The condition is written in the expanded
// { label, value } form, which reads back unchanged. It renames the chain,
// which must plan only the name, then imports both. Neither workflow is
// bound, and the template is never updated.
func TestAccWorkflowApproval(t *testing.T) {
	groupID := testAccOptionalEnv(t, "GROUP_ID")
	prefix := "tf-acc-" + acctest.RandString(8)
	config := func(chainName string) string {
		return providerConfig() + testAccCustomApp("app", prefix, "Workflow test") + fmt.Sprintf(`
resource "lumos_workflow" "chain" {
  app_id = lumos_app.app.id
  type   = "approval-chain"
  name   = %q
  steps = [{
    step = "approval_task"
    inputs = jsonencode({
      stage_number = 1
      assignees    = [{ uuid = %q, type = "group" }]
    })
  }]
}

resource "lumos_workflow" "template" {
  app_id      = lumos_app.app.id
  type        = "approval"
  name        = "%s-template"
  description = "Routes the Eng and IT teams to the chain"
  steps = [{
    step = "branch"
    inputs = jsonencode({
      branches = [{
        condition = {
          and = [{
            attribute = { label = "Team", value = "identity.team" }
            in        = [{ label = "Engineering", value = "Eng" }, { label = "IT", value = "IT" }]
          }]
        }
        workflow_id = lumos_workflow.chain.id
      }]
      fallback_workflow_id = lumos_workflow.chain.id
    })
  }]
}
`, chainName, groupID, prefix)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(prefix + "-chain"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("lumos_workflow.chain", "id"),
					resource.TestCheckResourceAttrPair("lumos_workflow.chain", "app_id", "lumos_app.app", "id"),
					resource.TestCheckNoResourceAttr("lumos_workflow.chain", "description"),
					resource.TestCheckResourceAttr("lumos_workflow.chain", "steps.0.step", "approval_task"),
					resource.TestCheckResourceAttr("lumos_workflow.template", "type", "approval"),
					resource.TestCheckResourceAttr("lumos_workflow.template", "steps.0.step", "branch"),
				),
			},
			{
				Config: config(prefix + "-renamed"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_workflow.chain", plancheck.ResourceActionUpdate),
						expectOnlyAttributesChange("lumos_workflow.chain", "name"),
						plancheck.ExpectResourceAction("lumos_workflow.template", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.TestCheckResourceAttr("lumos_workflow.chain", "name", prefix+"-renamed"),
			},
			{
				ResourceName:      "lumos_workflow.chain",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "lumos_workflow.template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
