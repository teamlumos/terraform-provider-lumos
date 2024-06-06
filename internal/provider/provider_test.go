package provider

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

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
