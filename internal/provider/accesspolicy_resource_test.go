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
)

// fakeAccessPolicyAPI serves access policy "policy1" in the API's shapes:
// apps and each app's entitlements come back in reverse order with labels,
// a missing is_preapproved or entitlements reads back as false or [], and
// an everyone policy reads back with the condition Lumos uses to match
// everyone. Once deleted is set, the policy is gone: reads return the API's
// 404 until a write creates it again.
type fakeAccessPolicyAPI struct {
	t       *testing.T
	mu      sync.Mutex
	writes  []map[string]any
	policy  map[string]any
	deleted bool
}

func (f *fakeAccessPolicyAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.URL.Path == "/v1/access-policies" && r.Method == http.MethodPost,
		r.URL.Path == "/v1/access-policies/policy1" && r.Method == http.MethodPut:
		body, _ := io.ReadAll(r.Body)
		var input map[string]any
		if err := json.Unmarshal(body, &input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		f.writes = append(f.writes, input)
		f.policy = f.stored(input)
		f.deleted = false
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusCreated)
		}
		json.NewEncoder(w).Encode(f.policy)
	case r.URL.Path == "/v1/access-policies/policy1" && r.Method == http.MethodGet && f.deleted:
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, `{"error":{"type":"not_found","code":"access_policy_not_found","detail":"Access policy not found."}}`)
	case r.URL.Path == "/v1/access-policies/policy1" && r.Method == http.MethodGet:
		json.NewEncoder(w).Encode(f.policy)
	case r.URL.Path == "/v1/access-policies/policy1" && r.Method == http.MethodDelete:
		w.WriteHeader(http.StatusNoContent)
	default:
		f.t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		http.Error(w, "unexpected request", http.StatusNotFound)
	}
}

func (f *fakeAccessPolicyAPI) stored(input map[string]any) map[string]any {
	apps := []any{}
	inputApps, _ := input["apps"].([]any)
	for i := len(inputApps) - 1; i >= 0; i-- {
		app := inputApps[i].(map[string]any)
		preapproved, _ := app["is_preapproved"].(bool)
		entitlements := []any{}
		refs, _ := app["entitlements"].([]any)
		for j := len(refs) - 1; j >= 0; j-- {
			id := refs[j].(map[string]any)["id"]
			entitlements = append(entitlements, map[string]any{"id": id, "label": "Label of " + id.(string)})
		}
		apps = append(apps, map[string]any{
			"id":             app["id"],
			"label":          "Label of " + app["id"].(string),
			"is_preapproved": preapproved,
			"entitlements":   entitlements,
		})
	}
	return map[string]any{
		"id":                     "policy1",
		"name":                   input["name"],
		"business_justification": input["business_justification"],
		"status":                 "draft",
		"is_enabled":             false,
		"is_everyone_condition":  true,
		"access_condition":       map[string]any{"attribute": map[string]any{"label": "Stable identifier", "value": "identity.stable_identifier"}, "exists": map[string]any{"label": "True", "value": true}},
		"apps":                   apps,
	}
}

// TestAccessPolicyGrants checks that grants read back in a different order
// plan no change, that an app grant written without entitlements is sent as
// [] and plans no change once it reads back as [], and that removing an
// app's entitlements from the config plans and sends the change to an
// app-level grant. It uses a fake API because the real API orders apps and
// entitlements by internal ids, so a test can't make them come back in an
// order other than the config's on demand.
func TestAccessPolicyGrants(t *testing.T) {
	api := &fakeAccessPolicyAPI{t: t}
	server := httptest.NewServer(api)
	defer server.Close()

	config := func(app2Entitlements string) string {
		return fakeProviderConfig(server.URL) + `
resource "lumos_access_policy" "test" {
  name                   = "Engineering tools"
  business_justification = "Engineers need these"
  is_everyone_condition  = true
  status                 = "draft"
  is_enabled             = false
  apps = [
    { id = "app1" },
    {
      id = "app2"
      ` + app2Entitlements + `
    },
  ]
}
`
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(`entitlements = [{ id = "ent1" }, { id = "ent2" }]`),
			},
			{
				Config: config(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_access_policy.test", plancheck.ResourceActionUpdate),
					},
				},
			},
		},
	})

	api.mu.Lock()
	defer api.mu.Unlock()
	if len(api.writes) != 2 {
		t.Fatalf("got %d writes, want a create and one update: %v", len(api.writes), api.writes)
	}
	for i, write := range api.writes {
		for _, app := range write["apps"].([]any) {
			app := app.(map[string]any)
			entitlements, ok := app["entitlements"].([]any)
			if !ok {
				t.Errorf("write %d sent app %v entitlements %v, want a list", i, app["id"], app["entitlements"])
				continue
			}
			want := 0
			if i == 0 && app["id"] == "app2" {
				want = 2
			}
			if len(entitlements) != want {
				t.Errorf("write %d sent app %v %d entitlements, want %d", i, app["id"], len(entitlements), want)
			}
		}
	}
}

// TestAccessPolicyDeletedOutsideTerraform checks that a policy whose read
// returns 404 is dropped from state on refresh and planned for creation,
// rather than failing the plan. It uses a fake API so the policy can
// disappear between steps without the test deleting it through the API
// outside Terraform.
func TestAccessPolicyDeletedOutsideTerraform(t *testing.T) {
	api := &fakeAccessPolicyAPI{t: t}
	server := httptest.NewServer(api)
	defer server.Close()

	config := fakeProviderConfig(server.URL) + `
resource "lumos_access_policy" "test" {
  name                   = "Engineering tools"
  business_justification = "Engineers need these"
  is_everyone_condition  = true
  status                 = "draft"
  is_enabled             = false
  apps                   = [{ id = "app1" }]
}
`

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				PreConfig: func() {
					api.mu.Lock()
					defer api.mu.Unlock()
					api.deleted = true
				},
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_access_policy.test", plancheck.ResourceActionCreate),
					},
				},
			},
		},
	})
}

// TestAccessPolicyRequiresStatusAndEnabled checks a policy config missing
// status or is_enabled fails at plan time: neither can be cleared, so
// leaving one out can't mean anything else.
func TestAccessPolicyRequiresStatusAndEnabled(t *testing.T) {
	for _, missing := range []string{"status", "is_enabled"} {
		t.Run(missing, func(t *testing.T) {
			set := map[string]string{"status": `"draft"`, "is_enabled": "false"}
			delete(set, missing)
			var lines []string
			for name, value := range set {
				lines = append(lines, fmt.Sprintf("  %s = %s", name, value))
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{{
					Config: fakeProviderConfig("http://127.0.0.1:1") + `
resource "lumos_access_policy" "test" {
  name                   = "Engineering tools"
  business_justification = "Engineers need these"
  is_everyone_condition  = true
  apps                   = [{ id = "app1" }]
` + strings.Join(lines, "\n") + `
}
`,
					PlanOnly:    true,
					ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, missing)),
				}},
			})
		})
	}
}
