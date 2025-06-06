// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/utils"
)

type AppApproversOutput struct {
	// Groups assigned as support request approvers.
	Groups []Group `json:"groups"`
	// Users assigned as support request approvers.
	Users []User `json:"users"`
}

func (a AppApproversOutput) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(a, "", false)
}

func (a *AppApproversOutput) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &a, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *AppApproversOutput) GetGroups() []Group {
	if o == nil {
		return nil
	}
	return o.Groups
}

func (o *AppApproversOutput) GetUsers() []User {
	if o == nil {
		return nil
	}
	return o.Users
}
