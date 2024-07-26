// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package types

import "github.com/hashicorp/terraform-plugin-framework/types"

type RequestablePermissionOutput struct {
	AppClassID    types.String                            `tfsdk:"app_class_id"`
	AppID         types.String                            `tfsdk:"app_id"`
	AppInstanceID types.String                            `tfsdk:"app_instance_id"`
	ID            types.String                            `tfsdk:"id"`
	Label         types.String                            `tfsdk:"label"`
	RequestConfig RequestablePermissionInputRequestConfig `tfsdk:"request_config"`
	Type          types.String                            `tfsdk:"type"`
}