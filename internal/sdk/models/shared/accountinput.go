// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type AccountInput struct {
	// The stable identifier of this account.
	UniqueIdentifier string `json:"unique_identifier"`
	// The email of this account.
	Email *string `json:"email,omitempty"`
	// The first name of the user.
	FirstName *string `json:"first_name,omitempty"`
	// The last name of the user.
	LastName *string `json:"last_name,omitempty"`
	// The datetime of last activity of the user.
	LastActivity *string `json:"last_activity,omitempty"`
	// The datetime of last login of the user.
	LastLogin *string `json:"last_login,omitempty"`
	// The status of the account.
	Status *AccountLifecycleStatus `json:"status,omitempty"`
	// The permissions of the account.
	Permissions []PermissionInput `json:"permissions,omitempty"`
}

func (o *AccountInput) GetUniqueIdentifier() string {
	if o == nil {
		return ""
	}
	return o.UniqueIdentifier
}

func (o *AccountInput) GetEmail() *string {
	if o == nil {
		return nil
	}
	return o.Email
}

func (o *AccountInput) GetFirstName() *string {
	if o == nil {
		return nil
	}
	return o.FirstName
}

func (o *AccountInput) GetLastName() *string {
	if o == nil {
		return nil
	}
	return o.LastName
}

func (o *AccountInput) GetLastActivity() *string {
	if o == nil {
		return nil
	}
	return o.LastActivity
}

func (o *AccountInput) GetLastLogin() *string {
	if o == nil {
		return nil
	}
	return o.LastLogin
}

func (o *AccountInput) GetStatus() *AccountLifecycleStatus {
	if o == nil {
		return nil
	}
	return o.Status
}

func (o *AccountInput) GetPermissions() []PermissionInput {
	if o == nil {
		return nil
	}
	return o.Permissions
}
