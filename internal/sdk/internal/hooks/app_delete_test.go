package hooks

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/config"
)

// The read-back needs a server that answers GET /v1/apps/{id}; the API
// can't be made to keep a record on demand, so a fake stands in for it.
func TestAppDeleteRecordKept(t *testing.T) {
	tests := []struct {
		name        string
		operationID string
		status      int
		readBack    func(w http.ResponseWriter)
		want        int
	}{
		{
			name: "disconnected but kept counts as deleted", operationID: "v1DeleteApp", status: http.StatusConflict,
			readBack: func(w http.ResponseWriter) { w.Write([]byte(`{"id":"123","disconnected":true}`)) },
			want:     http.StatusNoContent,
		},
		{
			name: "still connected stays a conflict", operationID: "v1DeleteApp", status: http.StatusConflict,
			readBack: func(w http.ResponseWriter) { w.Write([]byte(`{"id":"123","disconnected":false}`)) },
			want:     http.StatusConflict,
		},
		{
			name: "custom app kept unchanged stays a conflict", operationID: "v1DeleteApp", status: http.StatusConflict,
			readBack: func(w http.ResponseWriter) {
				w.Write([]byte(`{"id":"123","app_class_id":"custom","disconnected":false}`))
			},
			want: http.StatusConflict,
		},
		{
			name: "failed read-back stays a conflict", operationID: "v1DeleteApp", status: http.StatusConflict,
			readBack: func(w http.ResponseWriter) { w.WriteHeader(http.StatusInternalServerError) },
			want:     http.StatusConflict,
		},
		{name: "other status codes pass through", operationID: "v1DeleteApp", status: http.StatusNoContent, want: http.StatusNoContent},
		{name: "other operations pass through", operationID: "deleteEntitlement", status: http.StatusConflict, want: http.StatusConflict},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotAuth string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != "/v1/apps/123" {
					t.Errorf("unexpected read-back %s %s", r.Method, r.URL.Path)
				}
				gotAuth = r.Header.Get("Authorization")
				tt.readBack(w)
			}))
			defer server.Close()

			req, err := http.NewRequest(http.MethodDelete, server.URL+"/v1/apps/123", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", "Bearer test-token")
			res := &http.Response{StatusCode: tt.status, Header: http.Header{}, Body: http.NoBody, Request: req}
			hookCtx := AfterSuccessContext{HookContext: HookContext{
				OperationID:      tt.operationID,
				Context:          context.Background(),
				SDKConfiguration: config.SDKConfiguration{Client: server.Client()},
			}}

			got, err := New().AfterSuccess(hookCtx, res)
			if err != nil {
				t.Fatal(err)
			}
			if got.StatusCode != tt.want {
				t.Errorf("status = %d, want %d", got.StatusCode, tt.want)
			}
			if tt.readBack != nil && gotAuth != "Bearer test-token" {
				t.Errorf("read-back Authorization = %q, want the delete's", gotAuth)
			}
		})
	}
}
