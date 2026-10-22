package hooks

import (
	"net/http"
	"testing"
)

// Runs through New() so the test also fails if the hook isn't registered.
func TestRequireListFilter(t *testing.T) {
	tests := []struct {
		name        string
		operationID string
		rawQuery    string
		wantErr     string
	}{
		{name: "fails a user list without a search", operationID: "v1ListUsers", rawQuery: "page_size=100", wantErr: "lumos_users requires search"},
		{name: "fails an empty search", operationID: "v1ListUsers", rawQuery: "search=", wantErr: "lumos_users requires search"},
		{name: "passes a user search", operationID: "v1ListUsers", rawQuery: "search=alice"},
		{name: "fails an entitlement list filtered only by visibility", operationID: "listEntitlements", rawQuery: "visible_in_appstore=true", wantErr: "lumos_entitlements requires search"},
		{name: "passes an entitlement search", operationID: "listEntitlements", rawQuery: "search=Admin"},
		{name: "fails an account list without an app or email", operationID: "listAccounts", rawQuery: "status=ACTIVE", wantErr: "lumos_accounts requires app_id or email"},
		{name: "passes accounts of one app", operationID: "listAccounts", rawQuery: "app_id=123"},
		{name: "passes accounts by email", operationID: "listAccounts", rawQuery: "email=alice%40example.com"},
		{name: "fails a group list without a filter", operationID: "listGroups", rawQuery: "exact_match=true", wantErr: "lumos_groups requires name, integration_specific_id or app_id"},
		{name: "passes groups by name", operationID: "listGroups", rawQuery: "name=Engineering"},
		{name: "passes groups by integration id", operationID: "listGroups", rawQuery: "integration_specific_id=abc"},
		{name: "passes groups of one app", operationID: "listGroups", rawQuery: "app_id=123"},
		{name: "leaves other lists alone", operationID: "v1ListApps", rawQuery: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "https://api.example.com/v1/list?"+tt.rawQuery, nil)
			if err != nil {
				t.Fatal(err)
			}
			hookCtx := BeforeRequestContext{HookContext: HookContext{OperationID: tt.operationID}}

			got, err := New().BeforeRequest(hookCtx, req)
			if tt.wantErr == "" {
				if err != nil || got == nil {
					t.Fatalf("got (%v, %v), want the request sent", got, err)
				}
				return
			}
			if err == nil || err.Error() != tt.wantErr {
				t.Fatalf("error = %v, want %q", err, tt.wantErr)
			}
		})
	}
}
