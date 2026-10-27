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
func TestAppCreateBody(t *testing.T) {
	tests := []struct {
		name        string
		operationID string
		body        string
		want        string
	}{
		{
			name:        `drops a custom app's app_class_id = "custom"`,
			operationID: "v1CreateApp",
			body:        `{"app_class_id":"custom","label":"Wiki","category":"Developers","description":"Docs"}`,
			want:        `{"label":"Wiki","category":"Developers","description":"Docs"}`,
		},
		{
			name:        "keeps an integration's app_class_id",
			operationID: "v1CreateApp",
			body:        `{"app_class_id":"slack.com","auth":{"token":"x"}}`,
			want:        `{"app_class_id":"slack.com","auth":{"token":"x"}}`,
		},
		{
			name:        "leaves a custom app without app_class_id alone",
			operationID: "v1CreateApp",
			body:        `{"label":"Wiki","category":"Developers","description":"Docs"}`,
			want:        `{"label":"Wiki","category":"Developers","description":"Docs"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "https://api.example.com/v1/apps", strings.NewReader(tt.body))
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
