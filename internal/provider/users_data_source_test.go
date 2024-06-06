package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUsersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
					data "lumos_users" "test_users" {
						search_term = "%s"
					}`, os.Getenv("USER_EMAIL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lumos_users.test_users", "total", "1"),
					resource.TestCheckResourceAttr("data.lumos_users.test_users", "page", "1"),
					resource.TestCheckResourceAttr("data.lumos_users.test_users", "items.#", "1"),
					resource.TestCheckResourceAttr("data.lumos_users.test_users", "items.0.id", os.Getenv("USER_ID")),
				),
			},
		},
	})
}
