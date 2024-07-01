// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/utils"
)

// RequestablePermissionInputAccessRemovalInlineWebhook - A deprovisioning webhook can be optionally associated with this config.
type RequestablePermissionInputAccessRemovalInlineWebhook struct {
	// The ID of this inline webhook.
	ID string `json:"id"`
}

func (o *RequestablePermissionInputAccessRemovalInlineWebhook) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

// RequestablePermissionInputAllowedGroupsConfigType - The type of this allowed groups config, can be all groups or specific.
type RequestablePermissionInputAllowedGroupsConfigType string

const (
	RequestablePermissionInputAllowedGroupsConfigTypeAllGroups       RequestablePermissionInputAllowedGroupsConfigType = "ALL_GROUPS"
	RequestablePermissionInputAllowedGroupsConfigTypeSpecifiedGroups RequestablePermissionInputAllowedGroupsConfigType = "SPECIFIED_GROUPS"
)

func (e RequestablePermissionInputAllowedGroupsConfigType) ToPointer() *RequestablePermissionInputAllowedGroupsConfigType {
	return &e
}
func (e *RequestablePermissionInputAllowedGroupsConfigType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "ALL_GROUPS":
		fallthrough
	case "SPECIFIED_GROUPS":
		*e = RequestablePermissionInputAllowedGroupsConfigType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for RequestablePermissionInputAllowedGroupsConfigType: %v", v)
	}
}

// RequestablePermissionInputAllowedGroups - Refers to which group(s) can make requests to this permission.
type RequestablePermissionInputAllowedGroups struct {
	// The groups allowed to request this permission.
	Groups []BaseGroup `json:"groups,omitempty"`
	// The type of this allowed groups config, can be all groups or specific.
	Type *RequestablePermissionInputAllowedGroupsConfigType `default:"ALL_GROUPS" json:"type"`
}

func (r RequestablePermissionInputAllowedGroups) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(r, "", false)
}

func (r *RequestablePermissionInputAllowedGroups) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &r, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *RequestablePermissionInputAllowedGroups) GetGroups() []BaseGroup {
	if o == nil {
		return nil
	}
	return o.Groups
}

func (o *RequestablePermissionInputAllowedGroups) GetType() *RequestablePermissionInputAllowedGroupsConfigType {
	if o == nil {
		return nil
	}
	return o.Type
}

// RequestablePermissionInputAppStoreVisibilityOption - The appstore visibility of this request config.
type RequestablePermissionInputAppStoreVisibilityOption string

const (
	RequestablePermissionInputAppStoreVisibilityOptionHidden  RequestablePermissionInputAppStoreVisibilityOption = "HIDDEN"
	RequestablePermissionInputAppStoreVisibilityOptionVisible RequestablePermissionInputAppStoreVisibilityOption = "VISIBLE"
)

func (e RequestablePermissionInputAppStoreVisibilityOption) ToPointer() *RequestablePermissionInputAppStoreVisibilityOption {
	return &e
}
func (e *RequestablePermissionInputAppStoreVisibilityOption) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "HIDDEN":
		fallthrough
	case "VISIBLE":
		*e = RequestablePermissionInputAppStoreVisibilityOption(v)
		return nil
	default:
		return fmt.Errorf("invalid value for RequestablePermissionInputAppStoreVisibilityOption: %v", v)
	}
}

// RequestablePermissionInputApprovers - AppStore App approvers assigned.
type RequestablePermissionInputApprovers struct {
	// Groups assigned as support request approvers.
	Groups []BaseGroup `json:"groups,omitempty"`
	// Users assigned as support request approvers.
	Users []BaseUser `json:"users,omitempty"`
}

func (o *RequestablePermissionInputApprovers) GetGroups() []BaseGroup {
	if o == nil {
		return nil
	}
	return o.Groups
}

func (o *RequestablePermissionInputApprovers) GetUsers() []BaseUser {
	if o == nil {
		return nil
	}
	return o.Users
}

// RequestablePermissionInputApproversStage2 - AppStore App stage 2 approvers assigned.
type RequestablePermissionInputApproversStage2 struct {
	// Groups assigned as support request approvers.
	Groups []BaseGroup `json:"groups,omitempty"`
	// Users assigned as support request approvers.
	Users []BaseUser `json:"users,omitempty"`
}

func (o *RequestablePermissionInputApproversStage2) GetGroups() []BaseGroup {
	if o == nil {
		return nil
	}
	return o.Groups
}

func (o *RequestablePermissionInputApproversStage2) GetUsers() []BaseUser {
	if o == nil {
		return nil
	}
	return o.Users
}

// RequestablePermissionInputManagerApprovalOption - Manager approval can be configured as necessary to continue
type RequestablePermissionInputManagerApprovalOption string

const (
	RequestablePermissionInputManagerApprovalOptionNone            RequestablePermissionInputManagerApprovalOption = "NONE"
	RequestablePermissionInputManagerApprovalOptionInitialApproval RequestablePermissionInputManagerApprovalOption = "INITIAL_APPROVAL"
)

func (e RequestablePermissionInputManagerApprovalOption) ToPointer() *RequestablePermissionInputManagerApprovalOption {
	return &e
}
func (e *RequestablePermissionInputManagerApprovalOption) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "NONE":
		fallthrough
	case "INITIAL_APPROVAL":
		*e = RequestablePermissionInputManagerApprovalOption(v)
		return nil
	default:
		return fmt.Errorf("invalid value for RequestablePermissionInputManagerApprovalOption: %v", v)
	}
}

// RequestablePermissionInputRequestApprovalConfig - A request approval config can be optionally associated with this config
type RequestablePermissionInputRequestApprovalConfig struct {
	// AppStore App approvers assigned.
	Approvers *RequestablePermissionInputApprovers `json:"approvers,omitempty"`
	// AppStore App stage 2 approvers assigned.
	ApproversStage2 *RequestablePermissionInputApproversStage2 `json:"approvers_stage_2,omitempty"`
	// After the approval step, send a custom message to requesters. Note that the permission level approval message will override the App level approval message if custom_approval_message_override is set. Markdown for links and text formatting is supported.
	CustomApprovalMessage *string `json:"custom_approval_message,omitempty"`
	// Indicates if custom_approval_message is overridden.
	CustomApprovalMessageOverride *bool `json:"custom_approval_message_override,omitempty"`
	// Manager approval can be configured as necessary to continue
	ManagerApproval *RequestablePermissionInputManagerApprovalOption `default:"NONE" json:"manager_approval"`
	// Indicates if approval flow is overridden.
	RequestApprovalConfigOverride *bool `json:"request_approval_config_override,omitempty"`
	// The stages of this request approval.
	RequestApprovalStages []RequestApprovalStageInput `json:"request_approval_stages,omitempty"`
	// Only turn on when working with sensitive permissions to ensure a smooth employee experience.
	RequireAdditionalApproval *bool `json:"require_additional_approval,omitempty"`
}

func (r RequestablePermissionInputRequestApprovalConfig) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(r, "", false)
}

func (r *RequestablePermissionInputRequestApprovalConfig) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &r, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *RequestablePermissionInputRequestApprovalConfig) GetApprovers() *RequestablePermissionInputApprovers {
	if o == nil {
		return nil
	}
	return o.Approvers
}

func (o *RequestablePermissionInputRequestApprovalConfig) GetApproversStage2() *RequestablePermissionInputApproversStage2 {
	if o == nil {
		return nil
	}
	return o.ApproversStage2
}

func (o *RequestablePermissionInputRequestApprovalConfig) GetCustomApprovalMessage() *string {
	if o == nil {
		return nil
	}
	return o.CustomApprovalMessage
}

func (o *RequestablePermissionInputRequestApprovalConfig) GetCustomApprovalMessageOverride() *bool {
	if o == nil {
		return nil
	}
	return o.CustomApprovalMessageOverride
}

func (o *RequestablePermissionInputRequestApprovalConfig) GetManagerApproval() *RequestablePermissionInputManagerApprovalOption {
	if o == nil {
		return nil
	}
	return o.ManagerApproval
}

func (o *RequestablePermissionInputRequestApprovalConfig) GetRequestApprovalConfigOverride() *bool {
	if o == nil {
		return nil
	}
	return o.RequestApprovalConfigOverride
}

func (o *RequestablePermissionInputRequestApprovalConfig) GetRequestApprovalStages() []RequestApprovalStageInput {
	if o == nil {
		return nil
	}
	return o.RequestApprovalStages
}

func (o *RequestablePermissionInputRequestApprovalConfig) GetRequireAdditionalApproval() *bool {
	if o == nil {
		return nil
	}
	return o.RequireAdditionalApproval
}

// RequestablePermissionInputProvisioningGroup - The provisioning group optionally associated with this config.
type RequestablePermissionInputProvisioningGroup struct {
	// The ID of the app that sources this group.
	AppID *string `json:"app_id,omitempty"`
	// The ID of this group.
	ID *string `json:"id,omitempty"`
	// The ID of this group, specific to the integration.
	IntegrationSpecificID *string `json:"integration_specific_id,omitempty"`
}

func (o *RequestablePermissionInputProvisioningGroup) GetAppID() *string {
	if o == nil {
		return nil
	}
	return o.AppID
}

func (o *RequestablePermissionInputProvisioningGroup) GetID() *string {
	if o == nil {
		return nil
	}
	return o.ID
}

func (o *RequestablePermissionInputProvisioningGroup) GetIntegrationSpecificID() *string {
	if o == nil {
		return nil
	}
	return o.IntegrationSpecificID
}

// RequestablePermissionInputProvisioningWebhook - The provisioning webhook optionally associated with this config.
type RequestablePermissionInputProvisioningWebhook struct {
	// The ID of this inline webhook.
	ID string `json:"id"`
}

func (o *RequestablePermissionInputProvisioningWebhook) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

// RequestablePermissionInputRequestFulfillmentConfig - A request fulfillment config can be optionally associated with this config
type RequestablePermissionInputRequestFulfillmentConfig struct {
	// The manual instructions that go along.
	ManualInstructions *string `json:"manual_instructions,omitempty"`
	// Whether manual steps are needed.
	ManualStepsNeeded *bool `json:"manual_steps_needed,omitempty"`
	// The provisioning group optionally associated with this config.
	ProvisioningGroup *RequestablePermissionInputProvisioningGroup `json:"provisioning_group,omitempty"`
	// The provisioning webhook optionally associated with this config.
	ProvisioningWebhook *RequestablePermissionInputProvisioningWebhook `json:"provisioning_webhook,omitempty"`
	// If enabled, users can request an app for a selected duration. After expiry, Lumos will automatically remove user's access.
	TimeBasedAccess []TimeBasedAccessOptions `json:"time_based_access,omitempty"`
	// Indicates if time based access is overriden.
	TimeBasedAccessOverride *bool `json:"time_based_access_override,omitempty"`
}

func (o *RequestablePermissionInputRequestFulfillmentConfig) GetManualInstructions() *string {
	if o == nil {
		return nil
	}
	return o.ManualInstructions
}

func (o *RequestablePermissionInputRequestFulfillmentConfig) GetManualStepsNeeded() *bool {
	if o == nil {
		return nil
	}
	return o.ManualStepsNeeded
}

func (o *RequestablePermissionInputRequestFulfillmentConfig) GetProvisioningGroup() *RequestablePermissionInputProvisioningGroup {
	if o == nil {
		return nil
	}
	return o.ProvisioningGroup
}

func (o *RequestablePermissionInputRequestFulfillmentConfig) GetProvisioningWebhook() *RequestablePermissionInputProvisioningWebhook {
	if o == nil {
		return nil
	}
	return o.ProvisioningWebhook
}

func (o *RequestablePermissionInputRequestFulfillmentConfig) GetTimeBasedAccess() []TimeBasedAccessOptions {
	if o == nil {
		return nil
	}
	return o.TimeBasedAccess
}

func (o *RequestablePermissionInputRequestFulfillmentConfig) GetTimeBasedAccessOverride() *bool {
	if o == nil {
		return nil
	}
	return o.TimeBasedAccessOverride
}

// RequestablePermissionInputRequestValidationInlineWebhook - A request validation webhook can be optionally associated with this config.
type RequestablePermissionInputRequestValidationInlineWebhook struct {
	// The ID of this inline webhook.
	ID string `json:"id"`
}

func (o *RequestablePermissionInputRequestValidationInlineWebhook) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

// RequestablePermissionInputRequestConfig - The request config associated with this requestable permission.
type RequestablePermissionInputRequestConfig struct {
	// A deprovisioning webhook can be optionally associated with this config.
	AccessRemovalInlineWebhook *RequestablePermissionInputAccessRemovalInlineWebhook `json:"access_removal_inline_webhook,omitempty"`
	// Refers to which group(s) can make requests to this permission.
	AllowedGroups *RequestablePermissionInputAllowedGroups `json:"allowed_groups,omitempty"`
	// Indicates if allowed groups is overriden from the app-level settings.
	AllowedGroupsOverride *bool `json:"allowed_groups_override,omitempty"`
	// The appstore visibility of this request config.
	AppstoreVisibility *RequestablePermissionInputAppStoreVisibilityOption `default:"HIDDEN" json:"appstore_visibility"`
	// A request approval config can be optionally associated with this config
	RequestApprovalConfig *RequestablePermissionInputRequestApprovalConfig `json:"request_approval_config,omitempty"`
	// A request fulfillment config can be optionally associated with this config
	RequestFulfillmentConfig *RequestablePermissionInputRequestFulfillmentConfig `json:"request_fulfillment_config,omitempty"`
	// A request validation webhook can be optionally associated with this config.
	RequestValidationInlineWebhook *RequestablePermissionInputRequestValidationInlineWebhook `json:"request_validation_inline_webhook,omitempty"`
}

func (r RequestablePermissionInputRequestConfig) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(r, "", false)
}

func (r *RequestablePermissionInputRequestConfig) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &r, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *RequestablePermissionInputRequestConfig) GetAccessRemovalInlineWebhook() *RequestablePermissionInputAccessRemovalInlineWebhook {
	if o == nil {
		return nil
	}
	return o.AccessRemovalInlineWebhook
}

func (o *RequestablePermissionInputRequestConfig) GetAllowedGroups() *RequestablePermissionInputAllowedGroups {
	if o == nil {
		return nil
	}
	return o.AllowedGroups
}

func (o *RequestablePermissionInputRequestConfig) GetAllowedGroupsOverride() *bool {
	if o == nil {
		return nil
	}
	return o.AllowedGroupsOverride
}

func (o *RequestablePermissionInputRequestConfig) GetAppstoreVisibility() *RequestablePermissionInputAppStoreVisibilityOption {
	if o == nil {
		return nil
	}
	return o.AppstoreVisibility
}

func (o *RequestablePermissionInputRequestConfig) GetRequestApprovalConfig() *RequestablePermissionInputRequestApprovalConfig {
	if o == nil {
		return nil
	}
	return o.RequestApprovalConfig
}

func (o *RequestablePermissionInputRequestConfig) GetRequestFulfillmentConfig() *RequestablePermissionInputRequestFulfillmentConfig {
	if o == nil {
		return nil
	}
	return o.RequestFulfillmentConfig
}

func (o *RequestablePermissionInputRequestConfig) GetRequestValidationInlineWebhook() *RequestablePermissionInputRequestValidationInlineWebhook {
	if o == nil {
		return nil
	}
	return o.RequestValidationInlineWebhook
}

type RequestablePermissionInput struct {
	// The ID of the service associated with this requestable permission.
	AppClassID *string `json:"app_class_id,omitempty"`
	// The ID of the app associated with this requestable permission.
	AppID string `json:"app_id"`
	// Optionally, an app has an identifer associated with it's particular instance.
	AppInstanceID *string `json:"app_instance_id,omitempty"`
	// The label of this requestable permission.
	Label string `json:"label"`
	// The request config associated with this requestable permission.
	RequestConfig *RequestablePermissionInputRequestConfig `json:"request_config,omitempty"`
}

func (o *RequestablePermissionInput) GetAppClassID() *string {
	if o == nil {
		return nil
	}
	return o.AppClassID
}

func (o *RequestablePermissionInput) GetAppID() string {
	if o == nil {
		return ""
	}
	return o.AppID
}

func (o *RequestablePermissionInput) GetAppInstanceID() *string {
	if o == nil {
		return nil
	}
	return o.AppInstanceID
}

func (o *RequestablePermissionInput) GetLabel() string {
	if o == nil {
		return ""
	}
	return o.Label
}

func (o *RequestablePermissionInput) GetRequestConfig() *RequestablePermissionInputRequestConfig {
	if o == nil {
		return nil
	}
	return o.RequestConfig
}
