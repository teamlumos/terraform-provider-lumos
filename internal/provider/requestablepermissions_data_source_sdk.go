// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	tfTypes "github.com/teamlumos/terraform-provider-lumos/internal/provider/types"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
)

func (r *RequestablePermissionsDataSourceModel) RefreshFromSharedPageRequestablePermissionOutput(resp *shared.PageRequestablePermissionOutput) {
	if resp != nil {
		r.Items = []tfTypes.RequestablePermissionOutput{}
		if len(r.Items) > len(resp.Items) {
			r.Items = r.Items[:len(resp.Items)]
		}
		for itemsCount, itemsItem := range resp.Items {
			var items1 tfTypes.RequestablePermissionOutput
			items1.AppClassID = types.StringValue(itemsItem.AppClassID)
			items1.AppID = types.StringValue(itemsItem.AppID)
			items1.AppInstanceID = types.StringValue(itemsItem.AppInstanceID)
			items1.ID = types.StringPointerValue(itemsItem.ID)
			items1.Label = types.StringValue(itemsItem.Label)
			if itemsItem.RequestConfig.AccessRemovalInlineWebhook == nil {
				items1.RequestConfig.AccessRemovalInlineWebhook = nil
			} else {
				items1.RequestConfig.AccessRemovalInlineWebhook = &tfTypes.BaseInlineWebhook{}
				items1.RequestConfig.AccessRemovalInlineWebhook.Description = types.StringPointerValue(itemsItem.RequestConfig.AccessRemovalInlineWebhook.Description)
				items1.RequestConfig.AccessRemovalInlineWebhook.HookType = types.StringValue(string(itemsItem.RequestConfig.AccessRemovalInlineWebhook.HookType))
				items1.RequestConfig.AccessRemovalInlineWebhook.ID = types.StringValue(itemsItem.RequestConfig.AccessRemovalInlineWebhook.ID)
				items1.RequestConfig.AccessRemovalInlineWebhook.Name = types.StringValue(itemsItem.RequestConfig.AccessRemovalInlineWebhook.Name)
			}
			if itemsItem.RequestConfig.AllowedGroups == nil {
				items1.RequestConfig.AllowedGroups = nil
			} else {
				items1.RequestConfig.AllowedGroups = &tfTypes.AllowedGroupsConfigInput{}
				items1.RequestConfig.AllowedGroups.Groups = []tfTypes.Group{}
				for groupsCount, groupsItem := range itemsItem.RequestConfig.AllowedGroups.Groups {
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
					if groupsCount+1 > len(items1.RequestConfig.AllowedGroups.Groups) {
						items1.RequestConfig.AllowedGroups.Groups = append(items1.RequestConfig.AllowedGroups.Groups, groups1)
					} else {
						items1.RequestConfig.AllowedGroups.Groups[groupsCount].AppID = groups1.AppID
						items1.RequestConfig.AllowedGroups.Groups[groupsCount].Description = groups1.Description
						items1.RequestConfig.AllowedGroups.Groups[groupsCount].GroupLifecycle = groups1.GroupLifecycle
						items1.RequestConfig.AllowedGroups.Groups[groupsCount].ID = groups1.ID
						items1.RequestConfig.AllowedGroups.Groups[groupsCount].IntegrationSpecificID = groups1.IntegrationSpecificID
						items1.RequestConfig.AllowedGroups.Groups[groupsCount].Name = groups1.Name
						items1.RequestConfig.AllowedGroups.Groups[groupsCount].SourceAppID = groups1.SourceAppID
					}
				}
				if itemsItem.RequestConfig.AllowedGroups.Type != nil {
					items1.RequestConfig.AllowedGroups.Type = types.StringValue(string(*itemsItem.RequestConfig.AllowedGroups.Type))
				} else {
					items1.RequestConfig.AllowedGroups.Type = types.StringNull()
				}
			}
			items1.RequestConfig.AllowedGroupsOverride = types.BoolPointerValue(itemsItem.RequestConfig.AllowedGroupsOverride)
			if itemsItem.RequestConfig.AppstoreVisibility != nil {
				items1.RequestConfig.AppstoreVisibility = types.StringValue(string(*itemsItem.RequestConfig.AppstoreVisibility))
			} else {
				items1.RequestConfig.AppstoreVisibility = types.StringNull()
			}
			if itemsItem.RequestConfig.RequestApprovalConfig == nil {
				items1.RequestConfig.RequestApprovalConfig = nil
			} else {
				items1.RequestConfig.RequestApprovalConfig = &tfTypes.RequestApprovalConfigOutput{}
				if itemsItem.RequestConfig.RequestApprovalConfig.Approvers == nil {
					items1.RequestConfig.RequestApprovalConfig.Approvers = nil
				} else {
					items1.RequestConfig.RequestApprovalConfig.Approvers = &tfTypes.AppApproversInput{}
					items1.RequestConfig.RequestApprovalConfig.Approvers.Groups = []tfTypes.Group{}
					for groupsCount1, groupsItem1 := range itemsItem.RequestConfig.RequestApprovalConfig.Approvers.Groups {
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
						if groupsCount1+1 > len(items1.RequestConfig.RequestApprovalConfig.Approvers.Groups) {
							items1.RequestConfig.RequestApprovalConfig.Approvers.Groups = append(items1.RequestConfig.RequestApprovalConfig.Approvers.Groups, groups3)
						} else {
							items1.RequestConfig.RequestApprovalConfig.Approvers.Groups[groupsCount1].AppID = groups3.AppID
							items1.RequestConfig.RequestApprovalConfig.Approvers.Groups[groupsCount1].Description = groups3.Description
							items1.RequestConfig.RequestApprovalConfig.Approvers.Groups[groupsCount1].GroupLifecycle = groups3.GroupLifecycle
							items1.RequestConfig.RequestApprovalConfig.Approvers.Groups[groupsCount1].ID = groups3.ID
							items1.RequestConfig.RequestApprovalConfig.Approvers.Groups[groupsCount1].IntegrationSpecificID = groups3.IntegrationSpecificID
							items1.RequestConfig.RequestApprovalConfig.Approvers.Groups[groupsCount1].Name = groups3.Name
							items1.RequestConfig.RequestApprovalConfig.Approvers.Groups[groupsCount1].SourceAppID = groups3.SourceAppID
						}
					}
					items1.RequestConfig.RequestApprovalConfig.Approvers.Users = []tfTypes.User{}
					for usersCount, usersItem := range itemsItem.RequestConfig.RequestApprovalConfig.Approvers.Users {
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
						if usersCount+1 > len(items1.RequestConfig.RequestApprovalConfig.Approvers.Users) {
							items1.RequestConfig.RequestApprovalConfig.Approvers.Users = append(items1.RequestConfig.RequestApprovalConfig.Approvers.Users, users1)
						} else {
							items1.RequestConfig.RequestApprovalConfig.Approvers.Users[usersCount].Email = users1.Email
							items1.RequestConfig.RequestApprovalConfig.Approvers.Users[usersCount].FamilyName = users1.FamilyName
							items1.RequestConfig.RequestApprovalConfig.Approvers.Users[usersCount].GivenName = users1.GivenName
							items1.RequestConfig.RequestApprovalConfig.Approvers.Users[usersCount].ID = users1.ID
							items1.RequestConfig.RequestApprovalConfig.Approvers.Users[usersCount].Status = users1.Status
						}
					}
				}
				if itemsItem.RequestConfig.RequestApprovalConfig.ApproversStage2 == nil {
					items1.RequestConfig.RequestApprovalConfig.ApproversStage2 = nil
				} else {
					items1.RequestConfig.RequestApprovalConfig.ApproversStage2 = &tfTypes.AppApproversInput{}
					items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups = []tfTypes.Group{}
					for groupsCount2, groupsItem2 := range itemsItem.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups {
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
						if groupsCount2+1 > len(items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups) {
							items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups = append(items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups, groups5)
						} else {
							items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups[groupsCount2].AppID = groups5.AppID
							items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups[groupsCount2].Description = groups5.Description
							items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups[groupsCount2].GroupLifecycle = groups5.GroupLifecycle
							items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups[groupsCount2].ID = groups5.ID
							items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups[groupsCount2].IntegrationSpecificID = groups5.IntegrationSpecificID
							items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups[groupsCount2].Name = groups5.Name
							items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Groups[groupsCount2].SourceAppID = groups5.SourceAppID
						}
					}
					items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Users = []tfTypes.User{}
					for usersCount1, usersItem1 := range itemsItem.RequestConfig.RequestApprovalConfig.ApproversStage2.Users {
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
						if usersCount1+1 > len(items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Users) {
							items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Users = append(items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Users, users3)
						} else {
							items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Users[usersCount1].Email = users3.Email
							items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Users[usersCount1].FamilyName = users3.FamilyName
							items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Users[usersCount1].GivenName = users3.GivenName
							items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Users[usersCount1].ID = users3.ID
							items1.RequestConfig.RequestApprovalConfig.ApproversStage2.Users[usersCount1].Status = users3.Status
						}
					}
				}
				items1.RequestConfig.RequestApprovalConfig.CustomApprovalMessage = types.StringPointerValue(itemsItem.RequestConfig.RequestApprovalConfig.CustomApprovalMessage)
				items1.RequestConfig.RequestApprovalConfig.CustomApprovalMessageOverride = types.BoolPointerValue(itemsItem.RequestConfig.RequestApprovalConfig.CustomApprovalMessageOverride)
				if itemsItem.RequestConfig.RequestApprovalConfig.ManagerApproval != nil {
					items1.RequestConfig.RequestApprovalConfig.ManagerApproval = types.StringValue(string(*itemsItem.RequestConfig.RequestApprovalConfig.ManagerApproval))
				} else {
					items1.RequestConfig.RequestApprovalConfig.ManagerApproval = types.StringNull()
				}
				items1.RequestConfig.RequestApprovalConfig.RequestApprovalConfigOverride = types.BoolPointerValue(itemsItem.RequestConfig.RequestApprovalConfig.RequestApprovalConfigOverride)
				items1.RequestConfig.RequestApprovalConfig.RequireAdditionalApproval = types.BoolPointerValue(itemsItem.RequestConfig.RequestApprovalConfig.RequireAdditionalApproval)
			}
			if itemsItem.RequestConfig.RequestFulfillmentConfig == nil {
				items1.RequestConfig.RequestFulfillmentConfig = nil
			} else {
				items1.RequestConfig.RequestFulfillmentConfig = &tfTypes.RequestFulfillmentConfigInput{}
				items1.RequestConfig.RequestFulfillmentConfig.ManualInstructions = types.StringPointerValue(itemsItem.RequestConfig.RequestFulfillmentConfig.ManualInstructions)
				items1.RequestConfig.RequestFulfillmentConfig.ManualStepsNeeded = types.BoolPointerValue(itemsItem.RequestConfig.RequestFulfillmentConfig.ManualStepsNeeded)
				if itemsItem.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup == nil {
					items1.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup = nil
				} else {
					items1.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup = &tfTypes.Group{}
					items1.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.AppID = types.StringPointerValue(itemsItem.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.AppID)
					items1.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.Description = types.StringPointerValue(itemsItem.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.Description)
					if itemsItem.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.GroupLifecycle != nil {
						items1.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.GroupLifecycle = types.StringValue(string(*itemsItem.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.GroupLifecycle))
					} else {
						items1.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.GroupLifecycle = types.StringNull()
					}
					items1.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.ID = types.StringPointerValue(itemsItem.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.ID)
					items1.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.IntegrationSpecificID = types.StringPointerValue(itemsItem.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.IntegrationSpecificID)
					items1.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.Name = types.StringPointerValue(itemsItem.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.Name)
					items1.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.SourceAppID = types.StringPointerValue(itemsItem.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.SourceAppID)
				}
				if itemsItem.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook == nil {
					items1.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook = nil
				} else {
					items1.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook = &tfTypes.BaseInlineWebhook{}
					items1.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook.Description = types.StringPointerValue(itemsItem.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook.Description)
					items1.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook.HookType = types.StringValue(string(itemsItem.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook.HookType))
					items1.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook.ID = types.StringValue(itemsItem.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook.ID)
					items1.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook.Name = types.StringValue(itemsItem.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook.Name)
				}
				items1.RequestConfig.RequestFulfillmentConfig.TimeBasedAccess = []types.String{}
				for _, v := range itemsItem.RequestConfig.RequestFulfillmentConfig.TimeBasedAccess {
					items1.RequestConfig.RequestFulfillmentConfig.TimeBasedAccess = append(items1.RequestConfig.RequestFulfillmentConfig.TimeBasedAccess, types.StringValue(string(v)))
				}
				items1.RequestConfig.RequestFulfillmentConfig.TimeBasedAccessOverride = types.BoolPointerValue(itemsItem.RequestConfig.RequestFulfillmentConfig.TimeBasedAccessOverride)
			}
			if itemsItem.RequestConfig.RequestValidationInlineWebhook == nil {
				items1.RequestConfig.RequestValidationInlineWebhook = nil
			} else {
				items1.RequestConfig.RequestValidationInlineWebhook = &tfTypes.BaseInlineWebhook{}
				items1.RequestConfig.RequestValidationInlineWebhook.Description = types.StringPointerValue(itemsItem.RequestConfig.RequestValidationInlineWebhook.Description)
				items1.RequestConfig.RequestValidationInlineWebhook.HookType = types.StringValue(string(itemsItem.RequestConfig.RequestValidationInlineWebhook.HookType))
				items1.RequestConfig.RequestValidationInlineWebhook.ID = types.StringValue(itemsItem.RequestConfig.RequestValidationInlineWebhook.ID)
				items1.RequestConfig.RequestValidationInlineWebhook.Name = types.StringValue(itemsItem.RequestConfig.RequestValidationInlineWebhook.Name)
			}
			if itemsItem.Type != nil {
				items1.Type = types.StringValue(string(*itemsItem.Type))
			} else {
				items1.Type = types.StringNull()
			}
			if itemsCount+1 > len(r.Items) {
				r.Items = append(r.Items, items1)
			} else {
				r.Items[itemsCount].AppClassID = items1.AppClassID
				r.Items[itemsCount].AppID = items1.AppID
				r.Items[itemsCount].AppInstanceID = items1.AppInstanceID
				r.Items[itemsCount].ID = items1.ID
				r.Items[itemsCount].Label = items1.Label
				r.Items[itemsCount].RequestConfig = items1.RequestConfig
				r.Items[itemsCount].Type = items1.Type
			}
		}
		r.Page = types.Int64PointerValue(resp.Page)
		r.Pages = types.Int64PointerValue(resp.Pages)
		r.Size = types.Int64PointerValue(resp.Size)
		r.Total = types.Int64PointerValue(resp.Total)
	}
}
