// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type AddAppToAppStoreInput struct {
	// AppStore App instructions that are shown to the requester.
	CustomRequestInstructions *string                               `json:"custom_request_instructions,omitempty"`
	RequestFlow               *AppStoreAppSettingsRequestFlowInput  `json:"request_flow,omitempty"`
	Provisioning              *AppStoreAppSettingsProvisioningInput `json:"provisioning,omitempty"`
	// The ID of the app to add to the app store.
	AppID string `json:"app_id"`
}

func (o *AddAppToAppStoreInput) GetCustomRequestInstructions() *string {
	if o == nil {
		return nil
	}
	return o.CustomRequestInstructions
}

func (o *AddAppToAppStoreInput) GetRequestFlow() *AppStoreAppSettingsRequestFlowInput {
	if o == nil {
		return nil
	}
	return o.RequestFlow
}

func (o *AddAppToAppStoreInput) GetProvisioning() *AppStoreAppSettingsProvisioningInput {
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
