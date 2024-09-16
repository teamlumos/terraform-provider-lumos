// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type AppStoreAppSettingsProvisioningInput struct {
	// If enabled, Approvers must choose a group to provision the user to for access requests.
	GroupsProvisioning *GroupProvisioningOption `json:"groups_provisioning,omitempty"`
	// If enabled, users can request an app for a selected duration. After expiry, Lumos will automatically remove user's access.
	TimeBasedAccess []TimeBasedAccessOptions `json:"time_based_access,omitempty"`
	// Whether the app is configured to allow users to request multiple permissions in a single request
	AllowMultiplePermissionSelection *bool `json:"allow_multiple_permission_selection,omitempty"`
	// If enabled, Lumos will notify the App Admin after initial access is granted to perform additional manual steps. Note that if this option is enabled, this action must be confirmed by the App Admin in order to resolve the request.
	ManualStepsNeeded *bool `json:"manual_steps_needed,omitempty"`
	// Only Available if manual steps is active. During the provisioning step, Lumos will send a custom message to app admins explaining how to provision a user to the app. Markdown for links and text formatting is supported.
	CustomProvisioningInstructions *string `json:"custom_provisioning_instructions,omitempty"`
	// The provisioning webhook optionally associated with this app.
	ProvisioningWebhook *BaseInlineWebhook `json:"provisioning_webhook,omitempty"`
	// A deprovisioning webhook can be optionally associated with this app.
	AccessRemovalInlineWebhook *BaseInlineWebhook `json:"access_removal_inline_webhook,omitempty"`
}

func (o *AppStoreAppSettingsProvisioningInput) GetGroupsProvisioning() *GroupProvisioningOption {
	if o == nil {
		return nil
	}
	return o.GroupsProvisioning
}

func (o *AppStoreAppSettingsProvisioningInput) GetTimeBasedAccess() []TimeBasedAccessOptions {
	if o == nil {
		return nil
	}
	return o.TimeBasedAccess
}

func (o *AppStoreAppSettingsProvisioningInput) GetAllowMultiplePermissionSelection() *bool {
	if o == nil {
		return nil
	}
	return o.AllowMultiplePermissionSelection
}

func (o *AppStoreAppSettingsProvisioningInput) GetManualStepsNeeded() *bool {
	if o == nil {
		return nil
	}
	return o.ManualStepsNeeded
}

func (o *AppStoreAppSettingsProvisioningInput) GetCustomProvisioningInstructions() *string {
	if o == nil {
		return nil
	}
	return o.CustomProvisioningInstructions
}

func (o *AppStoreAppSettingsProvisioningInput) GetProvisioningWebhook() *BaseInlineWebhook {
	if o == nil {
		return nil
	}
	return o.ProvisioningWebhook
}

func (o *AppStoreAppSettingsProvisioningInput) GetAccessRemovalInlineWebhook() *BaseInlineWebhook {
	if o == nil {
		return nil
	}
	return o.AccessRemovalInlineWebhook
}
