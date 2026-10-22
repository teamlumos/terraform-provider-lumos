package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestLookupDataSourcesRequireFilter checks that a lookup of a large
// collection with no narrowing filter fails at plan without calling the API.
func TestLookupDataSourcesRequireFilter(t *testing.T) {
	tests := []struct {
		dataSource string
		config     string
		wantErr    string
	}{
		{dataSource: "lumos_users", config: `exact_match = true`, wantErr: "lumos_users requires search"},
		{dataSource: "lumos_entitlements", config: `visible_in_appstore = true`, wantErr: "lumos_entitlements requires search"},
		{dataSource: "lumos_accounts", config: `status = ["active"]`, wantErr: "lumos_accounts requires app_id or email"},
		{dataSource: "lumos_groups", config: `exact_match = true`, wantErr: "lumos_groups requires name, integration_specific_id or app_id"},
	}
	for _, tt := range tests {
		t.Run(tt.dataSource, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				http.Error(w, "unexpected "+r.Method+" "+r.URL.String(), http.StatusTeapot)
			}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{{
					Config: fmt.Sprintf(`
provider "lumos" {
  server_url  = %q
  http_bearer = "test-token"
}

data %q "test" {
  %s
}
`, server.URL, tt.dataSource, tt.config),
					ExpectError: regexp.MustCompile(regexp.QuoteMeta(tt.wantErr)),
				}},
			})
			if n := calls.Load(); n != 0 {
				t.Fatalf("the API was called %d times, want none", n)
			}
		})
	}
}
