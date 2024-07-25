// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

// UserLifecycleStatus - An enumeration.
type UserLifecycleStatus string

const (
	UserLifecycleStatusStaged    UserLifecycleStatus = "STAGED"
	UserLifecycleStatusActive    UserLifecycleStatus = "ACTIVE"
	UserLifecycleStatusSuspended UserLifecycleStatus = "SUSPENDED"
	UserLifecycleStatusInactive  UserLifecycleStatus = "INACTIVE"
)

func (e UserLifecycleStatus) ToPointer() *UserLifecycleStatus {
	return &e
}
func (e *UserLifecycleStatus) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "STAGED":
		fallthrough
	case "ACTIVE":
		fallthrough
	case "SUSPENDED":
		fallthrough
	case "INACTIVE":
		*e = UserLifecycleStatus(v)
		return nil
	default:
		return fmt.Errorf("invalid value for UserLifecycleStatus: %v", v)
	}
}
