// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	tfTypes "github.com/teamlumos/terraform-provider-lumos/internal/provider/types"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
)

func (r *PreApprovalRuleResourceModel) ToSharedPreApprovalRuleInput() *shared.PreApprovalRuleInput {
	appID := r.AppID.ValueString()
	justification := r.Justification.ValueString()
	var preapprovalWebhooks []shared.BaseInlineWebhook = []shared.BaseInlineWebhook{}
	for _, preapprovalWebhooksItem := range r.PreapprovalWebhooks {
		id := preapprovalWebhooksItem.ID.ValueString()
		preapprovalWebhooks = append(preapprovalWebhooks, shared.BaseInlineWebhook{
			ID: id,
		})
	}
	var preapprovedGroups []shared.BaseGroup = []shared.BaseGroup{}
	for _, preapprovedGroupsItem := range r.PreapprovedGroups {
		appId1 := new(string)
		if !preapprovedGroupsItem.AppID.IsUnknown() && !preapprovedGroupsItem.AppID.IsNull() {
			*appId1 = preapprovedGroupsItem.AppID.ValueString()
		} else {
			appId1 = nil
		}
		id1 := new(string)
		if !preapprovedGroupsItem.ID.IsUnknown() && !preapprovedGroupsItem.ID.IsNull() {
			*id1 = preapprovedGroupsItem.ID.ValueString()
		} else {
			id1 = nil
		}
		integrationSpecificID := new(string)
		if !preapprovedGroupsItem.IntegrationSpecificID.IsUnknown() && !preapprovedGroupsItem.IntegrationSpecificID.IsNull() {
			*integrationSpecificID = preapprovedGroupsItem.IntegrationSpecificID.ValueString()
		} else {
			integrationSpecificID = nil
		}
		preapprovedGroups = append(preapprovedGroups, shared.BaseGroup{
			AppID:                 appId1,
			ID:                    id1,
			IntegrationSpecificID: integrationSpecificID,
		})
	}
	var preapprovedPermissions []shared.RequestablePermissionBase = []shared.RequestablePermissionBase{}
	for _, preapprovedPermissionsItem := range r.PreapprovedPermissions {
		id2 := new(string)
		if !preapprovedPermissionsItem.ID.IsUnknown() && !preapprovedPermissionsItem.ID.IsNull() {
			*id2 = preapprovedPermissionsItem.ID.ValueString()
		} else {
			id2 = nil
		}
		preapprovedPermissions = append(preapprovedPermissions, shared.RequestablePermissionBase{
			ID: id2,
		})
	}
	var timeBasedAccess []shared.TimeBasedAccessOptions = []shared.TimeBasedAccessOptions{}
	for _, timeBasedAccessItem := range r.TimeBasedAccess {
		timeBasedAccess = append(timeBasedAccess, shared.TimeBasedAccessOptions(timeBasedAccessItem.ValueString()))
	}
	out := shared.PreApprovalRuleInput{
		AppID:                  appID,
		Justification:          justification,
		PreapprovalWebhooks:    preapprovalWebhooks,
		PreapprovedGroups:      preapprovedGroups,
		PreapprovedPermissions: preapprovedPermissions,
		TimeBasedAccess:        timeBasedAccess,
	}
	return &out
}

func (r *PreApprovalRuleResourceModel) RefreshFromSharedPreApprovalRuleOutput(resp *shared.PreApprovalRuleOutput) {
	if resp != nil {
		r.AppClassID = types.StringValue(resp.AppClassID)
		r.AppID = types.StringValue(resp.AppID)
		r.AppInstanceID = types.StringValue(resp.AppInstanceID)
		r.ID = types.StringPointerValue(resp.ID)
		r.Justification = types.StringValue(resp.Justification)
		r.PreapprovalWebhooks = []tfTypes.AddAppToAppStoreInputAccessRemovalInlineWebhook{}
		if len(r.PreapprovalWebhooks) > len(resp.PreapprovalWebhooks) {
			r.PreapprovalWebhooks = r.PreapprovalWebhooks[:len(resp.PreapprovalWebhooks)]
		}
		for preapprovalWebhooksCount, preapprovalWebhooksItem := range resp.PreapprovalWebhooks {
			var preapprovalWebhooks1 tfTypes.AddAppToAppStoreInputAccessRemovalInlineWebhook
			preapprovalWebhooks1.Description = types.StringPointerValue(preapprovalWebhooksItem.Description)
			preapprovalWebhooks1.HookType = types.StringValue(string(preapprovalWebhooksItem.HookType))
			preapprovalWebhooks1.ID = types.StringValue(preapprovalWebhooksItem.ID)
			preapprovalWebhooks1.Name = types.StringValue(preapprovalWebhooksItem.Name)
			if preapprovalWebhooksCount+1 > len(r.PreapprovalWebhooks) {
				r.PreapprovalWebhooks = append(r.PreapprovalWebhooks, preapprovalWebhooks1)
			} else {
				r.PreapprovalWebhooks[preapprovalWebhooksCount].Description = preapprovalWebhooks1.Description
				r.PreapprovalWebhooks[preapprovalWebhooksCount].HookType = preapprovalWebhooks1.HookType
				r.PreapprovalWebhooks[preapprovalWebhooksCount].ID = preapprovalWebhooks1.ID
				r.PreapprovalWebhooks[preapprovalWebhooksCount].Name = preapprovalWebhooks1.Name
			}
		}
		r.PreapprovedGroups = []tfTypes.Group{}
		if len(r.PreapprovedGroups) > len(resp.PreapprovedGroups) {
			r.PreapprovedGroups = r.PreapprovedGroups[:len(resp.PreapprovedGroups)]
		}
		for preapprovedGroupsCount, preapprovedGroupsItem := range resp.PreapprovedGroups {
			var preapprovedGroups1 tfTypes.Group
			preapprovedGroups1.AppID = types.StringPointerValue(preapprovedGroupsItem.AppID)
			preapprovedGroups1.Description = types.StringPointerValue(preapprovedGroupsItem.Description)
			if preapprovedGroupsItem.GroupLifecycle != nil {
				preapprovedGroups1.GroupLifecycle = types.StringValue(string(*preapprovedGroupsItem.GroupLifecycle))
			} else {
				preapprovedGroups1.GroupLifecycle = types.StringNull()
			}
			preapprovedGroups1.ID = types.StringPointerValue(preapprovedGroupsItem.ID)
			preapprovedGroups1.IntegrationSpecificID = types.StringPointerValue(preapprovedGroupsItem.IntegrationSpecificID)
			preapprovedGroups1.Name = types.StringPointerValue(preapprovedGroupsItem.Name)
			preapprovedGroups1.SourceAppID = types.StringPointerValue(preapprovedGroupsItem.SourceAppID)
			if preapprovedGroupsCount+1 > len(r.PreapprovedGroups) {
				r.PreapprovedGroups = append(r.PreapprovedGroups, preapprovedGroups1)
			} else {
				r.PreapprovedGroups[preapprovedGroupsCount].AppID = preapprovedGroups1.AppID
				r.PreapprovedGroups[preapprovedGroupsCount].Description = preapprovedGroups1.Description
				r.PreapprovedGroups[preapprovedGroupsCount].GroupLifecycle = preapprovedGroups1.GroupLifecycle
				r.PreapprovedGroups[preapprovedGroupsCount].ID = preapprovedGroups1.ID
				r.PreapprovedGroups[preapprovedGroupsCount].IntegrationSpecificID = preapprovedGroups1.IntegrationSpecificID
				r.PreapprovedGroups[preapprovedGroupsCount].Name = preapprovedGroups1.Name
				r.PreapprovedGroups[preapprovedGroupsCount].SourceAppID = preapprovedGroups1.SourceAppID
			}
		}
		r.PreapprovedPermissions = []tfTypes.RequestablePermissionBase{}
		if len(r.PreapprovedPermissions) > len(resp.PreapprovedPermissions) {
			r.PreapprovedPermissions = r.PreapprovedPermissions[:len(resp.PreapprovedPermissions)]
		}
		for preapprovedPermissionsCount, preapprovedPermissionsItem := range resp.PreapprovedPermissions {
			var preapprovedPermissions1 tfTypes.RequestablePermissionBase
			preapprovedPermissions1.AppClassID = types.StringValue(preapprovedPermissionsItem.AppClassID)
			preapprovedPermissions1.AppID = types.StringValue(preapprovedPermissionsItem.AppID)
			preapprovedPermissions1.AppInstanceID = types.StringValue(preapprovedPermissionsItem.AppInstanceID)
			preapprovedPermissions1.ID = types.StringPointerValue(preapprovedPermissionsItem.ID)
			preapprovedPermissions1.Label = types.StringValue(preapprovedPermissionsItem.Label)
			if preapprovedPermissionsItem.Type != nil {
				preapprovedPermissions1.Type = types.StringValue(string(*preapprovedPermissionsItem.Type))
			} else {
				preapprovedPermissions1.Type = types.StringNull()
			}
			if preapprovedPermissionsCount+1 > len(r.PreapprovedPermissions) {
				r.PreapprovedPermissions = append(r.PreapprovedPermissions, preapprovedPermissions1)
			} else {
				r.PreapprovedPermissions[preapprovedPermissionsCount].AppClassID = preapprovedPermissions1.AppClassID
				r.PreapprovedPermissions[preapprovedPermissionsCount].AppID = preapprovedPermissions1.AppID
				r.PreapprovedPermissions[preapprovedPermissionsCount].AppInstanceID = preapprovedPermissions1.AppInstanceID
				r.PreapprovedPermissions[preapprovedPermissionsCount].ID = preapprovedPermissions1.ID
				r.PreapprovedPermissions[preapprovedPermissionsCount].Label = preapprovedPermissions1.Label
				r.PreapprovedPermissions[preapprovedPermissionsCount].Type = preapprovedPermissions1.Type
			}
		}
		r.TimeBasedAccess = []types.String{}
		for _, v := range resp.TimeBasedAccess {
			r.TimeBasedAccess = append(r.TimeBasedAccess, types.StringValue(string(v)))
		}
	}
}

func (r *PreApprovalRuleResourceModel) ToSharedPreApprovalRuleUpdateInput() *shared.PreApprovalRuleUpdateInput {
	justification := r.Justification.ValueString()
	var preapprovalWebhooks []shared.BaseInlineWebhook = []shared.BaseInlineWebhook{}
	for _, preapprovalWebhooksItem := range r.PreapprovalWebhooks {
		id := preapprovalWebhooksItem.ID.ValueString()
		preapprovalWebhooks = append(preapprovalWebhooks, shared.BaseInlineWebhook{
			ID: id,
		})
	}
	var preapprovedGroups []shared.BaseGroup = []shared.BaseGroup{}
	for _, preapprovedGroupsItem := range r.PreapprovedGroups {
		appID := new(string)
		if !preapprovedGroupsItem.AppID.IsUnknown() && !preapprovedGroupsItem.AppID.IsNull() {
			*appID = preapprovedGroupsItem.AppID.ValueString()
		} else {
			appID = nil
		}
		id1 := new(string)
		if !preapprovedGroupsItem.ID.IsUnknown() && !preapprovedGroupsItem.ID.IsNull() {
			*id1 = preapprovedGroupsItem.ID.ValueString()
		} else {
			id1 = nil
		}
		integrationSpecificID := new(string)
		if !preapprovedGroupsItem.IntegrationSpecificID.IsUnknown() && !preapprovedGroupsItem.IntegrationSpecificID.IsNull() {
			*integrationSpecificID = preapprovedGroupsItem.IntegrationSpecificID.ValueString()
		} else {
			integrationSpecificID = nil
		}
		preapprovedGroups = append(preapprovedGroups, shared.BaseGroup{
			AppID:                 appID,
			ID:                    id1,
			IntegrationSpecificID: integrationSpecificID,
		})
	}
	var preapprovedPermissions []shared.RequestablePermissionBase = []shared.RequestablePermissionBase{}
	for _, preapprovedPermissionsItem := range r.PreapprovedPermissions {
		id2 := new(string)
		if !preapprovedPermissionsItem.ID.IsUnknown() && !preapprovedPermissionsItem.ID.IsNull() {
			*id2 = preapprovedPermissionsItem.ID.ValueString()
		} else {
			id2 = nil
		}
		preapprovedPermissions = append(preapprovedPermissions, shared.RequestablePermissionBase{
			ID: id2,
		})
	}
	var timeBasedAccess []shared.TimeBasedAccessOptions = []shared.TimeBasedAccessOptions{}
	for _, timeBasedAccessItem := range r.TimeBasedAccess {
		timeBasedAccess = append(timeBasedAccess, shared.TimeBasedAccessOptions(timeBasedAccessItem.ValueString()))
	}
	out := shared.PreApprovalRuleUpdateInput{
		Justification:          justification,
		PreapprovalWebhooks:    preapprovalWebhooks,
		PreapprovedGroups:      preapprovedGroups,
		PreapprovedPermissions: preapprovedPermissions,
		TimeBasedAccess:        timeBasedAccess,
	}
	return &out
}
