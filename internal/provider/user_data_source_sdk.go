// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
)

func (r *UserDataSourceModel) RefreshFromSharedUser(resp *shared.User) {
	if resp != nil {
		r.Email = types.StringPointerValue(resp.Email)
		r.FamilyName = types.StringPointerValue(resp.FamilyName)
		r.GivenName = types.StringPointerValue(resp.GivenName)
		r.ID = types.StringValue(resp.ID)
		if resp.Status != nil {
			r.Status = types.StringValue(string(*resp.Status))
		} else {
			r.Status = types.StringNull()
		}
	}
}
