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
// without it, removing a clearable setting from the config would leave it in
// place.
func TestClearUnsetFields(t *testing.T) {
	tests := []struct {
		name        string
		operationID string
		body        string
		want        string
	}{
		{
			name:        "sends null for an unset description and approval message",
			operationID: "updateEntitlement",
			body:        `{"label":"Admin"}`,
			want:        `{"label":"Admin","description_override":null,"custom_approval_message_override":null}`,
		},
		{
			name:        "keeps a configured description and approval message",
			operationID: "updateEntitlement",
			body:        `{"label":"Admin","description_override":"Ask first.","custom_approval_message_override":"Approved."}`,
			want:        `{"label":"Admin","description_override":"Ask first.","custom_approval_message_override":"Approved."}`,
		},
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
			name:        "sends null for each unset app setting and keeps the rest",
			operationID: "v1PatchAppStoreSettings",
			body:        `{"in_app_store":true,"redirect_url":"https://example.com","time_based_access":{"options":["1 day"]}}`,
			want:        `{"in_app_store":true,"redirect_url":"https://example.com","time_based_access":{"options":["1 day"],"default":null},"request_instructions":null,"request_validation_webhook_uuid":null,"default_entitlement":null}`,
		},
		{
			name:        "keeps configured app settings",
			operationID: "v1PatchAppStoreSettings",
			body:        `{"request_instructions":"Ask first.","request_validation_webhook_uuid":"w1","redirect_url":"https://example.com","default_entitlement":{"id":"e1"},"time_based_access":{"default":"1 day"}}`,
			want:        `{"request_instructions":"Ask first.","request_validation_webhook_uuid":"w1","redirect_url":"https://example.com","default_entitlement":{"id":"e1"},"time_based_access":{"default":"1 day"}}`,
		},
		{
			name:        "leaves time-based access out when the body has none",
			operationID: "v1PatchAppStoreSettings",
			body:        `{}`,
			want:        `{"request_instructions":null,"request_validation_webhook_uuid":null,"redirect_url":null,"default_entitlement":null}`,
		},
		{
			name:        "leaves other operations alone",
			operationID: "v1GetAppStoreSettings",
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
