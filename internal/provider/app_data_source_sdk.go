// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
)

func (r *AppDataSourceModel) RefreshFromSharedApp(resp *shared.App) {
	if resp != nil {
		r.AllowMultiplePermissionSelection = types.BoolValue(resp.AllowMultiplePermissionSelection)
		r.AppClassID = types.StringValue(resp.AppClassID)
		r.ID = types.StringValue(resp.ID)
		r.InstanceID = types.StringValue(resp.InstanceID)
		r.LogoURL = types.StringPointerValue(resp.LogoURL)
		r.RequestInstructions = types.StringPointerValue(resp.RequestInstructions)
		r.Sources = []types.String{}
		for _, v := range resp.Sources {
			r.Sources = append(r.Sources, types.StringValue(string(v)))
		}
		r.Status = types.StringValue(string(resp.Status))
		r.UserFriendlyLabel = types.StringValue(resp.UserFriendlyLabel)
		r.WebsiteURL = types.StringPointerValue(resp.WebsiteURL)
	}
}
