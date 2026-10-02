package provider

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
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

func approverChangeConfig(label, stage1UserID, stage1GroupID, stage2UserID, stage2GroupID string) string {
	return providerConfig() + fmt.Sprintf(`
	data "lumos_user" "stage_1" {
		user_id = "%[3]s"
	}
	data "lumos_group" "stage_1" {
		group_id = "%[4]s"
	}
	data "lumos_user" "stage_2" {
		user_id = "%[5]s"
	}
	data "lumos_group" "stage_2" {
		group_id = "%[6]s"
	}
	resource "lumos_requestable_permission" "test_approver_change" {
		app_id = "%[1]s"
		label  = "%[2]s"
		request_config = {
			appstore_visibility = "HIDDEN"
			request_approval_config = {
				request_approval_config_override = true
				manager_approval                 = "NONE"
				require_additional_approval      = true
				approvers = {
					users  = [{ id = "%[3]s" }]
					groups = [{ id = "%[4]s" }]
				}
				approvers_stage_2 = {
					users  = [{ id = "%[5]s" }]
					groups = [{ id = "%[6]s" }]
				}
			}
		}
	}`, os.Getenv("APP_ID"), label, stage1UserID, stage1GroupID, stage2UserID, stage2GroupID)
}

// The read-only user and group attributes under each approvers stage must
// describe the approver the id belongs to. They are only populated in state by
// a refresh: apply copies the plan's approver sets into state, so attributes
// that were unknown in the plan stay null until the next read.
func checkApproverDisplayFieldsMatch(resourceName string) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{}
	for _, stage := range []string{"approvers", "approvers_stage_2"} {
		prefix := "request_config.request_approval_config." + stage
		dataStage := "stage_1"
		if stage == "approvers_stage_2" {
			dataStage = "stage_2"
		}
		for _, attr := range []string{"id", "email", "given_name", "family_name", "status"} {
			checks = append(checks, resource.TestCheckTypeSetElemAttrPair(resourceName, prefix+".users.*."+attr, "data.lumos_user."+dataStage, attr))
		}
		for _, attr := range []string{"id", "name"} {
			checks = append(checks, resource.TestCheckTypeSetElemAttrPair(resourceName, prefix+".groups.*."+attr, "data.lumos_group."+dataStage, attr))
		}
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func expectApproverDisplayFieldsUnknown(resourceName string) []plancheck.PlanCheck {
	checks := []plancheck.PlanCheck{}
	for _, stage := range []string{"approvers", "approvers_stage_2"} {
		stagePath := tfjsonpath.New("request_config").AtMapKey("request_approval_config").AtMapKey(stage)
		for _, attr := range []string{"email", "given_name", "family_name", "status"} {
			checks = append(checks, plancheck.ExpectUnknownValue(resourceName, stagePath.AtMapKey("users").AtSliceIndex(0).AtMapKey(attr)))
		}
		checks = append(checks, plancheck.ExpectUnknownValue(resourceName, stagePath.AtMapKey("groups").AtSliceIndex(0).AtMapKey("name")))
	}
	return checks
}

func TestAccPermissionResourceApproverChange(t *testing.T) {
	resourceName := "lumos_requestable_permission.test_approver_change"
	userA, groupA := os.Getenv("USER_ID"), os.Getenv("GROUP_ID")
	userB, groupB := otherActiveUserID(t, userA), otherGroupID(t, groupA)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: approverChangeConfig("Approver change", userA, groupA, userB, groupB),
			},
			{
				RefreshState: true,
				Check:        checkApproverDisplayFieldsMatch(resourceName),
			},
			// Swapping the approver ids must not carry the previous approver's
			// name and email into the plan.
			{
				Config: approverChangeConfig("Approver change", userB, groupB, userA, groupA),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: append(
						[]plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate)},
						expectApproverDisplayFieldsUnknown(resourceName)...,
					),
				},
			},
			{
				RefreshState: true,
				Check:        checkApproverDisplayFieldsMatch(resourceName),
			},
			// An update that leaves the approvers alone keeps them intact.
			{
				Config: approverChangeConfig("Approver change renamed", userB, groupB, userA, groupA),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate)},
				},
			},
			{
				RefreshState: true,
				Check:        checkApproverDisplayFieldsMatch(resourceName),
			},
		},
	})
}
