// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

type ManagerApprovalOption string

const (
	ManagerApprovalOptionNone            ManagerApprovalOption = "NONE"
	ManagerApprovalOptionInitialApproval ManagerApprovalOption = "INITIAL_APPROVAL"
)

func (e ManagerApprovalOption) ToPointer() *ManagerApprovalOption {
	return &e
}
func (e *ManagerApprovalOption) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "NONE":
		fallthrough
	case "INITIAL_APPROVAL":
		*e = ManagerApprovalOption(v)
		return nil
	default:
		return fmt.Errorf("invalid value for ManagerApprovalOption: %v", v)
	}
}
