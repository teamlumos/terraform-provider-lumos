package provider

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"regexp"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// These tests cover the documented way to leave some settings to the Lumos
// UI: lifecycle { ignore_changes }. Terraform replaces each ignored
// attribute's config value with its state value before the provider plans,
// so a setting the config leaves out isn't cleared. Each test also runs the
// same config without the lifecycle block, which must clear the setting, so
// the checks can tell the two apart.

const removedAppStoreSettings = `
removed {
  from = lumos_app_appstore_settings.slack

  lifecycle {
    destroy = false
  }
}
`

const importAppStoreSettings = `
import {
  to = lumos_app_appstore_settings.slack
  id = "app1"
}
`

// ignoredAppStoreSettings is a lumos_app_appstore_settings config that sets
// only the required switches, ignoring the attributes in ignore.
func ignoredAppStoreSettings(serverURL, catalogVisibility string, multiple bool, ignore string) string {
	lifecycle := ""
	if ignore != "" {
		lifecycle = fmt.Sprintf("lifecycle {\n    ignore_changes = [%s]\n  }", ignore)
	}
	visibility := ""
	if catalogVisibility != "" {
		visibility = fmt.Sprintf("catalog_visibility = %q", catalogVisibility)
	}
	return fakeProviderConfig(serverURL) + fmt.Sprintf(`
resource "lumos_app_appstore_settings" "slack" {
  app_id                              = "app1"
  in_app_store                        = true
  %s
  allow_permission_selection          = true
  allow_multiple_permission_selection = %t
  %s
}
`, visibility, multiple, lifecycle)
}

// sentJSON encodes a value from a recorded PATCH body for comparison, with
// string lists sorted, since sets have no order. The fake API turns recorded
// time-based access option labels into option objects; this turns them back.
func sentJSON(v any) string {
	switch v := v.(type) {
	case []any:
		var labels []string
		for _, e := range v {
			s, ok := e.(string)
			if !ok {
				break
			}
			labels = append(labels, s)
		}
		if len(labels) == len(v) && len(v) > 0 {
			sort.Strings(labels)
			encoded, _ := json.Marshal(labels)
			return string(encoded)
		}
	case map[string]any:
		if options, ok := v["options"].([]any); ok {
			labels := make([]any, 0, len(options))
			for _, o := range options {
				if option, ok := o.(map[string]any); ok {
					labels = append(labels, option["label"])
				} else {
					labels = append(labels, o)
				}
			}
			encoded, _ := json.Marshal(map[string]any{"options": labels, "default": v["default"]})
			return string(encoded)
		}
	}
	encoded, _ := json.Marshal(v)
	return string(encoded)
}

// TestAppStoreSettingsIgnoreChanges imports settings with an ignored setting
// the config leaves out, then has the UI change it, then makes an unrelated
// change. The import and the UI change must plan nothing, and the update
// must send the value Lumos has, not clear it. Without the lifecycle block,
// the import must plan the clear and the update send it.
func TestAppStoreSettingsIgnoreChanges(t *testing.T) {
	timeBasedAccess := func(defaultLabel any, labels ...string) map[string]any {
		var options []any
		for _, l := range labels {
			options = append(options, map[string]any{"label": l, "value_seconds": nil})
		}
		return map[string]any{"options": options, "default": defaultLabel}
	}
	tests := []struct {
		name    string
		attr    string
		initial any
		fromUI  any
		// wantSent is what the unrelated update sends for attr, and
		// wantCleared what the update sends without the lifecycle block. An
		// empty wantCleared skips that run: the setting has nothing to clear.
		wantSent    string
		wantCleared string
		// stateAttr, if set, is checked to hold stateInitial after the import
		// and stateFromUI after the UI change and the update.
		stateAttr    string
		stateInitial string
		stateFromUI  string
	}{
		{
			name: "request_instructions", attr: "request_instructions",
			initial: "Ask first.", fromUI: "Ask your manager.",
			wantSent: `"Ask your manager."`, wantCleared: `null`,
			stateAttr: "request_instructions", stateInitial: "Ask first.", stateFromUI: "Ask your manager.",
		},
		{
			name: "redirect_url", attr: "redirect_url",
			initial: "https://a.example.com", fromUI: "https://b.example.com",
			wantSent: `"https://b.example.com"`, wantCleared: `null`,
			stateAttr: "redirect_url", stateInitial: "https://a.example.com", stateFromUI: "https://b.example.com",
		},
		{
			name: "allowed_requester_group_uuids", attr: "allowed_requester_group_uuids",
			initial: []any{"group1"}, fromUI: []any{"group1", "group2"},
			wantSent: `["group1","group2"]`, wantCleared: `[]`,
			stateAttr: "allowed_requester_group_uuids.#", stateInitial: "1", stateFromUI: "2",
		},
		{
			name: "time_based_access", attr: "time_based_access",
			initial: timeBasedAccess("1 day", "1 day"), fromUI: timeBasedAccess("1 week", "1 day", "1 week"),
			wantSent: `{"default":"1 week","options":["1 day","1 week"]}`, wantCleared: `{"default":null,"options":[]}`,
			stateAttr: "time_based_access.default", stateInitial: "1 day", stateFromUI: "1 week",
		},
		{
			// Required attributes can be ignored too: the config still sets
			// one, but a UI change isn't reverted.
			name: "catalog_visibility", attr: "catalog_visibility",
			initial: "visible", fromUI: "search-only",
			wantSent: `"search-only"`, wantCleared: `"visible"`,
			stateAttr: "catalog_visibility", stateInitial: "visible", stateFromUI: "search-only",
		},
		{
			name: "request_instructions unset in Lumos", attr: "request_instructions",
			wantSent: `null`,
		},
		{
			name: "allowed_requester_group_uuids unset in Lumos", attr: "allowed_requester_group_uuids",
			wantSent: `[]`,
		},
		{
			name: "time_based_access unset in Lumos", attr: "time_based_access",
			initial: timeBasedAccess(nil), fromUI: timeBasedAccess(nil),
			wantSent: `{"default":null,"options":[]}`,
		},
	}
	const addr = "lumos_app_appstore_settings.slack"

	newAPI := func(t *testing.T, attr string, value any) *fakeAPI {
		api := newFakeAPI(t)
		// Only the setting under test has a value, so the config can leave
		// the rest out without clearing anything.
		api.settings["allowed_requester_group_uuids"] = nil
		api.settings["time_based_access"] = map[string]any{"options": nil, "default": nil}
		api.settings[attr] = value
		return api
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(t, tt.attr, tt.initial)
			server := httptest.NewServer(api)
			defer server.Close()

			state := func(want string) resource.TestCheckFunc {
				if tt.stateAttr == "" {
					return resource.ComposeTestCheckFunc()
				}
				return resource.TestCheckResourceAttr(addr, tt.stateAttr, want)
			}
			patchesBeforeUpdate := -1
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: ignoredAppStoreSettings(server.URL, "visible", false, tt.attr) + importAppStoreSettings,
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectResourceAction(addr, plancheck.ResourceActionNoop),
							},
						},
						Check: state(tt.stateInitial),
					},
					{
						PreConfig: func() {
							api.mu.Lock()
							defer api.mu.Unlock()
							api.settings[tt.attr] = tt.fromUI
							patchesBeforeUpdate = len(api.patches)
						},
						Config: ignoredAppStoreSettings(server.URL, "visible", false, tt.attr),
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
						},
						Check: state(tt.stateFromUI),
					},
					{
						Config: ignoredAppStoreSettings(server.URL, "visible", true, tt.attr),
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								expectOnlyAttributesChange(addr, "allow_multiple_permission_selection"),
							},
						},
						Check: state(tt.stateFromUI),
					},
					{Config: fakeProviderConfig(server.URL) + removedAppStoreSettings},
				},
			})

			api.mu.Lock()
			defer api.mu.Unlock()
			if patchesBeforeUpdate != 0 || len(api.patches) != 1 {
				t.Fatalf("got %d PATCHes, %d of them before the update; want only the update: %v", len(api.patches), patchesBeforeUpdate, api.patches)
			}
			if got := sentJSON(api.patches[0][tt.attr]); got != tt.wantSent {
				t.Errorf("update sent %s %s, want %s", tt.attr, got, tt.wantSent)
			}
		})

		if tt.wantCleared == "" {
			continue
		}
		t.Run(tt.name+" without lifecycle", func(t *testing.T) {
			api := newAPI(t, tt.attr, tt.fromUI)
			server := httptest.NewServer(api)
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: ignoredAppStoreSettings(server.URL, "visible", false, "") + importAppStoreSettings,
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectResourceAction(addr, plancheck.ResourceActionUpdate),
								expectOnlyAttributesChange(addr, tt.attr),
							},
						},
					},
					{Config: fakeProviderConfig(server.URL) + removedAppStoreSettings},
				},
			})

			api.mu.Lock()
			defer api.mu.Unlock()
			if len(api.patches) != 1 {
				t.Fatalf("got %d PATCHes, want the update: %v", len(api.patches), api.patches)
			}
			if got, sent := api.patches[0][tt.attr]; !sent || sentJSON(got) != tt.wantCleared {
				t.Errorf("update sent %s %s (present: %t), want %s", tt.attr, sentJSON(got), sent, tt.wantCleared)
			}
		})
	}
}

// TestAppStoreSettingsIgnoreChangesOnCreate creates the settings resource
// for an app whose request instructions were set in the UI, ignoring them.
// Terraform has no state to take the ignored value from on create, so the
// create clears them, which is why the docs say to import instead.
func TestAppStoreSettingsIgnoreChangesOnCreate(t *testing.T) {
	api := newFakeAPI(t)
	api.settings["allowed_requester_group_uuids"] = nil
	api.settings["time_based_access"] = map[string]any{"options": nil, "default": nil}
	api.settings["request_instructions"] = "Ask first."
	server := httptest.NewServer(api)
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ignoredAppStoreSettings(server.URL, "visible", false, "request_instructions"),
				Check:  resource.TestCheckNoResourceAttr("lumos_app_appstore_settings.slack", "request_instructions"),
			},
			{Config: fakeProviderConfig(server.URL) + removedAppStoreSettings},
		},
	})

	api.mu.Lock()
	defer api.mu.Unlock()
	if got, sent := api.patches[0]["request_instructions"]; !sent || got != nil {
		t.Errorf("create sent request_instructions %v (present: %t), want null", got, sent)
	}
}

// TestAppStoreSettingsIgnoreChangesStillRequired checks that ignoring a
// required attribute doesn't let the config leave it out.
func TestAppStoreSettingsIgnoreChangesStillRequired(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:      ignoredAppStoreSettings("http://127.0.0.1:1", "", false, "catalog_visibility"),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(`The argument "catalog_visibility" is required`),
		}},
	})
}

// TestEntitlementAppStoreSettingsIgnoreChanges imports an entitlement's
// settings with the group override ignored and left out, has the UI change
// the override, then hides the entitlement. The update must send the
// override Lumos has, not null, which would make the entitlement inherit
// the app's groups. Without the lifecycle block, the import must send null.
func TestEntitlementAppStoreSettingsIgnoreChanges(t *testing.T) {
	const addr = "lumos_entitlement_appstore_settings.admin"
	config := func(serverURL string, visible, ignore bool) string {
		lifecycle := ""
		if ignore {
			lifecycle = "lifecycle {\n    ignore_changes = [allowed_requester_group_uuids_override]\n  }"
		}
		return fakeProviderConfig(serverURL) + fmt.Sprintf(`
resource "lumos_entitlement_appstore_settings" "admin" {
  entitlement_id      = "ent1"
  visible_in_appstore = %t
  %s
}
`, visible, lifecycle)
	}
	const importSettings = `
import {
  to = lumos_entitlement_appstore_settings.admin
  id = "ent1"
}
`
	const removed = `
removed {
  from = lumos_entitlement_appstore_settings.admin

  lifecycle {
    destroy = false
  }
}
`

	t.Run("ignored", func(t *testing.T) {
		api := newFakeEntitlementAPI(t, "okta")
		server := httptest.NewServer(api)
		defer server.Close()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: config(server.URL, true, true) + importSettings,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(addr, plancheck.ResourceActionNoop),
						},
					},
					Check: resource.TestCheckResourceAttr(addr, "allowed_requester_group_uuids_override.#", "1"),
				},
				{
					PreConfig: func() {
						api.mu.Lock()
						defer api.mu.Unlock()
						api.settings["allowed_requester_group_uuids_override"] = []any{"group1", "group2"}
					},
					Config: config(server.URL, true, true),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
					},
				},
				{
					Config: config(server.URL, false, true),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							expectOnlyAttributesChange(addr, "visible_in_appstore"),
						},
					},
					Check: resource.TestCheckResourceAttr(addr, "allowed_requester_group_uuids_override.#", "2"),
				},
				{Config: fakeProviderConfig(server.URL) + removed},
			},
		})

		api.mu.Lock()
		defer api.mu.Unlock()
		if len(api.patches) != 1 {
			t.Fatalf("got %d PATCHes, want the update: %v", len(api.patches), api.patches)
		}
		if got := sentJSON(api.patches[0]["allowed_requester_group_uuids_override"]); got != `["group1","group2"]` {
			t.Errorf("update sent allowed_requester_group_uuids_override %s, want [group1 group2]", got)
		}
	})

	t.Run("without lifecycle", func(t *testing.T) {
		api := newFakeEntitlementAPI(t, "okta")
		server := httptest.NewServer(api)
		defer server.Close()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: config(server.URL, true, false) + importSettings,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							expectOnlyAttributesChange(addr, "allowed_requester_group_uuids_override"),
						},
					},
				},
				{Config: fakeProviderConfig(server.URL) + removed},
			},
		})

		api.mu.Lock()
		defer api.mu.Unlock()
		if got, sent := api.patches[0]["allowed_requester_group_uuids_override"]; !sent || got != nil {
			t.Errorf("update sent allowed_requester_group_uuids_override %v (present: %t), want null", got, sent)
		}
	})
}

// TestAppCustomIgnoreChangesLogo creates a custom app with logo_url_override
// ignored and left out, has the UI set a logo, then edits the description.
// The UI change must plan nothing, and the update must send the UI's logo,
// not null. Without the lifecycle block, the update must clear the logo.
func TestAppCustomIgnoreChangesLogo(t *testing.T) {
	const addr = "lumos_app.wiki"
	const uiLogo = "https://example.com/ui-logo.png"
	config := func(serverURL, description string, ignore bool) string {
		lifecycle := ""
		if ignore {
			lifecycle = "lifecycle {\n    ignore_changes = [logo_url_override]\n  }"
		}
		return fakeProviderConfig(serverURL) + fmt.Sprintf(`
resource "lumos_app" "wiki" {
  label_override       = "Wiki"
  category_override    = "Developers"
  description_override = %q
  %s
}
`, description, lifecycle)
	}

	for _, ignore := range []bool{true, false} {
		t.Run(fmt.Sprintf("ignored=%t", ignore), func(t *testing.T) {
			api := &fakeCustomAppAPI{t: t}
			server := httptest.NewServer(api)
			defer server.Close()
			setLogoInUI := func() {
				api.mu.Lock()
				defer api.mu.Unlock()
				api.app["logo_url_override"] = uiLogo
			}

			steps := []resource.TestStep{{Config: config(server.URL, "Docs", ignore)}}
			if ignore {
				steps = append(steps, resource.TestStep{
					PreConfig: setLogoInUI,
					Config:    config(server.URL, "Docs", ignore),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(addr, "logo_url_override", uiLogo),
						resource.TestCheckResourceAttr(addr, "logo_url", uiLogo),
					),
				})
			}
			wantChanges := []string{"description_override", "description"}
			if !ignore {
				wantChanges = append(wantChanges, "logo_url_override", "logo_url")
			}
			steps = append(steps, resource.TestStep{
				PreConfig: setLogoInUI,
				Config:    config(server.URL, "Docs, edited", ignore),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{expectOnlyAttributesChange(addr, wantChanges...)},
				},
			})
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps:                    steps,
			})

			api.mu.Lock()
			defer api.mu.Unlock()
			if len(api.patches) != 1 {
				t.Fatalf("got %d PATCHes, want the update: %v", len(api.patches), api.patches)
			}
			got, sent := api.patches[0]["logo_url_override"]
			if ignore && got != uiLogo {
				t.Errorf("update sent logo_url_override %v (present: %t), want %q", got, sent, uiLogo)
			}
			if !ignore && (!sent || got != nil) {
				t.Errorf("update sent logo_url_override %v (present: %t), want null", got, sent)
			}
		})
	}
}
