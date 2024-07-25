// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/utils"
)

// RequestablePermissionOutputPermissionType - The type of this requestable permission.
type RequestablePermissionOutputPermissionType string

const (
	RequestablePermissionOutputPermissionTypeSynced RequestablePermissionOutputPermissionType = "SYNCED"
	RequestablePermissionOutputPermissionTypeNative RequestablePermissionOutputPermissionType = "NATIVE"
)

func (e RequestablePermissionOutputPermissionType) ToPointer() *RequestablePermissionOutputPermissionType {
	return &e
}
func (e *RequestablePermissionOutputPermissionType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "SYNCED":
		fallthrough
	case "NATIVE":
		*e = RequestablePermissionOutputPermissionType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for RequestablePermissionOutputPermissionType: %v", v)
	}
}

// AppStoreVisibilityOption - The appstore visibility of this request config.
type AppStoreVisibilityOption string

const (
	AppStoreVisibilityOptionHidden  AppStoreVisibilityOption = "HIDDEN"
	AppStoreVisibilityOptionVisible AppStoreVisibilityOption = "VISIBLE"
)

func (e AppStoreVisibilityOption) ToPointer() *AppStoreVisibilityOption {
	return &e
}
func (e *AppStoreVisibilityOption) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "HIDDEN":
		fallthrough
	case "VISIBLE":
		*e = AppStoreVisibilityOption(v)
		return nil
	default:
		return fmt.Errorf("invalid value for AppStoreVisibilityOption: %v", v)
	}
}

// RequestablePermissionOutputAllowedGroupsConfigType - The type of this allowed groups config, can be all groups or specific.
type RequestablePermissionOutputAllowedGroupsConfigType string

const (
	RequestablePermissionOutputAllowedGroupsConfigTypeAllGroups       RequestablePermissionOutputAllowedGroupsConfigType = "ALL_GROUPS"
	RequestablePermissionOutputAllowedGroupsConfigTypeSpecifiedGroups RequestablePermissionOutputAllowedGroupsConfigType = "SPECIFIED_GROUPS"
)

func (e RequestablePermissionOutputAllowedGroupsConfigType) ToPointer() *RequestablePermissionOutputAllowedGroupsConfigType {
	return &e
}
func (e *RequestablePermissionOutputAllowedGroupsConfigType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "ALL_GROUPS":
		fallthrough
	case "SPECIFIED_GROUPS":
		*e = RequestablePermissionOutputAllowedGroupsConfigType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for RequestablePermissionOutputAllowedGroupsConfigType: %v", v)
	}
}

// RequestablePermissionOutputAllowedGroups - The allowed groups config associated with this config.
type RequestablePermissionOutputAllowedGroups struct {
	// The type of this allowed groups config, can be all groups or specific.
	Type *RequestablePermissionOutputAllowedGroupsConfigType `default:"ALL_GROUPS" json:"type"`
	// The groups allowed to request this permission.
	Groups []Group `json:"groups,omitempty"`
}

func (r RequestablePermissionOutputAllowedGroups) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(r, "", false)
}

func (r *RequestablePermissionOutputAllowedGroups) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &r, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *RequestablePermissionOutputAllowedGroups) GetType() *RequestablePermissionOutputAllowedGroupsConfigType {
	if o == nil {
		return nil
	}
	return o.Type
}

func (o *RequestablePermissionOutputAllowedGroups) GetGroups() []Group {
	if o == nil {
		return nil
	}
	return o.Groups
}

// ManagerApprovalOption - Manager approval can be configured as necessary to continue
type ManagerApprovalOption string

const (
	ManagerApprovalOptionNone            ManagerApprovalOption = "NONE"
	ManagerApprovalOptionInitialApproval ManagerApprovalOption = "INITIAL_APPROVAL"
)

func (e ManagerApprovalOption) ToPointer() *ManagerApprovalOption {
	return &e
}
func (e *ManagerApprovalOption) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "NONE":
		fallthrough
	case "INITIAL_APPROVAL":
		*e = ManagerApprovalOption(v)
		return nil
	default:
		return fmt.Errorf("invalid value for ManagerApprovalOption: %v", v)
	}
}

// RequestablePermissionOutputApprovers - AppStore App approvers assigned.
type RequestablePermissionOutputApprovers struct {
	// Groups assigned as support request approvers.
	Groups []Group `json:"groups,omitempty"`
	// Users assigned as support request approvers.
	Users []User `json:"users,omitempty"`
}

func (o *RequestablePermissionOutputApprovers) GetGroups() []Group {
	if o == nil {
		return nil
	}
	return o.Groups
}

func (o *RequestablePermissionOutputApprovers) GetUsers() []User {
	if o == nil {
		return nil
	}
	return o.Users
}

// RequestablePermissionOutputApproversStage2 - AppStore App stage 2 approvers assigned.
type RequestablePermissionOutputApproversStage2 struct {
	// Groups assigned as support request approvers.
	Groups []Group `json:"groups,omitempty"`
	// Users assigned as support request approvers.
	Users []User `json:"users,omitempty"`
}

func (o *RequestablePermissionOutputApproversStage2) GetGroups() []Group {
	if o == nil {
		return nil
	}
	return o.Groups
}

func (o *RequestablePermissionOutputApproversStage2) GetUsers() []User {
	if o == nil {
		return nil
	}
	return o.Users
}

// RequestApprovalConfig - A request approval config can be optionally associated with this config
type RequestApprovalConfig struct {
	// Indicates if approval flow is overridden.
	RequestApprovalConfigOverride *bool `json:"request_approval_config_override,omitempty"`
	// Manager approval can be configured as necessary to continue
	ManagerApproval *ManagerApprovalOption `default:"NONE" json:"manager_approval"`
	// Only turn on when working with sensitive permissions to ensure a smooth employee experience.
	RequireAdditionalApproval *bool `json:"require_additional_approval,omitempty"`
	// After the approval step, send a custom message to requesters. Note that the permission level approval message will override the App level approval message if custom_approval_message_override is set. Markdown for links and text formatting is supported.
	CustomApprovalMessage *string `json:"custom_approval_message,omitempty"`
	// Indicates if custom_approval_message is overridden.
	CustomApprovalMessageOverride *bool `json:"custom_approval_message_override,omitempty"`
	// AppStore App approvers assigned.
	Approvers *RequestablePermissionOutputApprovers `json:"approvers,omitempty"`
	// AppStore App stage 2 approvers assigned.
	ApproversStage2 *RequestablePermissionOutputApproversStage2 `json:"approvers_stage_2,omitempty"`
	// The stages of this request approval.
	RequestApprovalStages []RequestApprovalStageOutput `json:"request_approval_stages,omitempty"`
}

func (r RequestApprovalConfig) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(r, "", false)
}

func (r *RequestApprovalConfig) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &r, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *RequestApprovalConfig) GetRequestApprovalConfigOverride() *bool {
	if o == nil {
		return nil
	}
	return o.RequestApprovalConfigOverride
}

func (o *RequestApprovalConfig) GetManagerApproval() *ManagerApprovalOption {
	if o == nil {
		return nil
	}
	return o.ManagerApproval
}

func (o *RequestApprovalConfig) GetRequireAdditionalApproval() *bool {
	if o == nil {
		return nil
	}
	return o.RequireAdditionalApproval
}

func (o *RequestApprovalConfig) GetCustomApprovalMessage() *string {
	if o == nil {
		return nil
	}
	return o.CustomApprovalMessage
}

func (o *RequestApprovalConfig) GetCustomApprovalMessageOverride() *bool {
	if o == nil {
		return nil
	}
	return o.CustomApprovalMessageOverride
}

func (o *RequestApprovalConfig) GetApprovers() *RequestablePermissionOutputApprovers {
	if o == nil {
		return nil
	}
	return o.Approvers
}

func (o *RequestApprovalConfig) GetApproversStage2() *RequestablePermissionOutputApproversStage2 {
	if o == nil {
		return nil
	}
	return o.ApproversStage2
}

func (o *RequestApprovalConfig) GetRequestApprovalStages() []RequestApprovalStageOutput {
	if o == nil {
		return nil
	}
	return o.RequestApprovalStages
}

// RequestablePermissionOutputLifecycle - The lifecycle of this group.
type RequestablePermissionOutputLifecycle string

const (
	RequestablePermissionOutputLifecycleSynced RequestablePermissionOutputLifecycle = "SYNCED"
	RequestablePermissionOutputLifecycleNative RequestablePermissionOutputLifecycle = "NATIVE"
)

func (e RequestablePermissionOutputLifecycle) ToPointer() *RequestablePermissionOutputLifecycle {
	return &e
}
func (e *RequestablePermissionOutputLifecycle) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "SYNCED":
		fallthrough
	case "NATIVE":
		*e = RequestablePermissionOutputLifecycle(v)
		return nil
	default:
		return fmt.Errorf("invalid value for RequestablePermissionOutputLifecycle: %v", v)
	}
}

// ProvisioningGroup - The provisioning group optionally assocated with this config.
type ProvisioningGroup struct {
	// The ID of this group.
	ID *string `json:"id,omitempty"`
	// The ID of the app that sources this group.
	AppID *string `json:"app_id,omitempty"`
	// The ID of this group, specific to the integration.
	IntegrationSpecificID *string `json:"integration_specific_id,omitempty"`
	// The name of this group.
	Name *string `json:"name,omitempty"`
	// The description of this group.
	Description *string `json:"description,omitempty"`
	// The lifecycle of this group.
	GroupLifecycle *RequestablePermissionOutputLifecycle `default:"SYNCED" json:"group_lifecycle"`
	// The ID of the app that sources this group.
	SourceAppID *string `json:"source_app_id,omitempty"`
}

func (p ProvisioningGroup) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(p, "", false)
}

func (p *ProvisioningGroup) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &p, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *ProvisioningGroup) GetID() *string {
	if o == nil {
		return nil
	}
	return o.ID
}

func (o *ProvisioningGroup) GetAppID() *string {
	if o == nil {
		return nil
	}
	return o.AppID
}

func (o *ProvisioningGroup) GetIntegrationSpecificID() *string {
	if o == nil {
		return nil
	}
	return o.IntegrationSpecificID
}

func (o *ProvisioningGroup) GetName() *string {
	if o == nil {
		return nil
	}
	return o.Name
}

func (o *ProvisioningGroup) GetDescription() *string {
	if o == nil {
		return nil
	}
	return o.Description
}

func (o *ProvisioningGroup) GetGroupLifecycle() *RequestablePermissionOutputLifecycle {
	if o == nil {
		return nil
	}
	return o.GroupLifecycle
}

func (o *ProvisioningGroup) GetSourceAppID() *string {
	if o == nil {
		return nil
	}
	return o.SourceAppID
}

// RequestablePermissionOutputProvisioningWebhook - The provisioning webhook optionally associated with this config.
type RequestablePermissionOutputProvisioningWebhook struct {
	// The ID of this inline webhook.
	ID string `json:"id"`
	// The type of this inline webhook.
	HookType InlineWebhookType `json:"hook_type"`
	// The name of this inline webhook.
	Name string `json:"name"`
	// The description of this inline webhook.
	Description *string `json:"description,omitempty"`
}

func (o *RequestablePermissionOutputProvisioningWebhook) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *RequestablePermissionOutputProvisioningWebhook) GetHookType() InlineWebhookType {
	if o == nil {
		return InlineWebhookType("")
	}
	return o.HookType
}

func (o *RequestablePermissionOutputProvisioningWebhook) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

func (o *RequestablePermissionOutputProvisioningWebhook) GetDescription() *string {
	if o == nil {
		return nil
	}
	return o.Description
}

// RequestFulfillmentConfig - A request fulfillment config can be optionally associated with this config
type RequestFulfillmentConfig struct {
	// Whether manual steps are needed.
	ManualStepsNeeded *bool `json:"manual_steps_needed,omitempty"`
	// The manual instructions that go along.
	ManualInstructions *string `json:"manual_instructions,omitempty"`
	// If enabled, users can request an app for a selected duration. After expiry, Lumos will automatically remove user's access.
	TimeBasedAccess []TimeBasedAccessOptions `json:"time_based_access,omitempty"`
	// Indicates if time based access is overriden.
	TimeBasedAccessOverride *bool `json:"time_based_access_override,omitempty"`
	// The provisioning group optionally assocated with this config.
	ProvisioningGroup *ProvisioningGroup `json:"provisioning_group,omitempty"`
	// The provisioning webhook optionally associated with this config.
	ProvisioningWebhook *RequestablePermissionOutputProvisioningWebhook `json:"provisioning_webhook,omitempty"`
}

func (o *RequestFulfillmentConfig) GetManualStepsNeeded() *bool {
	if o == nil {
		return nil
	}
	return o.ManualStepsNeeded
}

func (o *RequestFulfillmentConfig) GetManualInstructions() *string {
	if o == nil {
		return nil
	}
	return o.ManualInstructions
}

func (o *RequestFulfillmentConfig) GetTimeBasedAccess() []TimeBasedAccessOptions {
	if o == nil {
		return nil
	}
	return o.TimeBasedAccess
}

func (o *RequestFulfillmentConfig) GetTimeBasedAccessOverride() *bool {
	if o == nil {
		return nil
	}
	return o.TimeBasedAccessOverride
}

func (o *RequestFulfillmentConfig) GetProvisioningGroup() *ProvisioningGroup {
	if o == nil {
		return nil
	}
	return o.ProvisioningGroup
}

func (o *RequestFulfillmentConfig) GetProvisioningWebhook() *RequestablePermissionOutputProvisioningWebhook {
	if o == nil {
		return nil
	}
	return o.ProvisioningWebhook
}

// RequestablePermissionOutputAccessRemovalInlineWebhook - A deprovisioning webhook can be optionally associated with this config.
type RequestablePermissionOutputAccessRemovalInlineWebhook struct {
	// The ID of this inline webhook.
	ID string `json:"id"`
	// The type of this inline webhook.
	HookType InlineWebhookType `json:"hook_type"`
	// The name of this inline webhook.
	Name string `json:"name"`
	// The description of this inline webhook.
	Description *string `json:"description,omitempty"`
}

func (o *RequestablePermissionOutputAccessRemovalInlineWebhook) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *RequestablePermissionOutputAccessRemovalInlineWebhook) GetHookType() InlineWebhookType {
	if o == nil {
		return InlineWebhookType("")
	}
	return o.HookType
}

func (o *RequestablePermissionOutputAccessRemovalInlineWebhook) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

func (o *RequestablePermissionOutputAccessRemovalInlineWebhook) GetDescription() *string {
	if o == nil {
		return nil
	}
	return o.Description
}

// RequestablePermissionOutputRequestValidationInlineWebhook - A request validation webhook can be optionally associated with this config.
type RequestablePermissionOutputRequestValidationInlineWebhook struct {
	// The ID of this inline webhook.
	ID string `json:"id"`
	// The type of this inline webhook.
	HookType InlineWebhookType `json:"hook_type"`
	// The name of this inline webhook.
	Name string `json:"name"`
	// The description of this inline webhook.
	Description *string `json:"description,omitempty"`
}

func (o *RequestablePermissionOutputRequestValidationInlineWebhook) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *RequestablePermissionOutputRequestValidationInlineWebhook) GetHookType() InlineWebhookType {
	if o == nil {
		return InlineWebhookType("")
	}
	return o.HookType
}

func (o *RequestablePermissionOutputRequestValidationInlineWebhook) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

func (o *RequestablePermissionOutputRequestValidationInlineWebhook) GetDescription() *string {
	if o == nil {
		return nil
	}
	return o.Description
}

// RequestConfig - The request config associated with this requestable permission.
type RequestConfig struct {
	// The appstore visibility of this request config.
	AppstoreVisibility *AppStoreVisibilityOption `default:"HIDDEN" json:"appstore_visibility"`
	// Indicates if allowed groups is overriden from the app-level settings.
	AllowedGroupsOverride *bool `json:"allowed_groups_override,omitempty"`
	// The allowed groups config associated with this config.
	AllowedGroups *RequestablePermissionOutputAllowedGroups `json:"allowed_groups,omitempty"`
	// A request approval config can be optionally associated with this config
	RequestApprovalConfig *RequestApprovalConfig `json:"request_approval_config,omitempty"`
	// A request fulfillment config can be optionally associated with this config
	RequestFulfillmentConfig *RequestFulfillmentConfig `json:"request_fulfillment_config,omitempty"`
	// A deprovisioning webhook can be optionally associated with this config.
	AccessRemovalInlineWebhook *RequestablePermissionOutputAccessRemovalInlineWebhook `json:"access_removal_inline_webhook,omitempty"`
	// A request validation webhook can be optionally associated with this config.
	RequestValidationInlineWebhook *RequestablePermissionOutputRequestValidationInlineWebhook `json:"request_validation_inline_webhook,omitempty"`
}

func (r RequestConfig) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(r, "", false)
}

func (r *RequestConfig) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &r, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *RequestConfig) GetAppstoreVisibility() *AppStoreVisibilityOption {
	if o == nil {
		return nil
	}
	return o.AppstoreVisibility
}

func (o *RequestConfig) GetAllowedGroupsOverride() *bool {
	if o == nil {
		return nil
	}
	return o.AllowedGroupsOverride
}

func (o *RequestConfig) GetAllowedGroups() *RequestablePermissionOutputAllowedGroups {
	if o == nil {
		return nil
	}
	return o.AllowedGroups
}

func (o *RequestConfig) GetRequestApprovalConfig() *RequestApprovalConfig {
	if o == nil {
		return nil
	}
	return o.RequestApprovalConfig
}

func (o *RequestConfig) GetRequestFulfillmentConfig() *RequestFulfillmentConfig {
	if o == nil {
		return nil
	}
	return o.RequestFulfillmentConfig
}

func (o *RequestConfig) GetAccessRemovalInlineWebhook() *RequestablePermissionOutputAccessRemovalInlineWebhook {
	if o == nil {
		return nil
	}
	return o.AccessRemovalInlineWebhook
}

func (o *RequestConfig) GetRequestValidationInlineWebhook() *RequestablePermissionOutputRequestValidationInlineWebhook {
	if o == nil {
		return nil
	}
	return o.RequestValidationInlineWebhook
}

type RequestablePermissionOutput struct {
	// The ID of this requestable permission.
	ID *string `json:"id,omitempty"`
	// The type of this requestable permission.
	Type *RequestablePermissionOutputPermissionType `json:"type,omitempty"`
	// The label of this requestable permission.
	Label string `json:"label"`
	// The ID of the app associated with this requestable permission.
	AppID string `json:"app_id"`
	// The non-unique ID of the service associated with this requestable permission. Depending on how it is sourced in Lumos, this may be the app's name, website,  or other identifier.
	AppClassID string `json:"app_class_id"`
	// The ID of the instance associated with this requestable permission.
	AppInstanceID string `json:"app_instance_id"`
	// The request config associated with this requestable permission.
	RequestConfig RequestConfig `json:"request_config"`
}

func (o *RequestablePermissionOutput) GetID() *string {
	if o == nil {
		return nil
	}
	return o.ID
}

func (o *RequestablePermissionOutput) GetType() *RequestablePermissionOutputPermissionType {
	if o == nil {
		return nil
	}
	return o.Type
}

func (o *RequestablePermissionOutput) GetLabel() string {
	if o == nil {
		return ""
	}
	return o.Label
}

func (o *RequestablePermissionOutput) GetAppID() string {
	if o == nil {
		return ""
	}
	return o.AppID
}

func (o *RequestablePermissionOutput) GetAppClassID() string {
	if o == nil {
		return ""
	}
	return o.AppClassID
}

func (o *RequestablePermissionOutput) GetAppInstanceID() string {
	if o == nil {
		return ""
	}
	return o.AppInstanceID
}

func (o *RequestablePermissionOutput) GetRequestConfig() RequestConfig {
	if o == nil {
		return RequestConfig{}
	}
	return o.RequestConfig
}
