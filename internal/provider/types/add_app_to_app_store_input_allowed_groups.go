// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package types

import "github.com/hashicorp/terraform-plugin-framework/types"

type AddAppToAppStoreInputAllowedGroups struct {
	Groups []Group      `tfsdk:"groups"`
	Type   types.String `tfsdk:"type"`
}
