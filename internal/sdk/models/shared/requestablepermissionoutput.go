// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

type TypePermissionType string

const (
	TypePermissionTypeSynced TypePermissionType = "SYNCED"
	TypePermissionTypeNative TypePermissionType = "NATIVE"
)

func (e TypePermissionType) ToPointer() *TypePermissionType {
	return &e
}
func (e *TypePermissionType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "SYNCED":
		fallthrough
	case "NATIVE":
		*e = TypePermissionType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for TypePermissionType: %v", v)
	}
}

type RequestablePermissionOutput struct {
	// The ID of this requestable permission.
	ID *string `json:"id,omitempty"`
	// The type of this requestable permission.
	Type *TypePermissionType `json:"type,omitempty"`
	// The label of this requestable permission.
	Label string `json:"label"`
	// The ID of the app associated with this requestable permission.
	AppID string `json:"app_id"`
	// The non-unique ID of the service associated with this requestable permission. Depending on how it is sourced in Lumos, this may be the app's name, website,  or other identifier.
	AppClassID string `json:"app_class_id"`
	// The ID of the instance associated with this requestable permission. This may be an empty string.
	AppInstanceID string `json:"app_instance_id"`
	// The request config associated with this requestable permission.
	RequestConfig RequestConfigOutput `json:"request_config"`
}

func (o *RequestablePermissionOutput) GetID() *string {
	if o == nil {
		return nil
	}
	return o.ID
}

func (o *RequestablePermissionOutput) GetType() *TypePermissionType {
	if o == nil {
		return nil
	}
	return o.Type
}

func (o *RequestablePermissionOutput) GetLabel() string {
	if o == nil {
		return ""
	}
	return o.Label
}

func (o *RequestablePermissionOutput) GetAppID() string {
	if o == nil {
		return ""
	}
	return o.AppID
}

func (o *RequestablePermissionOutput) GetAppClassID() string {
	if o == nil {
		return ""
	}
	return o.AppClassID
}

func (o *RequestablePermissionOutput) GetAppInstanceID() string {
	if o == nil {
		return ""
	}
	return o.AppInstanceID
}

func (o *RequestablePermissionOutput) GetRequestConfig() RequestConfigOutput {
	if o == nil {
		return RequestConfigOutput{}
	}
	return o.RequestConfig
}
