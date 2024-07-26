// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

// DomainAppStatus - An enumeration.
type DomainAppStatus string

const (
	DomainAppStatusDiscovered  DomainAppStatus = "DISCOVERED"
	DomainAppStatusNeedsReview DomainAppStatus = "NEEDS_REVIEW"
	DomainAppStatusApproved    DomainAppStatus = "APPROVED"
	DomainAppStatusBlocklisted DomainAppStatus = "BLOCKLISTED"
	DomainAppStatusDeprecated  DomainAppStatus = "DEPRECATED"
)

func (e DomainAppStatus) ToPointer() *DomainAppStatus {
	return &e
}
func (e *DomainAppStatus) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "DISCOVERED":
		fallthrough
	case "NEEDS_REVIEW":
		fallthrough
	case "APPROVED":
		fallthrough
	case "BLOCKLISTED":
		fallthrough
	case "DEPRECATED":
		*e = DomainAppStatus(v)
		return nil
	default:
		return fmt.Errorf("invalid value for DomainAppStatus: %v", v)
	}
}
