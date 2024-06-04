// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package types

import "github.com/hashicorp/terraform-plugin-framework/types"

type RequestApprovalConfig struct {
	CustomApprovalMessage         types.String                 `tfsdk:"custom_approval_message"`
	CustomApprovalMessageOverride types.Bool                   `tfsdk:"custom_approval_message_override"`
	ManagerApproval               types.String                 `tfsdk:"manager_approval"`
	RequestApprovalConfigOverride types.Bool                   `tfsdk:"request_approval_config_override"`
	RequestApprovalStages         []RequestApprovalStageOutput `tfsdk:"request_approval_stages"`
	RequireAdditionalApproval     types.Bool                   `tfsdk:"require_additional_approval"`
}