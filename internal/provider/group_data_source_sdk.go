// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
)

func (r *GroupDataSourceModel) RefreshFromSharedGroup(resp *shared.Group) {
	if resp != nil {
		r.AppID = types.StringPointerValue(resp.AppID)
		r.Description = types.StringPointerValue(resp.Description)
		if resp.GroupLifecycle != nil {
			r.GroupLifecycle = types.StringValue(string(*resp.GroupLifecycle))
		} else {
			r.GroupLifecycle = types.StringNull()
		}
		r.ID = types.StringPointerValue(resp.ID)
		r.IntegrationSpecificID = types.StringPointerValue(resp.IntegrationSpecificID)
		r.Name = types.StringPointerValue(resp.Name)
		r.SourceAppID = types.StringPointerValue(resp.SourceAppID)
	}
}