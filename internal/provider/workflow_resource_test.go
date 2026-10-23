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

// fakeWorkflowAPI stores workflow "wf1" the way the API does: step inputs
// set to null are left out when read back, an omitted description is
// stored and read back as "", each step reads back with a name, and an
// update whose app_id differs from the stored one is refused.
type fakeWorkflowAPI struct {
	t        *testing.T
	mu       sync.Mutex
	workflow map[string]any
	bodies   []map[string]any
}

func (f *fakeWorkflowAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	switch {
	case r.URL.Path == "/v1/workflows" && r.Method == http.MethodPost,
		r.URL.Path == "/v1/workflows/wf1" && r.Method == http.MethodPut:
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		if err := json.Unmarshal(raw, &body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		f.bodies = append(f.bodies, body)
		if r.Method == http.MethodPut && body["app_id"] != f.workflow["app_id"] {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, `{"error":{"code":"workflow.app_reassignment_not_allowed","message":"app_id may not be updated directly after creation."}}`)
			return
		}
		steps := []any{}
		for _, s := range body["steps"].([]any) {
			step := s.(map[string]any)
			inputs := map[string]any{}
			if in, ok := step["inputs"].(map[string]any); ok {
				for k, v := range in {
					if v != nil {
						inputs[k] = v
					}
				}
			}
			steps = append(steps, map[string]any{"step": step["step"], "name": "Approval", "inputs": inputs})
		}
		description, _ := body["description"].(string)
		f.workflow = map[string]any{
			"id": "wf1", "name": body["name"], "description": description,
			"type": body["type"], "app_id": body["app_id"], "steps": steps,
		}
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusCreated)
		}
		json.NewEncoder(w).Encode(f.workflow)
	case r.URL.Path == "/v1/workflows/wf1" && r.Method == http.MethodGet:
		if f.workflow == nil {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, `{"error":{"code":"not_found","message":"Workflow not found"}}`)
			return
		}
		json.NewEncoder(w).Encode(f.workflow)
	case r.URL.Path == "/v1/workflows/wf1" && r.Method == http.MethodDelete:
		f.workflow = nil
		w.WriteHeader(http.StatusNoContent)
	default:
		f.t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		http.Error(w, "unexpected request", http.StatusNotFound)
	}
}

// workflowConfig is an approval chain with one step. extra goes inside the
// resource block, for attributes like description.
func workflowConfig(serverURL, name, step, inputs, extra string) string {
	return fakeProviderConfig(serverURL) + fmt.Sprintf(`
resource "lumos_workflow" "chain" {
  name   = %q
  type   = "approval-chain"
  app_id = "app1"
  %s
  steps = [{
    step   = %q
    inputs = jsonencode(%s)
  }]
}
`, name, extra, step, inputs)
}

const approvalInputs = `{
      stage_number = 1
      assignees    = [{ uuid = "group1", type = "group" }]
      instructions = null
    }`

// TestWorkflowNullInputKeys checks that a step input set to null in the
// config matches the API leaving that key out: apply succeeds with no
// "inconsistent result" error, the next plan is empty, and a real change
// still plans an update. It uses a fake API so the check covers only how the
// provider compares inputs.
func TestWorkflowNullInputKeys(t *testing.T) {
	fake := &fakeWorkflowAPI{t: t}
	server := httptest.NewServer(fake)
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: workflowConfig(server.URL, "IT approval", "approval_task", approvalInputs, "")},
			{
				Config: workflowConfig(server.URL, "IT approval", "approval_task", `{
      stage_number = 2
      assignees    = [{ uuid = "group1", type = "group" }]
      instructions = null
    }`, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_workflow.chain", plancheck.ResourceActionUpdate),
					},
				},
			},
		},
	})

	// The SDK sent inputs as a JSON object, not a string.
	for _, body := range fake.bodies {
		inputs := body["steps"].([]any)[0].(map[string]any)["inputs"]
		if _, ok := inputs.(map[string]any); !ok {
			t.Errorf("inputs sent as %T, want a JSON object", inputs)
		}
	}
}

// TestWorkflowDescription checks that a config without a description
// matches the "" the API stores for it, before and after one is set, and
// that description = "" fails at plan time. It uses a fake API so the check
// covers only the provider's handling of "".
func TestWorkflowDescription(t *testing.T) {
	fake := &fakeWorkflowAPI{t: t}
	server := httptest.NewServer(fake)
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: workflowConfig(server.URL, "IT approval", "approval_task", approvalInputs, ""),
				Check:  resource.TestCheckNoResourceAttr("lumos_workflow.chain", "description"),
			},
			{
				Config: workflowConfig(server.URL, "IT approval", "approval_task", approvalInputs, `description = "One stage"`),
				Check:  resource.TestCheckResourceAttr("lumos_workflow.chain", "description", "One stage"),
			},
			{
				Config: workflowConfig(server.URL, "IT approval", "approval_task", approvalInputs, ""),
				Check:  resource.TestCheckNoResourceAttr("lumos_workflow.chain", "description"),
			},
		},
	})
}

// TestWorkflowEmptyDescription checks that description = "" fails at plan
// time, since the API would read it back as no description. The check runs
// before any API call, so the server URL is never contacted.
func TestWorkflowEmptyDescription(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:      workflowConfig("http://127.0.0.1:1", "IT approval", "approval_task", approvalInputs, `description = ""`),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value Length.*description`),
		}},
	})
}

// TestWorkflowRenamePlan checks that renaming a workflow plans only the
// name, with nothing else shown as known after apply. It uses a fake API so
// the check covers only the provider's plan.
func TestWorkflowRenamePlan(t *testing.T) {
	fake := &fakeWorkflowAPI{t: t}
	server := httptest.NewServer(fake)
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: workflowConfig(server.URL, "IT approval", "approval_task", approvalInputs, `description = "One stage"`)},
			{
				Config: workflowConfig(server.URL, "Security approval", "approval_task", approvalInputs, `description = "One stage"`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						expectOnlyAttributesChange("lumos_workflow.chain", "name"),
					},
				},
			},
		},
	})
}

// TestWorkflowAppAssignedByAPI checks a workflow created without app_id,
// then given an app by the API, as binding does. Lumos can't clear that
// app_id, so a config that still leaves it out fails at plan time, saying
// which app to add; once the config sets it, the next update plans only the
// rename and sends that app_id. It uses a fake API so the app can be
// assigned between steps.
func TestWorkflowAppAssignedByAPI(t *testing.T) {
	fake := &fakeWorkflowAPI{t: t}
	server := httptest.NewServer(fake)
	defer server.Close()

	config := func(name, appID string) string {
		return fakeProviderConfig(server.URL) + fmt.Sprintf(`
resource "lumos_workflow" "provisioning" {
  name  = %q
  type  = "account-provisioning"
  %s
  steps = [{
    step   = "wait"
    inputs = jsonencode({ duration = "P1D" })
  }]
}
`, name, appID)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config("Provision", ""),
				Check:  resource.TestCheckNoResourceAttr("lumos_workflow.provisioning", "app_id"),
			},
			{
				PreConfig: func() {
					fake.mu.Lock()
					defer fake.mu.Unlock()
					fake.workflow["app_id"] = "app1"
				},
				Config:      config("Provision accounts", ""),
				ExpectError: regexp.MustCompile(`(?s)app_id can't be cleared.*Add\s+app_id = "app1"`),
			},
			{
				Config: config("Provision accounts", `app_id = "app1"`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("lumos_workflow.provisioning", plancheck.ResourceActionUpdate),
						expectOnlyAttributesChange("lumos_workflow.provisioning", "name"),
					},
				},
				Check: resource.TestCheckResourceAttr("lumos_workflow.provisioning", "app_id", "app1"),
			},
		},
	})
}

// TestWorkflowImport checks that importing a workflow by ID reads back
// every attribute the config set. Its inputs have no null keys, since an
// import reads inputs as the API returns them. It uses a fake API so the
// read-back covers only the provider's conversions.
func TestWorkflowImport(t *testing.T) {
	fake := &fakeWorkflowAPI{t: t}
	server := httptest.NewServer(fake)
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: workflowConfig(server.URL, "IT approval", "approval_task", `{
      stage_number = 1
      assignees    = [{ uuid = "group1", type = "group" }]
    }`, `description = "One stage"`)},
			{
				ResourceName:      "lumos_workflow.chain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestWorkflowStepRequired checks that a step without its step slug fails at
// plan time. The check runs before any API call, so the server URL is never
// contacted.
func TestWorkflowStepRequired(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: fakeProviderConfig("http://127.0.0.1:1") + `
resource "lumos_workflow" "chain" {
  name   = "IT approval"
  type   = "approval-chain"
  app_id = "app1"
  steps  = [{ inputs = jsonencode({ stage_number = 1 }) }]
}
`,
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(`(?s)Missing Attribute Value.*steps\[0\]\.step`),
		}},
	})
}
