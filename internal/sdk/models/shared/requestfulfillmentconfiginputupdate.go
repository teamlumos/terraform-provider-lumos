// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type RequestFulfillmentConfigInputUpdate struct {
	// Whether manual steps are needed.
	ManualStepsNeeded *bool `json:"manual_steps_needed,omitempty"`
	// The manual instructions that go along.
	ManualInstructions *string `json:"manual_instructions,omitempty"`
	// If enabled, users can request an app for a selected duration. After expiry, Lumos will automatically remove user's access.
	TimeBasedAccess []TimeBasedAccessOptions `json:"time_based_access,omitempty"`
	// Indicates if time based access is overriden.
	TimeBasedAccessOverride *bool `json:"time_based_access_override,omitempty"`
	// The provisioning webhook optionally associated with this config.
	ProvisioningWebhook *BaseInlineWebhook `json:"provisioning_webhook,omitempty"`
}

func (o *RequestFulfillmentConfigInputUpdate) GetManualStepsNeeded() *bool {
	if o == nil {
		return nil
	}
	return o.ManualStepsNeeded
}

func (o *RequestFulfillmentConfigInputUpdate) GetManualInstructions() *string {
	if o == nil {
		return nil
	}
	return o.ManualInstructions
}

func (o *RequestFulfillmentConfigInputUpdate) GetTimeBasedAccess() []TimeBasedAccessOptions {
	if o == nil {
		return nil
	}
	return o.TimeBasedAccess
}

func (o *RequestFulfillmentConfigInputUpdate) GetTimeBasedAccessOverride() *bool {
	if o == nil {
		return nil
	}
	return o.TimeBasedAccessOverride
}

func (o *RequestFulfillmentConfigInputUpdate) GetProvisioningWebhook() *BaseInlineWebhook {
	if o == nil {
		return nil
	}
	return o.ProvisioningWebhook
}
