package provider

import (
	"encoding/json"
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

// fakeSlotAPI serves app and entitlement workflow slots the way the API
// does: PUT binds a workflow and returns it, GET returns the bound workflow
// or 404, and DELETE on the entitlement approval override turns it off with
// 204 (or deleteStatus, when set).
type fakeSlotAPI struct {
	t            *testing.T
	mu           sync.Mutex
	bound        map[string]string
	requests     []string
	deleteStatus int
}

func newFakeSlotAPI(t *testing.T) *fakeSlotAPI {
	return &fakeSlotAPI{t: t, bound: map[string]string{}}
}

func (f *fakeSlotAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.requests = append(f.requests, r.Method+" "+r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	if !strings.HasPrefix(r.URL.Path, "/v1/apps/") && !strings.HasPrefix(r.URL.Path, "/v1/entitlements/") {
		f.t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		http.Error(w, "unexpected request", http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodPut:
		var body struct {
			WorkflowID string `json:"workflow_id"`
		}
		raw, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(raw, &body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		f.bound[r.URL.Path] = body.WorkflowID
		json.NewEncoder(w).Encode(map[string]any{"id": body.WorkflowID, "type": "approval", "steps": []any{}})
	case http.MethodGet:
		id, ok := f.bound[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, `{"error":{"code":"not_found","message":"No workflow bound"}}`)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"id": id, "type": "approval", "steps": []any{}})
	case http.MethodDelete:
		delete(f.bound, r.URL.Path)
		status := f.deleteStatus
		if status == 0 {
			status = http.StatusNoContent
		}
		w.WriteHeader(status)
		if status == http.StatusNotFound {
			io.WriteString(w, `{"error":{"code":"not_found","message":"Entitlement has no request configuration"}}`)
		}
	}
}

func (f *fakeSlotAPI) requestsWith(method string) []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []string
	for _, r := range f.requests {
		if strings.HasPrefix(r, method+" ") {
			out = append(out, r)
		}
	}
	return out
}

func appApprovalBindingConfig(serverURL, appID, workflowID string) string {
	return fakeProviderConfig(serverURL) + `
resource "lumos_app_approval_workflow" "test" {
  app_id      = "` + appID + `"
  workflow_id = "` + workflowID + `"
}
`
}

const removedAppApprovalBinding = `
removed {
  from = lumos_app_approval_workflow.test
  lifecycle {
    destroy = false
  }
}
`

// TestWorkflowBindingDestroyRefused checks that removing an app workflow
// binding from the config fails the plan, since Lumos can't empty the slot,
// and that a removed block with destroy = false stops managing it without a
// request to Lumos. It uses a fake API so the slot's state can be checked
// directly.
func TestWorkflowBindingDestroyRefused(t *testing.T) {
	fake := newFakeSlotAPI(t)
	server := httptest.NewServer(fake)
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: appApprovalBindingConfig(server.URL, "app1", "wf1")},
			{
				Config:      fakeProviderConfig(server.URL),
				ExpectError: regexp.MustCompile(`(?s)Workflow binding can't be removed.*removed block`),
			},
			{Config: fakeProviderConfig(server.URL) + removedAppApprovalBinding},
		},
	})

	if got := fake.bound["/v1/apps/app1/approval-workflow"]; got != "wf1" {
		t.Errorf("slot holds %q after the binding was removed from state, want wf1 still bound", got)
	}
	if puts := fake.requestsWith(http.MethodPut); len(puts) != 1 {
		t.Errorf("sent %v, want only the first bind", puts)
	}
	if deletes := fake.requestsWith(http.MethodDelete); len(deletes) != 0 {
		t.Errorf("sent %v, want no DELETE", deletes)
	}
}

// TestWorkflowBindingMoveRefused checks that changing app_id, which would
// replace the binding and leave the workflow bound to the old app, fails
// the plan. It uses a fake API so the check covers only the provider's plan.
func TestWorkflowBindingMoveRefused(t *testing.T) {
	fake := newFakeSlotAPI(t)
	server := httptest.NewServer(fake)
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: appApprovalBindingConfig(server.URL, "app1", "wf1")},
			{
				Config:      appApprovalBindingConfig(server.URL, "app2", "wf1"),
				ExpectError: regexp.MustCompile(`(?s)Workflow binding can't be moved.*app2`),
			},
			{
				Config: appApprovalBindingConfig(server.URL, "app1", "wf2"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_app_approval_workflow.test", plancheck.ResourceActionUpdate),
					},
				},
			},
			{Config: fakeProviderConfig(server.URL) + removedAppApprovalBinding},
		},
	})

	if _, ok := fake.bound["/v1/apps/app2/approval-workflow"]; ok {
		t.Error("bound a workflow to app2, want the move refused")
	}
	if got := fake.bound["/v1/apps/app1/approval-workflow"]; got != "wf2" {
		t.Errorf("app1 holds %q, want wf2 after the rebind", got)
	}
}

// unknownOwnerBindingConfig binds wf1 to an app whose ID comes from a
// terraform_data resource, so it isn't known until apply whenever owner
// changes, like an app created or replaced in the same apply.
func unknownOwnerBindingConfig(serverURL, owner string) string {
	return fakeProviderConfig(serverURL) + `
resource "terraform_data" "owner" {
  input = "` + owner + `"
}

resource "lumos_app_approval_workflow" "test" {
  app_id      = terraform_data.owner.output
  workflow_id = "wf1"
}
`
}

// TestWorkflowBindingUnknownOwner checks a binding whose app ID isn't known
// until apply. Creating it with the owner in one apply works. Repointing an
// existing binding at an owner that isn't known yet fails the plan, since
// the plan can't tell whether the binding moves, and a move would leave the
// workflow bound to the old app; no bind request is sent. It uses a fake API
// so the slots can be checked directly.
func TestWorkflowBindingUnknownOwner(t *testing.T) {
	fake := newFakeSlotAPI(t)
	server := httptest.NewServer(fake)
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: unknownOwnerBindingConfig(server.URL, "app1"),
				Check:  resource.TestCheckResourceAttr("lumos_app_approval_workflow.test", "app_id", "app1"),
			},
			{
				Config:      unknownOwnerBindingConfig(server.URL, "app2"),
				ExpectError: regexp.MustCompile(`(?s)Workflow binding can't be moved.*isn't known until apply`),
			},
			{
				Config: fakeProviderConfig(server.URL) + `
resource "terraform_data" "owner" {
  input = "app1"
}
` + removedAppApprovalBinding,
			},
		},
	})

	if got := fake.bound["/v1/apps/app1/approval-workflow"]; got != "wf1" {
		t.Errorf("app1 holds %q, want wf1 from the first apply", got)
	}
	if _, ok := fake.bound["/v1/apps/app2/approval-workflow"]; ok {
		t.Error("bound a workflow to app2, want the repoint refused")
	}
	if puts := fake.requestsWith(http.MethodPut); len(puts) != 1 {
		t.Errorf("sent %v, want only the first bind", puts)
	}
}

// TestWorkflowBindingReadAndImport checks that a binding reads its
// workflow_id back, imports by the app's ID, and plans to bind again when
// the slot reads back empty (404). It uses a fake API so the slot can be
// emptied between steps.
func TestWorkflowBindingReadAndImport(t *testing.T) {
	fake := newFakeSlotAPI(t)
	server := httptest.NewServer(fake)
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: appApprovalBindingConfig(server.URL, "app1", "wf1"),
				Check:  resource.TestCheckResourceAttr("lumos_app_approval_workflow.test", "workflow_id", "wf1"),
			},
			{
				ResourceName:                         "lumos_app_approval_workflow.test",
				ImportState:                          true,
				ImportStateId:                        "app1",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "app_id",
			},
			{
				PreConfig: func() {
					fake.mu.Lock()
					defer fake.mu.Unlock()
					delete(fake.bound, "/v1/apps/app1/approval-workflow")
				},
				Config: appApprovalBindingConfig(server.URL, "app1", "wf1"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_app_approval_workflow.test", plancheck.ResourceActionCreate),
					},
				},
			},
			{Config: fakeProviderConfig(server.URL) + removedAppApprovalBinding},
		},
	})
}

// TestEntitlementApprovalWorkflowDestroy checks that destroying the
// entitlement approval override binding sends DELETE, which turns the
// override off, and that destroy also succeeds when the API answers 404
// (the entitlement has no request configuration). It uses a fake API so
// the 404 can be produced on demand.
func TestEntitlementApprovalWorkflowDestroy(t *testing.T) {
	for _, status := range []int{http.StatusNoContent, http.StatusNotFound} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			fake := newFakeSlotAPI(t)
			fake.deleteStatus = status
			server := httptest.NewServer(fake)
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{Config: fakeProviderConfig(server.URL) + `
resource "lumos_entitlement_approval_workflow" "test" {
  entitlement_id = "ent1"
  workflow_id    = "wf1"
}
`},
					{Config: fakeProviderConfig(server.URL)},
				},
			})

			want := []string{"DELETE /v1/entitlements/ent1/approval-workflow-override"}
			if got := fake.requestsWith(http.MethodDelete); len(got) != 1 || got[0] != want[0] {
				t.Errorf("sent %v, want %v", got, want)
			}
		})
	}
}
