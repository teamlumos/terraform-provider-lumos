// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

type AppStoreVisibilityOption string

const (
	AppStoreVisibilityOptionHidden  AppStoreVisibilityOption = "HIDDEN"
	AppStoreVisibilityOptionVisible AppStoreVisibilityOption = "VISIBLE"
)

func (e AppStoreVisibilityOption) ToPointer() *AppStoreVisibilityOption {
	return &e
}
func (e *AppStoreVisibilityOption) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "HIDDEN":
		fallthrough
	case "VISIBLE":
		*e = AppStoreVisibilityOption(v)
		return nil
	default:
		return fmt.Errorf("invalid value for AppStoreVisibilityOption: %v", v)
	}
}
