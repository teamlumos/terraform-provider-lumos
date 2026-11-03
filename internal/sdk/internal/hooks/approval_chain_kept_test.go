package hooks

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestApprovalChainKept(t *testing.T) {
	chainBody := `{"error":{"detail":"An approval chain can't be deleted on its own. Edit the approval workflow that uses it instead.","type":"validation","code":"workflow.approval_chain_delete_not_supported","errors":[]}}`
	tests := []struct {
		name        string
		operationID string
		status      int
		body        string
		want        int
	}{
		{name: "refused chain delete counts as deleted", operationID: "deleteWorkflow", status: http.StatusUnprocessableEntity, body: chainBody, want: http.StatusNoContent},
		{
			name: "other 422 stays an error", operationID: "deleteWorkflow", status: http.StatusUnprocessableEntity,
			body: `{"error":{"detail":"Invalid.","type":"validation","code":"bad_request","errors":[]}}`,
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "bound workflow stays a conflict", operationID: "deleteWorkflow", status: http.StatusConflict,
			body: `{"error":{"detail":"This workflow is still bound as an app's approval workflow.","type":"conflict","code":"workflow.bound_delete_not_allowed","errors":[]}}`,
			want: http.StatusConflict,
		},
		{
			name: "default workflow stays a conflict", operationID: "deleteWorkflow", status: http.StatusConflict,
			body: `{"error":{"detail":"This workflow is an app's default workflow, so it can't be deleted.","type":"conflict","code":"workflow.default_delete_not_allowed","errors":[]}}`,
			want: http.StatusConflict,
		},
		{name: "unreadable body stays an error", operationID: "deleteWorkflow", status: http.StatusUnprocessableEntity, body: "not json", want: http.StatusUnprocessableEntity},
		{name: "other operations pass through", operationID: "deleteEntitlement", status: http.StatusUnprocessableEntity, body: chainBody, want: http.StatusUnprocessableEntity},
		{name: "successful delete passes through", operationID: "deleteWorkflow", status: http.StatusNoContent, want: http.StatusNoContent},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, "https://api.example.com/v1/workflows/123", nil)
			if err != nil {
				t.Fatal(err)
			}
			res := &http.Response{StatusCode: tt.status, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(tt.body)), Request: req}
			hookCtx := AfterSuccessContext{HookContext: HookContext{OperationID: tt.operationID, Context: context.Background()}}

			got, err := approvalChainKept{}.AfterSuccess(hookCtx, res)
			if err != nil {
				t.Fatal(err)
			}
			if got.StatusCode != tt.want {
				t.Errorf("status = %d, want %d", got.StatusCode, tt.want)
			}
			if got.StatusCode == tt.status {
				if body, _ := io.ReadAll(got.Body); string(body) != tt.body {
					t.Errorf("body = %q, want the original %q", body, tt.body)
				}
			}
		})
	}
}
