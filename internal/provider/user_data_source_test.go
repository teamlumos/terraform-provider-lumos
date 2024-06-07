package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
				data "lumos_user" "test_user" {
					user_id = "%s"
				}`, os.Getenv("USER_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lumos_user.test_user", "email", os.Getenv("USER_EMAIL")),
					resource.TestCheckResourceAttr("data.lumos_user.test_user", "status", "ACTIVE"),
				),
			},
		},
	})
}
