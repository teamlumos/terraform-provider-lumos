// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import "github.com/hashicorp/terraform-plugin-framework/types"

type RequestConfigInputCreate struct {
	AccessRemovalInlineWebhook     *BaseInlineWebhook             `tfsdk:"access_removal_inline_webhook"`
	AllowedGroups                  *AllowedGroupsConfigInput      `tfsdk:"allowed_groups"`
	AllowedGroupsOverride          types.Bool                     `tfsdk:"allowed_groups_override"`
	AppstoreVisibility             types.String                   `tfsdk:"appstore_visibility"`
	RequestApprovalConfig          *RequestApprovalConfigInput    `tfsdk:"request_approval_config"`
	RequestFulfillmentConfig       *RequestFulfillmentConfigInput `tfsdk:"request_fulfillment_config"`
	RequestValidationInlineWebhook *BaseInlineWebhook             `tfsdk:"request_validation_inline_webhook"`
}
