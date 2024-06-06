package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestGroupsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
					data "lumos_groups" "test_groups" {
						name = "%s"
					}`, os.Getenv("GROUP_NAME")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lumos_groups.test_groups", "total", "1"),
					resource.TestCheckResourceAttr("data.lumos_groups.test_groups", "page", "1"),
					resource.TestCheckResourceAttr("data.lumos_groups.test_groups", "items.#", "1"),
					resource.TestCheckResourceAttr("data.lumos_groups.test_groups", "items.0.id", os.Getenv("GROUP_ID")),
				),
			},
		},
	})
}
