// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package types

import "github.com/hashicorp/terraform-plugin-framework/types"

type AddAppToAppStoreInputProvisioning struct {
	AccessRemovalInlineWebhook     *AddAppToAppStoreInputAccessRemovalInlineWebhook `tfsdk:"access_removal_inline_webhook"`
	CustomProvisioningInstructions types.String                                     `tfsdk:"custom_provisioning_instructions"`
	GroupsProvisioning             types.String                                     `tfsdk:"groups_provisioning"`
	ManualStepsNeeded              types.Bool                                       `tfsdk:"manual_steps_needed"`
	ProvisioningWebhook            *AddAppToAppStoreInputAccessRemovalInlineWebhook `tfsdk:"provisioning_webhook"`
	TimeBasedAccess                []types.String                                   `tfsdk:"time_based_access"`
}
