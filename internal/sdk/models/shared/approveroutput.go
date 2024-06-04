// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/utils"
)

// ApproverOutputUser - Optionally, the approver can be a user.
type ApproverOutputUser struct {
	// The ID of this user.
	ID string `json:"id"`
	// The email of this user.
	Email *string `json:"email,omitempty"`
	// The given name of this user.
	GivenName *string `json:"given_name,omitempty"`
	// The family name of this user.
	FamilyName *string `json:"family_name,omitempty"`
	// The status of this user.
	Status *UserLifecycleStatus `json:"status,omitempty"`
}

func (o *ApproverOutputUser) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *ApproverOutputUser) GetEmail() *string {
	if o == nil {
		return nil
	}
	return o.Email
}

func (o *ApproverOutputUser) GetGivenName() *string {
	if o == nil {
		return nil
	}
	return o.GivenName
}

func (o *ApproverOutputUser) GetFamilyName() *string {
	if o == nil {
		return nil
	}
	return o.FamilyName
}

func (o *ApproverOutputUser) GetStatus() *UserLifecycleStatus {
	if o == nil {
		return nil
	}
	return o.Status
}

// ApproverOutputLifecycle - The lifecycle of this group.
type ApproverOutputLifecycle string

const (
	ApproverOutputLifecycleSynced ApproverOutputLifecycle = "SYNCED"
	ApproverOutputLifecycleNative ApproverOutputLifecycle = "NATIVE"
)

func (e ApproverOutputLifecycle) ToPointer() *ApproverOutputLifecycle {
	return &e
}
func (e *ApproverOutputLifecycle) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "SYNCED":
		fallthrough
	case "NATIVE":
		*e = ApproverOutputLifecycle(v)
		return nil
	default:
		return fmt.Errorf("invalid value for ApproverOutputLifecycle: %v", v)
	}
}

// ApproverOutputGroup - Optionally, the approver can be a group.
type ApproverOutputGroup struct {
	// The ID of this group.
	ID *string `json:"id,omitempty"`
	// The ID of the app that owns this group.
	AppID *string `json:"app_id,omitempty"`
	// The ID of this group, specific to the integration.
	IntegrationSpecificID *string `json:"integration_specific_id,omitempty"`
	// The name of this group.
	Name *string `json:"name,omitempty"`
	// The description of this group.
	Description *string `json:"description,omitempty"`
	// The lifecycle of this group.
	GroupLifecycle *ApproverOutputLifecycle `default:"SYNCED" json:"group_lifecycle"`
	// The ID of the app that owns this group.
	SourceAppID *string `json:"source_app_id,omitempty"`
}

func (a ApproverOutputGroup) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(a, "", false)
}

func (a *ApproverOutputGroup) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &a, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *ApproverOutputGroup) GetID() *string {
	if o == nil {
		return nil
	}
	return o.ID
}

func (o *ApproverOutputGroup) GetAppID() *string {
	if o == nil {
		return nil
	}
	return o.AppID
}

func (o *ApproverOutputGroup) GetIntegrationSpecificID() *string {
	if o == nil {
		return nil
	}
	return o.IntegrationSpecificID
}

func (o *ApproverOutputGroup) GetName() *string {
	if o == nil {
		return nil
	}
	return o.Name
}

func (o *ApproverOutputGroup) GetDescription() *string {
	if o == nil {
		return nil
	}
	return o.Description
}

func (o *ApproverOutputGroup) GetGroupLifecycle() *ApproverOutputLifecycle {
	if o == nil {
		return nil
	}
	return o.GroupLifecycle
}

func (o *ApproverOutputGroup) GetSourceAppID() *string {
	if o == nil {
		return nil
	}
	return o.SourceAppID
}

type ApproverOutput struct {
	// The type of this approver.
	Type ApproverType `json:"type"`
	// Optionally, the approver can be a user.
	User *ApproverOutputUser `json:"user,omitempty"`
	// Optionally, the approver can be a group.
	Group *ApproverOutputGroup `json:"group,omitempty"`
}

func (o *ApproverOutput) GetType() ApproverType {
	if o == nil {
		return ApproverType("")
	}
	return o.Type
}

func (o *ApproverOutput) GetUser() *ApproverOutputUser {
	if o == nil {
		return nil
	}
	return o.User
}

func (o *ApproverOutput) GetGroup() *ApproverOutputGroup {
	if o == nil {
		return nil
	}
	return o.Group
}
