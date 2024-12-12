// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type RequestConfigOutput struct {
	// The appstore visibility of this request config.
	AppstoreVisibility *AppStoreVisibilityOption `json:"appstore_visibility,omitempty"`
	// Indicates if allowed groups is overriden from the app-level settings.
	AllowedGroupsOverride *bool `json:"allowed_groups_override,omitempty"`
	// The allowed groups config associated with this config.
	AllowedGroups *AllowedGroupsConfigOutput `json:"allowed_groups,omitempty"`
	// A request approval config can be optionally associated with this config
	RequestApprovalConfig *RequestApprovalConfigOutput `json:"request_approval_config,omitempty"`
	// A request fulfillment config can be optionally associated with this config
	RequestFulfillmentConfig *RequestFulfillmentConfigOutput `json:"request_fulfillment_config,omitempty"`
	// A deprovisioning webhook can be optionally associated with this config.
	AccessRemovalInlineWebhook *InlineWebhook `json:"access_removal_inline_webhook,omitempty"`
	// A request validation webhook can be optionally associated with this config.
	RequestValidationInlineWebhook *InlineWebhook `json:"request_validation_inline_webhook,omitempty"`
}

func (o *RequestConfigOutput) GetAppstoreVisibility() *AppStoreVisibilityOption {
	if o == nil {
		return nil
	}
	return o.AppstoreVisibility
}

func (o *RequestConfigOutput) GetAllowedGroupsOverride() *bool {
	if o == nil {
		return nil
	}
	return o.AllowedGroupsOverride
}

func (o *RequestConfigOutput) GetAllowedGroups() *AllowedGroupsConfigOutput {
	if o == nil {
		return nil
	}
	return o.AllowedGroups
}

func (o *RequestConfigOutput) GetRequestApprovalConfig() *RequestApprovalConfigOutput {
	if o == nil {
		return nil
	}
	return o.RequestApprovalConfig
}

func (o *RequestConfigOutput) GetRequestFulfillmentConfig() *RequestFulfillmentConfigOutput {
	if o == nil {
		return nil
	}
	return o.RequestFulfillmentConfig
}

func (o *RequestConfigOutput) GetAccessRemovalInlineWebhook() *InlineWebhook {
	if o == nil {
		return nil
	}
	return o.AccessRemovalInlineWebhook
}

func (o *RequestConfigOutput) GetRequestValidationInlineWebhook() *InlineWebhook {
	if o == nil {
		return nil
	}
	return o.RequestValidationInlineWebhook
}