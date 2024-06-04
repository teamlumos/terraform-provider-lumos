// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/utils"
)

// AddAppToAppStoreInputAllowedGroupsConfigType - The type of this allowed groups config, can be all groups or specific.
type AddAppToAppStoreInputAllowedGroupsConfigType string

const (
	AddAppToAppStoreInputAllowedGroupsConfigTypeAllGroups       AddAppToAppStoreInputAllowedGroupsConfigType = "ALL_GROUPS"
	AddAppToAppStoreInputAllowedGroupsConfigTypeSpecifiedGroups AddAppToAppStoreInputAllowedGroupsConfigType = "SPECIFIED_GROUPS"
)

func (e AddAppToAppStoreInputAllowedGroupsConfigType) ToPointer() *AddAppToAppStoreInputAllowedGroupsConfigType {
	return &e
}
func (e *AddAppToAppStoreInputAllowedGroupsConfigType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "ALL_GROUPS":
		fallthrough
	case "SPECIFIED_GROUPS":
		*e = AddAppToAppStoreInputAllowedGroupsConfigType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for AddAppToAppStoreInputAllowedGroupsConfigType: %v", v)
	}
}

// AddAppToAppStoreInputAllowedGroups - The allowed groups associated with this config.
type AddAppToAppStoreInputAllowedGroups struct {
	// The type of this allowed groups config, can be all groups or specific.
	Type *AddAppToAppStoreInputAllowedGroupsConfigType `default:"ALL_GROUPS" json:"type"`
	// The groups associated with this config.
	Groups []BaseGroup `json:"groups,omitempty"`
}

func (a AddAppToAppStoreInputAllowedGroups) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(a, "", false)
}

func (a *AddAppToAppStoreInputAllowedGroups) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &a, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *AddAppToAppStoreInputAllowedGroups) GetType() *AddAppToAppStoreInputAllowedGroupsConfigType {
	if o == nil {
		return nil
	}
	return o.Type
}

func (o *AddAppToAppStoreInputAllowedGroups) GetGroups() []BaseGroup {
	if o == nil {
		return nil
	}
	return o.Groups
}

// AddAppToAppStoreInputApprovers - AppStore App approvers assigned.
type AddAppToAppStoreInputApprovers struct {
	// Groups assigned as support request approvers.
	Groups []BaseGroup `json:"groups,omitempty"`
	// Users assigned as support request approvers.
	Users []BaseUser `json:"users,omitempty"`
}

func (o *AddAppToAppStoreInputApprovers) GetGroups() []BaseGroup {
	if o == nil {
		return nil
	}
	return o.Groups
}

func (o *AddAppToAppStoreInputApprovers) GetUsers() []BaseUser {
	if o == nil {
		return nil
	}
	return o.Users
}

// AddAppToAppStoreInputApproversStage2 - AppStore App stage 2 approvers assigned.
type AddAppToAppStoreInputApproversStage2 struct {
	// Groups assigned as support request approvers.
	Groups []BaseGroup `json:"groups,omitempty"`
	// Users assigned as support request approvers.
	Users []BaseUser `json:"users,omitempty"`
}

func (o *AddAppToAppStoreInputApproversStage2) GetGroups() []BaseGroup {
	if o == nil {
		return nil
	}
	return o.Groups
}

func (o *AddAppToAppStoreInputApproversStage2) GetUsers() []BaseUser {
	if o == nil {
		return nil
	}
	return o.Users
}

// AddAppToAppStoreInputAdmins - AppStore App admins assigned.
type AddAppToAppStoreInputAdmins struct {
	// Groups assigned as app admins.
	Groups []BaseGroup `json:"groups,omitempty"`
	// Users assigned as app admins.
	Users []BaseUser `json:"users,omitempty"`
}

func (o *AddAppToAppStoreInputAdmins) GetGroups() []BaseGroup {
	if o == nil {
		return nil
	}
	return o.Groups
}

func (o *AddAppToAppStoreInputAdmins) GetUsers() []BaseUser {
	if o == nil {
		return nil
	}
	return o.Users
}

// AddAppToAppStoreInputRequestValidationInlineWebhook - A request validation webhook can be optionally associated with this app.
type AddAppToAppStoreInputRequestValidationInlineWebhook struct {
	// The ID of this inline webhook.
	ID string `json:"id"`
}

func (o *AddAppToAppStoreInputRequestValidationInlineWebhook) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

// AddAppToAppStoreInputRequestFlow - Request flow configuration to request access to app.
type AddAppToAppStoreInputRequestFlow struct {
	// AppStore App visibility.
	Discoverability *AppStoreVisibility `json:"discoverability,omitempty"`
	// During the approval step, send a custom message to requesters. Markdown for links and text formatting is supported.
	CustomApprovalMessage *string `json:"custom_approval_message,omitempty"`
	// When a user makes an access request, require that their manager approves the request before moving on to additional approvals.
	RequireManagerApproval *bool `json:"require_manager_approval,omitempty"`
	// Only turn on when working with sensitive permissions to ensure a smooth employee experience.
	RequireAdditionalApproval *bool `json:"require_additional_approval,omitempty"`
	// The allowed groups associated with this config.
	AllowedGroups *AddAppToAppStoreInputAllowedGroups `json:"allowed_groups,omitempty"`
	// AppStore App approvers assigned.
	Approvers *AddAppToAppStoreInputApprovers `json:"approvers,omitempty"`
	// AppStore App stage 2 approvers assigned.
	ApproversStage2 *AddAppToAppStoreInputApproversStage2 `json:"approvers_stage_2,omitempty"`
	// AppStore App admins assigned.
	Admins *AddAppToAppStoreInputAdmins `json:"admins,omitempty"`
	// A request validation webhook can be optionally associated with this app.
	RequestValidationInlineWebhook *AddAppToAppStoreInputRequestValidationInlineWebhook `json:"request_validation_inline_webhook,omitempty"`
}

func (o *AddAppToAppStoreInputRequestFlow) GetDiscoverability() *AppStoreVisibility {
	if o == nil {
		return nil
	}
	return o.Discoverability
}

func (o *AddAppToAppStoreInputRequestFlow) GetCustomApprovalMessage() *string {
	if o == nil {
		return nil
	}
	return o.CustomApprovalMessage
}

func (o *AddAppToAppStoreInputRequestFlow) GetRequireManagerApproval() *bool {
	if o == nil {
		return nil
	}
	return o.RequireManagerApproval
}

func (o *AddAppToAppStoreInputRequestFlow) GetRequireAdditionalApproval() *bool {
	if o == nil {
		return nil
	}
	return o.RequireAdditionalApproval
}

func (o *AddAppToAppStoreInputRequestFlow) GetAllowedGroups() *AddAppToAppStoreInputAllowedGroups {
	if o == nil {
		return nil
	}
	return o.AllowedGroups
}

func (o *AddAppToAppStoreInputRequestFlow) GetApprovers() *AddAppToAppStoreInputApprovers {
	if o == nil {
		return nil
	}
	return o.Approvers
}

func (o *AddAppToAppStoreInputRequestFlow) GetApproversStage2() *AddAppToAppStoreInputApproversStage2 {
	if o == nil {
		return nil
	}
	return o.ApproversStage2
}

func (o *AddAppToAppStoreInputRequestFlow) GetAdmins() *AddAppToAppStoreInputAdmins {
	if o == nil {
		return nil
	}
	return o.Admins
}

func (o *AddAppToAppStoreInputRequestFlow) GetRequestValidationInlineWebhook() *AddAppToAppStoreInputRequestValidationInlineWebhook {
	if o == nil {
		return nil
	}
	return o.RequestValidationInlineWebhook
}

// AddAppToAppStoreInputProvisioningWebhook - The provisioning webhook optionally associated with this app.
type AddAppToAppStoreInputProvisioningWebhook struct {
	// The ID of this inline webhook.
	ID string `json:"id"`
}

func (o *AddAppToAppStoreInputProvisioningWebhook) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

// AddAppToAppStoreInputAccessRemovalInlineWebhook - An inactivity workflow can be optionally associated with this app.
type AddAppToAppStoreInputAccessRemovalInlineWebhook struct {
	// The ID of this inline webhook.
	ID string `json:"id"`
}

func (o *AddAppToAppStoreInputAccessRemovalInlineWebhook) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

// AddAppToAppStoreInputProvisioning - Provisioning flow configuration to request access to app.
type AddAppToAppStoreInputProvisioning struct {
	// If enabled, Approvers must choose a group to provision the user to for access requests.
	GroupsProvisioning *GroupProvisioningOption `json:"groups_provisioning,omitempty"`
	// If enabled, users can request an app for a selected duration. After expiry, Lumos will automatically remove user's access.
	TimeBasedAccess []TimeBasedAccessOptions `json:"time_based_access,omitempty"`
	// If enabled, Lumos will reach out to the App Admin after initial access is granted to perform additional manual steps. Note that if this option is enabled, this action must be confirmed by the App Admin in order to resolve the request.
	ManualStepsNeeded *bool `json:"manual_steps_needed,omitempty"`
	// Only Available if manual steps is active. During the provisioning step, send a custom message to app admins explaining how to provision a user to the app. Markdown for links and text formatting is supported.
	CustomProvisioningInstructions *string `json:"custom_provisioning_instructions,omitempty"`
	// The provisioning webhook optionally associated with this app.
	ProvisioningWebhook *AddAppToAppStoreInputProvisioningWebhook `json:"provisioning_webhook,omitempty"`
	// An inactivity workflow can be optionally associated with this app.
	AccessRemovalInlineWebhook *AddAppToAppStoreInputAccessRemovalInlineWebhook `json:"access_removal_inline_webhook,omitempty"`
}

func (o *AddAppToAppStoreInputProvisioning) GetGroupsProvisioning() *GroupProvisioningOption {
	if o == nil {
		return nil
	}
	return o.GroupsProvisioning
}

func (o *AddAppToAppStoreInputProvisioning) GetTimeBasedAccess() []TimeBasedAccessOptions {
	if o == nil {
		return nil
	}
	return o.TimeBasedAccess
}

func (o *AddAppToAppStoreInputProvisioning) GetManualStepsNeeded() *bool {
	if o == nil {
		return nil
	}
	return o.ManualStepsNeeded
}

func (o *AddAppToAppStoreInputProvisioning) GetCustomProvisioningInstructions() *string {
	if o == nil {
		return nil
	}
	return o.CustomProvisioningInstructions
}

func (o *AddAppToAppStoreInputProvisioning) GetProvisioningWebhook() *AddAppToAppStoreInputProvisioningWebhook {
	if o == nil {
		return nil
	}
	return o.ProvisioningWebhook
}

func (o *AddAppToAppStoreInputProvisioning) GetAccessRemovalInlineWebhook() *AddAppToAppStoreInputAccessRemovalInlineWebhook {
	if o == nil {
		return nil
	}
	return o.AccessRemovalInlineWebhook
}

type AddAppToAppStoreInput struct {
	// AppStore App instructions.
	CustomRequestInstructions *string `json:"custom_request_instructions,omitempty"`
	// Request flow configuration to request access to app.
	RequestFlow *AddAppToAppStoreInputRequestFlow `json:"request_flow,omitempty"`
	// Provisioning flow configuration to request access to app.
	Provisioning *AddAppToAppStoreInputProvisioning `json:"provisioning,omitempty"`
	// The ID of the app to add to the app store.
	AppID string `json:"app_id"`
}

func (o *AddAppToAppStoreInput) GetCustomRequestInstructions() *string {
	if o == nil {
		return nil
	}
	return o.CustomRequestInstructions
}

func (o *AddAppToAppStoreInput) GetRequestFlow() *AddAppToAppStoreInputRequestFlow {
	if o == nil {
		return nil
	}
	return o.RequestFlow
}

func (o *AddAppToAppStoreInput) GetProvisioning() *AddAppToAppStoreInputProvisioning {
	if o == nil {
		return nil
	}
	return o.Provisioning
}

func (o *AddAppToAppStoreInput) GetAppID() string {
	if o == nil {
		return ""
	}
	return o.AppID
}