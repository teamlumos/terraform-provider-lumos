package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPreapprovalRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig() + fmt.Sprintf(`
				resource "lumos_pre_approval_rule" "test" {
					app_id                  = "%s"
					justification           = "Justification"
					preapproved_permissions = [{ 
						id = "%s"
					}]
					preapproved_groups      = [{
						id = "%s"
					}]
					time_based_access       = ["Unlimited"]
				}`, os.Getenv("APP_ID"), os.Getenv("PERMISSION_ID"), os.Getenv("GROUP_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of coffees returned
					resource.TestCheckResourceAttr("lumos_pre_approval_rule.test", "justification", "Justification"),
					resource.TestCheckResourceAttr("lumos_pre_approval_rule.test", "app_id", os.Getenv("APP_ID")),
				),
			},
		},
	})
}
