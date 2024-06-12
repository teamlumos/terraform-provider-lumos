// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	tfTypes "github.com/teamlumos/terraform-provider-lumos/internal/provider/types"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
)

func (r *RequestablePermissionDataSourceModel) RefreshFromSharedRequestablePermissionOutput(resp *shared.RequestablePermissionOutput) {
	if resp != nil {
		r.AppClassID = types.StringValue(resp.AppClassID)
		r.AppID = types.StringValue(resp.AppID)
		r.AppInstanceID = types.StringValue(resp.AppInstanceID)
		r.ID = types.StringPointerValue(resp.ID)
		r.Label = types.StringValue(resp.Label)
		if resp.RequestConfig.AccessRemovalInlineWebhook == nil {
			r.RequestConfig.AccessRemovalInlineWebhook = nil
		} else {
			r.RequestConfig.AccessRemovalInlineWebhook = &tfTypes.AddAppToAppStoreInputAccessRemovalInlineWebhook{}
			r.RequestConfig.AccessRemovalInlineWebhook.Description = types.StringPointerValue(resp.RequestConfig.AccessRemovalInlineWebhook.Description)
			r.RequestConfig.AccessRemovalInlineWebhook.HookType = types.StringValue(string(resp.RequestConfig.AccessRemovalInlineWebhook.HookType))
			r.RequestConfig.AccessRemovalInlineWebhook.ID = types.StringValue(resp.RequestConfig.AccessRemovalInlineWebhook.ID)
			r.RequestConfig.AccessRemovalInlineWebhook.Name = types.StringValue(resp.RequestConfig.AccessRemovalInlineWebhook.Name)
		}
		if resp.RequestConfig.AllowedGroups == nil {
			r.RequestConfig.AllowedGroups = nil
		} else {
			r.RequestConfig.AllowedGroups = &tfTypes.AddAppToAppStoreInputAllowedGroups{}
			r.RequestConfig.AllowedGroups.Groups = []tfTypes.Group{}
			if len(r.RequestConfig.AllowedGroups.Groups) > len(resp.RequestConfig.AllowedGroups.Groups) {
				r.RequestConfig.AllowedGroups.Groups = r.RequestConfig.AllowedGroups.Groups[:len(resp.RequestConfig.AllowedGroups.Groups)]
			}
			for groupsCount, groupsItem := range resp.RequestConfig.AllowedGroups.Groups {
				var groups1 tfTypes.Group
				groups1.AppID = types.StringPointerValue(groupsItem.AppID)
				groups1.Description = types.StringPointerValue(groupsItem.Description)
				if groupsItem.GroupLifecycle != nil {
					groups1.GroupLifecycle = types.StringValue(string(*groupsItem.GroupLifecycle))
				} else {
					groups1.GroupLifecycle = types.StringNull()
				}
				groups1.ID = types.StringPointerValue(groupsItem.ID)
				groups1.IntegrationSpecificID = types.StringPointerValue(groupsItem.IntegrationSpecificID)
				groups1.Name = types.StringPointerValue(groupsItem.Name)
				groups1.SourceAppID = types.StringPointerValue(groupsItem.SourceAppID)
				if groupsCount+1 > len(r.RequestConfig.AllowedGroups.Groups) {
					r.RequestConfig.AllowedGroups.Groups = append(r.RequestConfig.AllowedGroups.Groups, groups1)
				} else {
					r.RequestConfig.AllowedGroups.Groups[groupsCount].AppID = groups1.AppID
					r.RequestConfig.AllowedGroups.Groups[groupsCount].Description = groups1.Description
					r.RequestConfig.AllowedGroups.Groups[groupsCount].GroupLifecycle = groups1.GroupLifecycle
					r.RequestConfig.AllowedGroups.Groups[groupsCount].ID = groups1.ID
					r.RequestConfig.AllowedGroups.Groups[groupsCount].IntegrationSpecificID = groups1.IntegrationSpecificID
					r.RequestConfig.AllowedGroups.Groups[groupsCount].Name = groups1.Name
					r.RequestConfig.AllowedGroups.Groups[groupsCount].SourceAppID = groups1.SourceAppID
				}
			}
			if resp.RequestConfig.AllowedGroups.Type != nil {
				r.RequestConfig.AllowedGroups.Type = types.StringValue(string(*resp.RequestConfig.AllowedGroups.Type))
			} else {
				r.RequestConfig.AllowedGroups.Type = types.StringNull()
			}
		}
		r.RequestConfig.AllowedGroupsOverride = types.BoolPointerValue(resp.RequestConfig.AllowedGroupsOverride)
		if resp.RequestConfig.AppstoreVisibility != nil {
			r.RequestConfig.AppstoreVisibility = types.StringValue(string(*resp.RequestConfig.AppstoreVisibility))
		} else {
			r.RequestConfig.AppstoreVisibility = types.StringNull()
		}
		if resp.RequestConfig.RequestApprovalConfig == nil {
			r.RequestConfig.RequestApprovalConfig = nil
		} else {
			r.RequestConfig.RequestApprovalConfig = &tfTypes.RequestablePermissionInputRequestApprovalConfig{}
			if resp.RequestConfig.RequestApprovalConfig.Approvers == nil {
				r.RequestConfig.RequestApprovalConfig.Approvers = nil
			} else {
				r.RequestConfig.RequestApprovalConfig.Approvers = &tfTypes.AddAppToAppStoreInputAdmins{}
				r.RequestConfig.RequestApprovalConfig.Approvers.Groups = []tfTypes.Group{}
				if len(r.RequestConfig.RequestApprovalConfig.Approvers.Groups) > len(resp.RequestConfig.RequestApprovalConfig.Approvers.Groups) {
					r.RequestConfig.RequestApprovalConfig.Approvers.Groups = r.RequestConfig.RequestApprovalConfig.Approvers.Groups[:len(resp.RequestConfig.RequestApprovalConfig.Approvers.Groups)]
				}
				for groupsCount1, groupsItem1 := range resp.RequestConfig.RequestApprovalConfig.Approvers.Groups {
					var groups3 tfTypes.Group
					groups3.AppID = types.StringPointerValue(groupsItem1.AppID)
					groups3.Description = types.StringPointerValue(groupsItem1.Description)
					if groupsItem1.GroupLifecycle != nil {
						groups3.GroupLifecycle = types.StringValue(string(*groupsItem1.GroupLifecycle))
					} else {
						groups3.GroupLifecycle = types.StringNull()
					}
					groups3.ID = types.StringPointerValue(groupsItem1.ID)
					groups3.IntegrationSpecificID = types.StringPointerValue(groupsItem1.IntegrationSpecificID)
					groups3.Name = types.StringPointerValue(groupsItem1.Name)
					groups3.SourceAppID = types.StringPointerValue(groupsItem1.SourceAppID)
					if groupsCount1+1 > len(r.RequestConfig.RequestApprovalConfig.Approvers.Groups) {
						r.RequestConfig.RequestApprovalConfig.Approvers.Groups = append(r.RequestConfig.RequestApprovalConfig.Approvers.Groups, groups3)
					} else {
						r.RequestConfig.RequestApprovalConfig.Approvers.Groups[groupsCount1].AppID = groups3.AppID
						r.RequestConfig.RequestApprovalConfig.Approvers.Groups[groupsCount1].Description = groups3.Description
						r.RequestConfig.RequestApprovalConfig.Approvers.Groups[groupsCount1].GroupLifecycle = groups3.GroupLifecycle
						r.RequestConfig.RequestApprovalConfig.Approvers.Groups[groupsCount1].ID = groups3.ID
						r.RequestConfig.RequestApprovalConfig.Approvers.Groups[groupsCount1].IntegrationSpecificID = groups3.IntegrationSpecificID
						r.RequestConfig.RequestApprovalConfig.Approvers.Groups[groupsCount1].Name = groups3.Name
						r.RequestConfig.RequestApprovalConfig.Approvers.Groups[groupsCount1].SourceAppID = groups3.SourceAppID
					}
				}
				r.RequestConfig.RequestApprovalConfig.Approvers.Users = []tfTypes.User{}
				if len(r.RequestConfig.RequestApprovalConfig.Approvers.Users) > len(resp.RequestConfig.RequestApprovalConfig.Approvers.Users) {
					r.RequestConfig.RequestApprovalConfig.Approvers.Users = r.RequestConfig.RequestApprovalConfig.Approvers.Users[:len(resp.RequestConfig.RequestApprovalConfig.Approvers.Users)]
				}
				for usersCount, usersItem := range resp.RequestConfig.RequestApprovalConfig.Approvers.Users {
					var users1 tfTypes.User
					users1.Email = types.StringPointerValue(usersItem.Email)
					users1.FamilyName = types.StringPointerValue(usersItem.FamilyName)
					users1.GivenName = types.StringPointerValue(usersItem.GivenName)
					users1.ID = types.StringValue(usersItem.ID)
					if usersItem.Status != nil {
						users1.Status = types.StringValue(string(*usersItem.Status))
					} else {
						users1.Status = types.StringNull()
					}
					if usersCount+1 > len(r.RequestConfig.RequestApprovalConfig.Approvers.Users) {
						r.RequestConfig.RequestApprovalConfig.Approvers.Users = append(r.RequestConfig.RequestApprovalConfig.Approvers.Users, users1)
					} else {
						r.RequestConfig.RequestApprovalConfig.Approvers.Users[usersCount].Email = users1.Email
						r.RequestConfig.RequestApprovalConfig.Approvers.Users[usersCount].FamilyName = users1.FamilyName
						r.RequestConfig.RequestApprovalConfig.Approvers.Users[usersCount].GivenName = users1.GivenName
						r.RequestConfig.RequestApprovalConfig.Approvers.Users[usersCount].ID = users1.ID
						r.RequestConfig.RequestApprovalConfig.Approvers.Users[usersCount].Status = users1.Status
					}
				}
			}
			if resp.RequestConfig.RequestApprovalConfig.ApproversStage2 == nil {
				r.RequestConfig.RequestApprovalConfig.ApproversStage2 = nil
			} else {
				r.RequestConfig.RequestApprovalConfig.ApproversStage2 = &tfTypes.AddAppToAppStoreInputAdmins{}
				r.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups = []tfTypes.Group{}
				if len(r.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups) > len(resp.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups) {
					r.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups = r.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups[:len(resp.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups)]
				}
				for groupsCount2, groupsItem2 := range resp.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups {
					var groups5 tfTypes.Group
					groups5.AppID = types.StringPointerValue(groupsItem2.AppID)
					groups5.Description = types.StringPointerValue(groupsItem2.Description)
					if groupsItem2.GroupLifecycle != nil {
						groups5.GroupLifecycle = types.StringValue(string(*groupsItem2.GroupLifecycle))
					} else {
						groups5.GroupLifecycle = types.StringNull()
					}
					groups5.ID = types.StringPointerValue(groupsItem2.ID)
					groups5.IntegrationSpecificID = types.StringPointerValue(groupsItem2.IntegrationSpecificID)
					groups5.Name = types.StringPointerValue(groupsItem2.Name)
					groups5.SourceAppID = types.StringPointerValue(groupsItem2.SourceAppID)
					if groupsCount2+1 > len(r.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups) {
						r.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups = append(r.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups, groups5)
					} else {
						r.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups[groupsCount2].AppID = groups5.AppID
						r.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups[groupsCount2].Description = groups5.Description
						r.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups[groupsCount2].GroupLifecycle = groups5.GroupLifecycle
						r.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups[groupsCount2].ID = groups5.ID
						r.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups[groupsCount2].IntegrationSpecificID = groups5.IntegrationSpecificID
						r.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups[groupsCount2].Name = groups5.Name
						r.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups[groupsCount2].SourceAppID = groups5.SourceAppID
					}
				}
				r.RequestConfig.RequestApprovalConfig.ApproversStage2.Users = []tfTypes.User{}
				if len(r.RequestConfig.RequestApprovalConfig.ApproversStage2.Users) > len(resp.RequestConfig.RequestApprovalConfig.ApproversStage2.Users) {
					r.RequestConfig.RequestApprovalConfig.ApproversStage2.Users = r.RequestConfig.RequestApprovalConfig.ApproversStage2.Users[:len(resp.RequestConfig.RequestApprovalConfig.ApproversStage2.Users)]
				}
				for usersCount1, usersItem1 := range resp.RequestConfig.RequestApprovalConfig.ApproversStage2.Users {
					var users3 tfTypes.User
					users3.Email = types.StringPointerValue(usersItem1.Email)
					users3.FamilyName = types.StringPointerValue(usersItem1.FamilyName)
					users3.GivenName = types.StringPointerValue(usersItem1.GivenName)
					users3.ID = types.StringValue(usersItem1.ID)
					if usersItem1.Status != nil {
						users3.Status = types.StringValue(string(*usersItem1.Status))
					} else {
						users3.Status = types.StringNull()
					}
					if usersCount1+1 > len(r.RequestConfig.RequestApprovalConfig.ApproversStage2.Users) {
						r.RequestConfig.RequestApprovalConfig.ApproversStage2.Users = append(r.RequestConfig.RequestApprovalConfig.ApproversStage2.Users, users3)
					} else {
						r.RequestConfig.RequestApprovalConfig.ApproversStage2.Users[usersCount1].Email = users3.Email
						r.RequestConfig.RequestApprovalConfig.ApproversStage2.Users[usersCount1].FamilyName = users3.FamilyName
						r.RequestConfig.RequestApprovalConfig.ApproversStage2.Users[usersCount1].GivenName = users3.GivenName
						r.RequestConfig.RequestApprovalConfig.ApproversStage2.Users[usersCount1].ID = users3.ID
						r.RequestConfig.RequestApprovalConfig.ApproversStage2.Users[usersCount1].Status = users3.Status
					}
				}
			}
			r.RequestConfig.RequestApprovalConfig.CustomApprovalMessage = types.StringPointerValue(resp.RequestConfig.RequestApprovalConfig.CustomApprovalMessage)
			r.RequestConfig.RequestApprovalConfig.CustomApprovalMessageOverride = types.BoolPointerValue(resp.RequestConfig.RequestApprovalConfig.CustomApprovalMessageOverride)
			if resp.RequestConfig.RequestApprovalConfig.ManagerApproval != nil {
				r.RequestConfig.RequestApprovalConfig.ManagerApproval = types.StringValue(string(*resp.RequestConfig.RequestApprovalConfig.ManagerApproval))
			} else {
				r.RequestConfig.RequestApprovalConfig.ManagerApproval = types.StringNull()
			}
			r.RequestConfig.RequestApprovalConfig.RequestApprovalConfigOverride = types.BoolPointerValue(resp.RequestConfig.RequestApprovalConfig.RequestApprovalConfigOverride)
			r.RequestConfig.RequestApprovalConfig.RequireAdditionalApproval = types.BoolPointerValue(resp.RequestConfig.RequestApprovalConfig.RequireAdditionalApproval)
		}
		if resp.RequestConfig.RequestFulfillmentConfig == nil {
			r.RequestConfig.RequestFulfillmentConfig = nil
		} else {
			r.RequestConfig.RequestFulfillmentConfig = &tfTypes.RequestablePermissionInputRequestFulfillmentConfig{}
			r.RequestConfig.RequestFulfillmentConfig.ManualInstructions = types.StringPointerValue(resp.RequestConfig.RequestFulfillmentConfig.ManualInstructions)
			r.RequestConfig.RequestFulfillmentConfig.ManualStepsNeeded = types.BoolPointerValue(resp.RequestConfig.RequestFulfillmentConfig.ManualStepsNeeded)
			if resp.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup == nil {
				r.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup = nil
			} else {
				r.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup = &tfTypes.Group{}
				r.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.AppID = types.StringPointerValue(resp.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.AppID)
				r.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.Description = types.StringPointerValue(resp.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.Description)
				if resp.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.GroupLifecycle != nil {
					r.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.GroupLifecycle = types.StringValue(string(*resp.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.GroupLifecycle))
				} else {
					r.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.GroupLifecycle = types.StringNull()
				}
				r.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.ID = types.StringPointerValue(resp.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.ID)
				r.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.IntegrationSpecificID = types.StringPointerValue(resp.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.IntegrationSpecificID)
				r.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.Name = types.StringPointerValue(resp.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.Name)
				r.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.SourceAppID = types.StringPointerValue(resp.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.SourceAppID)
			}
			if resp.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook == nil {
				r.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook = nil
			} else {
				r.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook = &tfTypes.AddAppToAppStoreInputAccessRemovalInlineWebhook{}
				r.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook.Description = types.StringPointerValue(resp.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook.Description)
				r.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook.HookType = types.StringValue(string(resp.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook.HookType))
				r.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook.ID = types.StringValue(resp.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook.ID)
				r.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook.Name = types.StringValue(resp.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook.Name)
			}
			r.RequestConfig.RequestFulfillmentConfig.TimeBasedAccess = []types.String{}
			for _, v := range resp.RequestConfig.RequestFulfillmentConfig.TimeBasedAccess {
				r.RequestConfig.RequestFulfillmentConfig.TimeBasedAccess = append(r.RequestConfig.RequestFulfillmentConfig.TimeBasedAccess, types.StringValue(string(v)))
			}
			r.RequestConfig.RequestFulfillmentConfig.TimeBasedAccessOverride = types.BoolPointerValue(resp.RequestConfig.RequestFulfillmentConfig.TimeBasedAccessOverride)
		}
		if resp.RequestConfig.RequestValidationInlineWebhook == nil {
			r.RequestConfig.RequestValidationInlineWebhook = nil
		} else {
			r.RequestConfig.RequestValidationInlineWebhook = &tfTypes.AddAppToAppStoreInputAccessRemovalInlineWebhook{}
			r.RequestConfig.RequestValidationInlineWebhook.Description = types.StringPointerValue(resp.RequestConfig.RequestValidationInlineWebhook.Description)
			r.RequestConfig.RequestValidationInlineWebhook.HookType = types.StringValue(string(resp.RequestConfig.RequestValidationInlineWebhook.HookType))
			r.RequestConfig.RequestValidationInlineWebhook.ID = types.StringValue(resp.RequestConfig.RequestValidationInlineWebhook.ID)
			r.RequestConfig.RequestValidationInlineWebhook.Name = types.StringValue(resp.RequestConfig.RequestValidationInlineWebhook.Name)
		}
		if resp.Type != nil {
			r.Type = types.StringValue(string(*resp.Type))
		} else {
			r.Type = types.StringNull()
		}
	}
}
