package hooks

import (
	"net/http"
	"reflect"
	"testing"
)

func TestAppStoreSettingsExpand(t *testing.T) {
	tests := []struct {
		name        string
		operationID string
		rawQuery    string
		want        []string
	}{
		{name: "adds expand to an app read", operationID: "v1GetApp", rawQuery: "", want: []string{"appstore_settings"}},
		{name: "keeps other expansions", operationID: "v1GetApp", rawQuery: "expand=admins", want: []string{"admins", "appstore_settings"}},
		{name: "doesn't repeat it", operationID: "v1GetApp", rawQuery: "expand=appstore_settings", want: []string{"appstore_settings"}},
		{name: "leaves other operations alone", operationID: "v1ListApps", rawQuery: "", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "https://api.example.com/v1/apps/123?"+tt.rawQuery, nil)
			if err != nil {
				t.Fatal(err)
			}
			hookCtx := BeforeRequestContext{HookContext: HookContext{OperationID: tt.operationID}}

			got, err := appStoreSettingsExpand{}.BeforeRequest(hookCtx, req)
			if err != nil {
				t.Fatal(err)
			}
			if expand := got.URL.Query()["expand"]; !reflect.DeepEqual(expand, tt.want) {
				t.Errorf("expand = %v, want %v", expand, tt.want)
			}
		})
	}
}
