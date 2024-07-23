// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type CreateAccessRequestInput struct {
	AppID string `json:"app_id"`
	// User to request access for. If omitted, request for the current user.
	TargetUserID *string `json:"target_user_id,omitempty"`
	// Why the user needs access.
	Note string `json:"note"`
	// Once granted, how long the access should last. Omit for permanent access.
	ExpirationInSeconds *int64 `json:"expiration_in_seconds,omitempty"`
	// The IDs of the requestable permissions the user is requesting access to. Omit this to just request access to the app.
	RequestablePermissionIds []string `json:"requestable_permission_ids,omitempty"`
}

func (o *CreateAccessRequestInput) GetAppID() string {
	if o == nil {
		return ""
	}
	return o.AppID
}

func (o *CreateAccessRequestInput) GetTargetUserID() *string {
	if o == nil {
		return nil
	}
	return o.TargetUserID
}

func (o *CreateAccessRequestInput) GetNote() string {
	if o == nil {
		return ""
	}
	return o.Note
}

func (o *CreateAccessRequestInput) GetExpirationInSeconds() *int64 {
	if o == nil {
		return nil
	}
	return o.ExpirationInSeconds
}

func (o *CreateAccessRequestInput) GetRequestablePermissionIds() []string {
	if o == nil {
		return nil
	}
	return o.RequestablePermissionIds
}
