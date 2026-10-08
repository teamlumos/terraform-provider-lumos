package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func appStoreAppConfig(approverUserID, approvalMessage string) string {
	return providerConfig() + fmt.Sprintf(`
	data "lumos_user" "approver" {
		user_id = "%[1]s"
	}
	resource "lumos_app" "test_app" {
		name                 = "Terraform Testing - Appstore app"
		description          = "Terraform Testing Description"
		category             = "Developers"
		website_url          = "https://www.example.com"
		logo_url             = "https://www.example.com/logo.png"
		request_instructions = "Custom request instructions"
	}
	resource "lumos_app_store_app" "test_app_store_app" {
		app_id        = lumos_app.test_app.id
		request_flow = {
			custom_approval_message = "%[3]s"
			approvers = {
				users = [{ id = "%[1]s" }]
			}
			admins = {
				users = [{ id = "%[2]s" }]
			}
		}
	}`, approverUserID, os.Getenv("USER_ID"), approvalMessage)
}

// userDisplayFieldChecks applies check to the read-only display fields of the
// first user in the users set at usersPath.
func userDisplayFieldChecks(resourceName string, usersPath tfjsonpath.Path, check func(string, tfjsonpath.Path) plancheck.PlanCheck) []plancheck.PlanCheck {
	checks := []plancheck.PlanCheck{}
	for _, attr := range []string{"email", "given_name", "family_name", "status"} {
		checks = append(checks, check(resourceName, usersPath.AtSliceIndex(0).AtMapKey(attr)))
	}
	return checks
}

func TestAccAppStoreAppResource(t *testing.T) {
	resourceName := "lumos_app_store_app.test_app_store_app"
	secondApproverID := otherActiveUserID(t, os.Getenv("USER_ID"))
	approverUsersPath := tfjsonpath.New("request_flow").AtMapKey("approvers").AtMapKey("users")
	checkApproverMatchesUser := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckTypeSetElemAttrPair(resourceName, "request_flow.approvers.users.*.id", "data.lumos_user.approver", "id"),
		resource.TestCheckTypeSetElemAttrPair(resourceName, "request_flow.approvers.users.*.email", "data.lumos_user.approver", "email"),
		resource.TestCheckTypeSetElemAttrPair(resourceName, "request_flow.approvers.users.*.given_name", "data.lumos_user.approver", "given_name"),
		resource.TestCheckTypeSetElemAttrPair(resourceName, "request_flow.approvers.users.*.family_name", "data.lumos_user.approver", "family_name"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: appStoreAppConfig(os.Getenv("USER_ID"), "Approved"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_friendly_label", "Terraform Testing - Appstore app"),
					resource.TestCheckResourceAttr(resourceName, "instance_id", "custom_app_import"),
					resource.TestCheckResourceAttr(resourceName, "status", "APPROVED"),
				),
			},
			{
				RefreshState: true,
				Check:        checkApproverMatchesUser,
			},
			// Swapping the approver id must not carry the previous approver's
			// name and email into the plan.
			{
				Config: appStoreAppConfig(secondApproverID, "Approved"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: append(
						[]plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate)},
						userDisplayFieldChecks(resourceName, approverUsersPath, plancheck.ExpectUnknownValue)...,
					),
				},
			},
			{
				RefreshState: true,
				Check:        checkApproverMatchesUser,
			},
			// An update that leaves the approvers alone does not plan their
			// display fields as "(known after apply)".
			{
				Config: appStoreAppConfig(secondApproverID, "Approved, updated"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: append(
						[]plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate)},
						userDisplayFieldChecks(resourceName, approverUsersPath, expectNotUnknownValue)...,
					),
				},
			},
			{
				RefreshState: true,
				Check:        checkApproverMatchesUser,
			},
		},
	})
}
