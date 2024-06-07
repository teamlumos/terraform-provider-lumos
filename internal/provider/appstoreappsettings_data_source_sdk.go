// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	tfTypes "github.com/teamlumos/terraform-provider-lumos/internal/provider/types"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
)

func (r *AppStoreAppSettingsDataSourceModel) RefreshFromSharedAppStoreAppSettingsOutput(resp *shared.AppStoreAppSettingsOutput) {
	if resp != nil {
		r.CustomRequestInstructions = types.StringPointerValue(resp.CustomRequestInstructions)
		if resp.Provisioning == nil {
			r.Provisioning = nil
		} else {
			r.Provisioning = &tfTypes.AddAppToAppStoreInputProvisioning{}
			if resp.Provisioning.AccessRemovalInlineWebhook == nil {
				r.Provisioning.AccessRemovalInlineWebhook = nil
			} else {
				r.Provisioning.AccessRemovalInlineWebhook = &tfTypes.AddAppToAppStoreInputAccessRemovalInlineWebhook{}
				r.Provisioning.AccessRemovalInlineWebhook.Description = types.StringPointerValue(resp.Provisioning.AccessRemovalInlineWebhook.Description)
				r.Provisioning.AccessRemovalInlineWebhook.HookType = types.StringValue(string(resp.Provisioning.AccessRemovalInlineWebhook.HookType))
				r.Provisioning.AccessRemovalInlineWebhook.ID = types.StringValue(resp.Provisioning.AccessRemovalInlineWebhook.ID)
				r.Provisioning.AccessRemovalInlineWebhook.Name = types.StringValue(resp.Provisioning.AccessRemovalInlineWebhook.Name)
			}
			r.Provisioning.AllowMultiplePermissionSelection = types.BoolPointerValue(resp.Provisioning.AllowMultiplePermissionSelection)
			r.Provisioning.CustomProvisioningInstructions = types.StringPointerValue(resp.Provisioning.CustomProvisioningInstructions)
			if resp.Provisioning.GroupsProvisioning != nil {
				r.Provisioning.GroupsProvisioning = types.StringValue(string(*resp.Provisioning.GroupsProvisioning))
			} else {
				r.Provisioning.GroupsProvisioning = types.StringNull()
			}
			r.Provisioning.ManualStepsNeeded = types.BoolPointerValue(resp.Provisioning.ManualStepsNeeded)
			if resp.Provisioning.ProvisioningWebhook == nil {
				r.Provisioning.ProvisioningWebhook = nil
			} else {
				r.Provisioning.ProvisioningWebhook = &tfTypes.AddAppToAppStoreInputAccessRemovalInlineWebhook{}
				r.Provisioning.ProvisioningWebhook.Description = types.StringPointerValue(resp.Provisioning.ProvisioningWebhook.Description)
				r.Provisioning.ProvisioningWebhook.HookType = types.StringValue(string(resp.Provisioning.ProvisioningWebhook.HookType))
				r.Provisioning.ProvisioningWebhook.ID = types.StringValue(resp.Provisioning.ProvisioningWebhook.ID)
				r.Provisioning.ProvisioningWebhook.Name = types.StringValue(resp.Provisioning.ProvisioningWebhook.Name)
			}
			r.Provisioning.TimeBasedAccess = []types.String{}
			for _, v := range resp.Provisioning.TimeBasedAccess {
				r.Provisioning.TimeBasedAccess = append(r.Provisioning.TimeBasedAccess, types.StringValue(string(v)))
			}
		}
		if resp.RequestFlow == nil {
			r.RequestFlow = nil
		} else {
			r.RequestFlow = &tfTypes.AddAppToAppStoreInputRequestFlow{}
			if resp.RequestFlow.Admins == nil {
				r.RequestFlow.Admins = nil
			} else {
				r.RequestFlow.Admins = &tfTypes.AddAppToAppStoreInputAdmins{}
				r.RequestFlow.Admins.Groups = []tfTypes.Group{}
				if len(r.RequestFlow.Admins.Groups) > len(resp.RequestFlow.Admins.Groups) {
					r.RequestFlow.Admins.Groups = r.RequestFlow.Admins.Groups[:len(resp.RequestFlow.Admins.Groups)]
				}
				for groupsCount, groupsItem := range resp.RequestFlow.Admins.Groups {
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
					if groupsCount+1 > len(r.RequestFlow.Admins.Groups) {
						r.RequestFlow.Admins.Groups = append(r.RequestFlow.Admins.Groups, groups1)
					} else {
						r.RequestFlow.Admins.Groups[groupsCount].AppID = groups1.AppID
						r.RequestFlow.Admins.Groups[groupsCount].Description = groups1.Description
						r.RequestFlow.Admins.Groups[groupsCount].GroupLifecycle = groups1.GroupLifecycle
						r.RequestFlow.Admins.Groups[groupsCount].ID = groups1.ID
						r.RequestFlow.Admins.Groups[groupsCount].IntegrationSpecificID = groups1.IntegrationSpecificID
						r.RequestFlow.Admins.Groups[groupsCount].Name = groups1.Name
						r.RequestFlow.Admins.Groups[groupsCount].SourceAppID = groups1.SourceAppID
					}
				}
				r.RequestFlow.Admins.Users = []tfTypes.User{}
				if len(r.RequestFlow.Admins.Users) > len(resp.RequestFlow.Admins.Users) {
					r.RequestFlow.Admins.Users = r.RequestFlow.Admins.Users[:len(resp.RequestFlow.Admins.Users)]
				}
				for usersCount, usersItem := range resp.RequestFlow.Admins.Users {
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
					if usersCount+1 > len(r.RequestFlow.Admins.Users) {
						r.RequestFlow.Admins.Users = append(r.RequestFlow.Admins.Users, users1)
					} else {
						r.RequestFlow.Admins.Users[usersCount].Email = users1.Email
						r.RequestFlow.Admins.Users[usersCount].FamilyName = users1.FamilyName
						r.RequestFlow.Admins.Users[usersCount].GivenName = users1.GivenName
						r.RequestFlow.Admins.Users[usersCount].ID = users1.ID
						r.RequestFlow.Admins.Users[usersCount].Status = users1.Status
					}
				}
			}
			if resp.RequestFlow.AllowedGroups == nil {
				r.RequestFlow.AllowedGroups = nil
			} else {
				r.RequestFlow.AllowedGroups = &tfTypes.AddAppToAppStoreInputAllowedGroups{}
				r.RequestFlow.AllowedGroups.Groups = []tfTypes.Group{}
				if len(r.RequestFlow.AllowedGroups.Groups) > len(resp.RequestFlow.AllowedGroups.Groups) {
					r.RequestFlow.AllowedGroups.Groups = r.RequestFlow.AllowedGroups.Groups[:len(resp.RequestFlow.AllowedGroups.Groups)]
				}
				for groupsCount1, groupsItem1 := range resp.RequestFlow.AllowedGroups.Groups {
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
					if groupsCount1+1 > len(r.RequestFlow.AllowedGroups.Groups) {
						r.RequestFlow.AllowedGroups.Groups = append(r.RequestFlow.AllowedGroups.Groups, groups3)
					} else {
						r.RequestFlow.AllowedGroups.Groups[groupsCount1].AppID = groups3.AppID
						r.RequestFlow.AllowedGroups.Groups[groupsCount1].Description = groups3.Description
						r.RequestFlow.AllowedGroups.Groups[groupsCount1].GroupLifecycle = groups3.GroupLifecycle
						r.RequestFlow.AllowedGroups.Groups[groupsCount1].ID = groups3.ID
						r.RequestFlow.AllowedGroups.Groups[groupsCount1].IntegrationSpecificID = groups3.IntegrationSpecificID
						r.RequestFlow.AllowedGroups.Groups[groupsCount1].Name = groups3.Name
						r.RequestFlow.AllowedGroups.Groups[groupsCount1].SourceAppID = groups3.SourceAppID
					}
				}
				if resp.RequestFlow.AllowedGroups.Type != nil {
					r.RequestFlow.AllowedGroups.Type = types.StringValue(string(*resp.RequestFlow.AllowedGroups.Type))
				} else {
					r.RequestFlow.AllowedGroups.Type = types.StringNull()
				}
			}
			if resp.RequestFlow.Approvers == nil {
				r.RequestFlow.Approvers = nil
			} else {
				r.RequestFlow.Approvers = &tfTypes.AddAppToAppStoreInputAdmins{}
				r.RequestFlow.Approvers.Groups = []tfTypes.Group{}
				if len(r.RequestFlow.Approvers.Groups) > len(resp.RequestFlow.Approvers.Groups) {
					r.RequestFlow.Approvers.Groups = r.RequestFlow.Approvers.Groups[:len(resp.RequestFlow.Approvers.Groups)]
				}
				for groupsCount2, groupsItem2 := range resp.RequestFlow.Approvers.Groups {
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
					if groupsCount2+1 > len(r.RequestFlow.Approvers.Groups) {
						r.RequestFlow.Approvers.Groups = append(r.RequestFlow.Approvers.Groups, groups5)
					} else {
						r.RequestFlow.Approvers.Groups[groupsCount2].AppID = groups5.AppID
						r.RequestFlow.Approvers.Groups[groupsCount2].Description = groups5.Description
						r.RequestFlow.Approvers.Groups[groupsCount2].GroupLifecycle = groups5.GroupLifecycle
						r.RequestFlow.Approvers.Groups[groupsCount2].ID = groups5.ID
						r.RequestFlow.Approvers.Groups[groupsCount2].IntegrationSpecificID = groups5.IntegrationSpecificID
						r.RequestFlow.Approvers.Groups[groupsCount2].Name = groups5.Name
						r.RequestFlow.Approvers.Groups[groupsCount2].SourceAppID = groups5.SourceAppID
					}
				}
				r.RequestFlow.Approvers.Users = []tfTypes.User{}
				if len(r.RequestFlow.Approvers.Users) > len(resp.RequestFlow.Approvers.Users) {
					r.RequestFlow.Approvers.Users = r.RequestFlow.Approvers.Users[:len(resp.RequestFlow.Approvers.Users)]
				}
				for usersCount1, usersItem1 := range resp.RequestFlow.Approvers.Users {
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
					if usersCount1+1 > len(r.RequestFlow.Approvers.Users) {
						r.RequestFlow.Approvers.Users = append(r.RequestFlow.Approvers.Users, users3)
					} else {
						r.RequestFlow.Approvers.Users[usersCount1].Email = users3.Email
						r.RequestFlow.Approvers.Users[usersCount1].FamilyName = users3.FamilyName
						r.RequestFlow.Approvers.Users[usersCount1].GivenName = users3.GivenName
						r.RequestFlow.Approvers.Users[usersCount1].ID = users3.ID
						r.RequestFlow.Approvers.Users[usersCount1].Status = users3.Status
					}
				}
			}
			if resp.RequestFlow.ApproversStage2 == nil {
				r.RequestFlow.ApproversStage2 = nil
			} else {
				r.RequestFlow.ApproversStage2 = &tfTypes.AddAppToAppStoreInputAdmins{}
				r.RequestFlow.ApproversStage2.Groups = []tfTypes.Group{}
				if len(r.RequestFlow.ApproversStage2.Groups) > len(resp.RequestFlow.ApproversStage2.Groups) {
					r.RequestFlow.ApproversStage2.Groups = r.RequestFlow.ApproversStage2.Groups[:len(resp.RequestFlow.ApproversStage2.Groups)]
				}
				for groupsCount3, groupsItem3 := range resp.RequestFlow.ApproversStage2.Groups {
					var groups7 tfTypes.Group
					groups7.AppID = types.StringPointerValue(groupsItem3.AppID)
					groups7.Description = types.StringPointerValue(groupsItem3.Description)
					if groupsItem3.GroupLifecycle != nil {
						groups7.GroupLifecycle = types.StringValue(string(*groupsItem3.GroupLifecycle))
					} else {
						groups7.GroupLifecycle = types.StringNull()
					}
					groups7.ID = types.StringPointerValue(groupsItem3.ID)
					groups7.IntegrationSpecificID = types.StringPointerValue(groupsItem3.IntegrationSpecificID)
					groups7.Name = types.StringPointerValue(groupsItem3.Name)
					groups7.SourceAppID = types.StringPointerValue(groupsItem3.SourceAppID)
					if groupsCount3+1 > len(r.RequestFlow.ApproversStage2.Groups) {
						r.RequestFlow.ApproversStage2.Groups = append(r.RequestFlow.ApproversStage2.Groups, groups7)
					} else {
						r.RequestFlow.ApproversStage2.Groups[groupsCount3].AppID = groups7.AppID
						r.RequestFlow.ApproversStage2.Groups[groupsCount3].Description = groups7.Description
						r.RequestFlow.ApproversStage2.Groups[groupsCount3].GroupLifecycle = groups7.GroupLifecycle
						r.RequestFlow.ApproversStage2.Groups[groupsCount3].ID = groups7.ID
						r.RequestFlow.ApproversStage2.Groups[groupsCount3].IntegrationSpecificID = groups7.IntegrationSpecificID
						r.RequestFlow.ApproversStage2.Groups[groupsCount3].Name = groups7.Name
						r.RequestFlow.ApproversStage2.Groups[groupsCount3].SourceAppID = groups7.SourceAppID
					}
				}
				r.RequestFlow.ApproversStage2.Users = []tfTypes.User{}
				if len(r.RequestFlow.ApproversStage2.Users) > len(resp.RequestFlow.ApproversStage2.Users) {
					r.RequestFlow.ApproversStage2.Users = r.RequestFlow.ApproversStage2.Users[:len(resp.RequestFlow.ApproversStage2.Users)]
				}
				for usersCount2, usersItem2 := range resp.RequestFlow.ApproversStage2.Users {
					var users5 tfTypes.User
					users5.Email = types.StringPointerValue(usersItem2.Email)
					users5.FamilyName = types.StringPointerValue(usersItem2.FamilyName)
					users5.GivenName = types.StringPointerValue(usersItem2.GivenName)
					users5.ID = types.StringValue(usersItem2.ID)
					if usersItem2.Status != nil {
						users5.Status = types.StringValue(string(*usersItem2.Status))
					} else {
						users5.Status = types.StringNull()
					}
					if usersCount2+1 > len(r.RequestFlow.ApproversStage2.Users) {
						r.RequestFlow.ApproversStage2.Users = append(r.RequestFlow.ApproversStage2.Users, users5)
					} else {
						r.RequestFlow.ApproversStage2.Users[usersCount2].Email = users5.Email
						r.RequestFlow.ApproversStage2.Users[usersCount2].FamilyName = users5.FamilyName
						r.RequestFlow.ApproversStage2.Users[usersCount2].GivenName = users5.GivenName
						r.RequestFlow.ApproversStage2.Users[usersCount2].ID = users5.ID
						r.RequestFlow.ApproversStage2.Users[usersCount2].Status = users5.Status
					}
				}
			}
			r.RequestFlow.CustomApprovalMessage = types.StringPointerValue(resp.RequestFlow.CustomApprovalMessage)
			if resp.RequestFlow.Discoverability != nil {
				r.RequestFlow.Discoverability = types.StringValue(string(*resp.RequestFlow.Discoverability))
			} else {
				r.RequestFlow.Discoverability = types.StringNull()
			}
			if resp.RequestFlow.RequestValidationInlineWebhook == nil {
				r.RequestFlow.RequestValidationInlineWebhook = nil
			} else {
				r.RequestFlow.RequestValidationInlineWebhook = &tfTypes.AddAppToAppStoreInputAccessRemovalInlineWebhook{}
				r.RequestFlow.RequestValidationInlineWebhook.Description = types.StringPointerValue(resp.RequestFlow.RequestValidationInlineWebhook.Description)
				r.RequestFlow.RequestValidationInlineWebhook.HookType = types.StringValue(string(resp.RequestFlow.RequestValidationInlineWebhook.HookType))
				r.RequestFlow.RequestValidationInlineWebhook.ID = types.StringValue(resp.RequestFlow.RequestValidationInlineWebhook.ID)
				r.RequestFlow.RequestValidationInlineWebhook.Name = types.StringValue(resp.RequestFlow.RequestValidationInlineWebhook.Name)
			}
			r.RequestFlow.RequireAdditionalApproval = types.BoolPointerValue(resp.RequestFlow.RequireAdditionalApproval)
			r.RequestFlow.RequireManagerApproval = types.BoolPointerValue(resp.RequestFlow.RequireManagerApproval)
		}
	}
}
