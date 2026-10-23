package hooks

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// Runs through New() so the test also fails if the hook isn't registered.
func TestAppUpdateBody(t *testing.T) {
	tests := []struct {
		name        string
		operationID string
		body        string
		want        string
	}{
		{
			name:        "sends null for a custom app's unset logo and website",
			operationID: "v1UpdateApp",
			body:        `{"app_class_id":"custom","label":"Wiki","category":"Developers","description":"Docs"}`,
			want:        `{"label":"Wiki","category":"Developers","description":"Docs","logo_url":null,"website_url":null}`,
		},
		{
			name:        "keeps a custom app's configured logo and website",
			operationID: "v1UpdateApp",
			body:        `{"app_class_id":"custom","label":"Wiki","category":"Developers","description":"Docs","logo_url":"https://example.com/logo.png","website_url":"https://example.com"}`,
			want:        `{"label":"Wiki","category":"Developers","description":"Docs","logo_url":"https://example.com/logo.png","website_url":"https://example.com"}`,
		},
		{
			name:        "drops a connected app's metadata and keeps its credentials",
			operationID: "v1UpdateApp",
			body:        `{"app_class_id":"slack.com","auth":{"token":"x"},"label":"Slack","category":"Communication","description":"Team chat","logo_url":"https://example.com/slack.png","website_url":"https://slack.com"}`,
			want:        `{"auth":{"token":"x"}}`,
		},
		{
			name:        "leaves other operations alone",
			operationID: "v1CreateApp",
			body:        `{"app_class_id":"slack.com","label":"Slack"}`,
			want:        `{"app_class_id":"slack.com","label":"Slack"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPatch, "https://api.example.com/v1/apps/123", strings.NewReader(tt.body))
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
		})
	}
}
