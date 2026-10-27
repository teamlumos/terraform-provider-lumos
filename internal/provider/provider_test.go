package provider

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"sort"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// testAccPreCheck fails an acceptance test up front when the test domain's
// credentials are missing, instead of letting the provider fail mid-run.
func testAccPreCheck(t *testing.T) {
	t.Helper()
	for _, name := range []string{"LUMOS_API_TOKEN", "LUMOS_API_URL"} {
		if os.Getenv(name) == "" {
			t.Fatalf("%s must be set for acceptance tests", name)
		}
	}
}

func providerConfig() string {
	return fmt.Sprintf(`
provider "lumos" {
  http_bearer    = "%s"
  server_url     = "%s"
}
`, os.Getenv("LUMOS_API_TOKEN"), os.Getenv("LUMOS_API_URL"))
}

// expectOnlyAttributesChange fails a plan where the resource at addr changes
// any top-level attribute other than attrs, counting an attribute that
// becomes known after apply as a change. It is how a test checks that a
// one-field edit plans only that field.
func expectOnlyAttributesChange(addr string, attrs ...string) plancheck.PlanCheck {
	return onlyAttributesChange{addr: addr, attrs: attrs}
}

type onlyAttributesChange struct {
	addr  string
	attrs []string
}

func (c onlyAttributesChange) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	allowed := map[string]bool{}
	for _, a := range c.attrs {
		allowed[a] = true
	}
	for _, rc := range req.Plan.ResourceChanges {
		if rc.Address != c.addr {
			continue
		}
		before, _ := rc.Change.Before.(map[string]any)
		after, _ := rc.Change.After.(map[string]any)
		unknown, _ := rc.Change.AfterUnknown.(map[string]any)

		var unexpected []string
		for name := range mergedKeys(before, after, unknown) {
			if allowed[name] {
				continue
			}
			if !reflect.DeepEqual(before[name], after[name]) || containsTrue(unknown[name]) {
				unexpected = append(unexpected, name)
			}
		}
		if len(unexpected) > 0 {
			sort.Strings(unexpected)
			resp.Error = fmt.Errorf("%s: plan changes %v, want only %v", c.addr, unexpected, c.attrs)
		}
		return
	}
	resp.Error = fmt.Errorf("%s: not in the plan", c.addr)
}

func mergedKeys(maps ...map[string]any) map[string]bool {
	keys := map[string]bool{}
	for _, m := range maps {
		for k := range m {
			keys[k] = true
		}
	}
	return keys
}

// containsTrue reports whether an after_unknown value marks anything
// unknown; nested attributes report unknown parts as nested true values.
func containsTrue(v any) bool {
	switch v := v.(type) {
	case bool:
		return v
	case map[string]any:
		for _, x := range v {
			if containsTrue(x) {
				return true
			}
		}
	case []any:
		for _, x := range v {
			if containsTrue(x) {
				return true
			}
		}
	}
	return false
}

// testAccOptionalEnv returns an optional test-domain setting, or skips the
// test when it isn't set. These name things a test can't create through
// the provider, like an app with admins or a domain's time-based access
// labels.
func testAccOptionalEnv(t *testing.T, name string) string {
	t.Helper()
	value := os.Getenv(name)
	if value == "" {
		t.Skipf("%s isn't set", name)
	}
	return value
}

// testAccAppStoreAppMu runs the tests that list the app in
// LUMOS_TEST_APPSTORE_APP_ID one at a time: each leaves the app unlisted,
// which would break another test still relying on the listing.
var testAccAppStoreAppMu sync.Mutex

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"lumos": providerserver.NewProtocol6WithError(New("test")()),
	}
)
