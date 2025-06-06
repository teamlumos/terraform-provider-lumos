// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/utils"
)

type AllowedGroupsConfigOutput struct {
	Type *AllowedGroupsConfigType `json:"type,omitempty"`
	// The groups allowed to request this permission.
	Groups []Group `json:"groups"`
}

func (a AllowedGroupsConfigOutput) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(a, "", false)
}

func (a *AllowedGroupsConfigOutput) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &a, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *AllowedGroupsConfigOutput) GetType() *AllowedGroupsConfigType {
	if o == nil {
		return nil
	}
	return o.Type
}

func (o *AllowedGroupsConfigOutput) GetGroups() []Group {
	if o == nil {
		return nil
	}
	return o.Groups
}
