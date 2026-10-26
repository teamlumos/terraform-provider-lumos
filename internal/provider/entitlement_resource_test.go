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
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// fakeEntitlementAPI serves entitlement "ent1" and its App Store settings
// in the API's shapes, applying settings PATCHes as a merge patch. It keeps
// an empty group list as sent, so a test can check that the provider holds
// [] and null apart. Like the API, it answers 404 for the settings while
// the entitlement's app "app1" isn't in the App Store, and serves that
// app's settings for in_app_store.
type fakeEntitlementAPI struct {
	t            *testing.T
	mu           sync.Mutex
	requests     []string
	patches      []map[string]any
	entitlement  map[string]any
	settings     map[string]any
	deleteStatus int
	deleteBody   string

	appUnlisted     bool
	entitlementGone bool
	settingsGone    bool
}

func newFakeEntitlementAPI(t *testing.T, source string) *fakeEntitlementAPI {
	return &fakeEntitlementAPI{
		t: t,
		entitlement: map[string]any{
			"id":                               "ent1",
			"app_id":                           "app1",
			"label":                            "Super Admin",
			"entitlement_type_label":           "Role",
			"source":                           source,
			"description":                      nil,
			"approval_workflow_override":       nil,
			"provisioning_workflow":            nil,
			"deprovisioning_workflow":          nil,
			"custom_approval_message_override": nil,
		},
		settings: map[string]any{
			"visible_in_appstore":               true,
			"allowed_requester_groups_override": []any{map[string]any{"id": "group1"}},
			"time_based_access_override":        nil,
			"request_validation_webhook":        nil,
		},
	}
}

func (f *fakeEntitlementAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.requests = append(f.requests, r.Method+" "+r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	switch {
	case f.entitlementGone && strings.HasPrefix(r.URL.Path, "/v1/entitlements/ent1"),
		(f.appUnlisted || f.settingsGone) && r.URL.Path == "/v1/entitlements/ent1/appstore-settings":
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, `{"detail":"AppStore settings are not available for this entitlement."}`)
	case r.URL.Path == "/v1/apps/app1/appstore-settings" && r.Method == http.MethodGet:
		json.NewEncoder(w).Encode(map[string]any{"in_app_store": !f.appUnlisted})
	case r.URL.Path == "/v1/entitlements/ent1" && r.Method == http.MethodGet:
		entitlement := map[string]any{"appstore_settings": f.settings}
		for k, v := range f.entitlement {
			entitlement[k] = v
		}
		json.NewEncoder(w).Encode(entitlement)
	case r.URL.Path == "/v1/entitlements/ent1" && r.Method == http.MethodDelete:
		w.WriteHeader(f.deleteStatus)
		io.WriteString(w, f.deleteBody)
	case r.URL.Path == "/v1/entitlements/ent1/appstore-settings" && r.Method == http.MethodGet:
		json.NewEncoder(w).Encode(f.settings)
	case r.URL.Path == "/v1/entitlements/ent1/appstore-settings" && r.Method == http.MethodPatch:
		body, _ := io.ReadAll(r.Body)
		var patch map[string]any
		if err := json.Unmarshal(body, &patch); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		f.patches = append(f.patches, patch)
		for k, v := range patch {
			switch k {
			case "visible_in_appstore":
				if v != nil {
					f.settings[k] = v
				}
			case "allowed_requester_group_uuids_override":
				if v == nil {
					f.settings["allowed_requester_groups_override"] = nil
					continue
				}
				refs := []any{}
				for _, id := range v.([]any) {
					refs = append(refs, map[string]any{"id": id})
				}
				f.settings["allowed_requester_groups_override"] = refs
			case "time_based_access_override":
				f.settings[k] = withOptionObjects(v)
			case "request_validation_webhook_uuid":
				if v == nil {
					f.settings["request_validation_webhook"] = nil
					continue
				}
				f.settings["request_validation_webhook"] = map[string]any{"id": v}
			}
		}
		json.NewEncoder(w).Encode(f.settings)
	default:
		f.t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		http.Error(w, "unexpected request", http.StatusNotFound)
	}
}

// TestEntitlementAppStoreSettingsOnSyncedEntitlement imports an entitlement
// synced from an integration and its App Store settings, changes the group
// override from a list to [] to omitted, checking what the provider sends
// and keeps in state each time, then stops managing both with removed
// blocks, which must not call DELETE. It uses a fake API because a synced
// entitlement in an App Store app can't be set up on demand in a shared test
// domain.
func TestEntitlementAppStoreSettingsOnSyncedEntitlement(t *testing.T) {
	api := newFakeEntitlementAPI(t, "okta")
	server := httptest.NewServer(api)
	defer server.Close()

	config := func(groups string) string {
		return fakeProviderConfig(server.URL) + fmt.Sprintf(`
resource "lumos_entitlement" "admin" {
  app_id = "app1"
  label  = "Super Admin"
}

resource "lumos_entitlement_appstore_settings" "admin" {
  entitlement_id      = lumos_entitlement.admin.id
  visible_in_appstore = true
  %s
}
`, groups)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Matching config: the import plans no change, which also
				// shows the group refs read back as UUIDs.
				Config: config(`allowed_requester_group_uuids_override = ["group1"]`) + `
import {
  to = lumos_entitlement.admin
  id = "ent1"
}

import {
  to = lumos_entitlement_appstore_settings.admin
  id = "ent1"
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_entitlement.admin", plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction("lumos_entitlement_appstore_settings.admin", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.TestCheckResourceAttr("lumos_entitlement.admin", "appstore_settings.allowed_requester_group_uuids_override.#", "1"),
			},
			{
				// [] is sent as [] and kept in state, not turned into null.
				Config: config(`allowed_requester_group_uuids_override = []`),
				Check:  resource.TestCheckResourceAttr("lumos_entitlement_appstore_settings.admin", "allowed_requester_group_uuids_override.#", "0"),
			},
			{
				// Omitted, the provider sends null and the override reads
				// back as unset.
				Config: config(""),
				Check:  resource.TestCheckNoResourceAttr("lumos_entitlement_appstore_settings.admin", "allowed_requester_group_uuids_override"),
			},
			{
				Config: fakeProviderConfig(server.URL) + `
removed {
  from = lumos_entitlement_appstore_settings.admin

  lifecycle {
    destroy = false
  }
}

removed {
  from = lumos_entitlement.admin

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
	if len(api.patches) != 2 {
		t.Fatalf("got %d settings PATCHes, want one per group change: %v", len(api.patches), api.patches)
	}
	if got, ok := api.patches[0]["allowed_requester_group_uuids_override"].([]any); !ok || len(got) != 0 {
		t.Errorf("setting [] sent allowed_requester_group_uuids_override %v, want []", api.patches[0]["allowed_requester_group_uuids_override"])
	}
	if got, sent := api.patches[1]["allowed_requester_group_uuids_override"]; !sent || got != nil {
		t.Errorf("omitting the override sent %v (present: %v), want null", got, sent)
	}
	for _, r := range api.requests {
		if r == "DELETE /v1/entitlements/ent1" {
			t.Errorf("removing lumos_entitlement.admin with destroy = false called DELETE")
		}
	}
}

// TestEntitlementDestroyRefused checks that when the API refuses a delete,
// destroy fails with the API's error message, and that a removed block then
// stops managing the entitlement without another DELETE. It uses a fake API
// because a shared test domain can't produce an entitlement someone holds,
// or a synced one, on demand.
func TestEntitlementDestroyRefused(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		status     int
		body       string
		wantErrMsg string
	}{
		{
			name:       "held by an account",
			source:     "manual",
			status:     http.StatusConflict,
			body:       `{"detail":"Entitlement is still held by 2 active account(s); deletion would orphan their access.","code":"entitlement_in_use"}`,
			wantErrMsg: `still held by 2 active account`,
		},
		{
			name:       "synced from an integration",
			source:     "okta",
			status:     http.StatusUnprocessableEntity,
			body:       `{"detail":"Only entitlements from a manual source can be deleted."}`,
			wantErrMsg: `Only entitlements from a manual source can be deleted`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newFakeEntitlementAPI(t, tt.source)
			api.deleteStatus = tt.status
			api.deleteBody = tt.body
			server := httptest.NewServer(api)
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: fakeProviderConfig(server.URL) + `
import {
  to = lumos_entitlement.admin
  id = "ent1"
}

resource "lumos_entitlement" "admin" {
  app_id = "app1"
  label  = "Super Admin"
}
`,
					},
					{
						Config:      fakeProviderConfig(server.URL),
						ExpectError: regexp.MustCompile(tt.wantErrMsg),
					},
					{
						Config: fakeProviderConfig(server.URL) + `
removed {
  from = lumos_entitlement.admin

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
			deletes := 0
			for _, r := range api.requests {
				if r == "DELETE /v1/entitlements/ent1" {
					deletes++
				}
			}
			if deletes != 1 {
				t.Errorf("got %d DELETEs, want 1: the refused destroy, and none for the removed block", deletes)
			}
		})
	}
}

// entitlementAppStoreSettingsConfig manages ent1's App Store settings and
// nothing else, so a test can change what the fake answers for them alone.
func entitlementAppStoreSettingsConfig(serverURL string) string {
	return fakeProviderConfig(serverURL) + `
resource "lumos_entitlement_appstore_settings" "admin" {
  entitlement_id      = "ent1"
  visible_in_appstore = true
  allowed_requester_group_uuids_override = ["group1"]
}
`
}

const entitlementAppStoreSettingsImport = `
import {
  to = lumos_entitlement_appstore_settings.admin
  id = "ent1"
}
`

// TestEntitlementAppStoreSettingsAppUnlisted checks that while the
// entitlement's app isn't in the App Store, which makes the API answer 404
// for settings that still exist, refresh fails with a message saying so and
// keeps the settings in state, so listing the app again plans no change. It
// uses a fake API to switch the app's listing between steps (GOV-4617).
func TestEntitlementAppStoreSettingsAppUnlisted(t *testing.T) {
	api := newFakeEntitlementAPI(t, "manual")
	server := httptest.NewServer(api)
	defer server.Close()
	setUnlisted := func(unlisted bool) func() {
		return func() {
			api.mu.Lock()
			defer api.mu.Unlock()
			api.appUnlisted = unlisted
		}
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: entitlementAppStoreSettingsConfig(server.URL) + entitlementAppStoreSettingsImport,
			},
			{
				PreConfig:   setUnlisted(true),
				Config:      entitlementAppStoreSettingsConfig(server.URL),
				ExpectError: regexp.MustCompile(`the entitlement's app \(app1\) isn't in the App Store: list\s+the\s+app\s+again,\s+or\s+stop\s+managing\s+this\s+resource\s+with\s+a\s+removed\s+block\s+that\s+sets\s+destroy\s+=\s+false`),
			},
			{
				PreConfig: setUnlisted(false),
				Config:    entitlementAppStoreSettingsConfig(server.URL),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_entitlement_appstore_settings.admin", plancheck.ResourceActionNoop),
					},
				},
			},
			{
				// Still unlisted, a removed block stops managing the
				// settings: removing the resource from the config would plan
				// a destroy, which refreshes it first, and the API refuses
				// writes while the app is unlisted anyway.
				PreConfig: setUnlisted(true),
				Config: fakeProviderConfig(server.URL) + `
removed {
  from = lumos_entitlement_appstore_settings.admin

  lifecycle {
    destroy = false
  }
}
`,
			},
		},
	})
}

// TestEntitlementAppStoreSettingsDeleted checks that a 404 that means the
// settings are really gone still drops them from state, so the next plan
// creates them: the entitlement was deleted, or it has no settings although
// its app is in the App Store. It uses a fake API because neither can be
// set up on demand for an entitlement Terraform already manages.
func TestEntitlementAppStoreSettingsDeleted(t *testing.T) {
	tests := []struct {
		name string
		gone func(*fakeEntitlementAPI)
	}{
		{name: "entitlement deleted", gone: func(f *fakeEntitlementAPI) { f.entitlementGone = true }},
		{name: "settings gone from a listed app", gone: func(f *fakeEntitlementAPI) { f.settingsGone = true }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newFakeEntitlementAPI(t, "manual")
			server := httptest.NewServer(api)
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: entitlementAppStoreSettingsConfig(server.URL) + entitlementAppStoreSettingsImport,
					},
					{
						// The refresh below runs on the last config, which
						// mustn't still import.
						Config: entitlementAppStoreSettingsConfig(server.URL),
					},
					{
						PreConfig: func() {
							api.mu.Lock()
							defer api.mu.Unlock()
							tt.gone(api)
						},
						RefreshState:       true,
						ExpectNonEmptyPlan: true,
						Check: func(s *terraform.State) error {
							if _, ok := s.RootModule().Resources["lumos_entitlement_appstore_settings.admin"]; ok {
								return fmt.Errorf("lumos_entitlement_appstore_settings.admin is still in state after a refresh that read 404")
							}
							return nil
						},
					},
				},
			})
		})
	}
}
