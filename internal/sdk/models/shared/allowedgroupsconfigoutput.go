// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/utils"
)

// AllowedGroupsConfigType - The type of this allowed groups config, can be all groups or specific.
type AllowedGroupsConfigType string

const (
	AllowedGroupsConfigTypeAllGroups       AllowedGroupsConfigType = "ALL_GROUPS"
	AllowedGroupsConfigTypeSpecifiedGroups AllowedGroupsConfigType = "SPECIFIED_GROUPS"
)

func (e AllowedGroupsConfigType) ToPointer() *AllowedGroupsConfigType {
	return &e
}
func (e *AllowedGroupsConfigType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "ALL_GROUPS":
		fallthrough
	case "SPECIFIED_GROUPS":
		*e = AllowedGroupsConfigType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for AllowedGroupsConfigType: %v", v)
	}
}

type AllowedGroupsConfigOutput struct {
	// The type of this allowed groups config, can be all groups or specific.
	Type *AllowedGroupsConfigType `default:"ALL_GROUPS" json:"type"`
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
