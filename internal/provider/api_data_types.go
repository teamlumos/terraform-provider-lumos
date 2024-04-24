package provider

import (
	"strings"

	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type lumosAPIResponse interface {
}

type lumosAPIAppStoreRequestablePermissionResponse struct {
	Id            string                                     `json:"id"`
	Type          string                                     `json:"type"`
	Label         string                                     `json:"label"`
	AppId         string                                     `json:"app_id"`
	AppClassId    string                                     `json:"app_class_id"`
	AppInstanceId string                                     `json:"app_instance_id"`
	RequestConfig appStoreRequestablePermissionRequestConfig `json:"request_config"`
}

var appStoreVisibilityValueMap = make(map[bool]string)

func init() {
	appStoreVisibilityValueMap[true] = "VISIBLE"
	appStoreVisibilityValueMap[false] = "HIDDEN"
}

var appStoreVisibilityKeyMap = make(map[string]bool)

func init() {
	appStoreVisibilityKeyMap["VISIBLE"] = true
	appStoreVisibilityKeyMap["HIDDEN"] = false
}

var allowedGroupsValueMap = make(map[bool]string)

func init() {
	allowedGroupsValueMap[true] = "ALL_GROUPS"
	allowedGroupsValueMap[false] = "SPECIFIED_GROUPS"
}

var managerApprovalValueMap = make(map[bool]string)

func init() {
	managerApprovalValueMap[true] = "INITIAL_APPROVAL"
	managerApprovalValueMap[false] = "NONE"
}

var managerApprovalKeyMap = make(map[string]bool)

func init() {
	managerApprovalKeyMap["INITIAL_APPROVAL"] = true
	managerApprovalKeyMap["NONE"] = false
}

type appStoreRequestablePermissionRequestConfig struct {
	AppStoreVisibility             string                                                `json:"appstore_visibility"`
	AllowedGroups                  appStoreRequestablePermissionConfigAllowedGroups      `json:"allowed_groups"`
	RequestApprovalConfig          appStoreRequestablePermissionRequestApprovalConfig    `json:"request_approval_config"`
	RequestFulfillmentConfig       appStoreRequestablePermissionRequestFulfillmentConfig `json:"request_fulfillment_config"`
	AccessRemovalInlineWebhook     lumosAPIInlineWebhook                                 `json:"access_removal_inline_webhook"`
	RequestValidationInlineWebhook lumosAPIInlineWebhook                                 `json:"request_validation_inline_webhook"`
}

type appStoreRequestablePermissionConfigAllowedGroups struct {
	Type   string          `json:"type"`
	Groups []lumosAPIGroup `json:"groups"`
}

type appStoreRequestablePermissionRequestApprovalConfig struct {
	ManagerApproval           string                                              `json:"manager_approval"`
	RequireAdditionalApproval bool                                                `json:"require_additional_approval"`
	RequestApprovalStages     []appStoreRequestablePermissionManagerApprovalStage `json:"request_approval_stages"`
}

type appStoreRequestablePermissionManagerApprovalStage struct {
	Approvers []appStoreRequestablePermissionApprovers `json:"approvers"`
}

type appStoreRequestablePermissionApprovers struct {
	Type  string        `json:"type"`
	User  lumosAPIUser  `json:"user"`
	Group lumosAPIGroup `json:"group"`
}

type appStoreRequestablePermissionRequestFulfillmentConfig struct {
	ManualStepsNeeded   bool                  `json:"manual_steps_needed"`
	ManualInstructions  string                `json:"manual_instructions"`
	TimeBasedAccess     []string              `json:"time_based_access"`
	ProvisioningGroup   lumosAPIGroup         `json:"provisioning_group"`
	ProvisioningWebhook lumosAPIInlineWebhook `json:"provisioning_webhook"`
}

type lumosAPIInlineWebhook struct {
	Id string `json:"id"`
}

type lumosAPIGroup struct {
	Id                    string `json:"id"`
	AppId                 string `json:"app_id"`
	IntegrationSpecificId string `json:"integration_specific_id"`
}

type lumosAPIUser struct {
	Id string `json:"id"`
}

func buildAppStoreAppRequestablePermissionPayload(p requestablePermissionResourceModel) map[string]interface{} {
	request_config := buildAppStoreRequestablePermissionRequestConfig(p)

	// Build the payload from calculated inputs
	payload := map[string]interface{}{
		"label":  p.Label.ValueString(),
		"app_id": p.AppId.ValueString(),
	}

	if request_config != nil {
		payload["request_config"] = request_config
	}

	return payload
}

func parseAppStoreRequestablePermissionUserApprovers(user_approvers basetypes.SetValue, approvers_list []map[string]interface{}) []map[string]interface{} {
	for _, id := range user_approvers.Elements() {
		trimmedID := strings.Trim(id.String(), "\"")
		approver_user := map[string]interface{}{
			"type": "USER",
			"user": map[string]interface{}{
				"id": trimmedID,
			},
		}
		approvers_list = append(approvers_list, approver_user)
	}
	return approvers_list
}

func parseAppStoreRequestablePermissionGroupApprovers(group_approvers basetypes.SetValue, approvers_list []map[string]interface{}) []map[string]interface{} {
	for _, id := range group_approvers.Elements() {
		trimmedID := strings.Trim(id.String(), "\"")
		approver_group := map[string]interface{}{
			"type": "GROUP",
			"group": map[string]interface{}{
				"id": trimmedID,
			},
		}
		approvers_list = append(approvers_list, approver_group)
	}
	return approvers_list
}

func buildAppStoreRequestablePermissionAllowedGroups(p requestablePermissionResourceModel) map[string]interface{} {
	// Set Request Config -> Allowed Groups
	if !p.AllowedGroups.IsNull() {
		allowed_groups_type_bool := len(p.AllowedGroups.Elements()) == 0

		var groups []map[string]string

		for _, id := range p.AllowedGroups.Elements() {
			trimmedID := strings.Trim(id.String(), "\"")
			dict := map[string]string{"id": trimmedID}
			groups = append(groups, dict)
		}
		return map[string]interface{}{
			"type":   allowedGroupsValueMap[allowed_groups_type_bool],
			"groups": groups,
		}
	}

	return nil
}

func buildAppStoreRequestablePermissionApprovalConfig(p requestablePermissionResourceModel) map[string]interface{} {
	request_approval_config := map[string]interface{}{}

	// Set Request Config -> Approval Config -> Manager Approval
	if !p.ManagerApprovalRequired.IsNull() {
		key := p.ManagerApprovalRequired.ValueBool()
		request_approval_config["manager_approval"] = managerApprovalValueMap[key]
	}

	request_approval_config["require_additional_approval"] = p.RequireAdditionalApproval.ValueBool()

	// Set Request Config -> Approval Config -> Request Approval Stages
	if !p.ApproverGroupsStage1.IsNull() || !p.ApproverUsersStage1.IsNull() {
		request_approval_stages := []map[string]interface{}{}
		list_stage_1_approvers := []map[string]interface{}{}

		if !p.ApproverGroupsStage1.IsNull() {
			list_stage_1_approvers = parseAppStoreRequestablePermissionGroupApprovers(p.ApproverGroupsStage1, list_stage_1_approvers)
		}
		if !p.ApproverUsersStage1.IsNull() {
			list_stage_1_approvers = parseAppStoreRequestablePermissionUserApprovers(p.ApproverUsersStage1, list_stage_1_approvers)
		}

		approvers_stage_1 := map[string]interface{}{
			"approvers": list_stage_1_approvers,
		}
		request_approval_stages = append(request_approval_stages, approvers_stage_1)

		// approvers stage 2 is only valid if stage 1 is defined
		list_stage_2_approvers := []map[string]interface{}{}
		if !p.ApproverGroupsStage2.IsNull() {
			list_stage_2_approvers = parseAppStoreRequestablePermissionGroupApprovers(p.ApproverGroupsStage2, list_stage_2_approvers)
		}
		if !p.ApproverUsersStage2.IsNull() {
			list_stage_2_approvers = parseAppStoreRequestablePermissionUserApprovers(p.ApproverUsersStage2, list_stage_2_approvers)
		}

		if len(list_stage_2_approvers) != 0 {
			approvers_stage_2 := map[string]interface{}{
				"approvers": list_stage_2_approvers,
			}
			request_approval_stages = append(request_approval_stages, approvers_stage_2)
		}

		// Set approvers in payload
		request_approval_config["request_approval_stages"] = request_approval_stages

	}

	if len(request_approval_config) != 0 {
		return request_approval_config
	}

	return nil
}

func buildAppStoreAppRequestablePermissionFulfillmentConfig(p requestablePermissionResourceModel) map[string]interface{} {
	// Set Request Config -> Fulfullment Config
	request_fulfillment_config := map[string]interface{}{}

	// Set Request Config -> Fulfullment Config -> Manual Steps Needed
	if !p.ManualStepsNeeded.IsNull() {
		request_fulfillment_config["manual_steps_needed"] = p.ManualStepsNeeded.ValueBool()
	}

	// Set Request Config -> Fulfullment Config -> Manual Instrusctions
	if !p.ManualInstructions.IsNull() {
		request_fulfillment_config["manual_instructions"] = p.ManualInstructions.ValueString()
	}

	// Set Request Config -> Fulfullment Config -> Time Based Access
	if !p.TimeBasedAccessOptions.IsNull() {
		var options []string

		for _, option := range p.TimeBasedAccessOptions.Elements() {
			trimmedOption := strings.Trim(option.String(), "\"")
			options = append(options, trimmedOption)
		}
		request_fulfillment_config["time_based_access"] = options
	}

	// Set Request Config -> Fulfullment Config -> Provisioning Group
	if !p.ProvisioningGroup.Id.IsNull() {
		request_fulfillment_config["provisioning_group"] = map[string]string{
			"id": p.ProvisioningGroup.Id.ValueString(),
		}
	}
	// Set Request Config -> Fulfullment Config -> Provisioning Webhook
	if !p.ProvisioningInlineWebhook.IsNull() {
		request_fulfillment_config["provisioning_webhook"] = map[string]string{
			"id": p.ProvisioningInlineWebhook.ValueString(),
		}
	}

	if len(request_fulfillment_config) != 0 {
		return request_fulfillment_config
	}

	return nil

}

func buildAppStoreRequestablePermissionRequestConfig(p requestablePermissionResourceModel) map[string]interface{} {
	request_config := map[string]interface{}{}

	if !p.VisibleInAppStore.IsNull() {
		key := p.VisibleInAppStore.ValueBool()
		request_config["appstore_visibility"] = appStoreVisibilityValueMap[key]
	}

	allowed_groups := buildAppStoreRequestablePermissionAllowedGroups(p)
	if allowed_groups != nil {
		request_config["allowed_groups"] = allowed_groups
	}

	request_approval_config := buildAppStoreRequestablePermissionApprovalConfig(p)
	if request_approval_config != nil {
		request_config["request_approval_config"] = request_approval_config
	}

	request_fulfillment_config := buildAppStoreAppRequestablePermissionFulfillmentConfig(p)
	if request_fulfillment_config != nil {
		request_config["request_fulfillment_config"] = request_fulfillment_config
	}

	if !p.AccessRemovalInlineWebhook.IsNull() {
		request_config["access_removal_inline_webhook"] = map[string]string{
			"id": p.AccessRemovalInlineWebhook.ValueString(),
		}
	}

	if !p.RequestValidationInlineWebhook.IsNull() {
		request_config["request_validation_inline_webhook"] = map[string]string{
			"id": p.RequestValidationInlineWebhook.ValueString(),
		}
	}

	if len(request_config) != 0 {
		return request_config
	}

	return nil
}

func setRequestablePermissionResourceModelFromLumosAPIRequestablePermission(ctx context.Context, permissionResource requestablePermissionResourceModel, lumosApiPermission lumosAPIAppStoreRequestablePermissionResponse) requestablePermissionResourceModel {
	permissionResource.Id = types.StringValue(lumosApiPermission.Id)
	permissionResource.AppId = types.StringValue(lumosApiPermission.AppId)
	permissionResource.Label = types.StringValue(lumosApiPermission.Label)
	permissionResource.PermissionType = types.StringValue(lumosApiPermission.Type)
	permissionResource.VisibleInAppStore = types.BoolValue(appStoreVisibilityKeyMap[lumosApiPermission.RequestConfig.AppStoreVisibility])

	var group_ids []string

	for _, group := range lumosApiPermission.RequestConfig.AllowedGroups.Groups {
		group_ids = append(group_ids, group.Id)
	}
	if len(group_ids) > 0 {
		permissionResource.AllowedGroups, _ = types.SetValueFrom(ctx, types.StringType, group_ids)
	}

	permissionResource.ManagerApprovalRequired = types.BoolValue(managerApprovalKeyMap[lumosApiPermission.RequestConfig.RequestApprovalConfig.ManagerApproval])
	permissionResource.RequireAdditionalApproval = types.BoolValue(lumosApiPermission.RequestConfig.RequestApprovalConfig.RequireAdditionalApproval)

	var approvers []appStoreRequestablePermissionManagerApprovalStage = lumosApiPermission.RequestConfig.RequestApprovalConfig.RequestApprovalStages

	var approvers_stage_1 []appStoreRequestablePermissionApprovers
	var approvers_stage_2 []appStoreRequestablePermissionApprovers

	if len(approvers) > 0 {
		approvers_stage_1 = approvers[0].Approvers
	}
	if len(approvers) > 1 {
		approvers_stage_2 = approvers[1].Approvers
	}

	if len(approvers_stage_1) > 0 {
		var user_ids_approvers_stage_1 []string
		var group_ids_approvers_stage_1 []string

		for _, approver := range approvers_stage_1 {
			if approver.Type == "USER" {
				user_ids_approvers_stage_1 = append(user_ids_approvers_stage_1, approver.User.Id)
			}
			if approver.Type == "GROUP" {
				group_ids_approvers_stage_1 = append(group_ids_approvers_stage_1, approver.Group.Id)
			}
		}

		if len(user_ids_approvers_stage_1) > 0 {
			permissionResource.ApproverUsersStage1, _ = types.SetValueFrom(ctx, types.StringType, user_ids_approvers_stage_1)
		}

		if len(group_ids_approvers_stage_1) > 0 {
			permissionResource.ApproverGroupsStage1, _ = types.SetValueFrom(ctx, types.StringType, group_ids_approvers_stage_1)
		}

	}

	if len(approvers_stage_2) > 0 {
		var user_ids_approvers_stage_2 []string
		var group_ids_approvers_stage_2 []string
		for _, approver := range approvers_stage_2 {
			if approver.Type == "USER" {
				user_ids_approvers_stage_2 = append(user_ids_approvers_stage_2, approver.User.Id)
			}
			if approver.Type == "GROUP" {
				group_ids_approvers_stage_2 = append(group_ids_approvers_stage_2, approver.Group.Id)
			}
		}

		if len(user_ids_approvers_stage_2) > 0 {
			permissionResource.ApproverUsersStage2, _ = types.SetValueFrom(ctx, types.StringType, user_ids_approvers_stage_2)
		}

		if len(group_ids_approvers_stage_2) > 0 {
			permissionResource.ApproverGroupsStage2, _ = types.SetValueFrom(ctx, types.StringType, group_ids_approvers_stage_2)
		}

	}

	if lumosApiPermission.RequestConfig.AccessRemovalInlineWebhook.Id != "" {
		permissionResource.AccessRemovalInlineWebhook = types.StringValue(lumosApiPermission.RequestConfig.AccessRemovalInlineWebhook.Id)
	}

	if lumosApiPermission.RequestConfig.RequestValidationInlineWebhook.Id != "" {
		permissionResource.RequestValidationInlineWebhook = types.StringValue(lumosApiPermission.RequestConfig.RequestValidationInlineWebhook.Id)
	}

	if lumosApiPermission.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.Id != "" {
		permissionResource.ProvisioningGroup.Id = types.StringValue(lumosApiPermission.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.Id)
	} else if lumosApiPermission.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.AppId != "" && lumosApiPermission.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.IntegrationSpecificId != "" {
		permissionResource.ProvisioningGroup.AppId = types.StringValue(lumosApiPermission.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.AppId)
		permissionResource.ProvisioningGroup.IntegrationSpecificId = types.StringValue(lumosApiPermission.RequestConfig.RequestFulfillmentConfig.ProvisioningGroup.IntegrationSpecificId)
	}

	if lumosApiPermission.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook.Id != "" {
		permissionResource.ProvisioningInlineWebhook = types.StringValue(lumosApiPermission.RequestConfig.RequestFulfillmentConfig.ProvisioningWebhook.Id)
	}

	permissionResource.ManualStepsNeeded = types.BoolValue(lumosApiPermission.RequestConfig.RequestFulfillmentConfig.ManualStepsNeeded)

	if lumosApiPermission.RequestConfig.RequestFulfillmentConfig.ManualInstructions != "" {
		permissionResource.ManualInstructions = types.StringValue(lumosApiPermission.RequestConfig.RequestFulfillmentConfig.ManualInstructions)
	}
	permissionResource.TimeBasedAccessOptions, _ = types.SetValueFrom(ctx, types.StringType, lumosApiPermission.RequestConfig.RequestFulfillmentConfig.TimeBasedAccess)

	return permissionResource
}
