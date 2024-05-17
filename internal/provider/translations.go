package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/teamlumos/lumos-go-sdk/lumos"
	"github.com/teamlumos/terraform-provider-lumos/internal/generated/resource_pre_approval_rule"
)

func PreApprovalRuleTerraformToApi(ctx context.Context, input *resource_pre_approval_rule.PreApprovalRuleModel) (*lumos.PreApprovalRuleInput, error) {

	preApprovalRuleInput := lumos.PreApprovalRuleInput{}
	if !input.AppClassId.IsNull() && !input.AppClassId.IsUnknown() {
		preApprovalRuleInput.AppClassId = input.AppClassId.ValueStringPointer()
	}
	if !input.AppId.IsNull() && !input.AppId.IsUnknown() {
		preApprovalRuleInput.AppId = input.AppId.ValueStringPointer()

	}
	if !input.AppInstanceId.IsNull() && !input.AppInstanceId.IsUnknown() {
		preApprovalRuleInput.AppInstanceId = input.AppInstanceId.ValueStringPointer()

	}
	if !input.AppClassId.IsNull() && !input.AppClassId.IsUnknown() {
		preApprovalRuleInput.AppClassId = input.AppClassId.ValueStringPointer()
	}
	if !input.Justification.IsNull() && !input.Justification.IsUnknown() {
		preApprovalRuleInput.Justification = input.Justification.ValueStringPointer()
	}
	if !input.PreapprovedPermissions.IsNull() && !input.PreapprovedPermissions.IsUnknown() {
		var permissionElements = input.PreapprovedPermissions.Elements()
		preApprovedPermissions := make([]lumos.RequestablePermissionBase, len(permissionElements))
		for i, elem := range permissionElements {
			groupValue, ok := elem.(resource_pre_approval_rule.PreapprovedPermissionsValue)
			if !ok {
				return nil, fmt.Errorf("Unable to convert RequestablePermissionBase to RequestablePermissionBase")
			}

			preApprovedPermissions[i] = lumos.RequestablePermissionBase{
				Id: groupValue.Id.ValueStringPointer(),
			}
		}
		preApprovalRuleInput.PreapprovedPermissions = &preApprovedPermissions
	} else {
		preApprovedPermissions := make([]lumos.RequestablePermissionBase, 0)
		preApprovalRuleInput.PreapprovedPermissions = &preApprovedPermissions
	}
	if !input.PreapprovedGroups.IsNull() && !input.PreapprovedGroups.IsUnknown() {
		var groupElements = input.PreapprovedGroups.Elements()
		preApprovedGroups := make([]lumos.BaseGroup, len(groupElements))
		for i, elem := range groupElements {
			groupValue, ok := elem.(resource_pre_approval_rule.PreapprovedGroupsValue)
			if !ok {
				return nil, fmt.Errorf("Unable to convert BaseGroup to PreapprovedGroupsValue")
			}

			/* The caller must specify either (Id) or (AppId, IntegrationSpecificId) */
			if !groupValue.Id.IsNull() && !groupValue.Id.IsUnknown() {
				preApprovedGroups[i] = lumos.BaseGroup{
					Id:                    groupValue.Id.ValueStringPointer(),
					AppId:                 nil,
					IntegrationSpecificId: nil,
				}
			} else if !groupValue.AppId.IsNull() && !groupValue.AppId.IsUnknown() && !groupValue.IntegrationSpecificId.IsNull() && !groupValue.IntegrationSpecificId.IsUnknown() {
				preApprovedGroups[i] = lumos.BaseGroup{
					Id:                    nil,
					AppId:                 groupValue.AppId.ValueStringPointer(),
					IntegrationSpecificId: groupValue.IntegrationSpecificId.ValueStringPointer(),
				}
			} else {
				return nil, fmt.Errorf("Either (Id) or (AppId, IntegrationSpecificId) must be specified")
			}
		}
		preApprovalRuleInput.PreapprovedGroups = &preApprovedGroups
	}
	if !input.PreapprovalWebhooks.IsNull() && !input.PreapprovalWebhooks.IsUnknown() {
		var webhookElements = input.PreapprovalWebhooks.Elements()
		preApprovalWebhooks := make([]lumos.BaseInlineWebhook, len(webhookElements))
		for i, elem := range webhookElements {
			groupValue, ok := elem.(resource_pre_approval_rule.PreapprovalWebhooksValue)
			if !ok {
				return nil, fmt.Errorf("Unable to convert PreapprovalWebhooksValue to BaseInlineWebhook")
			}

			preApprovalWebhooks[i] = lumos.BaseInlineWebhook{
				Id: groupValue.Id.ValueString(),
			}
		}
		preApprovalRuleInput.PreapprovalWebhooks = &preApprovalWebhooks
	} else {
		preApprovalWebhooks := make([]lumos.BaseInlineWebhook, 0)
		preApprovalRuleInput.PreapprovalWebhooks = &preApprovalWebhooks
	}

	return &preApprovalRuleInput, nil
}

func PreApprovalRuleApiToTerraform(ctx context.Context, apiResponse lumos.PreApprovalRuleOutput, terraformModel resource_pre_approval_rule.PreApprovalRuleModel) resource_pre_approval_rule.PreApprovalRuleModel {
	terraformModel.Id = types.StringPointerValue(apiResponse.Id)
	terraformModel.AppClassId = types.StringPointerValue(apiResponse.AppClassId)
	terraformModel.AppId = types.StringPointerValue(apiResponse.AppId)
	terraformModel.AppInstanceId = types.StringPointerValue(apiResponse.AppInstanceId)
	terraformModel.Justification = types.StringValue(apiResponse.Justification)

	var timeBasedAccess = []string{}
	for i, timeBasedAccessOption := range *apiResponse.TimeBasedAccess {
		timeBasedAccess[i] = string(timeBasedAccessOption)
	}
	terraformModel.TimeBasedAccess, _ = types.ListValueFrom(ctx, types.StringType, timeBasedAccess)

	var preApprovedGroups []resource_pre_approval_rule.PreapprovedGroupsValue
	terraformModel.PreapprovedGroups.ElementsAs(ctx, &preApprovedGroups, false)
	for i, group := range *apiResponse.PreapprovedGroups {
		preApprovedGroups[i].Id = types.StringPointerValue(group.Id)
		preApprovedGroups[i].AppId = types.StringPointerValue(group.AppId)
		preApprovedGroups[i].IntegrationSpecificId = types.StringPointerValue(group.IntegrationSpecificId)
		preApprovedGroups[i].Name = types.StringPointerValue(group.Name)
		preApprovedGroups[i].Description = types.StringPointerValue(group.Description)
		preApprovedGroups[i].Lifecycle = types.StringPointerValue((*string)(group.Lifecycle))
		preApprovedGroups[i].SourceAppId = types.StringPointerValue(group.SourceAppId)
	}
	if len(preApprovedGroups) > 0 {
		terraformModel.PreapprovedGroups, _ = types.ListValueFrom(ctx, resource_pre_approval_rule.PreapprovedGroupsValue{}.Type(ctx), preApprovedGroups)
	} else {
		terraformModel.PreapprovedGroups, _ = types.ListValueFrom(ctx, resource_pre_approval_rule.PreapprovedGroupsValue{}.Type(ctx), make([]resource_pre_approval_rule.PreapprovedGroupsValue, 0))
	}

	var preApprovedPermissions []resource_pre_approval_rule.PreapprovedPermissionsValue
	terraformModel.PreapprovedPermissions.ElementsAs(ctx, &preApprovedPermissions, false)
	for i, permission := range *apiResponse.PreapprovedPermissions {
		preApprovedPermissions[i].Id = types.StringPointerValue(permission.Id)
		preApprovedPermissions[i].Label = types.StringValue(permission.Label)
		preApprovedPermissions[i].PreapprovedPermissionsType = types.StringPointerValue((*string)(permission.Type))
		preApprovedPermissions[i].AppId = types.StringValue(permission.AppId)
		preApprovedPermissions[i].AppClassId = types.StringValue(permission.AppClassId)
		preApprovedPermissions[i].AppInstanceId = types.StringValue(permission.AppInstanceId)
	}
	if len(preApprovedPermissions) > 0 {
		terraformModel.PreapprovedPermissions, _ = types.ListValueFrom(ctx, resource_pre_approval_rule.PreapprovedPermissionsValue{}.Type(ctx), preApprovedPermissions)
	} else {
		terraformModel.PreapprovedPermissions, _ = types.ListValueFrom(ctx, resource_pre_approval_rule.PreapprovedPermissionsValue{}.Type(ctx), make([]resource_pre_approval_rule.PreapprovedPermissionsValue, 0))

	}

	var preApprovalWebhooks []resource_pre_approval_rule.PreapprovalWebhooksValue
	terraformModel.PreapprovalWebhooks.ElementsAs(ctx, &preApprovalWebhooks, false)
	for i, webhook := range *apiResponse.PreapprovalWebhooks {
		preApprovalWebhooks[i].Description = types.StringPointerValue(webhook.Description)
		preApprovalWebhooks[i].HookType = types.StringValue(string(webhook.HookType))
		preApprovalWebhooks[i].Id = types.StringValue(string(webhook.Id))
		preApprovalWebhooks[i].Name = types.StringValue(string(webhook.Name))
	}
	if len(preApprovalWebhooks) > 0 {
		terraformModel.PreapprovalWebhooks, _ = types.ListValueFrom(ctx, resource_pre_approval_rule.PreapprovalWebhooksValue{}.Type(ctx), preApprovalWebhooks)
	} else {
		terraformModel.PreapprovalWebhooks, _ = types.ListValueFrom(ctx, resource_pre_approval_rule.PreapprovalWebhooksValue{}.Type(ctx), make([]resource_pre_approval_rule.PreapprovalWebhooksValue, 0))
	}

	return terraformModel
}
