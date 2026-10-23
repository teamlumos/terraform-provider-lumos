package hooks

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

// Runs through New() so the test also fails if the hook isn't registered.
func TestRateLimited(t *testing.T) {
	tests := []struct {
		name        string
		operationID string
		status      int
		retryAfter  string
		wantErr     string
	}{
		{name: "passes a successful response", operationID: "v1GetCurrentUser", status: http.StatusOK},
		{name: "leaves other errors alone", operationID: "v1GetCurrentUser", status: http.StatusConflict},
		{
			name: "fails a 429 with its Retry-After", operationID: "v1GetCurrentUser", status: http.StatusTooManyRequests, retryAfter: "300",
			wantErr: `v1GetCurrentUser: the Lumos API is still rate limiting this request (status 429, Retry-After: 300) after the provider retried it for up to 2 minutes. Re-run later, or lower terraform -parallelism. Response: {"detail":"slow down"}`,
		},
		{
			name: "fails a 429 without Retry-After", operationID: "v1GetCurrentUser", status: http.StatusTooManyRequests,
			wantErr: `v1GetCurrentUser: the Lumos API is still rate limiting this request (status 429, Retry-After: not sent) after the provider retried it for up to 2 minutes. Re-run later, or lower terraform -parallelism. Response: {"detail":"slow down"}`,
		},
		{
			name: "fails a rate-limited list page before listPageError does", operationID: "listGroups", status: http.StatusTooManyRequests, retryAfter: "1",
			wantErr: `listGroups: the Lumos API is still rate limiting this request (status 429, Retry-After: 1)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &http.Response{
				StatusCode: tt.status,
				Header:     http.Header{},
				Body:       io.NopCloser(strings.NewReader(`{"detail":"slow down"}`)),
			}
			if tt.retryAfter != "" {
				res.Header.Set("Retry-After", tt.retryAfter)
			}
			hookCtx := AfterSuccessContext{HookContext: HookContext{OperationID: tt.operationID}}

			got, err := New().AfterSuccess(hookCtx, res)
			if tt.wantErr == "" {
				if err != nil || got != res {
					t.Fatalf("got (%v, %v), want the response unchanged", got, err)
				}
				return
			}
			if err == nil || !strings.HasPrefix(err.Error(), tt.wantErr) {
				t.Fatalf("error = %v, want prefix %q", err, tt.wantErr)
			}
		})
	}
}
