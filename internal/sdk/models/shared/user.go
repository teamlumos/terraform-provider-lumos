// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type User struct {
	// The email of this user.
	Email *string `json:"email,omitempty"`
	// The family name of this user.
	FamilyName *string `json:"family_name,omitempty"`
	// The given name of this user.
	GivenName *string `json:"given_name,omitempty"`
	// The ID of this user.
	ID string `json:"id"`
	// The status of this user.
	Status *UserLifecycleStatus `json:"status,omitempty"`
}

func (o *User) GetEmail() *string {
	if o == nil {
		return nil
	}
	return o.Email
}

func (o *User) GetFamilyName() *string {
	if o == nil {
		return nil
	}
	return o.FamilyName
}

func (o *User) GetGivenName() *string {
	if o == nil {
		return nil
	}
	return o.GivenName
}

func (o *User) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *User) GetStatus() *UserLifecycleStatus {
	if o == nil {
		return nil
	}
	return o.Status
}
