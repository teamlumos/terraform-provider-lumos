package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func fakeProviderConfig(serverURL string) string {
	return fmt.Sprintf(`
provider "lumos" {
  server_url  = %q
  http_bearer = "test-token"
}
`, serverURL)
}

// TestAppConflictingFields checks a config that both connects an
// integration and sets custom-app fields fails at plan time. The check runs
// before any API call, so the server URL is never contacted.
func TestAppConflictingFields(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: fakeProviderConfig("http://127.0.0.1:1") + `
resource "lumos_app" "mixed" {
  app_class_id = "okta.com"
  auth         = jsonencode({ api_key = "x" })
  description  = "Our Okta tenant"
}
`,
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Combination.*app_class_id`),
		}},
	})
}

// fakeAPI records every request it receives and serves App Store settings
// for app "app1", applying PATCHes the way the API's merge patch does.
type fakeAPI struct {
	t        *testing.T
	mu       sync.Mutex
	requests []string
	patches  []map[string]any
	settings map[string]any
	app      map[string]any
}

func newFakeAPI(t *testing.T) *fakeAPI {
	settings := map[string]any{
		"in_app_store":                        true,
		"catalog_visibility":                  "visible",
		"allow_permission_selection":          true,
		"allow_multiple_permission_selection": false,
		"allowed_requester_group_uuids":       []any{"group1"},
		"request_instructions":                nil,
		"request_validation_webhook_uuid":     nil,
		"redirect_url":                        nil,
		"default_entitlement":                 nil,
		"custom_intake_fields":                nil,
		// The API returns options as objects; the provider shows labels.
		"time_based_access": map[string]any{
			"options": []any{map[string]any{"label": "1 day", "value_seconds": 86400}},
			"default": "1 day",
		},
	}
	return &fakeAPI{
		t:        t,
		settings: settings,
		app: map[string]any{
			"id":                              "app1",
			"label":                           "Slack",
			"app_class_id":                    "slack.com",
			"instance_id":                     "acme.slack.com",
			"category":                        "Communication",
			"description":                     "Team chat",
			"status":                          "APPROVED",
			"disconnected":                    false,
			"sync_status":                     "SUCCESS",
			"approval_workflow":               nil,
			"default_provisioning_workflow":   nil,
			"default_deprovisioning_workflow": nil,
		},
	}
}

func (f *fakeAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.requests = append(f.requests, r.Method+" "+r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.URL.Path == "/v1/apps/app1/appstore-settings" && r.Method == http.MethodGet:
		json.NewEncoder(w).Encode(f.settings)
	case r.URL.Path == "/v1/apps/app1/appstore-settings" && r.Method == http.MethodPatch:
		body, _ := io.ReadAll(r.Body)
		var patch map[string]any
		if err := json.Unmarshal(body, &patch); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		f.patches = append(f.patches, patch)
		for k, v := range patch {
			if k == "time_based_access" {
				// Like the API, only the sub-fields present change.
				current := f.settings[k].(map[string]any)
				for sub, subValue := range withOptionObjects(v).(map[string]any) {
					current[sub] = subValue
				}
				continue
			}
			f.settings[k] = v
		}
		json.NewEncoder(w).Encode(f.settings)
	case r.URL.Path == "/v1/apps/app1" && r.Method == http.MethodGet:
		app := map[string]any{"appstore_settings": f.settings}
		for k, v := range f.app {
			app[k] = v
		}
		json.NewEncoder(w).Encode(app)
	default:
		f.t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		http.Error(w, "unexpected request", http.StatusNotFound)
	}
}

// withOptionObjects turns written option labels into the objects the API
// reads back.
func withOptionObjects(v any) any {
	tba, ok := v.(map[string]any)
	if !ok {
		return v
	}
	if labels, ok := tba["options"].([]any); ok {
		options := make([]any, 0, len(labels))
		for _, l := range labels {
			options = append(options, map[string]any{"label": l, "value_seconds": nil})
		}
		tba["options"] = options
	}
	return tba
}

func (f *fakeAPI) sawRequest(want string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, r := range f.requests {
		if r == want {
			return true
		}
	}
	return false
}

// TestAppStoreSettingsOnImportedApp manages the App Store settings of an app
// Terraform doesn't own: import, a change that also removes the request
// validation webhook, then destroy. It uses a fake API because a
// UI-connected app with known settings and a validation webhook can't be set
// up on demand in a shared test domain.
func TestAppStoreSettingsOnImportedApp(t *testing.T) {
	api := newFakeAPI(t)
	api.settings["request_validation_webhook_uuid"] = "webhook1"
	server := httptest.NewServer(api)
	defer server.Close()

	settings := func(visibility, webhook string) string {
		return fakeProviderConfig(server.URL) + fmt.Sprintf(`
resource "lumos_app_appstore_settings" "slack" {
  app_id             = "app1"
  in_app_store       = true
  catalog_visibility = %q
  time_based_access = {
    options = ["1 day"]
    default = "1 day"
  }
  %s
}
`, visibility, webhook)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Matching config: the import plans no change, which also
				// shows time-based access options read back as labels.
				Config: settings("visible", `request_validation_webhook_uuid = "webhook1"`) + `
import {
  to = lumos_app_appstore_settings.slack
  id = "app1"
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_app_appstore_settings.slack", plancheck.ResourceActionNoop),
					},
				},
			},
			{
				Config: settings("search_only", ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_app_appstore_settings.slack", plancheck.ResourceActionUpdate),
						expectOnlyAttributesChange("lumos_app_appstore_settings.slack", "catalog_visibility", "request_validation_webhook_uuid"),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_app_appstore_settings.slack", "catalog_visibility", "search_only"),
					resource.TestCheckNoResourceAttr("lumos_app_appstore_settings.slack", "request_validation_webhook_uuid"),
				),
			},
		},
	})

	api.mu.Lock()
	defer api.mu.Unlock()
	if len(api.patches) != 2 {
		t.Fatalf("got %d PATCHes, want the update and the destroy: %v", len(api.patches), api.patches)
	}
	if got := api.patches[0]["catalog_visibility"]; got != "search_only" {
		t.Errorf("update sent catalog_visibility %v, want search_only", got)
	}
	if got, sent := api.patches[0]["request_validation_webhook_uuid"]; !sent || got != nil {
		t.Errorf("update sent request_validation_webhook_uuid %v (sent: %t), want null", got, sent)
	}
	// Settings the config leaves out are sent back with their current value.
	if got := fmt.Sprint(api.patches[0]["allowed_requester_group_uuids"]); got != "[group1]" {
		t.Errorf("update sent allowed_requester_group_uuids %s, want [group1]", got)
	}
	if got, want := api.patches[1], (map[string]any{"in_app_store": false}); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("destroy sent %v, want %v", got, want)
	}
}

// TestAppImportedFromUIRemovedWithoutDelete imports an app connected in the
// Lumos UI, then stops managing it with a removed block. The API refuses to
// delete such an app, so the documented way out must never call DELETE. It
// uses a fake API for the same reason as TestAppStoreSettingsOnImportedApp.
// It runs for a connected and a disconnected app: either way, the import
// plans no change.
func TestAppImportedFromUIRemovedWithoutDelete(t *testing.T) {
	for _, disconnected := range []bool{false, true} {
		t.Run(fmt.Sprintf("disconnected=%t", disconnected), func(t *testing.T) {
			testAppImportedFromUIRemovedWithoutDelete(t, disconnected)
		})
	}
}

func testAppImportedFromUIRemovedWithoutDelete(t *testing.T, disconnected bool) {
	api := newFakeAPI(t)
	api.app["disconnected"] = disconnected
	server := httptest.NewServer(api)
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fakeProviderConfig(server.URL) + `
import {
  to = lumos_app.slack
  id = "app1"
}

resource "lumos_app" "slack" {
  app_class_id = "slack.com"
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_app.slack", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_app.slack", "appstore_settings.in_app_store", "true"),
					resource.TestCheckResourceAttr("lumos_app.slack", "disconnected", fmt.Sprint(disconnected)),
				),
			},
			{
				Config: fakeProviderConfig(server.URL) + `
removed {
  from = lumos_app.slack

  lifecycle {
    destroy = false
  }
}
`,
			},
		},
	})

	if api.sawRequest("DELETE /v1/apps/app1") {
		t.Errorf("removing lumos_app.slack with destroy = false called DELETE")
	}
}
