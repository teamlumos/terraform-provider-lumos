// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

// InlineWebhookType - An enumeration.
type InlineWebhookType string

const (
	InlineWebhookTypePreApproval       InlineWebhookType = "PRE_APPROVAL"
	InlineWebhookTypeProvision         InlineWebhookType = "PROVISION"
	InlineWebhookTypeDeprovision       InlineWebhookType = "DEPROVISION"
	InlineWebhookTypeRequestValidation InlineWebhookType = "REQUEST_VALIDATION"
	InlineWebhookTypeSiem              InlineWebhookType = "SIEM"
)

func (e InlineWebhookType) ToPointer() *InlineWebhookType {
	return &e
}
func (e *InlineWebhookType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "PRE_APPROVAL":
		fallthrough
	case "PROVISION":
		fallthrough
	case "DEPROVISION":
		fallthrough
	case "REQUEST_VALIDATION":
		fallthrough
	case "SIEM":
		*e = InlineWebhookType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for InlineWebhookType: %v", v)
	}
}
