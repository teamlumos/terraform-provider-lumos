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
			name:        "sends null for a custom app's unset overrides",
			operationID: "v1UpdateApp",
			body:        `{"app_class_id":"custom","label_override":"Wiki","category_override":"Developers","description_override":"Docs"}`,
			want:        `{"label_override":"Wiki","category_override":"Developers","description_override":"Docs","logo_url_override":null,"website_url_override":null}`,
		},
		{
			name:        "keeps a custom app's configured overrides",
			operationID: "v1UpdateApp",
			body:        `{"app_class_id":"custom","label_override":"Wiki","category_override":"Developers","description_override":"Docs","logo_url_override":"https://example.com/logo.png","website_url_override":"https://example.com"}`,
			want:        `{"label_override":"Wiki","category_override":"Developers","description_override":"Docs","logo_url_override":"https://example.com/logo.png","website_url_override":"https://example.com"}`,
		},
		{
			name:        "sends a connected app's credentials and configured overrides, and null for the rest",
			operationID: "v1UpdateApp",
			body:        `{"app_class_id":"slack.com","auth":{"token":"x"},"label_override":"Slack"}`,
			want:        `{"auth":{"token":"x"},"label_override":"Slack","description_override":null,"category_override":null,"logo_url_override":null,"website_url_override":null}`,
		},
		{
			name:        "leaves other operations alone",
			operationID: "v1CreateApp",
			body:        `{"app_class_id":"slack.com","label_override":"Slack"}`,
			want:        `{"app_class_id":"slack.com","label_override":"Slack"}`,
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
