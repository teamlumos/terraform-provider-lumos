// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

type AttributeInput struct {
	// The unique identifier of the custom attribute
	UniqueIdentifier *string `json:"unique_identifier,omitempty"`
	// The type of the attribute.
	Type *string `json:"type,omitempty"`
	// The name of the attribute.
	Name *string `json:"name,omitempty"`
}

func (o *AttributeInput) GetUniqueIdentifier() *string {
	if o == nil {
		return nil
	}
	return o.UniqueIdentifier
}

func (o *AttributeInput) GetType() *string {
	if o == nil {
		return nil
	}
	return o.Type
}

func (o *AttributeInput) GetName() *string {
	if o == nil {
		return nil
	}
	return o.Name
}
