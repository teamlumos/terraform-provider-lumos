// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type ActivityRecordAppInput struct {
	// The ID of the app as it is identified in the source E.g. the ID that Okta uses to identify the app
	InstanceIdentifier string `json:"instance_identifier"`
}

func (o *ActivityRecordAppInput) GetInstanceIdentifier() string {
	if o == nil {
		return ""
	}
	return o.InstanceIdentifier
}
