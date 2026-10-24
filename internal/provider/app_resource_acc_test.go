package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCustomApp(name, label, description string) string {
	return fmt.Sprintf(`
resource "lumos_app" %q {
  label       = %q
  category    = "Developers"
  description = %q
}
`, name, label, description)
}

// TestAccAppCustom creates a custom app with a logo and website, edits only
// its description, removes the logo and website, then imports it. The edit
// must plan that one field and nothing else (GOV-4388), and removing the
// logo and website must clear them.
func TestAccAppCustom(t *testing.T) {
	label := "tf-acc-" + acctest.RandString(8)
	withLinks := func(description string) string {
		return fmt.Sprintf(`
resource "lumos_app" "test" {
  label       = %q
  category    = "Developers"
  description = %q
  logo_url    = "https://example.com/tf-acc-logo.png"
  website_url = "https://example.com/tf-acc"
}
`, label, description)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + withLinks("Created by an acceptance test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_app.test", "label", label),
					resource.TestCheckResourceAttr("lumos_app.test", "app_class_id", "custom"),
					resource.TestCheckResourceAttr("lumos_app.test", "website_url", "https://example.com/tf-acc"),
					resource.TestCheckResourceAttrSet("lumos_app.test", "id"),
				),
			},
			{
				Config: providerConfig() + withLinks("Edited by an acceptance test"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_app.test", plancheck.ResourceActionUpdate),
						expectOnlyAttributesChange("lumos_app.test", "description"),
					},
				},
				Check: resource.TestCheckResourceAttr("lumos_app.test", "description", "Edited by an acceptance test"),
			},
			{
				Config: providerConfig() + testAccCustomApp("test", label, "Edited by an acceptance test"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						expectOnlyAttributesChange("lumos_app.test", "logo_url", "website_url"),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("lumos_app.test", "logo_url"),
					resource.TestCheckNoResourceAttr("lumos_app.test", "website_url"),
				),
			},
			{
				ResourceName:      "lumos_app.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccAppStoreSettings manages the App Store settings of a custom app the
// test creates. The app stays unlisted: listing needs an admin, which the
// provider can't assign, so TestAccAppStoreListing covers that with an
// existing app. LUMOS_TEST_TBA_LABEL (one of the domain's time-based access
// labels) and GROUP_ID add those settings when set. The first app's default
// entitlement is one the test creates. Removing it from the config clears it
// before the post-test destroy deletes the entitlement and the app.
func TestAccAppStoreSettings(t *testing.T) {
	prefix := "tf-acc-" + acctest.RandString(8)
	tbaLabel := os.Getenv("LUMOS_TEST_TBA_LABEL")
	groupID := os.Getenv("GROUP_ID")

	apps := testAccCustomApp("first", prefix+"-first", "App Store settings test") +
		testAccCustomApp("second", prefix+"-second", "App Store settings test") + fmt.Sprintf(`
resource "lumos_entitlement" "default" {
  app_id = lumos_app.first.id
  label  = %q
}
`, prefix+"-default")
	// full sets every optional setting the test can: an intake field, request
	// instructions, redirect URL, default entitlement and, when their
	// variables are set, requester groups and time-based access. Without it,
	// the config sets only the required switches.
	settings := func(appRef, visibility string, full bool) string {
		var extra []string
		if full {
			extra = append(extra,
				`custom_intake_fields = [{
    reference_key = "justification"
    field_type    = "text"
    prompt        = "Why do you need access?"
    is_required   = true
  }]`,
				`request_instructions = "Set by an acceptance test."`,
				`redirect_url = "https://example.com/tf-acc"`,
				"default_entitlement = {\n    id = lumos_entitlement.default.id\n  }",
			)
			if tbaLabel != "" {
				extra = append(extra, fmt.Sprintf("time_based_access = {\n    options = [%q]\n    default = %[1]q\n  }", tbaLabel))
			}
			if groupID != "" {
				extra = append(extra, fmt.Sprintf("allowed_requester_group_uuids = [%q]", groupID))
			}
		}
		return fmt.Sprintf(`
resource "lumos_app_appstore_settings" "test" {
  app_id                              = %s.id
  in_app_store                        = false
  catalog_visibility                  = %q
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
  %s
}
`, appRef, visibility, strings.Join(extra, "\n  "))
	}

	steps := []resource.TestStep{
		{
			Config: providerConfig() + apps + settings("lumos_app.first", "visible", true),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("lumos_app_appstore_settings.test", "catalog_visibility", "visible"),
				resource.TestCheckResourceAttr("lumos_app_appstore_settings.test", "custom_intake_fields.#", "1"),
				resource.TestCheckResourceAttrPair("lumos_app_appstore_settings.test", "app_id", "lumos_app.first", "id"),
				resource.TestCheckResourceAttrPair("lumos_app_appstore_settings.test", "default_entitlement.id", "lumos_entitlement.default", "id"),
				resource.TestCheckResourceAttr("lumos_app_appstore_settings.test", "redirect_url", "https://example.com/tf-acc"),
			),
		},
		{
			Config: providerConfig() + apps + settings("lumos_app.first", "search-only", true),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					expectOnlyAttributesChange("lumos_app_appstore_settings.test", "catalog_visibility"),
				},
			},
			Check: resource.TestCheckResourceAttr("lumos_app_appstore_settings.test", "catalog_visibility", "search-only"),
		},
	}
	// Removing every optional setting clears it: lists and time-based access
	// options plan empty, the rest null. The checks run after the read, so
	// they also show that the API keeps empty time-based access options
	// empty instead of filling in the organization's. Clearing all of this
	// also lets the post-test destroy delete the app: rows that still point at
	// a custom app make the API keep it.
	clearedChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("lumos_app_appstore_settings.test", "custom_intake_fields.#", "0"),
		resource.TestCheckResourceAttr("lumos_app_appstore_settings.test", "allowed_requester_group_uuids.#", "0"),
		resource.TestCheckResourceAttr("lumos_app_appstore_settings.test", "time_based_access.options.#", "0"),
		resource.TestCheckNoResourceAttr("lumos_app_appstore_settings.test", "time_based_access.default"),
		resource.TestCheckNoResourceAttr("lumos_app_appstore_settings.test", "request_instructions"),
		resource.TestCheckNoResourceAttr("lumos_app_appstore_settings.test", "redirect_url"),
		resource.TestCheckNoResourceAttr("lumos_app_appstore_settings.test", "default_entitlement.id"),
	}
	steps = append(steps, resource.TestStep{
		Config: providerConfig() + apps + settings("lumos_app.first", "search-only", false),
		ConfigPlanChecks: resource.ConfigPlanChecks{
			PreApply: []plancheck.PlanCheck{
				expectOnlyAttributesChange("lumos_app_appstore_settings.test",
					"custom_intake_fields", "allowed_requester_group_uuids", "time_based_access",
					"request_instructions", "redirect_url", "default_entitlement"),
			},
		},
		Check: resource.ComposeAggregateTestCheckFunc(clearedChecks...),
	})
	steps = append(steps, resource.TestStep{
		// Settings belong to one app: pointing the resource at another app is
		// a replace, not an update.
		Config: providerConfig() + apps + settings("lumos_app.second", "search-only", false),
		ConfigPlanChecks: resource.ConfigPlanChecks{
			PreApply: []plancheck.PlanCheck{
				plancheck.ExpectResourceAction("lumos_app_appstore_settings.test", plancheck.ResourceActionReplace),
			},
		},
		Check: resource.TestCheckResourceAttrPair("lumos_app_appstore_settings.test", "app_id", "lumos_app.second", "id"),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

// TestAccAppStoreListing lists an existing app in the App Store, then moves
// its settings resource to a custom app the test creates. The replace has to
// take the old app out of the App Store. It needs
// LUMOS_TEST_APPSTORE_APP_ID: an app with an admin, reserved for tests,
// since the test leaves it unlisted. The custom app has no admin, so it
// stays unlisted throughout.
func TestAccAppStoreListing(t *testing.T) {
	appID := testAccOptionalEnv(t, "LUMOS_TEST_APPSTORE_APP_ID")
	label := "tf-acc-" + acctest.RandString(8)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNotInAppStore(appID),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "lumos_app_appstore_settings" "listed" {
  app_id                              = %q
  in_app_store                        = true
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
}
`, appID),
				Check: resource.TestCheckResourceAttr("lumos_app_appstore_settings.listed", "in_app_store", "true"),
			},
			{
				Config: providerConfig() + testAccCustomApp("unlisted", label, "App Store listing test") + `
resource "lumos_app_appstore_settings" "listed" {
  app_id                              = lumos_app.unlisted.id
  in_app_store                        = false
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_app_appstore_settings.listed", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("lumos_app_appstore_settings.listed", "app_id", "lumos_app.unlisted", "id"),
					testAccCheckNotInAppStore(appID),
				),
			},
		},
	})
}

// testAccCheckNotInAppStore reads the app's settings straight from the API,
// since the resource is gone from state by the time CheckDestroy runs.
func testAccCheckNotInAppStore(appID string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		req, err := http.NewRequest(http.MethodGet, strings.TrimSuffix(os.Getenv("LUMOS_API_URL"), "/")+"/v1/apps/"+appID+"/appstore-settings", nil)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+os.Getenv("LUMOS_API_TOKEN"))
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("reading app %s settings: status %d", appID, res.StatusCode)
		}
		var settings struct {
			InAppStore bool `json:"in_app_store"`
		}
		if err := json.NewDecoder(res.Body).Decode(&settings); err != nil {
			return err
		}
		if settings.InAppStore {
			return fmt.Errorf("app %s is still in the App Store after its settings were destroyed", appID)
		}
		return nil
	}
}

// TestAccAppStoreSettingsIgnoreChanges sets request instructions, standing
// in for an app owner setting them in the UI, then drops them from the config
// with ignore_changes and makes an unrelated change. The update must keep
// the instructions. Dropping the lifecycle block then clears them.
func TestAccAppStoreSettingsIgnoreChanges(t *testing.T) {
	const addr = "lumos_app_appstore_settings.test"
	app := testAccCustomApp("first", "tf-acc-"+acctest.RandString(8), "ignore_changes test")
	settings := func(multiple bool, extra string) string {
		return fmt.Sprintf(`
resource "lumos_app_appstore_settings" "test" {
  app_id                              = lumos_app.first.id
  in_app_store                        = false
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = %t
  %s
}
`, multiple, extra)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + app + settings(false, `request_instructions = "Set outside Terraform."`),
			},
			{
				Config: providerConfig() + app + settings(true, "lifecycle {\n    ignore_changes = [request_instructions]\n  }"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						expectOnlyAttributesChange(addr, "allow_multiple_permission_selection"),
					},
				},
				Check: resource.TestCheckResourceAttr(addr, "request_instructions", "Set outside Terraform."),
			},
			{
				Config: providerConfig() + app + settings(true, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						expectOnlyAttributesChange(addr, "request_instructions"),
					},
				},
				Check: resource.TestCheckNoResourceAttr(addr, "request_instructions"),
			},
		},
	})
}
