// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type Group struct {
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
	GroupLifecycle *Lifecycle `json:"group_lifecycle,omitempty"`
	// The ID of the app that sources this group.
	SourceAppID *string `json:"source_app_id,omitempty"`
}

func (o *Group) GetID() *string {
	if o == nil {
		return nil
	}
	return o.ID
}

func (o *Group) GetAppID() *string {
	if o == nil {
		return nil
	}
	return o.AppID
}

func (o *Group) GetIntegrationSpecificID() *string {
	if o == nil {
		return nil
	}
	return o.IntegrationSpecificID
}

func (o *Group) GetName() *string {
	if o == nil {
		return nil
	}
	return o.Name
}

func (o *Group) GetDescription() *string {
	if o == nil {
		return nil
	}
	return o.Description
}

func (o *Group) GetGroupLifecycle() *Lifecycle {
	if o == nil {
		return nil
	}
	return o.GroupLifecycle
}

func (o *Group) GetSourceAppID() *string {
	if o == nil {
		return nil
	}
	return o.SourceAppID
}
