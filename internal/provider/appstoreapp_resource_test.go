package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func appStoreAppConfig(approverUserID string) string {
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
			approvers = {
				users = [{ id = "%[1]s" }]
			}
			admins = {
				users = [{ id = "%[2]s" }]
			}
		}
	}`, approverUserID, os.Getenv("USER_ID"))
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
				Config: appStoreAppConfig(os.Getenv("USER_ID")),
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
				Config: appStoreAppConfig(secondApproverID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
						plancheck.ExpectUnknownValue(resourceName, approverUsersPath.AtSliceIndex(0).AtMapKey("email")),
						plancheck.ExpectUnknownValue(resourceName, approverUsersPath.AtSliceIndex(0).AtMapKey("given_name")),
						plancheck.ExpectUnknownValue(resourceName, approverUsersPath.AtSliceIndex(0).AtMapKey("family_name")),
					},
				},
			},
			{
				RefreshState: true,
				Check:        checkApproverMatchesUser,
			},
		},
	})
}
