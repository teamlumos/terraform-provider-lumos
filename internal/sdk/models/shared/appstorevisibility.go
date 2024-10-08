// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

type AppStoreVisibility string

const (
	AppStoreVisibilityFull    AppStoreVisibility = "FULL"
	AppStoreVisibilityLimited AppStoreVisibility = "LIMITED"
	AppStoreVisibilityNone    AppStoreVisibility = "NONE"
)

func (e AppStoreVisibility) ToPointer() *AppStoreVisibility {
	return &e
}
func (e *AppStoreVisibility) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "FULL":
		fallthrough
	case "LIMITED":
		fallthrough
	case "NONE":
		*e = AppStoreVisibility(v)
		return nil
	default:
		return fmt.Errorf("invalid value for AppStoreVisibility: %v", v)
	}
}
