// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

// SupportRequestStatus - An enumeration.
type SupportRequestStatus string

const (
	SupportRequestStatusPending                     SupportRequestStatus = "PENDING"
	SupportRequestStatusPendingManagerApproval      SupportRequestStatus = "PENDING_MANAGER_APPROVAL"
	SupportRequestStatusManagerApproved             SupportRequestStatus = "MANAGER_APPROVED"
	SupportRequestStatusManagerDenied               SupportRequestStatus = "MANAGER_DENIED"
	SupportRequestStatusPendingApproval             SupportRequestStatus = "PENDING_APPROVAL"
	SupportRequestStatusApproved                    SupportRequestStatus = "APPROVED"
	SupportRequestStatusDenied                      SupportRequestStatus = "DENIED"
	SupportRequestStatusExpired                     SupportRequestStatus = "EXPIRED"
	SupportRequestStatusCancelled                   SupportRequestStatus = "CANCELLED"
	SupportRequestStatusPendingProvisioning         SupportRequestStatus = "PENDING_PROVISIONING"
	SupportRequestStatusPendingManualProvisioning   SupportRequestStatus = "PENDING_MANUAL_PROVISIONING"
	SupportRequestStatusDeniedProvisioning          SupportRequestStatus = "DENIED_PROVISIONING"
	SupportRequestStatusProvisioned                 SupportRequestStatus = "PROVISIONED"
	SupportRequestStatusPendingManualDeprovisioning SupportRequestStatus = "PENDING_MANUAL_DEPROVISIONING"
	SupportRequestStatusTimeBasedExpired            SupportRequestStatus = "TIME_BASED_EXPIRED"
	SupportRequestStatusCompleted                   SupportRequestStatus = "COMPLETED"
	SupportRequestStatusReverting                   SupportRequestStatus = "REVERTING"
	SupportRequestStatusReverted                    SupportRequestStatus = "REVERTED"
)

func (e SupportRequestStatus) ToPointer() *SupportRequestStatus {
	return &e
}
func (e *SupportRequestStatus) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "PENDING":
		fallthrough
	case "PENDING_MANAGER_APPROVAL":
		fallthrough
	case "MANAGER_APPROVED":
		fallthrough
	case "MANAGER_DENIED":
		fallthrough
	case "PENDING_APPROVAL":
		fallthrough
	case "APPROVED":
		fallthrough
	case "DENIED":
		fallthrough
	case "EXPIRED":
		fallthrough
	case "CANCELLED":
		fallthrough
	case "PENDING_PROVISIONING":
		fallthrough
	case "PENDING_MANUAL_PROVISIONING":
		fallthrough
	case "DENIED_PROVISIONING":
		fallthrough
	case "PROVISIONED":
		fallthrough
	case "PENDING_MANUAL_DEPROVISIONING":
		fallthrough
	case "TIME_BASED_EXPIRED":
		fallthrough
	case "COMPLETED":
		fallthrough
	case "REVERTING":
		fallthrough
	case "REVERTED":
		*e = SupportRequestStatus(v)
		return nil
	default:
		return fmt.Errorf("invalid value for SupportRequestStatus: %v", v)
	}
}
