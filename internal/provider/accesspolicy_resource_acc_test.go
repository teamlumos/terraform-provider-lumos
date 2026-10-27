package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// TestAccAccessPolicy grants a custom app and a manual entitlement the test
// creates, moving the policy between the everyone condition and a condition
// written in the expanded form, then publishes it, renames it while
// published and imports it. Every step must plan no change once applied.
// The condition matches an address no one has, so the policy can't collide
// with another published policy's condition, and it stays disabled, so it
// grants nothing; the everyone condition is only ever on the draft.
func TestAccAccessPolicy(t *testing.T) {
	prefix := "tf-acc-" + acctest.RandString(8)
	condition := fmt.Sprintf(`jsonencode({
    attribute = { label = "Email", value = "identity.email" }
    equals    = { label = %[1]q, value = %[1]q }
  })`, prefix+"@example.com")

	config := func(name, audience, status, extra string) string {
		return providerConfig() + testAccCustomApp("app", prefix, "Access policy test") + fmt.Sprintf(`
resource "lumos_entitlement" "test" {
  app_id = lumos_app.app.id
  label  = %[1]q
}

resource "lumos_access_policy" "test" {
  name                   = %[2]q
  business_justification = "Terraform acceptance test"
  %[3]s
  status                 = %[4]q
  is_enabled             = false
  apps = [
    {
      id           = lumos_app.app.id
      entitlements = [{ id = lumos_entitlement.test.id }]
    },
  ]
}
%[5]s
`, prefix+"-admin", name, audience, status, extra)
	}
	everyone := "is_everyone_condition = true"
	byCondition := "access_condition = " + condition
	byName := `
data "lumos_access_policies" "by_name" {
  name = lumos_access_policy.test.name
}
`
	renamed := prefix + "-renamed"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(prefix, everyone, "draft", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_access_policy.test", "is_everyone_condition", "true"),
					resource.TestCheckNoResourceAttr("lumos_access_policy.test", "access_condition"),
					resource.TestCheckResourceAttr("lumos_access_policy.test", "status", "draft"),
					resource.TestCheckResourceAttr("lumos_access_policy.test", "is_enabled", "false"),
					resource.TestCheckResourceAttrSet("lumos_access_policy.test", "id"),
				),
			},
			{
				Config: config(prefix, byCondition, "draft", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_access_policy.test", "is_everyone_condition", "false"),
					resource.TestCheckResourceAttrSet("lumos_access_policy.test", "access_condition"),
				),
			},
			{
				// Back to everyone: the update must not send the previous
				// condition along with the flag.
				Config: config(prefix, everyone, "draft", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_access_policy.test", "is_everyone_condition", "true"),
					resource.TestCheckNoResourceAttr("lumos_access_policy.test", "access_condition"),
				),
			},
			{
				Config: config(prefix, byCondition, "published", byName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_access_policy.test", "status", "published"),
					resource.TestCheckResourceAttr("data.lumos_access_policies.by_name", "access_policies.#", "1"),
					resource.TestCheckResourceAttrPair("data.lumos_access_policies.by_name", "access_policies.0.access_condition", "lumos_access_policy.test", "access_condition"),
				),
			},
			{
				Config: config(renamed, byCondition, "published", ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_access_policy.test", plancheck.ResourceActionUpdate),
						expectOnlyAttributesChange("lumos_access_policy.test", "name"),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_access_policy.test", "name", renamed),
					resource.TestCheckResourceAttr("lumos_access_policy.test", "status", "published"),
				),
			},
			{
				ResourceName:      "lumos_access_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
