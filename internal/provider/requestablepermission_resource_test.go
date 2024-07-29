package provider

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPermissionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig() + fmt.Sprintf(`
				resource "lumos_requestable_permission" "test" {
					app_id      = "%s"
					label       = "Permission Label"
				}`, os.Getenv("APP_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of coffees returned
					resource.TestCheckResourceAttr("lumos_requestable_permission.test", "label", "Permission Label"),
					resource.TestCheckResourceAttr("lumos_requestable_permission.test", "type", "NATIVE"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lumos_requestable_permission.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ErrorCheck: func(err error) error {
			if err == nil {
				return errors.New("expected error but got none")
			}

			// some simple example with string matching
			if strings.Contains(err.Error(), "Invalid request approval config") {
				return nil
			}

			// return original error if no match
			return err
		},
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
				resource "lumos_requestable_permission" "test_manager_approval" {
					app_id      = "%s"
					label       = "Permission Label"
					request_config = {
						appstore_visibility = "HIDDEN"
						request_approval_config = {
						  approvers = {
							users  = []
							groups = []
						  },
						  manager_approval = "NONE"
						}
					  }
				}`, os.Getenv("APP_ID")),
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
				resource "lumos_requestable_permission" "test_approvers" {
					app_id      = "%s"
					label       = "Permission Label"
					request_config = {
						appstore_visibility = "HIDDEN"
						request_approval_config = {
						  approvers = {
							users  = [{ id = "%s" }]
							groups = []
						  },
						}
					  }
				}`, os.Getenv("APP_ID"), os.Getenv("USER_ID")),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ErrorCheck: func(err error) error {
			if err == nil {
				return errors.New("expected error but got none")
			}

			// some simple example with string matching
			if strings.Contains(err.Error(), "Invalid allowed groups") {
				return nil
			}

			// return original error if no match
			return err
		},
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
				resource "lumos_requestable_permission" "test_allowed_groups" {
					app_id      = "%s"
					label       = "Permission Label"
					request_config = {
						allowed_groups = {
							groups = [{ id = "123" }]
						}
					}
				}`, os.Getenv("APP_ID")),
			},
		},
	})
}
