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
// [] and null apart. A create makes "ent1" a manual entitlement, and an
// entitlement PATCH is applied as the API applies it: a synced
// entitlement's current label is accepted and ignored, and a new one is
// refused; description is read-only; and description reads back as
// description_override if set, otherwise as the integration's description,
// which a cleared override brings back only on the next sync. Like the API,
// it answers 404 for the settings while the entitlement's app "app1" isn't
// in the App Store, and serves that app's settings for in_app_store.
type fakeEntitlementAPI struct {
	t                      *testing.T
	mu                     sync.Mutex
	requests               []string
	patches                []map[string]any
	entitlementBodies      []map[string]any
	entitlement            map[string]any
	integrationDescription any
	settings               map[string]any
	deleteStatus           int
	deleteBody             string

	appUnlisted     bool
	entitlementGone bool
	settingsGone    bool
}

func newFakeEntitlementAPI(t *testing.T, source string) *fakeEntitlementAPI {
	var integrationDescription any
	if source != "manual" {
		integrationDescription = "All engineers"
	}
	return &fakeEntitlementAPI{
		t:                      t,
		integrationDescription: integrationDescription,
		entitlement: map[string]any{
			"id":                               "ent1",
			"app_id":                           "app1",
			"label":                            "Super Admin",
			"entitlement_type_label":           "Role",
			"source":                           source,
			"description":                      integrationDescription,
			"description_override":             nil,
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
	case r.URL.Path == "/v1/entitlements" && r.Method == http.MethodPost:
		body := f.decodeEntitlementBody(w, r)
		if body == nil {
			return
		}
		f.entitlement["app_id"] = body["app_id"]
		f.entitlement["label"] = body["label"]
		f.entitlement["source"] = "manual"
		f.entitlement["description_override"] = body["description_override"]
		f.entitlement["description"] = body["description_override"]
		f.entitlement["custom_approval_message_override"] = body["custom_approval_message_override"]
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(f.entitlement)
	case r.URL.Path == "/v1/entitlements/ent1" && r.Method == http.MethodPatch:
		body := f.decodeEntitlementBody(w, r)
		if body == nil {
			return
		}
		if _, ok := body["description"]; ok {
			http.Error(w, `{"detail":"description is read-only; set description_override."}`, http.StatusBadRequest)
			return
		}
		if label, ok := body["label"]; ok && label != f.entitlement["label"] && f.entitlement["source"] != "manual" {
			http.Error(w, `{"detail":"Only manual entitlements can be renamed."}`, http.StatusBadRequest)
			return
		}
		for k, v := range body {
			if k == "description_override" && v == nil && f.entitlement[k] == nil {
				continue
			}
			f.entitlement[k] = v
			if k == "description_override" {
				f.entitlement["description"] = v
			}
		}
		json.NewEncoder(w).Encode(f.entitlement)
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

// sync runs the integration's sync, which brings back the integration's
// description on an entitlement with no description_override.
func (f *fakeEntitlementAPI) sync() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.entitlement["description_override"] == nil {
		f.entitlement["description"] = f.integrationDescription
	}
}

// decodeEntitlementBody records an entitlement create or update body. It
// writes a 400 and returns nil if the body isn't a JSON object.
func (f *fakeEntitlementAPI) decodeEntitlementBody(w http.ResponseWriter, r *http.Request) map[string]any {
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	f.entitlementBodies = append(f.entitlementBodies, body)
	return body
}

// TestEntitlementApprovalMessage creates an entitlement with an approval
// message, changes the message, then leaves it out of the config, checking
// the request body each time. Leaving it out must send null, which the API
// takes as "inherit the app's message", and then plan no change.
func TestEntitlementApprovalMessage(t *testing.T) {
	api := newFakeEntitlementAPI(t, "manual")
	api.deleteStatus = http.StatusNoContent
	server := httptest.NewServer(api)
	defer server.Close()

	config := func(message string) string {
		return fakeProviderConfig(server.URL) + fmt.Sprintf(`
resource "lumos_entitlement" "admin" {
  app_id = "app1"
  label  = "Super Admin"
  %s
}
`, message)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(`custom_approval_message_override = "Approved."`),
				Check:  resource.TestCheckResourceAttr("lumos_entitlement.admin", "custom_approval_message_override", "Approved."),
			},
			{
				Config: config(`custom_approval_message_override = "Approved by IT."`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						expectOnlyAttributesChange("lumos_entitlement.admin", "custom_approval_message_override"),
					},
				},
				Check: resource.TestCheckResourceAttr("lumos_entitlement.admin", "custom_approval_message_override", "Approved by IT."),
			},
			{
				Config: config(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						expectOnlyAttributesChange("lumos_entitlement.admin", "custom_approval_message_override"),
					},
				},
				Check: resource.TestCheckNoResourceAttr("lumos_entitlement.admin", "custom_approval_message_override"),
			},
		},
	})

	api.mu.Lock()
	defer api.mu.Unlock()
	// The create, then one update per message change; destroy is a DELETE.
	if len(api.entitlementBodies) != 3 {
		t.Fatalf("got %d entitlement create and update bodies, want 3: %v", len(api.entitlementBodies), api.entitlementBodies)
	}
	if got := api.entitlementBodies[0]["custom_approval_message_override"]; got != "Approved." {
		t.Errorf("create sent custom_approval_message_override %v, want %q", got, "Approved.")
	}
	if got := api.entitlementBodies[1]["custom_approval_message_override"]; got != "Approved by IT." {
		t.Errorf("changing the message sent custom_approval_message_override %v, want %q", got, "Approved by IT.")
	}
	if got, sent := api.entitlementBodies[2]["custom_approval_message_override"]; !sent || got != nil {
		t.Errorf("omitting the message sent %v (present: %v), want null", got, sent)
	}
	for i, body := range api.entitlementBodies[1:] {
		if _, sent := body["description"]; sent {
			t.Errorf("update %d sent description %v; it's read-only, so sending it would mark it as set by you", i+1, body["description"])
		}
	}
}

// TestEntitlementApprovalMessageOnSyncedEntitlement imports an entitlement
// synced from an integration, sets its approval message, then leaves it out
// of the config. Each update resends the label, which the integration owns:
// the API accepts it unchanged, so a synced entitlement's message can be
// managed without renaming it.
func TestEntitlementApprovalMessageOnSyncedEntitlement(t *testing.T) {
	api := newFakeEntitlementAPI(t, "okta")
	server := httptest.NewServer(api)
	defer server.Close()

	config := func(message string) string {
		return fakeProviderConfig(server.URL) + fmt.Sprintf(`
resource "lumos_entitlement" "admin" {
  app_id = "app1"
  label  = "Super Admin"
  %s
}
`, message)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config("") + `
import {
  to = lumos_entitlement.admin
  id = "ent1"
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_entitlement.admin", plancheck.ResourceActionNoop),
					},
				},
			},
			{
				Config: config(`custom_approval_message_override = "Approved."`),
				Check:  resource.TestCheckResourceAttr("lumos_entitlement.admin", "custom_approval_message_override", "Approved."),
			},
			{
				Config: config(""),
				Check:  resource.TestCheckNoResourceAttr("lumos_entitlement.admin", "custom_approval_message_override"),
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
	if len(api.entitlementBodies) != 2 {
		t.Fatalf("got %d entitlement update bodies, want one per message change: %v", len(api.entitlementBodies), api.entitlementBodies)
	}
	if got := api.entitlementBodies[0]["custom_approval_message_override"]; got != "Approved." {
		t.Errorf("setting the message sent custom_approval_message_override %v, want %q", got, "Approved.")
	}
	if got, sent := api.entitlementBodies[1]["custom_approval_message_override"]; !sent || got != nil {
		t.Errorf("omitting the message sent %v (present: %v), want null", got, sent)
	}
}

// TestEntitlementDescription creates a manual entitlement with a
// description, changes it, then leaves it out of the config, checking the
// request body and the description read back each time. Leaving it out must
// send null, which leaves a manual entitlement with no description.
func TestEntitlementDescription(t *testing.T) {
	api := newFakeEntitlementAPI(t, "manual")
	api.deleteStatus = http.StatusNoContent
	server := httptest.NewServer(api)
	defer server.Close()

	config := func(description string) string {
		return fakeProviderConfig(server.URL) + fmt.Sprintf(`
resource "lumos_entitlement" "admin" {
  app_id = "app1"
  label  = "Super Admin"
  %s
}
`, description)
	}
	descriptionChanges := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			expectOnlyAttributesChange("lumos_entitlement.admin", "description_override", "description"),
		},
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(`description_override = "Ask #eng first."`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_entitlement.admin", "description_override", "Ask #eng first."),
					resource.TestCheckResourceAttr("lumos_entitlement.admin", "description", "Ask #eng first."),
				),
			},
			{
				Config:           config(`description_override = "Ask IT."`),
				ConfigPlanChecks: descriptionChanges,
				Check:            resource.TestCheckResourceAttr("lumos_entitlement.admin", "description", "Ask IT."),
			},
			{
				Config:           config(""),
				ConfigPlanChecks: descriptionChanges,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("lumos_entitlement.admin", "description_override"),
					resource.TestCheckNoResourceAttr("lumos_entitlement.admin", "description"),
				),
			},
		},
	})

	api.mu.Lock()
	defer api.mu.Unlock()
	if len(api.entitlementBodies) != 3 {
		t.Fatalf("got %d entitlement create and update bodies, want 3: %v", len(api.entitlementBodies), api.entitlementBodies)
	}
	if got := api.entitlementBodies[0]["description_override"]; got != "Ask #eng first." {
		t.Errorf("create sent description_override %v, want %q", got, "Ask #eng first.")
	}
	if got := api.entitlementBodies[1]["description_override"]; got != "Ask IT." {
		t.Errorf("changing the description sent description_override %v, want %q", got, "Ask IT.")
	}
	if got, sent := api.entitlementBodies[2]["description_override"]; !sent || got != nil {
		t.Errorf("omitting the description sent %v (present: %v), want null", got, sent)
	}
}

// TestEntitlementDescriptionOnSyncedEntitlement imports an entitlement
// synced from an integration, sets a description in place of the
// integration's, then leaves it out of the config. The integration's
// description returns on the next sync, which must plan no change.
func TestEntitlementDescriptionOnSyncedEntitlement(t *testing.T) {
	api := newFakeEntitlementAPI(t, "okta")
	server := httptest.NewServer(api)
	defer server.Close()

	config := func(description string) string {
		return fakeProviderConfig(server.URL) + fmt.Sprintf(`
resource "lumos_entitlement" "admin" {
  app_id = "app1"
  label  = "Super Admin"
  %s
}
`, description)
	}
	noChange := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction("lumos_entitlement.admin", plancheck.ResourceActionNoop),
		},
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config("") + `
import {
  to = lumos_entitlement.admin
  id = "ent1"
}
`,
				ConfigPlanChecks: noChange,
				Check:            resource.TestCheckResourceAttr("lumos_entitlement.admin", "description", "All engineers"),
			},
			{
				Config: config(`description_override = "Ask #eng first."`),
				Check:  resource.TestCheckResourceAttr("lumos_entitlement.admin", "description", "Ask #eng first."),
			},
			{
				// Cleared, the description reads null until the next sync.
				Config: config(""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("lumos_entitlement.admin", "description_override"),
					resource.TestCheckNoResourceAttr("lumos_entitlement.admin", "description"),
				),
			},
			{
				PreConfig:        api.sync,
				Config:           config(""),
				ConfigPlanChecks: noChange,
				Check:            resource.TestCheckResourceAttr("lumos_entitlement.admin", "description", "All engineers"),
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
	if len(api.entitlementBodies) != 2 {
		t.Fatalf("got %d entitlement update bodies, want one per description change: %v", len(api.entitlementBodies), api.entitlementBodies)
	}
	if got := api.entitlementBodies[0]["description_override"]; got != "Ask #eng first." {
		t.Errorf("setting the description sent description_override %v, want %q", got, "Ask #eng first.")
	}
	if got, sent := api.entitlementBodies[1]["description_override"]; !sent || got != nil {
		t.Errorf("omitting the description sent %v (present: %v), want null", got, sent)
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
