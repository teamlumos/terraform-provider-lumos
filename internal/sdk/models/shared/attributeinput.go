// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type AttributeInput struct {
	// The name of the attribute.
	Name *string `json:"name,omitempty"`
	// The type of the attribute.
	Type *string `json:"type,omitempty"`
	// The unique identifier of the custom attribute
	UniqueIdentifier *string `json:"unique_identifier,omitempty"`
}

func (o *AttributeInput) GetName() *string {
	if o == nil {
		return nil
	}
	return o.Name
}

func (o *AttributeInput) GetType() *string {
	if o == nil {
		return nil
	}
	return o.Type
}

func (o *AttributeInput) GetUniqueIdentifier() *string {
	if o == nil {
		return nil
	}
	return o.UniqueIdentifier
}
