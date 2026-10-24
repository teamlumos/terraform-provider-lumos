package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// appWithCustomAttributes is an app as the API returns it with
// expand=custom_attributes, with one attribute of each value shape.
const appWithCustomAttributes = `{
  "id": "app-1",
  "label": "Okta",
  "custom_attributes": {
    "Tier": {"type": "text", "value": "gold"},
    "Owners": {"type": "user", "value": [{"id": "user-1", "email": "ada@example.com"}]},
    "Renewal": {"type": "text", "value": "2027-01-31T00:00:00Z"},
    "Seats": {"type": "text", "value": 42},
    "Unset": {"type": "text", "value": null}
  }
}`

// TestAppCustomAttributes runs data "lumos_apps" and data "lumos_app" against
// a fake API and checks that both ask for custom attributes and put each
// value shape in state.
func TestAppCustomAttributes(t *testing.T) {
	tests := []struct {
		dataSource string
		config     string
		wantPath   string
		body       string
		prefix     string
	}{
		{
			dataSource: "lumos_apps",
			config:     `search = "Okta"`,
			wantPath:   "/v1/apps",
			body:       `{"apps":[` + appWithCustomAttributes + `],"next_page_token":null}`,
			prefix:     "apps.0.",
		},
		{
			dataSource: "lumos_app",
			config:     `id = "app-1"`,
			wantPath:   "/v1/apps/app-1",
			body:       appWithCustomAttributes,
			prefix:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.dataSource, func(t *testing.T) {
			var mu sync.Mutex
			var expands [][]string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != tt.wantPath {
					http.NotFound(w, r)
					return
				}
				mu.Lock()
				expands = append(expands, r.URL.Query()["expand"])
				mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, tt.body)
			}))
			defer server.Close()

			addr := "data." + tt.dataSource + ".test"
			attrs := tt.prefix + "custom_attributes."
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
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(addr, attrs+"%", "5"),
						resource.TestCheckResourceAttr(addr, attrs+"Tier.type", "text"),
						resource.TestCheckResourceAttr(addr, attrs+"Tier.value.str", "gold"),
						resource.TestCheckResourceAttr(addr, attrs+"Owners.type", "user"),
						resource.TestCheckResourceAttr(addr, attrs+"Owners.value.array_of_user.0.id", "user-1"),
						resource.TestCheckResourceAttr(addr, attrs+"Owners.value.array_of_user.0.email", "ada@example.com"),
						// A date-time is also a valid str, and the generated
						// decoder keeps str.
						resource.TestCheckResourceAttr(addr, attrs+"Renewal.value.str", "2027-01-31T00:00:00Z"),
						resource.TestCheckResourceAttr(addr, attrs+"Seats.value.integer", "42"),
						resource.TestCheckNoResourceAttr(addr, attrs+"Unset.value.str"),
					),
				}},
			})

			mu.Lock()
			defer mu.Unlock()
			for _, got := range expands {
				if !slices.Contains(got, "custom_attributes") {
					t.Fatalf("expand = %q on a read, want it to include custom_attributes", got)
				}
			}
		})
	}
}
