// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

// GroupProvisioningOption - An enumeration.
type GroupProvisioningOption string

const (
	GroupProvisioningOptionDirectToUser     GroupProvisioningOption = "DIRECT_TO_USER"
	GroupProvisioningOptionGroupsAndHidden  GroupProvisioningOption = "GROUPS_AND_HIDDEN"
	GroupProvisioningOptionGroupsAndVisible GroupProvisioningOption = "GROUPS_AND_VISIBLE"
)

func (e GroupProvisioningOption) ToPointer() *GroupProvisioningOption {
	return &e
}
func (e *GroupProvisioningOption) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "DIRECT_TO_USER":
		fallthrough
	case "GROUPS_AND_HIDDEN":
		fallthrough
	case "GROUPS_AND_VISIBLE":
		*e = GroupProvisioningOption(v)
		return nil
	default:
		return fmt.Errorf("invalid value for GroupProvisioningOption: %v", v)
	}
}
