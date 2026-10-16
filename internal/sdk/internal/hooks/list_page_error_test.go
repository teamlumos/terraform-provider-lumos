package hooks

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

// Runs through New() so the test also fails if the hook isn't registered.
func TestListPageError(t *testing.T) {
	tests := []struct {
		name        string
		operationID string
		status      int
		wantErr     string
	}{
		{name: "passes a successful page", operationID: "v1ListUsers", status: http.StatusOK},
		{name: "fails an error page with its status and body", operationID: "v1ListUsers", status: http.StatusBadRequest, wantErr: `v1ListUsers: status 400: {"detail":"page token expired"}`},
		{name: "leaves other operations alone", operationID: "v1GetUser", status: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &http.Response{
				StatusCode: tt.status,
				Body:       io.NopCloser(strings.NewReader(`{"detail":"page token expired"}`)),
			}
			hookCtx := AfterSuccessContext{HookContext: HookContext{OperationID: tt.operationID}}

			got, err := New().AfterSuccess(hookCtx, res)
			if tt.wantErr == "" {
				if err != nil || got != res {
					t.Fatalf("got (%v, %v), want the response unchanged", got, err)
				}
				return
			}
			if err == nil || err.Error() != tt.wantErr {
				t.Fatalf("error = %v, want %q", err, tt.wantErr)
			}
		})
	}
}
