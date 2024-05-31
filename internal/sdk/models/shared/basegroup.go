// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type BaseGroup struct {
	// The ID of this group.
	ID *string `json:"id,omitempty"`
	// The ID of the app that owns this group.
	AppID *string `json:"app_id,omitempty"`
	// The ID of this group, specific to the integration.
	IntegrationSpecificID *string `json:"integration_specific_id,omitempty"`
}

func (o *BaseGroup) GetID() *string {
	if o == nil {
		return nil
	}
	return o.ID
}

func (o *BaseGroup) GetAppID() *string {
	if o == nil {
		return nil
	}
	return o.AppID
}

func (o *BaseGroup) GetIntegrationSpecificID() *string {
	if o == nil {
		return nil
	}
	return o.IntegrationSpecificID
}
