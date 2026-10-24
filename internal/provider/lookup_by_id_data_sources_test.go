package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestLookupByIDDataSourcesRequireID checks that a by-ID lookup without an
// ID fails at plan.
func TestLookupByIDDataSourcesRequireID(t *testing.T) {
	for _, dataSource := range []string{"lumos_app", "lumos_entitlement", "lumos_user", "lumos_group"} {
		t.Run(dataSource, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{{
					Config: fmt.Sprintf(`
provider "lumos" {
  server_url  = "http://127.0.0.1:1"
  http_bearer = "test-token"
}

data %q "test" {}
`, dataSource),
					ExpectError: regexp.MustCompile(`The argument "id" is required`),
				}},
			})
		})
	}
}

// TestLookupByIDDataSources runs each by-ID data source against a fake API
// and checks it reads the object at its ID's path, with the expansions the
// provider asks for, and puts the response in state.
func TestLookupByIDDataSources(t *testing.T) {
	tests := []struct {
		dataSource string
		config     string
		wantPath   string
		wantQuery  string
		body       string
		wantAttrs  map[string]string
	}{
		{
			dataSource: "lumos_app",
			config:     `id = "app-1"`,
			wantPath:   "/v1/apps/app-1",
			wantQuery:  "expand=appstore_settings&expand=custom_attributes",
			body:       `{"id":"app-1","label":"Okta","appstore_settings":{"in_app_store":true}}`,
			wantAttrs:  map[string]string{"id": "app-1", "label": "Okta", "appstore_settings.in_app_store": "true"},
		},
		{
			dataSource: "lumos_entitlement",
			config:     `id = "ent-1"`,
			wantPath:   "/v1/entitlements/ent-1",
			wantQuery:  "expand=appstore_settings",
			body:       `{"id":"ent-1","app_id":"app-1","label":"Admin","appstore_settings":{"visible_in_appstore":true}}`,
			wantAttrs:  map[string]string{"id": "ent-1", "app_id": "app-1", "label": "Admin", "appstore_settings.visible_in_appstore": "true"},
		},
		{
			dataSource: "lumos_user",
			config:     `id = "user-1"`,
			wantPath:   "/v1/users/user-1",
			wantQuery:  "",
			body:       `{"id":"user-1","email":"ada@example.com"}`,
			wantAttrs:  map[string]string{"id": "user-1", "email": "ada@example.com"},
		},
		{
			dataSource: "lumos_user",
			config: `
  id     = "user-1"
  expand = ["custom_attributes"]
`,
			wantPath:  "/v1/users/user-1",
			wantQuery: "expand=custom_attributes",
			body:      `{"id":"user-1","email":"ada@example.com","custom_attributes":{"Team":{"type":"text","value":"Platform"}}}`,
			wantAttrs: map[string]string{"id": "user-1", "custom_attributes.Team.value.str": "Platform"},
		},
		{
			dataSource: "lumos_group",
			config:     `id = "group-1"`,
			wantPath:   "/v1/groups/group-1",
			wantQuery:  "",
			body:       `{"id":"group-1","name":"Engineering"}`,
			wantAttrs:  map[string]string{"id": "group-1", "name": "Engineering"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.dataSource+" "+tt.wantQuery, func(t *testing.T) {
			var mu sync.Mutex
			var requests []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				requests = append(requests, r.Method+" "+r.URL.Path+"?"+r.URL.RawQuery)
				mu.Unlock()
				if r.Method != http.MethodGet || r.URL.Path != tt.wantPath {
					http.NotFound(w, r)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, tt.body)
			}))
			defer server.Close()

			addr := "data." + tt.dataSource + ".test"
			var checks []resource.TestCheckFunc
			for k, v := range tt.wantAttrs {
				checks = append(checks, resource.TestCheckResourceAttr(addr, k, v))
			}
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
					Check: resource.ComposeAggregateTestCheckFunc(checks...),
				}},
			})

			// Terraform reads the data source more than once per step.
			want := "GET " + tt.wantPath + "?" + tt.wantQuery
			mu.Lock()
			defer mu.Unlock()
			for _, got := range requests {
				if got != want {
					t.Fatalf("requests %q, want only %q", requests, want)
				}
			}
		})
	}
}
