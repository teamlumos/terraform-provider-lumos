package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAppStoreAppDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig() + fmt.Sprintf(`
				data "lumos_app_store_app" "test" {
					app_id = "%s"
				}`, os.Getenv("APP_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lumos_app_store_app.test", "user_friendly_label", os.Getenv("APP_NAME")),
				),
			},
		},
	})
}
