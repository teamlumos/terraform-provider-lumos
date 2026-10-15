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

// TestAccAppCustom creates a custom app, edits only its description, then
// imports it. The edit must plan that one field and nothing else
// (GOV-4388).
func TestAccAppCustom(t *testing.T) {
	label := "tf-acc-" + acctest.RandString(8)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + testAccCustomApp("test", label, "Created by an acceptance test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_app.test", "label", label),
					resource.TestCheckResourceAttr("lumos_app.test", "app_class_id", "custom"),
					resource.TestCheckResourceAttrSet("lumos_app.test", "id"),
				),
			},
			{
				Config: providerConfig() + testAccCustomApp("test", label, "Edited by an acceptance test"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_app.test", plancheck.ResourceActionUpdate),
						expectOnlyAttributesChange("lumos_app.test", "description"),
					},
				},
				Check: resource.TestCheckResourceAttr("lumos_app.test", "description", "Edited by an acceptance test"),
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
// labels) and GROUP_ID add those settings when set.
func TestAccAppStoreSettings(t *testing.T) {
	prefix := "tf-acc-" + acctest.RandString(8)
	tbaLabel := os.Getenv("LUMOS_TEST_TBA_LABEL")
	groupID := os.Getenv("GROUP_ID")

	apps := testAccCustomApp("first", prefix+"-first", "App Store settings test") +
		testAccCustomApp("second", prefix+"-second", "App Store settings test")
	settings := func(appRef, visibility string, groups []string, intakeField bool) string {
		intake := "[]"
		if intakeField {
			intake = `[{
    reference_key = "justification"
    field_type    = "text"
    prompt        = "Why do you need access?"
    is_required   = true
  }]`
		}
		var extra []string
		if tbaLabel != "" {
			extra = append(extra, fmt.Sprintf("time_based_access = {\n    options = [%q]\n    default = %q\n  }", tbaLabel, tbaLabel))
		}
		if groups != nil {
			quoted := make([]string, len(groups))
			for i, g := range groups {
				quoted[i] = fmt.Sprintf("%q", g)
			}
			extra = append(extra, fmt.Sprintf("allowed_requester_group_uuids = [%s]", strings.Join(quoted, ", ")))
		}
		return fmt.Sprintf(`
resource "lumos_app_appstore_settings" "test" {
  app_id               = %s.id
  in_app_store         = false
  catalog_visibility   = %q
  request_instructions = "Set by an acceptance test."
  custom_intake_fields = %s
  %s
}
`, appRef, visibility, intake, strings.Join(extra, "\n  "))
	}
	var groups []string
	if groupID != "" {
		groups = []string{groupID}
	}

	steps := []resource.TestStep{
		{
			Config: providerConfig() + apps + settings("lumos_app.first", "visible", groups, true),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("lumos_app_appstore_settings.test", "catalog_visibility", "visible"),
				resource.TestCheckResourceAttr("lumos_app_appstore_settings.test", "custom_intake_fields.#", "1"),
				resource.TestCheckResourceAttrPair("lumos_app_appstore_settings.test", "app_id", "lumos_app.first", "id"),
			),
		},
		{
			Config: providerConfig() + apps + settings("lumos_app.first", "search_only", groups, true),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					expectOnlyAttributesChange("lumos_app_appstore_settings.test", "catalog_visibility"),
				},
			},
			Check: resource.TestCheckResourceAttr("lumos_app_appstore_settings.test", "catalog_visibility", "search_only"),
		},
	}
	// Empty lists are sent, so they clear the intake field and, when set, the
	// requester groups. Clearing them also lets the post-test destroy delete
	// the app: rows that still point at a custom app make the API keep it.
	var noGroups []string
	if groupID != "" {
		noGroups = []string{}
	}
	clearedChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("lumos_app_appstore_settings.test", "custom_intake_fields.#", "0"),
	}
	if groupID != "" {
		clearedChecks = append(clearedChecks, resource.TestCheckResourceAttr("lumos_app_appstore_settings.test", "allowed_requester_group_uuids.#", "0"))
	}
	steps = append(steps, resource.TestStep{
		Config: providerConfig() + apps + settings("lumos_app.first", "search_only", noGroups, false),
		Check:  resource.ComposeAggregateTestCheckFunc(clearedChecks...),
	})
	steps = append(steps, resource.TestStep{
		// Settings belong to one app: pointing the resource at another app is
		// a replace, not an update.
		Config: providerConfig() + apps + settings("lumos_app.second", "search_only", noGroups, false),
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
  app_id       = %q
  in_app_store = true
}
`, appID),
				Check: resource.TestCheckResourceAttr("lumos_app_appstore_settings.listed", "in_app_store", "true"),
			},
			{
				Config: providerConfig() + testAccCustomApp("unlisted", label, "App Store listing test") + `
resource "lumos_app_appstore_settings" "listed" {
  app_id       = lumos_app.unlisted.id
  in_app_store = false
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
