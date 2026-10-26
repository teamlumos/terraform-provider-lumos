package hooks

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/config"
)

// The check needs a server that answers the entitlement and app settings
// reads; a fake stands in so each answer can be set per case.
func TestEntitlementSettingsUnlisted(t *testing.T) {
	notFound := func(w http.ResponseWriter) { w.WriteHeader(http.StatusNotFound) }
	entitlement := func(w http.ResponseWriter) { w.Write([]byte(`{"id":"ent1","app_id":"app1"}`)) }
	tests := []struct {
		name        string
		operationID string
		status      int
		entitlement func(w http.ResponseWriter)
		appSettings func(w http.ResponseWriter)
		wantStatus  int
		wantErr     string
	}{
		{
			name: "unlisted app fails the read", operationID: "getEntitlementAppStoreSettings", status: http.StatusNotFound,
			entitlement: entitlement,
			appSettings: func(w http.ResponseWriter) { w.Write([]byte(`{"in_app_store":false}`)) },
			wantErr:     "the entitlement's app (app1) isn't in the App Store: list the app again, or stop managing this resource with a removed block that sets destroy = false",
		},
		{
			// The entitlement has no settings of its own.
			name: "listed app stays a 404", operationID: "getEntitlementAppStoreSettings", status: http.StatusNotFound,
			entitlement: entitlement,
			appSettings: func(w http.ResponseWriter) { w.Write([]byte(`{"in_app_store":true}`)) },
			wantStatus:  http.StatusNotFound,
		},
		{
			name: "deleted entitlement stays a 404", operationID: "getEntitlementAppStoreSettings", status: http.StatusNotFound,
			entitlement: notFound,
			wantStatus:  http.StatusNotFound,
		},
		{
			name: "deleted app stays a 404", operationID: "getEntitlementAppStoreSettings", status: http.StatusNotFound,
			entitlement: entitlement,
			appSettings: notFound,
			wantStatus:  http.StatusNotFound,
		},
		{
			name: "failed entitlement read fails the read", operationID: "getEntitlementAppStoreSettings", status: http.StatusNotFound,
			entitlement: func(w http.ResponseWriter) { w.WriteHeader(http.StatusInternalServerError) },
			wantErr:     "reading the entitlement to check its app failed: status 500",
		},
		{
			name: "failed app settings read fails the read", operationID: "getEntitlementAppStoreSettings", status: http.StatusNotFound,
			entitlement: entitlement,
			appSettings: func(w http.ResponseWriter) { w.WriteHeader(http.StatusForbidden) },
			wantErr:     "reading its app's (app1) App Store settings failed: status 403",
		},
		{name: "other status codes pass through", operationID: "getEntitlementAppStoreSettings", status: http.StatusOK, wantStatus: http.StatusOK},
		{name: "other operations pass through", operationID: "getEntitlement", status: http.StatusNotFound, wantStatus: http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
					t.Errorf("lookup Authorization = %q, want the settings read's", got)
				}
				var answer func(w http.ResponseWriter)
				switch {
				case r.Method == http.MethodGet && r.URL.Path == "/v1/entitlements/ent1" && r.URL.RawQuery == "":
					answer = tt.entitlement
				case r.Method == http.MethodGet && r.URL.Path == "/v1/apps/app1/appstore-settings":
					answer = tt.appSettings
				}
				if answer == nil {
					t.Errorf("unexpected lookup %s %s", r.Method, r.URL)
					w.WriteHeader(http.StatusTeapot)
					return
				}
				answer(w)
			}))
			defer server.Close()

			req, err := http.NewRequest(http.MethodGet, server.URL+"/v1/entitlements/ent1/appstore-settings", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", "Bearer test-token")
			body := `{"detail":"AppStore settings are not available for this entitlement."}`
			res := &http.Response{StatusCode: tt.status, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(body)), Request: req}
			hookCtx := AfterSuccessContext{HookContext: HookContext{
				OperationID:      tt.operationID,
				Context:          context.Background(),
				SDKConfiguration: config.SDKConfiguration{Client: server.Client()},
			}}

			got, err := New().AfterSuccess(hookCtx, res)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("err = %v, want one containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", got.StatusCode, tt.wantStatus)
			}
			if b, _ := io.ReadAll(got.Body); string(b) != body {
				t.Errorf("body = %q, want the original %q", b, body)
			}
		})
	}
}
