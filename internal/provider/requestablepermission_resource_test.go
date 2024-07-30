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
				resource "lumos_requestable_permission" "test_permission_with_overrides" {
					app_id      = "%s"
					label       = "Permission with Config Overrides"
					request_config = {
						allowed_groups_override = true
						allowed_groups = {
							groups = [
								{ 
									id = "%s"
								},
							]
							type = "SPECIFIED_GROUPS"
						}
						appstore_visibility = "HIDDEN"
						request_approval_config = {
							request_approval_config_override = true
							approvers = {
								users  = [{ id = "%s" }]
								groups = []
							}
							manager_approval = "NONE"
						}
					}
				}`, os.Getenv("APP_ID"), os.Getenv("GROUP_ID"), os.Getenv("USER_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_requestable_permission.test_permission_with_overrides", "label", "Permission with Config Overrides"),
					resource.TestCheckResourceAttr("lumos_requestable_permission.test_permission_with_overrides", "type", "NATIVE"),
				),
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
