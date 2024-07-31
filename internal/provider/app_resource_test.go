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
					name        = "Terraform Testing - Create app"
					description = "Terraform Testing Description"
					category    = "Developers"
					logo_url    = "https://www.example.com/logo.png"
					website_url = "https://www.example.com"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_app.test_app", "user_friendly_label", "Terraform Testing - Create app"),
					resource.TestCheckResourceAttr("lumos_app.test_app", "instance_id", "custom_app_import"),
					resource.TestCheckResourceAttr("lumos_app.test_app", "status", "APPROVED"),
					resource.TestCheckResourceAttr("lumos_app.test_app", "logo_url", "https://www.example.com/logo.png"),
					resource.TestCheckResourceAttr("lumos_app.test_app", "website_url", "https://www.example.com"),
				),
			},
			// Update testing
			{
				Config: providerConfig() + `
				resource "lumos_app" "test_app" {
					name        = "Terraform Testing - Update app"
					description = "New Description"
					category    = "Communication"
					logo_url    = "https://www.example1.com/logo.png"
					website_url = "https://www.example1.com"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_app.test_app", "user_friendly_label", "Terraform Testing - Update app"),
					resource.TestCheckResourceAttr("lumos_app.test_app", "instance_id", "custom_app_import"),
					resource.TestCheckResourceAttr("lumos_app.test_app", "status", "APPROVED"),
					resource.TestCheckResourceAttr("lumos_app.test_app", "logo_url", "https://www.example1.com/logo.png"),
					resource.TestCheckResourceAttr("lumos_app.test_app", "website_url", "https://www.example1.com"),
				),
			},
		},
	})
}
