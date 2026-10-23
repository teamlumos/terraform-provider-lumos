package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
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
			"status":                          "approved",
			"disconnected":                    false,
			"sync_status":                     "success",
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
					current[sub] = nilIfEmpty(subValue)
				}
				continue
			}
			f.settings[k] = nilIfEmpty(v)
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

// nilIfEmpty stores an empty list as null, since the API reads one back that
// way.
func nilIfEmpty(v any) any {
	if list, ok := v.([]any); ok && len(list) == 0 {
		return nil
	}
	return v
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
  app_id                              = "app1"
  in_app_store                        = true
  catalog_visibility                  = %q
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
  allowed_requester_group_uuids       = ["group1"]
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
				Config: settings("search-only", ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_app_appstore_settings.slack", plancheck.ResourceActionUpdate),
						expectOnlyAttributesChange("lumos_app_appstore_settings.slack", "catalog_visibility", "request_validation_webhook_uuid"),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_app_appstore_settings.slack", "catalog_visibility", "search-only"),
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
	if got := api.patches[0]["catalog_visibility"]; got != "search-only" {
		t.Errorf("update sent catalog_visibility %v, want search-only", got)
	}
	if got, sent := api.patches[0]["request_validation_webhook_uuid"]; !sent || got != nil {
		t.Errorf("update sent request_validation_webhook_uuid %v (sent: %t), want null", got, sent)
	}
	if got, want := api.patches[1], (map[string]any{"in_app_store": false}); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("destroy sent %v, want %v", got, want)
	}
}

// TestAppStoreSettingsInsertIntakeField adds an intake field without a
// description in front of one with a description. The new field must be sent
// without one, not with the description of the field that had its position
// before. It uses a fake API so the check covers only the request the
// provider builds.
func TestAppStoreSettingsInsertIntakeField(t *testing.T) {
	api := newFakeAPI(t)
	server := httptest.NewServer(api)
	defer server.Close()

	const why = `{
    reference_key = "why"
    field_type    = "text"
    prompt        = "Why do you need access?"
    is_required   = true
    description   = "Name the project."
  }`
	const ticket = `{
    reference_key = "ticket"
    field_type    = "text"
    prompt        = "Ticket link"
    is_required   = false
  }`
	settings := func(fields ...string) string {
		return fakeProviderConfig(server.URL) + fmt.Sprintf(`
resource "lumos_app_appstore_settings" "slack" {
  app_id                              = "app1"
  in_app_store                        = true
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
  custom_intake_fields                = [%s]
}
`, strings.Join(fields, ", "))
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: settings(why)},
			{
				Config: settings(ticket, why),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("lumos_app_appstore_settings.slack", "custom_intake_fields.0.description"),
					resource.TestCheckResourceAttr("lumos_app_appstore_settings.slack", "custom_intake_fields.1.description", "Name the project."),
				),
			},
		},
	})

	api.mu.Lock()
	defer api.mu.Unlock()
	if len(api.patches) < 2 {
		t.Fatalf("got %d PATCHes, want the create and the update: %v", len(api.patches), api.patches)
	}
	fields, _ := api.patches[1]["custom_intake_fields"].([]any)
	if len(fields) != 2 {
		t.Fatalf("update sent custom_intake_fields %v, want two fields", api.patches[1]["custom_intake_fields"])
	}
	if got := fields[0].(map[string]any)["description"]; got != nil {
		t.Errorf("update sent the new field with description %v, want none", got)
	}
}

// TestAppStoreSettingsRequiredSwitches checks that a settings config missing
// a setting the API can't clear fails at plan time, before any API call.
func TestAppStoreSettingsRequiredSwitches(t *testing.T) {
	switches := map[string]string{
		"in_app_store":                        "true",
		"catalog_visibility":                  `"visible"`,
		"allow_permission_selection":          "true",
		"allow_multiple_permission_selection": "false",
	}
	configs := map[string]string{
		"visible_in_appstore": `
resource "lumos_entitlement_appstore_settings" "admin" {
  entitlement_id = "ent1"
}
`,
	}
	for missing := range switches {
		var set []string
		for name, value := range switches {
			if name != missing {
				set = append(set, fmt.Sprintf("  %s = %s", name, value))
			}
		}
		configs[missing] = fmt.Sprintf(`
resource "lumos_app_appstore_settings" "slack" {
  app_id = "app1"
%s
}
`, strings.Join(set, "\n"))
	}

	for missing, config := range configs {
		t.Run(missing, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{{
					Config:      fakeProviderConfig("http://127.0.0.1:1") + config,
					PlanOnly:    true,
					ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, missing)),
				}},
			})
		})
	}
}

// TestAppStoreSettingsOmittedSettingsCleared imports settings that have
// every optional setting set, then removes them all from the config. The
// plan must show each one cleared, and the update must send only what the
// config and those cleared values say, never a value read from the API. A
// config that sets time-based access without options clears the options.
func TestAppStoreSettingsOmittedSettingsCleared(t *testing.T) {
	api := newFakeAPI(t)
	api.settings["request_instructions"] = "Ask first."
	api.settings["request_validation_webhook_uuid"] = "webhook1"
	api.settings["redirect_url"] = "https://example.com"
	api.settings["default_entitlement"] = map[string]any{"id": "ent1"}
	api.settings["custom_intake_fields"] = []any{map[string]any{
		"reference_key": "why", "field_type": "text", "prompt": "Why?", "is_required": true,
		"description": nil, "dropdown_options": nil,
	}}
	server := httptest.NewServer(api)
	defer server.Close()

	settings := func(optional string) string {
		return fakeProviderConfig(server.URL) + fmt.Sprintf(`
resource "lumos_app_appstore_settings" "slack" {
  app_id                              = "app1"
  in_app_store                        = true
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
  %s
}
`, optional)
	}
	const addr = "lumos_app_appstore_settings.slack"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: settings(`
  allowed_requester_group_uuids   = ["group1"]
  request_instructions            = "Ask first."
  request_validation_webhook_uuid = "webhook1"
  redirect_url                    = "https://example.com"
  default_entitlement             = { id = "ent1" }
  custom_intake_fields = [{
    reference_key = "why"
    field_type    = "text"
    prompt        = "Why?"
    is_required   = true
  }]
  time_based_access = {
    options = ["1 day"]
    default = "1 day"
  }
`) + `
import {
  to = lumos_app_appstore_settings.slack
  id = "app1"
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(addr, plancheck.ResourceActionNoop),
					},
				},
			},
			{
				Config: settings(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(addr, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("allowed_requester_group_uuids"), knownvalue.SetExact([]knownvalue.Check{})),
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("custom_intake_fields"), knownvalue.ListExact([]knownvalue.Check{})),
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("time_based_access"), knownvalue.ObjectExact(map[string]knownvalue.Check{
							"options": knownvalue.ListExact([]knownvalue.Check{}),
							"default": knownvalue.Null(),
						})),
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("request_instructions"), knownvalue.Null()),
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("default_entitlement"), knownvalue.Null()),
					},
				},
			},
			{
				Config: settings(`time_based_access = { default = "1 day" }`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						expectOnlyAttributesChange(addr, "time_based_access"),
					},
				},
			},
		},
	})

	api.mu.Lock()
	defer api.mu.Unlock()
	if len(api.patches) < 2 {
		t.Fatalf("got %d PATCHes, want both updates: %v", len(api.patches), api.patches)
	}
	want := map[string]string{
		"in_app_store":                        `true`,
		"catalog_visibility":                  `"visible"`,
		"allow_permission_selection":          `true`,
		"allow_multiple_permission_selection": `false`,
		"allowed_requester_group_uuids":       `[]`,
		"custom_intake_fields":                `[]`,
		"time_based_access":                   `{"default":null,"options":[]}`,
		"request_instructions":                `null`,
		"request_validation_webhook_uuid":     `null`,
		"redirect_url":                        `null`,
		"default_entitlement":                 `null`,
	}
	for i, wantTBA := range []string{`{"default":null,"options":[]}`, `{"default":"1 day","options":[]}`} {
		want["time_based_access"] = wantTBA
		got := map[string]string{}
		for k, v := range api.patches[i] {
			encoded, _ := json.Marshal(v)
			got[k] = string(encoded)
		}
		if fmt.Sprint(got) != fmt.Sprint(want) {
			t.Errorf("update %d sent %v, want %v", i+1, got, want)
		}
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
