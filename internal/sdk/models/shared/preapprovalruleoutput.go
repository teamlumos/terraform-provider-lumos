// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type PreApprovalRuleOutput struct {
	// The ID of this preapproval rule.
	ID *string `json:"id,omitempty"`
	// The justification of this preapproval rule.
	Justification string `json:"justification"`
	// The ID of the app associated with this pre-approval rule.
	AppID *string `json:"app_id,omitempty"`
	// The ID of the service associated with this pre-approval rule.
	AppClassID *string `json:"app_class_id,omitempty"`
	// Optionally, an app has an identifer associated with it's particular instance.
	AppInstanceID *string `json:"app_instance_id,omitempty"`
	// The preapproved groups of this preapproval rule.
	PreapprovedGroups []Group `json:"preapproved_groups,omitempty"`
	// The preapproved permissions of this preapproval rule.
	PreapprovedPermissions []RequestablePermissionBaseOutput `json:"preapproved_permissions,omitempty"`
	// The preapproval webhooks of this preapproval rule.
	PreapprovalWebhooks []InlineWebhook `json:"preapproval_webhooks,omitempty"`
	// Preapproval rule time access length,
	TimeBasedAccess []TimeBasedAccessOptions `json:"time_based_access,omitempty"`
}

func (o *PreApprovalRuleOutput) GetID() *string {
	if o == nil {
		return nil
	}
	return o.ID
}

func (o *PreApprovalRuleOutput) GetJustification() string {
	if o == nil {
		return ""
	}
	return o.Justification
}

func (o *PreApprovalRuleOutput) GetAppID() *string {
	if o == nil {
		return nil
	}
	return o.AppID
}

func (o *PreApprovalRuleOutput) GetAppClassID() *string {
	if o == nil {
		return nil
	}
	return o.AppClassID
}

func (o *PreApprovalRuleOutput) GetAppInstanceID() *string {
	if o == nil {
		return nil
	}
	return o.AppInstanceID
}

func (o *PreApprovalRuleOutput) GetPreapprovedGroups() []Group {
	if o == nil {
		return nil
	}
	return o.PreapprovedGroups
}

func (o *PreApprovalRuleOutput) GetPreapprovedPermissions() []RequestablePermissionBaseOutput {
	if o == nil {
		return nil
	}
	return o.PreapprovedPermissions
}

func (o *PreApprovalRuleOutput) GetPreapprovalWebhooks() []InlineWebhook {
	if o == nil {
		return nil
	}
	return o.PreapprovalWebhooks
}

func (o *PreApprovalRuleOutput) GetTimeBasedAccess() []TimeBasedAccessOptions {
	if o == nil {
		return nil
	}
	return o.TimeBasedAccess
}