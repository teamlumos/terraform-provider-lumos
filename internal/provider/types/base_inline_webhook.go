// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BaseInlineWebhook struct {
	Description types.String `tfsdk:"description"`
	HookType    types.String `tfsdk:"hook_type"`
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
}
