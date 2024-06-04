// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/utils"
)

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

// RequestablePermissionInputAllowedGroups - The allowed groups associated with this config.
type RequestablePermissionInputAllowedGroups struct {
	// The type of this allowed groups config, can be all groups or specific.
	Type *RequestablePermissionInputAllowedGroupsConfigType `default:"ALL_GROUPS" json:"type"`
	// The groups associated with this config.
	Groups []BaseGroup `json:"groups,omitempty"`
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

func (o *RequestablePermissionInputAllowedGroups) GetType() *RequestablePermissionInputAllowedGroupsConfigType {
	if o == nil {
		return nil
	}
	return o.Type
}

func (o *RequestablePermissionInputAllowedGroups) GetGroups() []BaseGroup {
	if o == nil {
		return nil
	}
	return o.Groups
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
	// Indicates if approval flow is overrided
	RequestApprovalConfigOverride *bool `json:"request_approval_config_override,omitempty"`
	// Manager approval can be configured as necessary to continue
	ManagerApproval *RequestablePermissionInputManagerApprovalOption `default:"NONE" json:"manager_approval"`
	// Only turn on when working with sensitive permissions to ensure a smooth employee experience.
	RequireAdditionalApproval *bool `json:"require_additional_approval,omitempty"`
	// During the approval step, send a custom message to requesters. Note that the Permission level approval message will override the App level approval message. Markdown for links and text formatting is supported.
	CustomApprovalMessage *string `json:"custom_approval_message,omitempty"`
	// Indicates if custom_approval_message is overrided
	CustomApprovalMessageOverride *bool `json:"custom_approval_message_override,omitempty"`
	// The stages of this request approval.
	RequestApprovalStages []RequestApprovalStageInput `json:"request_approval_stages,omitempty"`
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

func (o *RequestablePermissionInputRequestApprovalConfig) GetRequestApprovalConfigOverride() *bool {
	if o == nil {
		return nil
	}
	return o.RequestApprovalConfigOverride
}

func (o *RequestablePermissionInputRequestApprovalConfig) GetManagerApproval() *RequestablePermissionInputManagerApprovalOption {
	if o == nil {
		return nil
	}
	return o.ManagerApproval
}

func (o *RequestablePermissionInputRequestApprovalConfig) GetRequireAdditionalApproval() *bool {
	if o == nil {
		return nil
	}
	return o.RequireAdditionalApproval
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

func (o *RequestablePermissionInputRequestApprovalConfig) GetRequestApprovalStages() []RequestApprovalStageInput {
	if o == nil {
		return nil
	}
	return o.RequestApprovalStages
}

// RequestablePermissionInputAccessRemovalInlineWebhook - An inactivity workflow can be optionally associated with this config.
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

// RequestablePermissionInputProvisioningGroup - The provisioning group optionally associated with this config.
type RequestablePermissionInputProvisioningGroup struct {
	// The ID of this group.
	ID *string `json:"id,omitempty"`
	// The ID of the app that owns this group.
	AppID *string `json:"app_id,omitempty"`
	// The ID of this group, specific to the integration.
	IntegrationSpecificID *string `json:"integration_specific_id,omitempty"`
}

func (o *RequestablePermissionInputProvisioningGroup) GetID() *string {
	if o == nil {
		return nil
	}
	return o.ID
}

func (o *RequestablePermissionInputProvisioningGroup) GetAppID() *string {
	if o == nil {
		return nil
	}
	return o.AppID
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
	// Whether manual steps are needed.
	ManualStepsNeeded *bool `json:"manual_steps_needed,omitempty"`
	// The manual instructions that go along.
	ManualInstructions *string `json:"manual_instructions,omitempty"`
	// If enabled, users can request an app for a selected duration. After expiry, Lumos will automatically remove user's access.
	TimeBasedAccess []TimeBasedAccessOptions `json:"time_based_access,omitempty"`
	// Indicates if time based access is overrided
	TimeBasedAccessOverride *bool `json:"time_based_access_override,omitempty"`
	// The provisioning group optionally associated with this config.
	ProvisioningGroup *RequestablePermissionInputProvisioningGroup `json:"provisioning_group,omitempty"`
	// The provisioning webhook optionally associated with this config.
	ProvisioningWebhook *RequestablePermissionInputProvisioningWebhook `json:"provisioning_webhook,omitempty"`
}

func (o *RequestablePermissionInputRequestFulfillmentConfig) GetManualStepsNeeded() *bool {
	if o == nil {
		return nil
	}
	return o.ManualStepsNeeded
}

func (o *RequestablePermissionInputRequestFulfillmentConfig) GetManualInstructions() *string {
	if o == nil {
		return nil
	}
	return o.ManualInstructions
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

// RequestablePermissionInputRequestConfig - The request config associated with this requestable permission.
type RequestablePermissionInputRequestConfig struct {
	// The appstore visibility of this request config.
	AppstoreVisibility *RequestablePermissionInputAppStoreVisibilityOption `default:"HIDDEN" json:"appstore_visibility"`
	// Indicates if allowed groups is overrided
	AllowedGroupsOverride *bool `json:"allowed_groups_override,omitempty"`
	// The allowed groups associated with this config.
	AllowedGroups *RequestablePermissionInputAllowedGroups `json:"allowed_groups,omitempty"`
	// A request approval config can be optionally associated with this config
	RequestApprovalConfig *RequestablePermissionInputRequestApprovalConfig `json:"request_approval_config,omitempty"`
	// An inactivity workflow can be optionally associated with this config.
	AccessRemovalInlineWebhook *RequestablePermissionInputAccessRemovalInlineWebhook `json:"access_removal_inline_webhook,omitempty"`
	// A request validation webhook can be optionally associated with this config.
	RequestValidationInlineWebhook *RequestablePermissionInputRequestValidationInlineWebhook `json:"request_validation_inline_webhook,omitempty"`
	// A request fulfillment config can be optionally associated with this config
	RequestFulfillmentConfig *RequestablePermissionInputRequestFulfillmentConfig `json:"request_fulfillment_config,omitempty"`
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

func (o *RequestablePermissionInputRequestConfig) GetAppstoreVisibility() *RequestablePermissionInputAppStoreVisibilityOption {
	if o == nil {
		return nil
	}
	return o.AppstoreVisibility
}

func (o *RequestablePermissionInputRequestConfig) GetAllowedGroupsOverride() *bool {
	if o == nil {
		return nil
	}
	return o.AllowedGroupsOverride
}

func (o *RequestablePermissionInputRequestConfig) GetAllowedGroups() *RequestablePermissionInputAllowedGroups {
	if o == nil {
		return nil
	}
	return o.AllowedGroups
}

func (o *RequestablePermissionInputRequestConfig) GetRequestApprovalConfig() *RequestablePermissionInputRequestApprovalConfig {
	if o == nil {
		return nil
	}
	return o.RequestApprovalConfig
}

func (o *RequestablePermissionInputRequestConfig) GetAccessRemovalInlineWebhook() *RequestablePermissionInputAccessRemovalInlineWebhook {
	if o == nil {
		return nil
	}
	return o.AccessRemovalInlineWebhook
}

func (o *RequestablePermissionInputRequestConfig) GetRequestValidationInlineWebhook() *RequestablePermissionInputRequestValidationInlineWebhook {
	if o == nil {
		return nil
	}
	return o.RequestValidationInlineWebhook
}

func (o *RequestablePermissionInputRequestConfig) GetRequestFulfillmentConfig() *RequestablePermissionInputRequestFulfillmentConfig {
	if o == nil {
		return nil
	}
	return o.RequestFulfillmentConfig
}

type RequestablePermissionInput struct {
	// The ID of the app associated with this requestable permission.
	AppID *string `json:"app_id,omitempty"`
	// The ID of the service associated with this requestable permission.
	AppClassID *string `json:"app_class_id,omitempty"`
	// Optionally, an app has an identifer associated with it's particular instance.
	AppInstanceID *string `json:"app_instance_id,omitempty"`
	// The label of this requestable permission.
	Label string `json:"label"`
	// The request config associated with this requestable permission.
	RequestConfig *RequestablePermissionInputRequestConfig `json:"request_config,omitempty"`
}

func (o *RequestablePermissionInput) GetAppID() *string {
	if o == nil {
		return nil
	}
	return o.AppID
}

func (o *RequestablePermissionInput) GetAppClassID() *string {
	if o == nil {
		return nil
	}
	return o.AppClassID
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