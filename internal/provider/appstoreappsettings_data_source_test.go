package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAppStoreAppSettingsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
				data "lumos_app_store_app_settings" "test_app_settings" {
					id        = "%s"
				}`, os.Getenv("APP_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lumos_app_store_app_settings.test_app_settings", "custom_request_instructions", "Please only request apps that are relevant to your role."),
				),
			},
		},
	})
}
