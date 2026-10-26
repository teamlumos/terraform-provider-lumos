package provider

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccPreApprovalRule writes a rule on a custom app and a manual
// entitlement the test creates, changes its justification and condition,
// moves it through each way of choosing entitlements, reads it through both
// data sources and imports it. Every step must plan no change once applied.
// The conditions match an address no one has, so the rule pre-approves
// nothing. The second condition checks a start date to show that a
// datetime value reads back as written.
func TestAccPreApprovalRule(t *testing.T) {
	prefix := "tf-acc-" + acctest.RandString(8)
	email := fmt.Sprintf(`{
      attribute = { label = "Email", value = "identity.email" }
      equals    = { label = %[1]q, value = %[1]q }
    }`, prefix+"@example.com")
	byEmail := "jsonencode(" + email + ")"
	byEmailAndStartDate := `jsonencode({
    and = [
      ` + email + `,
      {
        attribute = { label = "Start date", value = "identity.start_date" }
        ">="      = { label = "2026-01-01T00:00:00Z", value = "2026-01-01T00:00:00Z" }
      },
    ]
  })`

	config := func(justification, condition, entitlements, extra string) string {
		return providerConfig() + testAccCustomApp("app", prefix, "Pre-approval rule test") + fmt.Sprintf(`
resource "lumos_entitlement" "test" {
  app_id = lumos_app.app.id
  label  = %[1]q
}

resource "lumos_pre_approval_rule" "test" {
  app_id        = lumos_app.app.id
  justification = %[2]q
  condition     = %[3]s
  %[4]s
}
%[5]s
`, prefix+"-admin", justification, condition, entitlements, extra)
	}
	includes := "request_includes_entitlements = [lumos_entitlement.test.id]"
	excludes := "request_does_not_include_entitlements = [lumos_entitlement.test.id]"
	anyEntitlement := "request_has_entitlements = true"
	dataSources := `
data "lumos_pre_approval_rule" "test" {
  app_id = lumos_pre_approval_rule.test.app_id
  id     = lumos_pre_approval_rule.test.id
}

data "lumos_pre_approval_rules" "test" {
  app_id = lumos_pre_approval_rule.test.app_id
}
`
	rule := "lumos_pre_approval_rule.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config("Terraform acceptance test", byEmail, includes, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(rule, "id"),
					resource.TestCheckResourceAttrPair(rule, "app_id", "lumos_app.app", "id"),
					resource.TestCheckTypeSetElemAttrPair(rule, "request_includes_entitlements.*", "lumos_entitlement.test", "id"),
					resource.TestCheckNoResourceAttr(rule, "request_has_entitlements"),
				),
			},
			{
				Config: config("Terraform acceptance test, edited", byEmailAndStartDate, includes, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rule, plancheck.ResourceActionUpdate),
						expectOnlyAttributesChange(rule, "justification", "condition"),
					},
				},
			},
			{
				Config: config("Terraform acceptance test, edited", byEmailAndStartDate, anyEntitlement, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rule, "request_has_entitlements", "true"),
					resource.TestCheckResourceAttr(rule, "request_includes_entitlements.#", "0"),
				),
			},
			{
				Config: config("Terraform acceptance test, edited", byEmailAndStartDate, excludes, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemAttrPair(rule, "request_does_not_include_entitlements.*", "lumos_entitlement.test", "id"),
					resource.TestCheckNoResourceAttr(rule, "request_has_entitlements"),
				),
			},
			{
				// No entitlements field: the update replaces the list, and
				// the rule covers only a request for the app itself.
				Config: config("Terraform acceptance test, edited", byEmailAndStartDate, "", dataSources),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rule, "request_has_entitlements", "false"),
					resource.TestCheckResourceAttr(rule, "request_does_not_include_entitlements.#", "0"),
					resource.TestCheckResourceAttrPair("data.lumos_pre_approval_rule.test", "condition", rule, "condition"),
					resource.TestCheckResourceAttr("data.lumos_pre_approval_rule.test", "request_has_entitlements", "false"),
					resource.TestCheckResourceAttr("data.lumos_pre_approval_rules.test", "pre_approval_rules.#", "1"),
					resource.TestCheckResourceAttrPair("data.lumos_pre_approval_rules.test", "pre_approval_rules.0.id", rule, "id"),
				),
			},
			{
				ResourceName:      rule,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					attributes := s.RootModule().Resources[rule].Primary.Attributes
					id, err := json.Marshal(map[string]string{"app_id": attributes["app_id"], "id": attributes["id"]})
					return string(id), err
				},
			},
		},
	})
}

// TestAccPreApprovalRuleWebhook writes a rule that pre-approves requests a
// pre-approval webhook approves, on a custom app the test creates, and
// reads it back through the data source. LUMOS_TEST_PREAPPROVAL_WEBHOOK_ID
// names the webhook, which the provider can't create.
func TestAccPreApprovalRuleWebhook(t *testing.T) {
	webhookID := testAccOptionalEnv(t, "LUMOS_TEST_PREAPPROVAL_WEBHOOK_ID")
	testAccPreApprovalRuleCondition(t, fmt.Sprintf(`jsonencode({
    attribute = { label = "Access check", value = "webhook.responses[\"%s\"]" }
    equals    = { label = "Pre-approved", value = "PREAPPROVED" }
  })`, webhookID))
}

// TestAccPreApprovalRuleOnCall writes a rule that pre-approves whoever is on
// call, on a custom app the test creates, and reads it back through the
// data source. LUMOS_TEST_ONCALL_SCHEDULE_ENTITLEMENT_ID names the on-call
// schedule's entitlement, which the provider can't create.
func TestAccPreApprovalRuleOnCall(t *testing.T) {
	scheduleID := testAccOptionalEnv(t, "LUMOS_TEST_ONCALL_SCHEDULE_ENTITLEMENT_ID")
	testAccPreApprovalRuleCondition(t, fmt.Sprintf(`jsonencode({
    attribute = { label = "On call", value = "identity.oncall[\"%s\"].on_call" }
    equals    = { label = "On call", value = true }
  })`, scheduleID))
}

// testAccPreApprovalRuleCondition writes a rule with condition on a custom
// app, checks the data source reads the condition back as written, and
// imports the rule. Every step must plan no change once applied.
func testAccPreApprovalRuleCondition(t *testing.T, condition string) {
	prefix := "tf-acc-" + acctest.RandString(8)
	rule := "lumos_pre_approval_rule.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + testAccCustomApp("app", prefix, "Pre-approval rule test") + `
resource "lumos_pre_approval_rule" "test" {
  app_id        = lumos_app.app.id
  justification = "Terraform acceptance test"
  condition     = ` + condition + `
}

data "lumos_pre_approval_rule" "test" {
  app_id = lumos_pre_approval_rule.test.app_id
  id     = lumos_pre_approval_rule.test.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.lumos_pre_approval_rule.test", "condition", rule, "condition"),
				),
			},
			{
				ResourceName:      rule,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					attributes := s.RootModule().Resources[rule].Primary.Attributes
					id, err := json.Marshal(map[string]string{"app_id": attributes["app_id"], "id": attributes["id"]})
					return string(id), err
				},
			},
		},
	})
}
