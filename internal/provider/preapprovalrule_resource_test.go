package provider

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"sort"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

const fakePreApprovalRulePath = "/v1/apps/app1/pre-approval-rules/rule1"

var preApprovalRuleEntitlementFields = []string{
	"request_includes_entitlements",
	"request_does_not_include_entitlements",
	"request_has_entitlements",
}

// fakePreApprovalRuleAPI serves pre-approval rule "rule1" of app "app1" in
// the API's shapes: entitlement IDs come back in reverse order, the
// condition comes back with a null for every operator it doesn't use, and a
// rule with no entitlements field reads back as request_has_entitlements =
// false. An update that sends any entitlements field replaces all three.
// Once deleted is set, reads return the API's 404; once unrepresentable is
// set, they return the 409 for a rule covering access with no entitlement.
type fakePreApprovalRuleAPI struct {
	t               *testing.T
	mu              sync.Mutex
	writes          []map[string]any
	rule            map[string]any
	deleted         bool
	unrepresentable bool
}

func (f *fakePreApprovalRuleAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.URL.Path == "/v1/apps/app1/pre-approval-rules" && r.Method == http.MethodPost,
		r.URL.Path == fakePreApprovalRulePath && r.Method == http.MethodPatch:
		body, _ := io.ReadAll(r.Body)
		var input map[string]any
		if err := json.Unmarshal(body, &input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		f.writes = append(f.writes, input)
		f.store(input, r.Method == http.MethodPost)
		f.deleted = false
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusCreated)
		}
		json.NewEncoder(w).Encode(f.rule)
	case r.URL.Path == fakePreApprovalRulePath && r.Method == http.MethodGet && f.deleted:
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, `{"error":{"type":"not_found","code":"pre_approval_rule_not_found","detail":"Pre-approval rule not found."}}`)
	case r.URL.Path == fakePreApprovalRulePath && r.Method == http.MethodGet && f.unrepresentable:
		w.WriteHeader(http.StatusConflict)
		io.WriteString(w, `{"error":{"type":"conflict","code":"conflict","detail":"This pre-approval rule covers access that has no entitlement ID."}}`)
	case r.URL.Path == fakePreApprovalRulePath && r.Method == http.MethodGet:
		json.NewEncoder(w).Encode(f.rule)
	case r.URL.Path == fakePreApprovalRulePath && r.Method == http.MethodDelete:
		w.WriteHeader(http.StatusNoContent)
	default:
		f.t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		http.Error(w, "unexpected request", http.StatusNotFound)
	}
}

func (f *fakePreApprovalRuleAPI) store(input map[string]any, create bool) {
	if create {
		f.rule = map[string]any{"id": "rule1", "app_id": "app1", "request_has_entitlements": false}
	}
	if justification, ok := input["justification"]; ok {
		f.rule["justification"] = justification
	}
	if condition, ok := input["condition"].(map[string]any); ok {
		stored := map[string]any{"equals": nil, "in": nil, "contains": nil}
		for k, v := range condition {
			stored[k] = v
		}
		f.rule["condition"] = stored
	}
	replaced := false
	for _, name := range preApprovalRuleEntitlementFields {
		if _, ok := input[name]; ok {
			replaced = true
		}
	}
	if !replaced {
		return
	}
	for _, name := range preApprovalRuleEntitlementFields {
		delete(f.rule, name)
	}
	for _, name := range preApprovalRuleEntitlementFields[:2] {
		ids, _ := input[name].([]any)
		if len(ids) == 0 {
			continue
		}
		reversed := []any{}
		for i := len(ids) - 1; i >= 0; i-- {
			reversed = append(reversed, ids[i])
		}
		f.rule[name] = reversed
		return
	}
	has, _ := input["request_has_entitlements"].(bool)
	f.rule["request_has_entitlements"] = has
}

func fakePreApprovalRuleConfig(serverURL, entitlements string) string {
	return fakeProviderConfig(serverURL) + `
resource "lumos_pre_approval_rule" "test" {
  app_id        = "app1"
  justification = "Engineers need it"
  condition = jsonencode({
    attribute = { label = "Email", value = "identity.email" }
    equals    = { label = "a@example.com", value = "a@example.com" }
  })
  ` + entitlements + `
}
`
}

// TestPreApprovalRuleEntitlements checks that entitlement IDs read back in
// a different order plan no change, and that each update sends exactly the
// one entitlements field the config means: a list, true, or false when the
// config sets none, which replaces the previous list. It uses a fake API
// because a test can't choose the order the real API returns entitlement
// IDs in.
func TestPreApprovalRuleEntitlements(t *testing.T) {
	api := &fakePreApprovalRuleAPI{t: t}
	server := httptest.NewServer(api)
	defer server.Close()

	updates := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction("lumos_pre_approval_rule.test", plancheck.ResourceActionUpdate),
		},
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fakePreApprovalRuleConfig(server.URL, `request_includes_entitlements = ["ent1", "ent2"]`),
			},
			{
				Config:           fakePreApprovalRuleConfig(server.URL, `request_has_entitlements = true`),
				ConfigPlanChecks: updates,
			},
			{
				Config:           fakePreApprovalRuleConfig(server.URL, `request_does_not_include_entitlements = ["ent1"]`),
				ConfigPlanChecks: updates,
			},
			{
				Config:           fakePreApprovalRuleConfig(server.URL, ""),
				ConfigPlanChecks: updates,
			},
		},
	})

	api.mu.Lock()
	defer api.mu.Unlock()
	want := []map[string]any{
		{"request_includes_entitlements": []any{"ent1", "ent2"}},
		{"request_has_entitlements": true},
		{"request_does_not_include_entitlements": []any{"ent1"}},
		{"request_has_entitlements": false},
	}
	if len(api.writes) != len(want) {
		t.Fatalf("got %d writes, want a create and %d updates: %v", len(api.writes), len(want)-1, api.writes)
	}
	for i, write := range api.writes {
		sent := map[string]any{}
		for _, name := range preApprovalRuleEntitlementFields {
			if v, ok := write[name]; ok {
				if ids, ok := v.([]any); ok {
					// A set is sent in no particular order.
					sort.Slice(ids, func(a, b int) bool { return ids[a].(string) < ids[b].(string) })
				}
				sent[name] = v
			}
		}
		if !reflect.DeepEqual(sent, want[i]) {
			t.Errorf("write %d sent entitlements fields %v, want %v", i, sent, want[i])
		}
	}
}

// TestPreApprovalRuleDeletedOutsideTerraform checks that a rule whose read
// returns 404 is dropped from state on refresh and planned for creation,
// rather than failing the plan. It uses a fake API so the rule can
// disappear between steps without the test deleting it through the API
// outside Terraform.
func TestPreApprovalRuleDeletedOutsideTerraform(t *testing.T) {
	api := &fakePreApprovalRuleAPI{t: t}
	server := httptest.NewServer(api)
	defer server.Close()

	config := fakePreApprovalRuleConfig(server.URL, "")
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
						plancheck.ExpectResourceAction("lumos_pre_approval_rule.test", plancheck.ResourceActionCreate),
					},
				},
			},
		},
	})
}

// TestPreApprovalRuleUnrepresentable checks that a refresh fails when the
// read returns 409, rather than dropping the rule from state and planning a
// second one. It uses a fake API because the API returns that 409 only for
// a rule covering access with no entitlement, which can't be created
// through the API.
func TestPreApprovalRuleUnrepresentable(t *testing.T) {
	api := &fakePreApprovalRuleAPI{t: t}
	server := httptest.NewServer(api)
	defer server.Close()

	config := fakePreApprovalRuleConfig(server.URL, "")
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
					api.unrepresentable = true
				},
				Config:      config,
				ExpectError: regexp.MustCompile(`unexpected response code 409`),
			},
			{
				PreConfig: func() {
					api.mu.Lock()
					defer api.mu.Unlock()
					api.unrepresentable = false
				},
				Config: config,
			},
		},
	})
}
