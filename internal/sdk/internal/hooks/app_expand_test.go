package hooks

import (
	"net/http"
	"reflect"
	"testing"
)

// Runs through New() so the test also fails if the hook isn't registered.
func TestReadExpand(t *testing.T) {
	tests := []struct {
		name        string
		operationID string
		rawQuery    string
		want        []string
	}{
		{name: "adds expand to an app read", operationID: "v1GetApp", rawQuery: "", want: []string{"appstore_settings", "custom_attributes"}},
		{name: "keeps other expansions", operationID: "v1GetApp", rawQuery: "expand=admins", want: []string{"admins", "appstore_settings", "custom_attributes"}},
		{name: "doesn't repeat them", operationID: "v1GetApp", rawQuery: "expand=custom_attributes&expand=appstore_settings", want: []string{"custom_attributes", "appstore_settings"}},
		{name: "adds expand to an entitlement read", operationID: "getEntitlement", rawQuery: "", want: []string{"appstore_settings"}},
		{name: "adds expand to an app list", operationID: "v1ListApps", rawQuery: "page_size=100", want: []string{"appstore_settings", "custom_attributes"}},
		{name: "leaves other operations alone", operationID: "listAppEntitlements", rawQuery: "", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "https://api.example.com/v1/apps/123?"+tt.rawQuery, nil)
			if err != nil {
				t.Fatal(err)
			}
			hookCtx := BeforeRequestContext{HookContext: HookContext{OperationID: tt.operationID}}

			got, err := New().BeforeRequest(hookCtx, req)
			if err != nil {
				t.Fatal(err)
			}
			if expand := got.URL.Query()["expand"]; !reflect.DeepEqual(expand, tt.want) {
				t.Errorf("expand = %v, want %v", expand, tt.want)
			}
		})
	}
}
