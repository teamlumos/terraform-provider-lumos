package hooks

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

// Runs through New() so the test also fails if the hook isn't registered:
// without it, the request would go to the trailing-slash path with an empty
// body and could succeed without taking the app out of the App Store.
func TestRemoveAppFromAppStore(t *testing.T) {
	tests := []struct {
		name        string
		operationID string
		wantPath    string
		wantBody    string
	}{
		{name: "sends in_app_store false to the settings endpoint", operationID: "v1RemoveAppFromAppStore", wantPath: "/v1/apps/123/appstore-settings", wantBody: `{"in_app_store":false}`},
		{name: "leaves other operations alone", operationID: "v1PatchAppStoreSettings", wantPath: "/v1/apps/123/appstore-settings/", wantBody: `{}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPatch, "https://api.example.com/v1/apps/123/appstore-settings/", strings.NewReader(`{}`))
			if err != nil {
				t.Fatal(err)
			}
			hookCtx := BeforeRequestContext{HookContext: HookContext{OperationID: tt.operationID}}

			got, err := New().BeforeRequest(hookCtx, req)
			if err != nil {
				t.Fatal(err)
			}
			if got.Method != http.MethodPatch {
				t.Errorf("method = %s, want PATCH", got.Method)
			}
			if got.URL.Path != tt.wantPath {
				t.Errorf("path = %s, want %s", got.URL.Path, tt.wantPath)
			}
			body, err := io.ReadAll(got.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(body) != tt.wantBody {
				t.Errorf("body = %s, want %s", body, tt.wantBody)
			}
		})
	}
}
