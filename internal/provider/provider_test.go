package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
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

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"lumos": providerserver.NewProtocol6WithError(New("test")()),
	}
)
