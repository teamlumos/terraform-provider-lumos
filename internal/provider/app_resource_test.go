package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAppResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig() + `
				resource "lumos_app" "test_app" {
					name        = "Terraform Testing"
					description = "Terraform Testing Description"
					category    = "Developers"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_app.test_app", "user_friendly_label", "Terraform Testing"),
					resource.TestCheckResourceAttr("lumos_app.test_app", "instance_id", "custom_app_import"),
					resource.TestCheckResourceAttr("lumos_app.test_app", "status", "APPROVED"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "lumos_app.test_app",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"category", "name", "description"},
			},
		},
	})
}
