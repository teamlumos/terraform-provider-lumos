package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRequestablePermissionDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig() + fmt.Sprintf(`
				data "lumos_requestable_permission" "test_permission" {
					id = "%s"
				}`, os.Getenv("PERMISSION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lumos_requestable_permission.test_permission", "id", os.Getenv("PERMISSION_ID")),
				),
			},
		},
	})
}
