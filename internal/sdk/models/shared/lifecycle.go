// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

type Lifecycle string

const (
	LifecycleSynced Lifecycle = "SYNCED"
	LifecycleNative Lifecycle = "NATIVE"
)

func (e Lifecycle) ToPointer() *Lifecycle {
	return &e
}
func (e *Lifecycle) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "SYNCED":
		fallthrough
	case "NATIVE":
		*e = Lifecycle(v)
		return nil
	default:
		return fmt.Errorf("invalid value for Lifecycle: %v", v)
	}
}
