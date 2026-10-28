package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"slices"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

// TestAppConflictingFields checks a custom app config that sets credentials
// fails at plan time, whether it leaves app_class_id out or sets it to
// "custom": the API won't take credentials for a custom app. The check runs
// before any API call, so the server URL is never contacted.
func TestAppConflictingFields(t *testing.T) {
	for name, appClassID := range map[string]string{"omitted": "", "custom": `app_class_id = "custom"`} {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{{
					Config: fakeProviderConfig("http://127.0.0.1:1") + fmt.Sprintf(`
resource "lumos_app" "wiki" {
  %s
  auth                 = jsonencode({ api_key = "x" })
  label_override       = "Wiki"
  category_override    = "Developers"
  description_override = "Docs"
}
`, appClassID),
					PlanOnly:    true,
					ExpectError: regexp.MustCompile(`(?s)Field needs an integration.*auth only applies`),
				}},
			})
		})
	}
}

// fakeAPI records every request it receives and serves App Store settings
// for app "app1", applying PATCHes the way the API's merge patch does. app
// holds the app's display defaults under label, logo_url and so on, and
// reads back as withEffectiveDisplay shows it. It reads the app's admins back
// the way the API does, users first, each sorted by id.
type fakeAPI struct {
	t        *testing.T
	mu       sync.Mutex
	requests []string
	patches  []map[string]any
	settings map[string]any
	app      map[string]any
	admins   []any
	// appPatches holds each lumos_app update body, apart from the settings
	// PATCHes in patches.
	appPatches []map[string]any
	// adminPuts holds each admins list written.
	adminPuts [][]any
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
			"logo_url":                        nil,
			"website_url":                     nil,
			"label_override":                  nil,
			"description_override":            nil,
			"category_override":               nil,
			"logo_url_override":               nil,
			"website_url_override":            nil,
			"status":                          "approved",
			"disconnected":                    false,
			"sync_status":                     "success",
			"approval_workflow":               nil,
			"default_provisioning_workflow":   nil,
			"default_deprovisioning_workflow": nil,
			"post_approval_message":           nil,
		},
		admins: []any{},
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
			if k == "custom_intake_fields" {
				v = requiredFirst(v)
			}
			f.settings[k] = nilIfEmpty(v)
		}
		json.NewEncoder(w).Encode(f.settings)
	case r.URL.Path == "/v1/apps/app1" && (r.Method == http.MethodGet || r.Method == http.MethodPatch):
		if r.Method == http.MethodPatch {
			var patch map[string]any
			if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			f.appPatches = append(f.appPatches, patch)
			// Like the API, settings or version need auth, and credentials
			// reconnect the integration and start a sync.
			_, hasAuth := patch["auth"]
			_, hasSettings := patch["settings"]
			_, hasVersion := patch["version"]
			if (hasSettings || hasVersion) && !hasAuth {
				http.Error(w, "auth is required with settings or version", http.StatusBadRequest)
				return
			}
			if hasAuth {
				f.app["disconnected"] = false
				f.app["sync_status"] = "syncing"
			}
			// The overrides change; credentials are never read back.
			for _, k := range appDisplayOverrides {
				if v, ok := patch[k]; ok {
					f.app[k] = v
				}
			}
			if message, ok := patch["post_approval_message"]; ok {
				f.app["post_approval_message"] = message
			}
		}
		app := withEffectiveDisplay(f.app)
		app["appstore_settings"] = f.settings
		json.NewEncoder(w).Encode(app)
	case r.URL.Path == "/v1/apps/app1/admins" && (r.Method == http.MethodGet || r.Method == http.MethodPut):
		if r.Method == http.MethodPut {
			var body struct {
				Admins []any `json:"admins"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Admins == nil {
				http.Error(w, fmt.Sprintf("want an admins list: %v", err), http.StatusUnprocessableEntity)
				return
			}
			f.adminPuts = append(f.adminPuts, body.Admins)
			f.admins = sortedAdmins(body.Admins)
		}
		json.NewEncoder(w).Encode(map[string]any{"admins": f.admins})
	default:
		f.t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		http.Error(w, "unexpected request", http.StatusNotFound)
	}
}

// appDisplayOverrides maps each display field to the override the API shows
// in it where one is set.
var appDisplayOverrides = map[string]string{
	"label":       "label_override",
	"description": "description_override",
	"category":    "category_override",
	"logo_url":    "logo_url_override",
	"website_url": "website_url_override",
}

// withEffectiveDisplay copies app, which holds the app's display defaults
// under label, logo_url and so on, and reads each display field as the API
// does: its override where one is set, otherwise the default.
func withEffectiveDisplay(app map[string]any) map[string]any {
	read := map[string]any{}
	for k, v := range app {
		read[k] = v
	}
	for field, override := range appDisplayOverrides {
		if v := app[override]; v != nil {
			read[field] = v
		}
	}
	return read
}

// sortedAdmins orders admins as the API reads them back: users first, each
// sorted by id.
func sortedAdmins(admins []any) []any {
	sorted := append([]any{}, admins...)
	sort.Slice(sorted, func(i, j int) bool {
		a, b := sorted[i].(map[string]any), sorted[j].(map[string]any)
		if a["type"] != b["type"] {
			return a["type"] == "user"
		}
		return a["id"].(string) < b["id"].(string)
	})
	return sorted
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

// requiredFirst reorders written intake fields required-first, since the API
// reads them back that way whatever order they were written in.
func requiredFirst(v any) any {
	fields, ok := v.([]any)
	if !ok {
		return v
	}
	sorted := slices.Clone(fields)
	slices.SortStableFunc(sorted, func(a, b any) int {
		aRequired, _ := a.(map[string]any)["is_required"].(bool)
		bRequired, _ := b.(map[string]any)["is_required"].(bool)
		switch {
		case aRequired == bRequired:
			return 0
		case aRequired:
			return -1
		}
		return 1
	})
	return sorted
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

// TestAppStoreSettingsInsertIntakeField adds an optional intake field without
// a description in front of a required one with a description, then edits the
// new field's prompt. The new field must be sent without a description, the
// API reading the fields back required-first must not plan a change, and the
// edit must change only the intake fields. It uses a fake API so the check
// covers only what the provider plans and sends.
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
	ticket := func(prompt string) string {
		return fmt.Sprintf(`{
    reference_key = "ticket"
    field_type    = "text"
    prompt        = %q
    is_required   = false
  }`, prompt)
	}
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

	const addr = "lumos_app_appstore_settings.slack"

	// Each step's post-apply plan must be empty, so the fields read back in
	// the API's order must not plan a change.
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: settings(why)},
			{
				Config: settings(ticket("Ticket link"), why),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs(addr, "custom_intake_fields.*", map[string]string{
						"reference_key": "ticket",
						"prompt":        "Ticket link",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(addr, "custom_intake_fields.*", map[string]string{
						"reference_key": "why",
						"description":   "Name the project.",
					}),
				),
			},
			{
				Config: settings(ticket("Link to the ticket"), why),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						expectOnlyAttributesChange(addr, "custom_intake_fields"),
					},
				},
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
	for _, f := range fields {
		field := f.(map[string]any)
		if field["reference_key"] == "ticket" && field["description"] != nil {
			t.Errorf("update sent the new field with description %v, want none", field["description"])
		}
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
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("custom_intake_fields"), knownvalue.SetExact([]knownvalue.Check{})),
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

// TestAppStoreSettingsAdmins creates settings with admins, reorders them,
// changes them, then removes them from the config. The admins must be
// written before the settings, since the API lists an app only once it has
// an admin; the API's order must never show as a change; and removing them
// must send an empty list. Destroy must leave them as they are.
func TestAppStoreSettingsAdmins(t *testing.T) {
	api := newFakeAPI(t)
	server := httptest.NewServer(api)
	defer server.Close()

	settings := func(inAppStore bool, admins ...string) string {
		list := ""
		if admins != nil {
			list = fmt.Sprintf("admins = [%s]", strings.Join(admins, ", "))
		}
		return fakeProviderConfig(server.URL) + fmt.Sprintf(`
resource "lumos_app_appstore_settings" "slack" {
  app_id                              = "app1"
  in_app_store                        = %t
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
  %s
}
`, inAppStore, list)
	}
	const addr = "lumos_app_appstore_settings.slack"
	const group = `{ id = "group1", type = "group" }`
	const userA = `{ id = "user-a", type = "user" }`
	const userB = `{ id = "user-b", type = "user" }`

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: settings(true, group, userB),
				Check:  resource.TestCheckResourceAttr(addr, "admins.#", "2"),
			},
			{
				// The API reads users first; the config lists the group first.
				Config: settings(true, userB, group),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(addr, plancheck.ResourceActionNoop),
					},
				},
			},
			{
				Config: settings(true, group, userB, userA),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						expectOnlyAttributesChange(addr, "admins"),
					},
				},
			},
			{
				// The real API keeps a listed app only while it has an admin,
				// so the config that removes them also unlists the app.
				Config: settings(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						expectOnlyAttributesChange(addr, "in_app_store", "admins"),
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("admins"), knownvalue.SetExact([]knownvalue.Check{})),
					},
				},
				Check: resource.TestCheckResourceAttr(addr, "admins.#", "0"),
			},
		},
	})

	api.mu.Lock()
	defer api.mu.Unlock()
	wantAdmins := []string{
		`[{"id":"group1","type":"group"},{"id":"user-b","type":"user"}]`,
		`[{"id":"group1","type":"group"},{"id":"user-a","type":"user"},{"id":"user-b","type":"user"}]`,
		`[]`,
	}
	var gotAdmins []string
	for _, put := range api.adminPuts {
		encoded, _ := json.Marshal(sortedByID(put))
		gotAdmins = append(gotAdmins, string(encoded))
	}
	if fmt.Sprint(gotAdmins) != fmt.Sprint(wantAdmins) {
		t.Errorf("admins writes sent %v, want %v", gotAdmins, wantAdmins)
	}
	// Each create or update writes the admins, then the settings; destroy
	// only unlists the app.
	var writes []string
	for _, r := range api.requests {
		if !strings.HasPrefix(r, "GET ") {
			writes = append(writes, r)
		}
	}
	var wantWrites []string
	for range 3 {
		wantWrites = append(wantWrites, "PUT /v1/apps/app1/admins", "PATCH /v1/apps/app1/appstore-settings")
	}
	wantWrites = append(wantWrites, "PATCH /v1/apps/app1/appstore-settings")
	if fmt.Sprint(writes) != fmt.Sprint(wantWrites) {
		t.Errorf("writes were %v, want %v", writes, wantWrites)
	}
}

// sortedByID orders admins by id alone, so a check doesn't depend on the
// order a set is sent in.
func sortedByID(admins []any) []any {
	sorted := append([]any{}, admins...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].(map[string]any)["id"].(string) < sorted[j].(map[string]any)["id"].(string)
	})
	return sorted
}

// TestAppStoreSettingsImportUIAdmins imports settings whose admins were set
// in the Lumos UI. A config that lists the same admins, in any order, must
// plan no change and write nothing. Leaving them out then plans to clear
// them.
func TestAppStoreSettingsImportUIAdmins(t *testing.T) {
	api := newFakeAPI(t)
	api.admins = sortedAdmins([]any{
		map[string]any{"id": "user-b", "type": "user"},
		map[string]any{"id": "group1", "type": "group"},
		map[string]any{"id": "user-a", "type": "user"},
	})
	server := httptest.NewServer(api)
	defer server.Close()

	settings := func(inAppStore bool, optional string) string {
		return fakeProviderConfig(server.URL) + fmt.Sprintf(`
resource "lumos_app_appstore_settings" "slack" {
  app_id                              = "app1"
  in_app_store                        = %t
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
  allowed_requester_group_uuids       = ["group1"]
  time_based_access = {
    options = ["1 day"]
    default = "1 day"
  }
  %s
}
`, inAppStore, optional)
	}
	const addr = "lumos_app_appstore_settings.slack"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: settings(true, `
  admins = [
    { id = "group1", type = "group" },
    { id = "user-a", type = "user" },
    { id = "user-b", type = "user" },
  ]
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
				Check: func(*terraform.State) error {
					api.mu.Lock()
					defer api.mu.Unlock()
					if len(api.adminPuts) != 0 {
						return fmt.Errorf("import wrote admins %v, want none", api.adminPuts)
					}
					return nil
				},
			},
			{
				// Lumos keeps a listed app only while it has an admin.
				Config: settings(false, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						expectOnlyAttributesChange(addr, "in_app_store", "admins"),
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("admins"), knownvalue.SetExact([]knownvalue.Check{})),
					},
				},
			},
		},
	})
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

// TestAppEmptyStringsRejected checks that "" fails at plan time, with a
// message to leave the field out, on the fields the API doesn't keep "" for.
// The check runs before any API call, so the server URL is never contacted.
func TestAppEmptyStringsRejected(t *testing.T) {
	configs := map[string]string{
		"logo_url_override": `
resource "lumos_app" "wiki" {
  label_override       = "Wiki"
  category_override    = "Developers"
  description_override = "Team wiki"
  logo_url_override    = ""
}
`,
		"website_url_override": `
resource "lumos_app" "wiki" {
  label_override       = "Wiki"
  category_override    = "Developers"
  description_override = "Team wiki"
  website_url_override = ""
}
`,
		"request_instructions": `
resource "lumos_app_appstore_settings" "slack" {
  app_id                              = "app1"
  in_app_store                        = true
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
  request_instructions                = ""
}
`,
		"description": `
resource "lumos_app_appstore_settings" "slack" {
  app_id                              = "app1"
  in_app_store                        = true
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
  custom_intake_fields = [{
    reference_key = "why"
    field_type    = "text"
    prompt        = "Why?"
    is_required   = true
    description   = ""
  }]
}
`,
	}
	for field, config := range configs {
		t.Run(field, func(t *testing.T) {
			field := regexp.QuoteMeta(field)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{{
					Config:      fakeProviderConfig("http://127.0.0.1:1") + config,
					PlanOnly:    true,
					ExpectError: regexp.MustCompile(fmt.Sprintf(`(?s)%s can't be "".*Leave\s+%s\s+out\s+to\s+clear\s+it`, field, field)),
				}},
			})
		})
	}
}

// TestAppIntegrationRemovalRejected imports a connected app, then drops
// app_class_id from its config. There's no empty value to clear it to, so
// the plan must fail rather than turn the integration into a custom app.
func TestAppIntegrationRemovalRejected(t *testing.T) {
	api := newFakeAPI(t)
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
  app_class_id   = "slack.com"
  label_override = "Slack (Acme)"
}
`,
			},
			{
				Config: fakeProviderConfig(server.URL) + `
resource "lumos_app" "slack" {
  label_override = "Slack (Acme)"
}
`,
				ExpectError: regexp.MustCompile(`app_class_id can't be removed`),
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
}

// fakeCustomAppAPI serves custom app "custom1", applying PATCHes as the
// API's merge patch does, and records each create and PATCH body. Like the
// API, it creates the app from the label, category and description
// overrides, which become its defaults too, and reads it back as
// withEffectiveDisplay shows it.
type fakeCustomAppAPI struct {
	t       *testing.T
	mu      sync.Mutex
	app     map[string]any
	creates []map[string]any
	patches []map[string]any
}

func (f *fakeCustomAppAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.URL.Path == "/v1/apps" && r.Method == http.MethodPost:
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		f.creates = append(f.creates, body)
		f.app = map[string]any{
			"id": "custom1", "app_class_id": "custom", "instance_id": "custom_app_import",
			"status": "approved", "disconnected": false, "sync_status": nil,
			"approval_workflow": nil, "default_provisioning_workflow": nil, "default_deprovisioning_workflow": nil,
			"appstore_settings": nil,
		}
		for field, override := range appDisplayOverrides {
			f.app[field] = body[override]
			f.app[override] = body[override]
		}
		f.app["post_approval_message"] = body["post_approval_message"]
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(withEffectiveDisplay(f.app))
	case r.URL.Path == "/v1/apps/custom1" && r.Method == http.MethodGet:
		json.NewEncoder(w).Encode(withEffectiveDisplay(f.app))
	case r.URL.Path == "/v1/apps/custom1" && r.Method == http.MethodPatch:
		var patch map[string]any
		if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		f.patches = append(f.patches, patch)
		for k, v := range patch {
			f.app[k] = v
		}
		json.NewEncoder(w).Encode(withEffectiveDisplay(f.app))
	case r.URL.Path == "/v1/apps/custom1" && r.Method == http.MethodDelete:
		w.WriteHeader(http.StatusNoContent)
	default:
		f.t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		http.Error(w, "unexpected request", http.StatusNotFound)
	}
}

// TestAppCustomClearsOverrides creates a custom app with a logo and website,
// then removes both from the config. The create must plan the label it
// shows. Removing the logo and website must plan their overrides cleared,
// and the update must send null for each, the overrides the config sets,
// and no app_class_id, which the API would reject for a custom app. The
// unset post-approval message is sent as null too.
func TestAppCustomClearsOverrides(t *testing.T) {
	api := &fakeCustomAppAPI{t: t}
	server := httptest.NewServer(api)
	defer server.Close()

	config := func(extra string) string {
		return fakeProviderConfig(server.URL) + fmt.Sprintf(`
resource "lumos_app" "wiki" {
  label_override       = "Wiki"
  category_override    = "Developers"
  description_override = "Docs"
  %s
}
`, extra)
	}
	const addr = "lumos_app.wiki"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(`
  logo_url_override    = "https://example.com/logo.png"
  website_url_override = "https://example.com"
`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("label"), knownvalue.StringExact("Wiki")),
					},
				},
				Check: resource.TestCheckResourceAttr(addr, "website_url", "https://example.com"),
			},
			{
				Config: config(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						expectOnlyAttributesChange(addr, "logo_url_override", "website_url_override", "logo_url", "website_url"),
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("logo_url_override"), knownvalue.Null()),
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("website_url_override"), knownvalue.Null()),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(addr, "logo_url_override"),
					resource.TestCheckNoResourceAttr(addr, "website_url_override"),
				),
			},
		},
	})

	api.mu.Lock()
	defer api.mu.Unlock()
	if len(api.patches) != 1 {
		t.Fatalf("got %d PATCHes, want one update: %v", len(api.patches), api.patches)
	}
	want := map[string]string{
		"label_override":       `"Wiki"`,
		"category_override":    `"Developers"`,
		"description_override": `"Docs"`,
		"logo_url_override":    `null`,
		"website_url_override": `null`,

		"post_approval_message": `null`,
	}
	if got := encodeFields(api.patches[0]); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("update sent %v, want %v", got, want)
	}
}

// TestAppCustomClassID creates a custom app whose config sets app_class_id =
// "custom", the value a custom app reads back with, then drops it from the
// config and puts it back. Neither form plans a change, and the create must
// leave app_class_id out, since the API takes any app_class_id on create as
// an integration to connect.
func TestAppCustomClassID(t *testing.T) {
	api := &fakeCustomAppAPI{t: t}
	server := httptest.NewServer(api)
	defer server.Close()

	config := func(appClassID string) string {
		return fakeProviderConfig(server.URL) + fmt.Sprintf(`
resource "lumos_app" "wiki" {
  %s
  label_override       = "Wiki"
  category_override    = "Developers"
  description_override = "Docs"
}
`, appClassID)
	}
	const addr = "lumos_app.wiki"
	noop := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(addr, plancheck.ResourceActionNoop)},
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(`app_class_id = "custom"`),
				Check:  resource.TestCheckResourceAttr(addr, "app_class_id", "custom"),
			},
			{Config: config(""), ConfigPlanChecks: noop},
			{Config: config(`app_class_id = "custom"`), ConfigPlanChecks: noop},
		},
	})

	api.mu.Lock()
	defer api.mu.Unlock()
	if len(api.creates) != 1 {
		t.Fatalf("got %d creates, want 1: %v", len(api.creates), api.creates)
	}
	if v, ok := api.creates[0]["app_class_id"]; ok {
		t.Errorf("create sent app_class_id = %v; want it left out: %v", v, api.creates[0])
	}
	if len(api.patches) != 0 {
		t.Errorf("got %d updates, want none: %v", len(api.patches), api.patches)
	}
}

// TestAppCustomClassIDImported imports a custom app created in the UI, which
// has no overrides, with a config that copies app_class_id = "custom" from
// its state, and expects no change. The API requires overrides only to
// create a custom app, so the config doesn't need them.
func TestAppCustomClassIDImported(t *testing.T) {
	api := &fakeCustomAppAPI{t: t, app: map[string]any{
		"id": "custom1", "app_class_id": "custom", "instance_id": "custom_app_import",
		"label": "Wiki", "category": "Developers", "description": "Docs",
		"logo_url": nil, "website_url": nil,
		"label_override": nil, "category_override": nil, "description_override": nil,
		"logo_url_override": nil, "website_url_override": nil,
		"status": "approved", "disconnected": false, "sync_status": nil,
		"approval_workflow": nil, "default_provisioning_workflow": nil, "default_deprovisioning_workflow": nil,
		"appstore_settings": nil,
	}}
	server := httptest.NewServer(api)
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fakeProviderConfig(server.URL) + `
import {
  to = lumos_app.wiki
  id = "custom1"
}

resource "lumos_app" "wiki" {
  app_class_id = "custom"
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_app.wiki", plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})

	api.mu.Lock()
	defer api.mu.Unlock()
	if len(api.patches) != 0 {
		t.Errorf("got %d updates, want none: %v", len(api.patches), api.patches)
	}
}

// TestAppPostApprovalMessage creates a custom app with a post-approval
// message, changes it, then removes it from the config. The create must send
// it, and removing it must plan null and send null, which removes it.
func TestAppPostApprovalMessage(t *testing.T) {
	api := &fakeCustomAppAPI{t: t}
	server := httptest.NewServer(api)
	defer server.Close()

	config := func(extra string) string {
		return fakeProviderConfig(server.URL) + fmt.Sprintf(`
resource "lumos_app" "wiki" {
  label_override       = "Wiki"
  category_override    = "Developers"
  description_override = "Docs"
  %s
}
`, extra)
	}
	const addr = "lumos_app.wiki"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(`post_approval_message = "Welcome aboard."`),
				Check:  resource.TestCheckResourceAttr(addr, "post_approval_message", "Welcome aboard."),
			},
			{
				Config: config(`post_approval_message = "Ask in #it for help."`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						expectOnlyAttributesChange(addr, "post_approval_message"),
					},
				},
				Check: resource.TestCheckResourceAttr(addr, "post_approval_message", "Ask in #it for help."),
			},
			{
				Config: config(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						expectOnlyAttributesChange(addr, "post_approval_message"),
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("post_approval_message"), knownvalue.Null()),
					},
				},
				Check: resource.TestCheckNoResourceAttr(addr, "post_approval_message"),
			},
		},
	})

	api.mu.Lock()
	defer api.mu.Unlock()
	if len(api.creates) != 1 {
		t.Fatalf("got %d creates, want 1: %v", len(api.creates), api.creates)
	}
	if got := api.creates[0]["post_approval_message"]; got != "Welcome aboard." {
		t.Errorf("create sent post_approval_message %v, want %q", got, "Welcome aboard.")
	}
	if len(api.patches) != 2 {
		t.Fatalf("got %d PATCHes, want the change and the removal: %v", len(api.patches), api.patches)
	}
	for i, want := range []any{"Ask in #it for help.", nil} {
		if got, sent := api.patches[i]["post_approval_message"]; !sent || got != want {
			t.Errorf("update %d sent post_approval_message %v (sent: %t), want %v", i+1, got, sent, want)
		}
	}
}

// TestAppImportedPostApprovalMessage imports a connected app whose
// post-approval message was set in the Lumos UI. A config with the same
// message plans no change; one that leaves it out plans to remove it.
func TestAppImportedPostApprovalMessage(t *testing.T) {
	api := newFakeAPI(t)
	api.app["post_approval_message"] = "Welcome aboard."
	server := httptest.NewServer(api)
	defer server.Close()

	config := func(extra string) string {
		return fakeProviderConfig(server.URL) + fmt.Sprintf(`
resource "lumos_app" "slack" {
  app_class_id = "slack.com"
  %s
}
`, extra)
	}
	const addr = "lumos_app.slack"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(`post_approval_message = "Welcome aboard."`) + `
import {
  to = lumos_app.slack
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
				Config: config(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						expectOnlyAttributesChange(addr, "post_approval_message"),
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("post_approval_message"), knownvalue.Null()),
					},
				},
				Check: resource.TestCheckNoResourceAttr(addr, "post_approval_message"),
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

	api.mu.Lock()
	defer api.mu.Unlock()
	if len(api.appPatches) != 1 {
		t.Fatalf("got %d app updates, want the removal: %v", len(api.appPatches), api.appPatches)
	}
	want := map[string]string{
		"label_override":        `null`,
		"category_override":     `null`,
		"description_override":  `null`,
		"logo_url_override":     `null`,
		"website_url_override":  `null`,
		"post_approval_message": `null`,
	}
	if got := encodeFields(api.appPatches[0]); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("removal sent %v, want %v", got, want)
	}
}

// TestAppConnectedCredentialChanges imports a connected app with a
// post-approval message its config keeps, renames it, then changes its
// auth, settings and version in turn. Importing and renaming must keep the
// connection: those updates send no credentials, and the fields a
// reconnect recomputes keep their values in the plan. Each credential
// change must plan those fields, and the display fields without an
// override, as unknown, and send auth, settings and version together. Every
// update sends the configured message and null for each override the
// config leaves out.
func TestAppConnectedCredentialChanges(t *testing.T) {
	api := newFakeAPI(t)
	api.app["post_approval_message"] = "Welcome aboard."
	server := httptest.NewServer(api)
	defer server.Close()

	config := func(credentials string) string {
		return fakeProviderConfig(server.URL) + fmt.Sprintf(`
resource "lumos_app" "slack" {
  app_class_id          = "slack.com"
  label_override        = "Slack (Acme)"
  post_approval_message = "Welcome aboard."
  %s
}
`, credentials)
	}
	const addr = "lumos_app.slack"
	reconnect := []string{
		"sync_status", "disconnected", "instance_id",
		"description", "category", "logo_url", "website_url",
	}
	keepsConnection := []plancheck.PlanCheck{
		plancheck.ExpectKnownValue(addr, tfjsonpath.New("sync_status"), knownvalue.StringExact("success")),
		plancheck.ExpectKnownValue(addr, tfjsonpath.New("disconnected"), knownvalue.Bool(false)),
		plancheck.ExpectKnownValue(addr, tfjsonpath.New("instance_id"), knownvalue.StringExact("acme.slack.com")),
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// An imported app has no auth in state, and its first
				// update records the configured auth without sending it.
				Config: strings.Replace(config(`auth = jsonencode({ token = "v1" })`), `label_override        = "Slack (Acme)"`, "", 1) + `
import {
  to = lumos_app.slack
  id = "app1"
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: append([]plancheck.PlanCheck{expectOnlyAttributesChange(addr, "auth")}, keepsConnection...),
				},
			},
			{
				Config: config(`auth = jsonencode({ token = "v1" })`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: append([]plancheck.PlanCheck{expectOnlyAttributesChange(addr, "label_override", "label")}, keepsConnection...),
				},
				Check: resource.TestCheckResourceAttr(addr, "sync_status", "success"),
			},
			{
				Config: config(`auth = jsonencode({ token = "v2" })`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{expectOnlyAttributesChange(addr, append(reconnect, "auth")...)},
				},
				Check: resource.TestCheckResourceAttr(addr, "sync_status", "syncing"),
			},
			{
				Config: config(`
  auth     = jsonencode({ token = "v2" })
  settings = jsonencode({ region = "us" })`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{expectOnlyAttributesChange(addr, append(reconnect, "settings")...)},
				},
			},
			{
				Config: config(`
  auth     = jsonencode({ token = "v2" })
  settings = jsonencode({ region = "us" })
  version  = "2"`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{expectOnlyAttributesChange(addr, append(reconnect, "version")...)},
				},
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

	api.mu.Lock()
	defer api.mu.Unlock()
	if len(api.appPatches) != 5 {
		t.Fatalf("got %d app updates, want 5: %v", len(api.appPatches), api.appPatches)
	}
	credentials := []map[string]string{
		{},
		{},
		{"auth": `{"token":"v2"}`},
		{"auth": `{"token":"v2"}`, "settings": `{"region":"us"}`},
		{"auth": `{"token":"v2"}`, "settings": `{"region":"us"}`, "version": `"2"`},
	}
	for i, sent := range credentials {
		want := map[string]string{
			"label_override":       `"Slack (Acme)"`,
			"category_override":    `null`,
			"description_override": `null`,
			"logo_url_override":    `null`,
			"website_url_override": `null`,

			"post_approval_message": `"Welcome aboard."`,
		}
		if i == 0 {
			want["label_override"] = `null`
		}
		for k, v := range sent {
			want[k] = v
		}
		if got := encodeFields(api.appPatches[i]); fmt.Sprint(got) != fmt.Sprint(want) {
			t.Errorf("update %d sent %v, want %v", i+1, got, want)
		}
	}
}

// TestAppDisconnectedUpdateReconnects imports an app disconnected in the
// Lumos UI and renames it. An update of a disconnected app sends its
// credentials, reconnecting it, so the plan changes disconnected to false
// and leaves the other fields a reconnect recomputes unknown.
func TestAppDisconnectedUpdateReconnects(t *testing.T) {
	api := newFakeAPI(t)
	api.app["disconnected"] = true
	api.app["sync_status"] = nil
	server := httptest.NewServer(api)
	defer server.Close()

	const addr = "lumos_app.slack"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fakeProviderConfig(server.URL) + `
resource "lumos_app" "slack" {
  app_class_id   = "slack.com"
  auth           = jsonencode({ token = "v1" })
  label_override = "Slack (Acme)"
}

import {
  to = lumos_app.slack
  id = "app1"
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("disconnected"), knownvalue.Bool(false)),
						plancheck.ExpectUnknownValue(addr, tfjsonpath.New("sync_status")),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "disconnected", "false"),
					resource.TestCheckResourceAttr(addr, "sync_status", "syncing"),
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

	api.mu.Lock()
	defer api.mu.Unlock()
	if len(api.appPatches) != 1 {
		t.Fatalf("got %d app updates, want 1: %v", len(api.appPatches), api.appPatches)
	}
	if got := encodeFields(api.appPatches[0])["auth"]; got != `{"token":"v1"}` {
		t.Errorf("update sent auth %s, want {\"token\":\"v1\"}", got)
	}
}

// TestAppDisconnectedPlansReconnect connects an integration, then
// disconnects it outside Terraform. With no change to the config, the plan
// must change disconnected to false and leave the other fields a reconnect
// recomputes unknown, and the update must send auth and settings. Once
// reconnected, the plan is empty.
func TestAppDisconnectedPlansReconnect(t *testing.T) {
	api := newFakeAPI(t)
	server := httptest.NewServer(api)
	defer server.Close()

	const addr = "lumos_app.slack"
	config := fakeProviderConfig(server.URL) + `
resource "lumos_app" "slack" {
  app_class_id = "slack.com"
  auth         = jsonencode({ token = "v1" })
  settings     = jsonencode({ region = "us" })
}
`
	disconnect := func() {
		api.mu.Lock()
		defer api.mu.Unlock()
		api.app["disconnected"] = true
		api.app["sync_status"] = nil
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config + `
import {
  to = lumos_app.slack
  id = "app1"
}
`,
			},
			{
				PreConfig: disconnect,
				Config:    config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						expectOnlyAttributesChange(addr,
							"disconnected", "sync_status", "instance_id",
							"label", "description", "category", "logo_url", "website_url"),
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("disconnected"), knownvalue.Bool(false)),
						plancheck.ExpectUnknownValue(addr, tfjsonpath.New("sync_status")),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "disconnected", "false"),
					resource.TestCheckResourceAttr(addr, "sync_status", "syncing"),
				),
			},
			{
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
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

	api.mu.Lock()
	defer api.mu.Unlock()
	if len(api.appPatches) != 2 {
		t.Fatalf("got %d app updates, want 2: %v", len(api.appPatches), api.appPatches)
	}
	if _, ok := api.appPatches[0]["auth"]; ok {
		t.Errorf("import update sent auth: %v", api.appPatches[0])
	}
	got := encodeFields(api.appPatches[1])
	if got["auth"] != `{"token":"v1"}` || got["settings"] != `{"region":"us"}` {
		t.Errorf("reconnect sent %v, want the configured auth and settings", got)
	}
}

// TestAppIgnoreAuthChanges connects an integration, then drops auth from
// the config and lists it in ignore_changes. Renaming the app must then
// plan no auth change and send no credentials. Once the app is disconnected
// outside Terraform, the plan reconnects it with the auth in state.
func TestAppIgnoreAuthChanges(t *testing.T) {
	api := newFakeAPI(t)
	server := httptest.NewServer(api)
	defer server.Close()

	const addr = "lumos_app.slack"
	config := func(body string) string {
		return fakeProviderConfig(server.URL) + `
resource "lumos_app" "slack" {
  app_class_id = "slack.com"
  settings     = jsonencode({ region = "us" })
` + body + `
}
`
	}
	ignored := config(`
  label_override = "Slack (Acme)"

  lifecycle {
    ignore_changes = [auth]
  }`)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(`  auth = jsonencode({ token = "v1" })`) + `
import {
  to = lumos_app.slack
  id = "app1"
}
`,
			},
			{
				Config: ignored,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{expectOnlyAttributesChange(addr, "label_override", "label")},
				},
			},
			{
				PreConfig: func() {
					api.mu.Lock()
					defer api.mu.Unlock()
					api.app["disconnected"] = true
				},
				Config: ignored,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("disconnected"), knownvalue.Bool(false)),
					},
				},
				Check: resource.TestCheckResourceAttr(addr, "disconnected", "false"),
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

	api.mu.Lock()
	defer api.mu.Unlock()
	if len(api.appPatches) != 3 {
		t.Fatalf("got %d app updates, want 3: %v", len(api.appPatches), api.appPatches)
	}
	for i, patch := range api.appPatches[:2] {
		if _, ok := patch["auth"]; ok {
			t.Errorf("update %d sent auth: %v", i+1, patch)
		}
	}
	if got := encodeFields(api.appPatches[2])["auth"]; got != `{"token":"v1"}` {
		t.Errorf("reconnect sent auth %s, want the auth in state", got)
	}
}

// TestAppDisconnectedCustomApp imports a custom app that reads back
// disconnected. A custom app has nothing to reconnect, so the plan is empty.
func TestAppDisconnectedCustomApp(t *testing.T) {
	api := &fakeCustomAppAPI{t: t, app: map[string]any{
		"id": "custom1", "app_class_id": "custom", "instance_id": "custom_app_import",
		"label": "Wiki", "category": "Developers", "description": "Docs",
		"logo_url": nil, "website_url": nil,
		"label_override": nil, "category_override": nil, "description_override": nil,
		"logo_url_override": nil, "website_url_override": nil,
		"status": "approved", "disconnected": true, "sync_status": nil,
		"approval_workflow": nil, "default_provisioning_workflow": nil, "default_deprovisioning_workflow": nil,
		"appstore_settings": nil,
	}}
	server := httptest.NewServer(api)
	defer server.Close()

	config := fakeProviderConfig(server.URL) + `
resource "lumos_app" "wiki" {
  app_class_id = "custom"
}
`
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config + `
import {
  to = lumos_app.wiki
  id = "custom1"
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_app.wiki", plancheck.ResourceActionNoop),
					},
				},
			},
			{
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})

	api.mu.Lock()
	defer api.mu.Unlock()
	if len(api.patches) != 0 {
		t.Errorf("got %d updates, want none: %v", len(api.patches), api.patches)
	}
}

// TestAppConnectedDisplayOverrides imports a connected app whose logo was
// set in the Lumos UI, with a config that renames it, sets a website and
// leaves the logo out. The plan must show the new name and website and the
// logo override cleared, while the category and description keep the
// integration's defaults. Removing the website from the config then clears
// it too.
func TestAppConnectedDisplayOverrides(t *testing.T) {
	api := newFakeAPI(t)
	api.app["logo_url_override"] = "https://example.com/slack.png"
	server := httptest.NewServer(api)
	defer server.Close()

	config := func(extra string) string {
		return fakeProviderConfig(server.URL) + fmt.Sprintf(`
resource "lumos_app" "slack" {
  app_class_id   = "slack.com"
  label_override = "Slack (Acme)"
  %s
}
`, extra)
	}
	const addr = "lumos_app.slack"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(`website_url_override = "https://acme.slack.com"`) + `
import {
  to = lumos_app.slack
  id = "app1"
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						expectOnlyAttributesChange(addr,
							"label_override", "logo_url_override", "website_url_override",
							"label", "logo_url", "website_url"),
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("label"), knownvalue.StringExact("Slack (Acme)")),
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("website_url"), knownvalue.StringExact("https://acme.slack.com")),
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("logo_url_override"), knownvalue.Null()),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "label", "Slack (Acme)"),
					resource.TestCheckResourceAttr(addr, "category", "Communication"),
					resource.TestCheckNoResourceAttr(addr, "logo_url"),
					resource.TestCheckResourceAttr(addr, "website_url", "https://acme.slack.com"),
				),
			},
			{
				Config: config(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						expectOnlyAttributesChange(addr, "website_url_override", "website_url"),
						plancheck.ExpectKnownValue(addr, tfjsonpath.New("website_url_override"), knownvalue.Null()),
					},
				},
				Check: resource.TestCheckNoResourceAttr(addr, "website_url"),
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

	api.mu.Lock()
	defer api.mu.Unlock()
	if len(api.appPatches) != 2 {
		t.Fatalf("got %d app updates, want 2: %v", len(api.appPatches), api.appPatches)
	}
	overrides := map[string]string{
		"label_override":       `"Slack (Acme)"`,
		"category_override":    `null`,
		"description_override": `null`,
		"logo_url_override":    `null`,

		"post_approval_message": `null`,
	}
	for i, website := range []string{`"https://acme.slack.com"`, `null`} {
		want := map[string]string{"website_url_override": website}
		for k, v := range overrides {
			want[k] = v
		}
		if got := encodeFields(api.appPatches[i]); fmt.Sprint(got) != fmt.Sprint(want) {
			t.Errorf("update %d sent %v, want %v", i+1, got, want)
		}
	}
}

// encodeFields JSON-encodes each value of a request body, for comparing
// bodies field by field.
func encodeFields(body map[string]any) map[string]string {
	encoded := map[string]string{}
	for k, v := range body {
		b, _ := json.Marshal(v)
		encoded[k] = string(b)
	}
	return encoded
}
