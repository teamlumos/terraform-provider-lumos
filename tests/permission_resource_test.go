package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPermissionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig() + fmt.Sprintf(`
				resource "lumos_requestable_permission" "test" {
					app_id      = "%s"
					label       = "Permission Label"
				}`, os.Getenv("APP_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of coffees returned
					resource.TestCheckResourceAttr("lumos_requestable_permission.test", "label", "Permission Label"),
					resource.TestCheckResourceAttr("lumos_requestable_permission.test", "type", "NATIVE"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lumos_requestable_permission.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
