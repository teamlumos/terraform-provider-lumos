package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// testAccCheckNotEmpty fails when the list attribute at key has no elements.
func testAccCheckNotEmpty(name, key string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith(name, key+".#", func(v string) error {
		if v == "0" {
			return fmt.Errorf("%s is empty", key)
		}
		return nil
	})
}

// TestAccLookupDataSources reads each lookup data source against the test
// domain. Lookups of existing objects use the test domain's settings (APP_ID,
// APP_NAME, GROUP_NAME, GROUP_ID, USER_ID) and skip when one isn't set; the
// rest only check that the read succeeds or finds something every domain has.
func TestAccLookupDataSources(t *testing.T) {
	tests := []struct {
		name   string
		env    []string
		config func(env map[string]string) string
		check  func(env map[string]string) resource.TestCheckFunc
	}{
		{
			name: "apps",
			env:  []string{"APP_ID", "APP_NAME"},
			config: func(env map[string]string) string {
				return fmt.Sprintf(`data "lumos_apps" "test" { search = %q }`, env["APP_NAME"])
			},
			check: func(env map[string]string) resource.TestCheckFunc {
				return resource.TestCheckTypeSetElemNestedAttrs("data.lumos_apps.test", "apps.*", map[string]string{"id": env["APP_ID"]})
			},
		},
		{
			// The tenant-wide search looks for the label of one of the app's
			// own entitlements, so it finds that entitlement again.
			name: "app and tenant-wide entitlements",
			env:  []string{"APP_ID"},
			config: func(env map[string]string) string {
				return fmt.Sprintf(`
data "lumos_app_entitlements" "test" { app_id = %q }

data "lumos_entitlements" "test" {
  search = data.lumos_app_entitlements.test.entitlements[0].label
}
`, env["APP_ID"])
			},
			check: func(env map[string]string) resource.TestCheckFunc {
				return resource.ComposeAggregateTestCheckFunc(
					testAccCheckNotEmpty("data.lumos_app_entitlements.test", "entitlements"),
					resource.TestCheckResourceAttr("data.lumos_app_entitlements.test", "entitlements.0.app_id", env["APP_ID"]),
					resource.TestCheckTypeSetElemAttrPair("data.lumos_entitlements.test", "entitlements.*.id", "data.lumos_app_entitlements.test", "entitlements.0.id"),
				)
			},
		},
		{
			name: "accounts",
			env:  []string{"APP_ID"},
			config: func(env map[string]string) string {
				return fmt.Sprintf(`data "lumos_accounts" "test" { app_id = %q }`, env["APP_ID"])
			},
			check: func(map[string]string) resource.TestCheckFunc {
				return resource.TestCheckResourceAttrSet("data.lumos_accounts.test", "accounts.#")
			},
		},
		{
			name: "user roles",
			env:  []string{"USER_ID"},
			config: func(env map[string]string) string {
				return fmt.Sprintf(`data "lumos_user_roles" "test" { user_id = %q }`, env["USER_ID"])
			},
			check: func(env map[string]string) resource.TestCheckFunc {
				return resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lumos_user_roles.test", "user_id", env["USER_ID"]),
					resource.TestCheckResourceAttrSet("data.lumos_user_roles.test", "roles.#"),
				)
			},
		},
		{
			name: "groups",
			env:  []string{"GROUP_ID", "GROUP_NAME"},
			config: func(env map[string]string) string {
				return fmt.Sprintf(`
data "lumos_groups" "test" {
  name        = %q
  exact_match = true
}
`, env["GROUP_NAME"])
			},
			check: func(env map[string]string) resource.TestCheckFunc {
				return resource.TestCheckTypeSetElemNestedAttrs("data.lumos_groups.test", "groups.*", map[string]string{"id": env["GROUP_ID"]})
			},
		},
		{
			name: "inline webhooks",
			config: func(map[string]string) string {
				return `data "lumos_inline_webhooks" "test" {}`
			},
			check: func(map[string]string) resource.TestCheckFunc {
				return resource.TestCheckResourceAttrSet("data.lumos_inline_webhooks.test", "inline_webhooks.#")
			},
		},
		{
			name: "app categories",
			config: func(map[string]string) string {
				return `data "lumos_app_categories" "test" {}`
			},
			check: func(map[string]string) resource.TestCheckFunc {
				return testAccCheckNotEmpty("data.lumos_app_categories.test", "app_categories")
			},
		},
		{
			// The connection schema is read for the first integration the
			// search finds, so the test doesn't depend on one integration's
			// identifier.
			name: "integrations and connection schema",
			config: func(map[string]string) string {
				return `
data "lumos_integrations" "test" { query = "Okta" }

data "lumos_integration_connection_schema" "test" {
  app_class_id = data.lumos_integrations.test.integrations[0].app_class_id
}
`
			},
			check: func(map[string]string) resource.TestCheckFunc {
				return resource.ComposeAggregateTestCheckFunc(
					testAccCheckNotEmpty("data.lumos_integrations.test", "integrations"),
					resource.TestCheckResourceAttrPair("data.lumos_integration_connection_schema.test", "app_class_id", "data.lumos_integrations.test", "integrations.0.app_class_id"),
					resource.TestCheckResourceAttrSet("data.lumos_integration_connection_schema.test", "auth_schema"),
				)
			},
		},
		{
			name: "custom attribute definitions",
			config: func(map[string]string) string {
				return `data "lumos_custom_attribute_definitions" "test" {}`
			},
			check: func(map[string]string) resource.TestCheckFunc {
				return resource.TestCheckResourceAttrSet("data.lumos_custom_attribute_definitions.test", "custom_attribute_definitions.#")
			},
		},
		{
			name: "workflow steps",
			env:  []string{"APP_ID"},
			config: func(env map[string]string) string {
				return fmt.Sprintf(`
data "lumos_workflow_steps" "test" {
  app_id = %q
  type   = "account-provisioning"
}
`, env["APP_ID"])
			},
			check: func(map[string]string) resource.TestCheckFunc {
				return testAccCheckNotEmpty("data.lumos_workflow_steps.test", "workflow_steps")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := map[string]string{}
			for _, name := range tt.env {
				env[name] = os.Getenv(name)
				if env[name] == "" {
					t.Skipf("%s isn't set", name)
				}
			}
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{{
					Config: providerConfig() + tt.config(env),
					Check:  tt.check(env),
				}},
			})
		})
	}
}
