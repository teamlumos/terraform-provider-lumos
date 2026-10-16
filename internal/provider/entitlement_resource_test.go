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

// fakeEntitlementAPI serves entitlement "ent1" and its App Store settings
// in the API's shapes, applying settings PATCHes as a merge patch. It keeps
// an empty group list as sent, so a test can check that the provider holds
// [] and null apart.
type fakeEntitlementAPI struct {
	t            *testing.T
	mu           sync.Mutex
	requests     []string
	patches      []map[string]any
	entitlement  map[string]any
	settings     map[string]any
	deleteStatus int
	deleteBody   string
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
	api := newFakeEntitlementAPI(t, "OKTA")
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
			source:     "MANUAL",
			status:     http.StatusConflict,
			body:       `{"detail":"Entitlement is still held by 2 active account(s); deletion would orphan their access.","code":"entitlement_in_use"}`,
			wantErrMsg: `still held by 2 active account`,
		},
		{
			name:       "synced from an integration",
			source:     "OKTA",
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
