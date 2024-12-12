// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

// AppStoreAppSettingsRequestFlowOutputApproversAppApproversOutput - AppStore App approvers assigned.
type AppStoreAppSettingsRequestFlowOutputApproversAppApproversOutput struct {
	// Groups assigned as support request approvers.
	Groups []Group `json:"groups,omitempty"`
	// Users assigned as support request approvers.
	Users []User `json:"users,omitempty"`
}

func (o *AppStoreAppSettingsRequestFlowOutputApproversAppApproversOutput) GetGroups() []Group {
	if o == nil {
		return nil
	}
	return o.Groups
}

func (o *AppStoreAppSettingsRequestFlowOutputApproversAppApproversOutput) GetUsers() []User {
	if o == nil {
		return nil
	}
	return o.Users
}

// AppStoreAppSettingsRequestFlowOutputAppApproversOutput - AppStore App stage 2 approvers assigned.
type AppStoreAppSettingsRequestFlowOutputAppApproversOutput struct {
	// Groups assigned as support request approvers.
	Groups []Group `json:"groups,omitempty"`
	// Users assigned as support request approvers.
	Users []User `json:"users,omitempty"`
}

func (o *AppStoreAppSettingsRequestFlowOutputAppApproversOutput) GetGroups() []Group {
	if o == nil {
		return nil
	}
	return o.Groups
}

func (o *AppStoreAppSettingsRequestFlowOutputAppApproversOutput) GetUsers() []User {
	if o == nil {
		return nil
	}
	return o.Users
}

// AppAdminsOutput - AppStore App admins assigned.
type AppAdminsOutput struct {
	// Groups assigned as app admins.
	Groups []Group `json:"groups,omitempty"`
	// Users assigned as app admins.
	Users []User `json:"users,omitempty"`
}

func (o *AppAdminsOutput) GetGroups() []Group {
	if o == nil {
		return nil
	}
	return o.Groups
}

func (o *AppAdminsOutput) GetUsers() []User {
	if o == nil {
		return nil
	}
	return o.Users
}

type AppStoreAppSettingsRequestFlowOutput struct {
	// AppStore App visibility.
	Discoverability *AppStoreVisibility `json:"discoverability,omitempty"`
	// After the approval step, send a custom message to requesters. Markdown for links and text formatting is supported.
	CustomApprovalMessage *string `json:"custom_approval_message,omitempty"`
	// When a user makes an access request, require that their manager approves the request before moving on to additional approvals.
	RequireManagerApproval *bool `json:"require_manager_approval,omitempty"`
	// Only turn on when working with sensitive permissions to ensure a smooth employee experience.
	RequireAdditionalApproval *bool `json:"require_additional_approval,omitempty"`
	// The allowed groups config associated with this config.
	AllowedGroups *AllowedGroupsConfigOutput `json:"allowed_groups,omitempty"`
	// AppStore App approvers assigned.
	Approvers *AppStoreAppSettingsRequestFlowOutputApproversAppApproversOutput `json:"approvers,omitempty"`
	// AppStore App stage 2 approvers assigned.
	ApproversStage2 *AppStoreAppSettingsRequestFlowOutputAppApproversOutput `json:"approvers_stage_2,omitempty"`
	// AppStore App admins assigned.
	Admins *AppAdminsOutput `json:"admins,omitempty"`
	// A request validation webhook can be optionally associated with this config.
	RequestValidationInlineWebhook *InlineWebhook `json:"request_validation_inline_webhook,omitempty"`
}

func (o *AppStoreAppSettingsRequestFlowOutput) GetDiscoverability() *AppStoreVisibility {
	if o == nil {
		return nil
	}
	return o.Discoverability
}

func (o *AppStoreAppSettingsRequestFlowOutput) GetCustomApprovalMessage() *string {
	if o == nil {
		return nil
	}
	return o.CustomApprovalMessage
}

func (o *AppStoreAppSettingsRequestFlowOutput) GetRequireManagerApproval() *bool {
	if o == nil {
		return nil
	}
	return o.RequireManagerApproval
}

func (o *AppStoreAppSettingsRequestFlowOutput) GetRequireAdditionalApproval() *bool {
	if o == nil {
		return nil
	}
	return o.RequireAdditionalApproval
}

func (o *AppStoreAppSettingsRequestFlowOutput) GetAllowedGroups() *AllowedGroupsConfigOutput {
	if o == nil {
		return nil
	}
	return o.AllowedGroups
}

func (o *AppStoreAppSettingsRequestFlowOutput) GetApprovers() *AppStoreAppSettingsRequestFlowOutputApproversAppApproversOutput {
	if o == nil {
		return nil
	}
	return o.Approvers
}

func (o *AppStoreAppSettingsRequestFlowOutput) GetApproversStage2() *AppStoreAppSettingsRequestFlowOutputAppApproversOutput {
	if o == nil {
		return nil
	}
	return o.ApproversStage2
}

func (o *AppStoreAppSettingsRequestFlowOutput) GetAdmins() *AppAdminsOutput {
	if o == nil {
		return nil
	}
	return o.Admins
}

func (o *AppStoreAppSettingsRequestFlowOutput) GetRequestValidationInlineWebhook() *InlineWebhook {
	if o == nil {
		return nil
	}
	return o.RequestValidationInlineWebhook
}