package hooks

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// Runs through New() so the test also fails if the hook isn't registered:
// without it, removing an override from the config would leave it in place.
func TestInheritUnsetEntitlementOverrides(t *testing.T) {
	tests := []struct {
		name        string
		operationID string
		body        string
		want        string
	}{
		{
			name:        "sends null for every unset override",
			operationID: "updateEntitlementAppStoreSettings",
			body:        `{}`,
			want:        `{"allowed_requester_group_uuids_override":null,"time_based_access_override":null,"request_validation_webhook_uuid":null}`,
		},
		{
			name:        "keeps an empty group list, which means anyone may request",
			operationID: "updateEntitlementAppStoreSettings",
			body:        `{"allowed_requester_group_uuids_override":[]}`,
			want:        `{"allowed_requester_group_uuids_override":[],"time_based_access_override":null,"request_validation_webhook_uuid":null}`,
		},
		{
			name:        "keeps configured overrides and doesn't add visible_in_appstore",
			operationID: "updateEntitlementAppStoreSettings",
			body:        `{"allowed_requester_group_uuids_override":["g1"],"time_based_access_override":{"options":["1 day"]},"request_validation_webhook_uuid":"w1"}`,
			want:        `{"allowed_requester_group_uuids_override":["g1"],"time_based_access_override":{"options":["1 day"]},"request_validation_webhook_uuid":"w1"}`,
		},
		{
			name:        "keeps visible_in_appstore when it's set",
			operationID: "updateEntitlementAppStoreSettings",
			body:        `{"visible_in_appstore":true,"request_validation_webhook_uuid":"w1"}`,
			want:        `{"visible_in_appstore":true,"allowed_requester_group_uuids_override":null,"time_based_access_override":null,"request_validation_webhook_uuid":"w1"}`,
		},
		{
			name:        "leaves other operations alone",
			operationID: "v1PatchAppStoreSettings",
			body:        `{}`,
			want:        `{}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPatch, "https://api.example.com/v1/entitlements/456/appstore-settings", strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}
			hookCtx := BeforeRequestContext{HookContext: HookContext{OperationID: tt.operationID}}

			got, err := New().BeforeRequest(hookCtx, req)
			if err != nil {
				t.Fatal(err)
			}
			body, err := io.ReadAll(got.Body)
			if err != nil {
				t.Fatal(err)
			}
			var gotFields, wantFields map[string]any
			if err := json.Unmarshal(body, &gotFields); err != nil {
				t.Fatalf("body %s: %v", body, err)
			}
			if err := json.Unmarshal([]byte(tt.want), &wantFields); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(gotFields, wantFields) {
				t.Errorf("body = %s, want %s", body, tt.want)
			}
			if got.ContentLength != int64(len(body)) {
				t.Errorf("ContentLength = %d, want %d", got.ContentLength, len(body))
			}
		})
	}
}
