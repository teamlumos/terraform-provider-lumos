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
			{
				Config: providerConfig() + fmt.Sprintf(`
				resource "lumos_requestable_permission" "test_permission_with_overrides_1" {
					app_id      = "%s"
					label       = "Permission with Config Overrides 1"
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
					resource.TestCheckResourceAttr("lumos_requestable_permission.test_permission_with_overrides_1", "label", "Permission with Config Overrides 1"),
					resource.TestCheckResourceAttr("lumos_requestable_permission.test_permission_with_overrides_1", "type", "NATIVE"),
				),
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
				resource "lumos_requestable_permission" "test_permission_with_overrides_2" {
					app_id      = "%s"
					label       = "Permission with Config Overrides 2"
					request_config = {
						allowed_groups_override = true
						allowed_groups = {
							type = "ALL_GROUPS"
						}
						appstore_visibility = "HIDDEN"
						request_approval_config = {
							request_approval_config_override = true
							approvers = {
								users  = [{ id = "%s" }]
								groups = []
							}
						}
					}
				}`, os.Getenv("APP_ID"), os.Getenv("USER_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_requestable_permission.test_permission_with_overrides_2", "label", "Permission with Config Overrides 2"),
					resource.TestCheckResourceAttr("lumos_requestable_permission.test_permission_with_overrides_2", "type", "NATIVE"),
				),
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
				resource "lumos_requestable_permission" "test_permission_with_overrides_3" {
					app_id      = "%s"
					label       = "Permission with Config Overrides 3"
					request_config = {
						allowed_groups_override = true
						allowed_groups = {
							groups = []
							type = "ALL_GROUPS"
						}
						appstore_visibility = "HIDDEN"
						request_approval_config = {
							request_approval_config_override = true
							manager_approval = "NONE"
						}
					}
				}`, os.Getenv("APP_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_requestable_permission.test_permission_with_overrides_3", "label", "Permission with Config Overrides 3"),
					resource.TestCheckResourceAttr("lumos_requestable_permission.test_permission_with_overrides_3", "type", "NATIVE"),
				),
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
				resource "lumos_requestable_permission" "test_permission_with_no_overrides_1" {
					app_id      = "%s"
					label       = "Permission with No Config Overrides 1"
					request_config = {
						appstore_visibility = "HIDDEN"
					}
				}`, os.Getenv("APP_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lumos_requestable_permission.test_permission_with_no_overrides_1", "label", "Permission with No Config Overrides 1"),
					resource.TestCheckResourceAttr("lumos_requestable_permission.test_permission_with_no_overrides_1", "type", "NATIVE"),
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
					label       = "manager approval"
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
					label       = "Approvers"
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
			if strings.Contains(err.Error(), "Invalid time based access") {
				return nil
			}

			// return original error if no match
			return err
		},
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
				resource "lumos_requestable_permission" "test_lowercase_time_based_access" {
					app_id      = "%s"
					label       = "Lowercase time based access"
					request_config = {
						request_fulfillment_config = {
							time_based_access_override = true
							time_based_access = ["unlimited"]
						}
					}
				}`, os.Getenv("APP_ID")),
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
			if strings.Contains(err.Error(), "Invalid time based access") {
				return nil
			}

			// return original error if no match
			return err
		},
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
				resource "lumos_requestable_permission" "test_uppercase_time_based_access" {
					app_id      = "%s"
					label       = "Upper case time based access"
					request_config = {
						request_fulfillment_config = {
							time_based_access_override = true
							time_based_access = ["2 Hours"]
						}
					}
				}`, os.Getenv("APP_ID")),
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
			if strings.Contains(err.Error(), "Invalid time based access") {
				return nil
			}

			// return original error if no match
			return err
		},
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
				resource "lumos_requestable_permission" "test_override_but_no_time_based_access" {
					app_id      = "%s"
					label       = "Override but no time based access"
					request_config = {
						request_fulfillment_config = {
							time_based_access_override = true
							time_based_access = []
						}
					}
				}`, os.Getenv("APP_ID")),
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
			if strings.Contains(err.Error(), "Invalid time based access") {
				return nil
			}

			// return original error if no match
			return err
		},
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
				resource "lumos_requestable_permission" "test_no_override_but_time_based_access" {
					app_id      = "%s"
					label       = "No override but time based access"
					request_config = {
						request_fulfillment_config = {
							time_based_access_override = false
							time_based_access = ["Unlimited"]
						}
					}
				}`, os.Getenv("APP_ID")),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ErrorCheck: func(err error) error {
			if err == nil {
				return nil
			}

			return err
		},
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
				resource "lumos_requestable_permission" "test_override_with_time_based_access" {
					app_id      = "%s"
					label       = "Override with time based access"
					request_config = {
						request_fulfillment_config = {
							time_based_access_override = true
							time_based_access = ["Unlimited"]
						}
					}
				}`, os.Getenv("APP_ID")),
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
				resource "lumos_requestable_permission" "test_allowed_groups_override" {
					app_id      = "%s"
					label       = "Invalid allowed groups permission"
					request_config = {
						allowed_groups = {
							groups = [{ id = "%s" }]
						}
					}
				}`, os.Getenv("APP_ID"), os.Getenv("GROUP_ID")),
			},
		},
	})
}