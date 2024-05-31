// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
	"github.com/speakeasy/terraform-provider-lumos/internal/sdk/internal/utils"
)

// AppStoreAppSettingsOutputAllowedGroupsConfigType - The type of this allowed groups config, can be all groups or specific.
type AppStoreAppSettingsOutputAllowedGroupsConfigType string

const (
	AppStoreAppSettingsOutputAllowedGroupsConfigTypeAllGroups       AppStoreAppSettingsOutputAllowedGroupsConfigType = "ALL_GROUPS"
	AppStoreAppSettingsOutputAllowedGroupsConfigTypeSpecifiedGroups AppStoreAppSettingsOutputAllowedGroupsConfigType = "SPECIFIED_GROUPS"
)

func (e AppStoreAppSettingsOutputAllowedGroupsConfigType) ToPointer() *AppStoreAppSettingsOutputAllowedGroupsConfigType {
	return &e
}
func (e *AppStoreAppSettingsOutputAllowedGroupsConfigType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "ALL_GROUPS":
		fallthrough
	case "SPECIFIED_GROUPS":
		*e = AppStoreAppSettingsOutputAllowedGroupsConfigType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for AppStoreAppSettingsOutputAllowedGroupsConfigType: %v", v)
	}
}

// AppStoreAppSettingsOutputAllowedGroups - The allowed groups config associated with this config.
type AppStoreAppSettingsOutputAllowedGroups struct {
	// The type of this allowed groups config, can be all groups or specific.
	Type *AppStoreAppSettingsOutputAllowedGroupsConfigType `default:"ALL_GROUPS" json:"type"`
	// The groups associated with this config.
	Groups []Group `json:"groups,omitempty"`
}

func (a AppStoreAppSettingsOutputAllowedGroups) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(a, "", false)
}

func (a *AppStoreAppSettingsOutputAllowedGroups) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &a, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *AppStoreAppSettingsOutputAllowedGroups) GetType() *AppStoreAppSettingsOutputAllowedGroupsConfigType {
	if o == nil {
		return nil
	}
	return o.Type
}

func (o *AppStoreAppSettingsOutputAllowedGroups) GetGroups() []Group {
	if o == nil {
		return nil
	}
	return o.Groups
}

// Approvers - AppStore App approvers assigned.
type Approvers struct {
	// Groups assigned as support request approvers.
	Groups []Group `json:"groups,omitempty"`
	// Users assigned as support request approvers.
	Users []User `json:"users,omitempty"`
}

func (o *Approvers) GetGroups() []Group {
	if o == nil {
		return nil
	}
	return o.Groups
}

func (o *Approvers) GetUsers() []User {
	if o == nil {
		return nil
	}
	return o.Users
}

// ApproversStage2 - AppStore App stage 2 approvers assigned.
type ApproversStage2 struct {
	// Groups assigned as support request approvers.
	Groups []Group `json:"groups,omitempty"`
	// Users assigned as support request approvers.
	Users []User `json:"users,omitempty"`
}

func (o *ApproversStage2) GetGroups() []Group {
	if o == nil {
		return nil
	}
	return o.Groups
}

func (o *ApproversStage2) GetUsers() []User {
	if o == nil {
		return nil
	}
	return o.Users
}

// Admins - AppStore App admins assigned.
type Admins struct {
	// Groups assigned as app admins.
	Groups []Group `json:"groups,omitempty"`
	// Users assigned as app admins.
	Users []User `json:"users,omitempty"`
}

func (o *Admins) GetGroups() []Group {
	if o == nil {
		return nil
	}
	return o.Groups
}

func (o *Admins) GetUsers() []User {
	if o == nil {
		return nil
	}
	return o.Users
}

// AppStoreAppSettingsOutputRequestValidationInlineWebhook - A request validation webhook can be optionally associated with this config.
type AppStoreAppSettingsOutputRequestValidationInlineWebhook struct {
	// The ID of this inline webhook.
	ID string `json:"id"`
	// The type of this inline webhook.
	HookType InlineWebhookType `json:"hook_type"`
	// The name of this inline webhook.
	Name string `json:"name"`
	// The description of this inline webhook.
	Description *string `json:"description,omitempty"`
}

func (o *AppStoreAppSettingsOutputRequestValidationInlineWebhook) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *AppStoreAppSettingsOutputRequestValidationInlineWebhook) GetHookType() InlineWebhookType {
	if o == nil {
		return InlineWebhookType("")
	}
	return o.HookType
}

func (o *AppStoreAppSettingsOutputRequestValidationInlineWebhook) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

func (o *AppStoreAppSettingsOutputRequestValidationInlineWebhook) GetDescription() *string {
	if o == nil {
		return nil
	}
	return o.Description
}

// RequestFlow - Request flow configuration to request access to app.
type RequestFlow struct {
	// AppStore App visibility.
	Discoverability *AppStoreVisibility `json:"discoverability,omitempty"`
	// During the approval step, send a custom message to requesters. Markdown for links and text formatting is supported.
	CustomApprovalMessage *string `json:"custom_approval_message,omitempty"`
	// When a user makes an access request, require that their manager approves the request before moving on to additional approvals.
	RequireManagerApproval *bool `json:"require_manager_approval,omitempty"`
	// Only turn on when working with sensitive permissions to ensure a smooth employee experience.
	RequireAdditionalApproval *bool `json:"require_additional_approval,omitempty"`
	// The allowed groups config associated with this config.
	AllowedGroups *AppStoreAppSettingsOutputAllowedGroups `json:"allowed_groups,omitempty"`
	// AppStore App approvers assigned.
	Approvers *Approvers `json:"approvers,omitempty"`
	// AppStore App stage 2 approvers assigned.
	ApproversStage2 *ApproversStage2 `json:"approvers_stage_2,omitempty"`
	// AppStore App admins assigned.
	Admins *Admins `json:"admins,omitempty"`
	// A request validation webhook can be optionally associated with this config.
	RequestValidationInlineWebhook *AppStoreAppSettingsOutputRequestValidationInlineWebhook `json:"request_validation_inline_webhook,omitempty"`
}

func (o *RequestFlow) GetDiscoverability() *AppStoreVisibility {
	if o == nil {
		return nil
	}
	return o.Discoverability
}

func (o *RequestFlow) GetCustomApprovalMessage() *string {
	if o == nil {
		return nil
	}
	return o.CustomApprovalMessage
}

func (o *RequestFlow) GetRequireManagerApproval() *bool {
	if o == nil {
		return nil
	}
	return o.RequireManagerApproval
}

func (o *RequestFlow) GetRequireAdditionalApproval() *bool {
	if o == nil {
		return nil
	}
	return o.RequireAdditionalApproval
}

func (o *RequestFlow) GetAllowedGroups() *AppStoreAppSettingsOutputAllowedGroups {
	if o == nil {
		return nil
	}
	return o.AllowedGroups
}

func (o *RequestFlow) GetApprovers() *Approvers {
	if o == nil {
		return nil
	}
	return o.Approvers
}

func (o *RequestFlow) GetApproversStage2() *ApproversStage2 {
	if o == nil {
		return nil
	}
	return o.ApproversStage2
}

func (o *RequestFlow) GetAdmins() *Admins {
	if o == nil {
		return nil
	}
	return o.Admins
}

func (o *RequestFlow) GetRequestValidationInlineWebhook() *AppStoreAppSettingsOutputRequestValidationInlineWebhook {
	if o == nil {
		return nil
	}
	return o.RequestValidationInlineWebhook
}

// ProvisioningWebhook - The provisioning webhook optionally associated with this config.
type ProvisioningWebhook struct {
	// The ID of this inline webhook.
	ID string `json:"id"`
	// The type of this inline webhook.
	HookType InlineWebhookType `json:"hook_type"`
	// The name of this inline webhook.
	Name string `json:"name"`
	// The description of this inline webhook.
	Description *string `json:"description,omitempty"`
}

func (o *ProvisioningWebhook) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *ProvisioningWebhook) GetHookType() InlineWebhookType {
	if o == nil {
		return InlineWebhookType("")
	}
	return o.HookType
}

func (o *ProvisioningWebhook) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

func (o *ProvisioningWebhook) GetDescription() *string {
	if o == nil {
		return nil
	}
	return o.Description
}

// AppStoreAppSettingsOutputAccessRemovalInlineWebhook - An inactivity workflow can be optionally associated with this config.
type AppStoreAppSettingsOutputAccessRemovalInlineWebhook struct {
	// The ID of this inline webhook.
	ID string `json:"id"`
	// The type of this inline webhook.
	HookType InlineWebhookType `json:"hook_type"`
	// The name of this inline webhook.
	Name string `json:"name"`
	// The description of this inline webhook.
	Description *string `json:"description,omitempty"`
}

func (o *AppStoreAppSettingsOutputAccessRemovalInlineWebhook) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *AppStoreAppSettingsOutputAccessRemovalInlineWebhook) GetHookType() InlineWebhookType {
	if o == nil {
		return InlineWebhookType("")
	}
	return o.HookType
}

func (o *AppStoreAppSettingsOutputAccessRemovalInlineWebhook) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

func (o *AppStoreAppSettingsOutputAccessRemovalInlineWebhook) GetDescription() *string {
	if o == nil {
		return nil
	}
	return o.Description
}

// Provisioning flow configuration to request access to app.
type Provisioning struct {
	// If enabled, Approvers must choose a group to provision the user to for access requests.
	GroupsProvisioning *GroupProvisioningOption `json:"groups_provisioning,omitempty"`
	// If enabled, users can request an app for a selected duration. After expiry, Lumos will automatically remove user's access.
	TimeBasedAccess []TimeBasedAccessOptions `json:"time_based_access,omitempty"`
	// If enabled, Lumos will reach out to the App Admin after initial access is granted to perform additional manual steps. Note that if this option is enabled, this action must be confirmed by the App Admin in order to resolve the request.
	ManualStepsNeeded *bool `json:"manual_steps_needed,omitempty"`
	// Only Available if manual steps is active. During the provisioning step, send a custom message to app admins explaining how to provision a user to the app. Markdown for links and text formatting is supported.
	CustomProvisioningInstructions *string `json:"custom_provisioning_instructions,omitempty"`
	// The provisioning webhook optionally associated with this config.
	ProvisioningWebhook *ProvisioningWebhook `json:"provisioning_webhook,omitempty"`
	// An inactivity workflow can be optionally associated with this config.
	AccessRemovalInlineWebhook *AppStoreAppSettingsOutputAccessRemovalInlineWebhook `json:"access_removal_inline_webhook,omitempty"`
}

func (o *Provisioning) GetGroupsProvisioning() *GroupProvisioningOption {
	if o == nil {
		return nil
	}
	return o.GroupsProvisioning
}

func (o *Provisioning) GetTimeBasedAccess() []TimeBasedAccessOptions {
	if o == nil {
		return nil
	}
	return o.TimeBasedAccess
}

func (o *Provisioning) GetManualStepsNeeded() *bool {
	if o == nil {
		return nil
	}
	return o.ManualStepsNeeded
}

func (o *Provisioning) GetCustomProvisioningInstructions() *string {
	if o == nil {
		return nil
	}
	return o.CustomProvisioningInstructions
}

func (o *Provisioning) GetProvisioningWebhook() *ProvisioningWebhook {
	if o == nil {
		return nil
	}
	return o.ProvisioningWebhook
}

func (o *Provisioning) GetAccessRemovalInlineWebhook() *AppStoreAppSettingsOutputAccessRemovalInlineWebhook {
	if o == nil {
		return nil
	}
	return o.AccessRemovalInlineWebhook
}

type AppStoreAppSettingsOutput struct {
	// AppStore App instructions.
	CustomRequestInstructions *string `json:"custom_request_instructions,omitempty"`
	// Request flow configuration to request access to app.
	RequestFlow *RequestFlow `json:"request_flow,omitempty"`
	// Provisioning flow configuration to request access to app.
	Provisioning *Provisioning `json:"provisioning,omitempty"`
}

func (o *AppStoreAppSettingsOutput) GetCustomRequestInstructions() *string {
	if o == nil {
		return nil
	}
	return o.CustomRequestInstructions
}

func (o *AppStoreAppSettingsOutput) GetRequestFlow() *RequestFlow {
	if o == nil {
		return nil
	}
	return o.RequestFlow
}

func (o *AppStoreAppSettingsOutput) GetProvisioning() *Provisioning {
	if o == nil {
		return nil
	}
	return o.Provisioning
}
