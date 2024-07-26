// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

// OrderCustomAttributeType - An enumeration.
type OrderCustomAttributeType string

const (
	OrderCustomAttributeTypeText OrderCustomAttributeType = "TEXT"
	OrderCustomAttributeTypeUser OrderCustomAttributeType = "USER"
)

func (e OrderCustomAttributeType) ToPointer() *OrderCustomAttributeType {
	return &e
}
func (e *OrderCustomAttributeType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "TEXT":
		fallthrough
	case "USER":
		*e = OrderCustomAttributeType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for OrderCustomAttributeType: %v", v)
	}
}
