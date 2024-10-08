// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

type VendorAgreementCustomAttributeType string

const (
	VendorAgreementCustomAttributeTypeText         VendorAgreementCustomAttributeType = "TEXT"
	VendorAgreementCustomAttributeTypeUser         VendorAgreementCustomAttributeType = "USER"
	VendorAgreementCustomAttributeTypeSingleSelect VendorAgreementCustomAttributeType = "SINGLE_SELECT"
)

func (e VendorAgreementCustomAttributeType) ToPointer() *VendorAgreementCustomAttributeType {
	return &e
}
func (e *VendorAgreementCustomAttributeType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "TEXT":
		fallthrough
	case "USER":
		fallthrough
	case "SINGLE_SELECT":
		*e = VendorAgreementCustomAttributeType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for VendorAgreementCustomAttributeType: %v", v)
	}
}
