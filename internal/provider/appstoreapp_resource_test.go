package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAppStoreAppResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig() + fmt.Sprintf(`
				resource "lumos_app" "test_app" {
					name        = "Terraform Testing - Appstore app"
					description = "Terraform Testing Description"
					category    = "Developers"
				}
				resource "lumos_app_store_app" "test_app_store_app" {
					app_id        = lumos_app.test_app.id
					request_flow = {
						approvers = {
							users = [{ id = "%s" }]
						}
						admins = {
							users = [{ id = "%s" }]
						}
					}
				}`, os.Getenv("USER_ID"), os.Getenv("USER_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_app_store_app.test_app_store_app", "user_friendly_label", "Terraform Testing - Appstore app"),
					resource.TestCheckResourceAttr("lumos_app_store_app.test_app_store_app", "instance_id", "custom_app_import"),
					resource.TestCheckResourceAttr("lumos_app_store_app.test_app_store_app", "status", "APPROVED"),
				),
			},
		},
	})
}
