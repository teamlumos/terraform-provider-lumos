// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

type FlowState string

const (
	FlowStateSuccess    FlowState = "SUCCESS"
	FlowStateFailure    FlowState = "FAILURE"
	FlowStateRunning    FlowState = "RUNNING"
	FlowStateNotStarted FlowState = "NOT_STARTED"
)

func (e FlowState) ToPointer() *FlowState {
	return &e
}
func (e *FlowState) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "SUCCESS":
		fallthrough
	case "FAILURE":
		fallthrough
	case "RUNNING":
		fallthrough
	case "NOT_STARTED":
		*e = FlowState(v)
		return nil
	default:
		return fmt.Errorf("invalid value for FlowState: %v", v)
	}
}
