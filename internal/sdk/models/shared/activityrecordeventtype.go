// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

// ActivityRecordEventType - An enumeration.
type ActivityRecordEventType string

const (
	ActivityRecordEventTypeLogin    ActivityRecordEventType = "LOGIN"
	ActivityRecordEventTypeActivity ActivityRecordEventType = "ACTIVITY"
)

func (e ActivityRecordEventType) ToPointer() *ActivityRecordEventType {
	return &e
}
func (e *ActivityRecordEventType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "LOGIN":
		fallthrough
	case "ACTIVITY":
		*e = ActivityRecordEventType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for ActivityRecordEventType: %v", v)
	}
}
