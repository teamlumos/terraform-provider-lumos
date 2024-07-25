// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

type AccountsUploadInput struct {
	// Accounts to upload.
	Accounts []AccountInput `json:"accounts,omitempty"`
	// The ID of the app to upload accounts to.
	AppID string `json:"app_id"`
}

func (o *AccountsUploadInput) GetAccounts() []AccountInput {
	if o == nil {
		return nil
	}
	return o.Accounts
}

func (o *AccountsUploadInput) GetAppID() string {
	if o == nil {
		return ""
	}
	return o.AppID
}
