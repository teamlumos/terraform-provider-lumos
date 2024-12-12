// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/utils"
)

// AllowedGroupsConfigInputAllowedGroupsConfigType - The type of this allowed groups config, can be all groups or specific.
type AllowedGroupsConfigInputAllowedGroupsConfigType string

const (
	AllowedGroupsConfigInputAllowedGroupsConfigTypeAllGroups       AllowedGroupsConfigInputAllowedGroupsConfigType = "ALL_GROUPS"
	AllowedGroupsConfigInputAllowedGroupsConfigTypeSpecifiedGroups AllowedGroupsConfigInputAllowedGroupsConfigType = "SPECIFIED_GROUPS"
)

func (e AllowedGroupsConfigInputAllowedGroupsConfigType) ToPointer() *AllowedGroupsConfigInputAllowedGroupsConfigType {
	return &e
}
func (e *AllowedGroupsConfigInputAllowedGroupsConfigType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "ALL_GROUPS":
		fallthrough
	case "SPECIFIED_GROUPS":
		*e = AllowedGroupsConfigInputAllowedGroupsConfigType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for AllowedGroupsConfigInputAllowedGroupsConfigType: %v", v)
	}
}

type AllowedGroupsConfigInput struct {
	// The type of this allowed groups config, can be all groups or specific.
	Type *AllowedGroupsConfigInputAllowedGroupsConfigType `default:"ALL_GROUPS" json:"type"`
	// The groups allowed to request this permission.
	Groups []BaseGroup `json:"groups,omitempty"`
}

func (a AllowedGroupsConfigInput) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(a, "", false)
}

func (a *AllowedGroupsConfigInput) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &a, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *AllowedGroupsConfigInput) GetType() *AllowedGroupsConfigInputAllowedGroupsConfigType {
	if o == nil {
		return nil
	}
	return o.Type
}

func (o *AllowedGroupsConfigInput) GetGroups() []BaseGroup {
	if o == nil {
		return nil
	}
	return o.Groups
}