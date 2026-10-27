package provider

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// TestAccEntitlementManual creates a manual entitlement on a custom app the
// test creates, renames it, then imports it. The rename must plan only the
// label (GOV-4388).
func TestAccEntitlementManual(t *testing.T) {
	prefix := "tf-acc-" + acctest.RandString(8)
	config := func(label string) string {
		return providerConfig() + testAccCustomApp("app", prefix, "Entitlement test") + fmt.Sprintf(`
resource "lumos_entitlement" "test" {
  app_id = lumos_app.app.id
  label  = %q
}
`, label)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(prefix + "-admin"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_entitlement.test", "label", prefix+"-admin"),
					resource.TestCheckResourceAttr("lumos_entitlement.test", "source", "manual"),
					resource.TestCheckResourceAttrPair("lumos_entitlement.test", "app_id", "lumos_app.app", "id"),
					resource.TestCheckResourceAttrSet("lumos_entitlement.test", "id"),
				),
			},
			{
				Config: config(prefix + "-owner"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_entitlement.test", plancheck.ResourceActionUpdate),
						expectOnlyAttributesChange("lumos_entitlement.test", "label"),
					},
				},
				Check: resource.TestCheckResourceAttr("lumos_entitlement.test", "label", prefix+"-owner"),
			},
			{
				ResourceName:      "lumos_entitlement.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccEntitlementAppStoreSettings lists the app in
// LUMOS_TEST_APPSTORE_APP_ID (settings are only available for App Store
// apps) with the token's user as its admin, creates an entitlement on it, and manages that entitlement's App
// Store settings. GROUP_ID adds a group override, which is then omitted to
// go back to inheriting the app's groups, and LUMOS_TEST_TBA_LABEL adds a
// time-based access override. The app is left unlisted, as
// TestAccAppStoreListing leaves it.
func TestAccEntitlementAppStoreSettings(t *testing.T) {
	appID := testAccOptionalEnv(t, "LUMOS_TEST_APPSTORE_APP_ID")
	groupID := os.Getenv("GROUP_ID")
	tbaLabel := os.Getenv("LUMOS_TEST_TBA_LABEL")
	label := "tf-acc-" + acctest.RandString(8)

	config := func(label string, visible bool, groups []string) string {
		override := ""
		if groups != nil {
			quoted := make([]string, len(groups))
			for i, g := range groups {
				quoted[i] = fmt.Sprintf("%q", g)
			}
			override = fmt.Sprintf("allowed_requester_group_uuids_override = [%s]", strings.Join(quoted, ", "))
		}
		if tbaLabel != "" {
			override += fmt.Sprintf("\n  time_based_access_override = {\n    options = [%q]\n  }", tbaLabel)
		}
		return providerConfig() + testAccCurrentUser + fmt.Sprintf(`
resource "lumos_app_appstore_settings" "listed" {
  app_id                              = %[1]q
  in_app_store                        = true
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
  %[5]s
}

resource "lumos_entitlement" "test" {
  app_id = %[1]q
  label  = %[2]q
}

resource "lumos_entitlement_appstore_settings" "test" {
  entitlement_id      = lumos_entitlement.test.id
  visible_in_appstore = %[3]t
  %[4]s

  depends_on = [lumos_app_appstore_settings.listed]
}
`, appID, label, visible, override, testAccAdmins)
	}
	var groups []string
	if groupID != "" {
		groups = []string{groupID}
	}

	renamed := label + "-renamed"
	steps := []resource.TestStep{
		{
			Config: config(label, true, groups),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("lumos_entitlement_appstore_settings.test", "visible_in_appstore", "true"),
				resource.TestCheckResourceAttrPair("lumos_entitlement_appstore_settings.test", "entitlement_id", "lumos_entitlement.test", "id"),
			),
		},
		{
			Config: config(label, false, groups),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					expectOnlyAttributesChange("lumos_entitlement_appstore_settings.test", "visible_in_appstore"),
				},
			},
			Check: resource.TestCheckResourceAttr("lumos_entitlement_appstore_settings.test", "visible_in_appstore", "false"),
		},
		{
			// Renaming an entitlement whose app is in the App Store must plan
			// only the label, including nothing in its read-only
			// appstore_settings (GOV-4388).
			Config: config(renamed, false, groups),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					expectOnlyAttributesChange("lumos_entitlement.test", "label"),
				},
			},
			Check: resource.TestCheckResourceAttr("lumos_entitlement.test", "label", renamed),
		},
	}
	if groupID != "" {
		steps[0].Check = resource.ComposeAggregateTestCheckFunc(steps[0].Check,
			resource.TestCheckTypeSetElemAttr("lumos_entitlement_appstore_settings.test", "allowed_requester_group_uuids_override.*", groupID),
		)
		steps = append(steps, resource.TestStep{
			// Omitted, the override goes back to inheriting the app's groups.
			Config: config(renamed, false, nil),
			Check:  resource.TestCheckNoResourceAttr("lumos_entitlement_appstore_settings.test", "allowed_requester_group_uuids_override"),
		})
	}

	// See TestAccAppStoreListing for why this locks after t.Parallel.
	t.Parallel()
	testAccAppStoreAppMu.Lock()
	defer testAccAppStoreAppMu.Unlock()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNotInAppStore(appID),
		Steps:                    steps,
	})
}

// TestAccEntitlementAppStoreSettingsAppUnlisted lists the app in
// LUMOS_TEST_APPSTORE_APP_ID, manages an entitlement's App Store settings,
// then unlists the app. The API then answers 404 for the settings, which
// still exist, and refresh must fail with a message saying why instead of
// dropping them from state (GOV-4617). A removed block then stops managing
// them, as the message says. The app is left unlisted, as
// TestAccAppStoreListing leaves it.
func TestAccEntitlementAppStoreSettingsAppUnlisted(t *testing.T) {
	appID := testAccOptionalEnv(t, "LUMOS_TEST_APPSTORE_APP_ID")
	label := "tf-acc-" + acctest.RandString(8)

	app := func(listed bool) string {
		return providerConfig() + fmt.Sprintf(`
resource "lumos_app_appstore_settings" "listed" {
  app_id                              = %[1]q
  in_app_store                        = %[2]t
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
}

resource "lumos_entitlement" "test" {
  app_id = %[1]q
  label  = %[3]q
}
`, appID, listed, label)
	}
	settings := `
resource "lumos_entitlement_appstore_settings" "test" {
  entitlement_id      = lumos_entitlement.test.id
  visible_in_appstore = true

  depends_on = [lumos_app_appstore_settings.listed]
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNotInAppStore(appID),
		Steps: []resource.TestStep{
			{
				Config: app(true) + settings,
				Check:  resource.TestCheckResourceAttr("lumos_entitlement_appstore_settings.test", "visible_in_appstore", "true"),
			},
			{
				// The apply unlists the app; the plan after it refreshes the
				// settings.
				Config:      app(false) + settings,
				ExpectError: regexp.MustCompile(`isn't in the App Store: list\s+the\s+app\s+again`),
			},
			{
				Config: app(false) + `
removed {
  from = lumos_entitlement_appstore_settings.test

  lifecycle {
    destroy = false
  }
}
`,
			},
		},
	})
}
