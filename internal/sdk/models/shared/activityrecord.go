// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/utils"
	"time"
)

// ActivityRecordAccount - Metadata that Lumos can use to match the activity record to a software account within Lumos.
type ActivityRecordAccount struct {
	// The external app's user ID for the account.
	ExternalID *string `json:"external_id,omitempty"`
	// The email associated with the account
	Email *string `json:"email,omitempty"`
}

func (o *ActivityRecordAccount) GetExternalID() *string {
	if o == nil {
		return nil
	}
	return o.ExternalID
}

func (o *ActivityRecordAccount) GetEmail() *string {
	if o == nil {
		return nil
	}
	return o.Email
}

// Event - Metadata about the event being uploaded.
type Event struct {
	// The type of event being uploaded.
	Type ActivityRecordEventType `json:"type"`
}

func (o *Event) GetType() ActivityRecordEventType {
	if o == nil {
		return ActivityRecordEventType("")
	}
	return o.Type
}

// ActivityRecordApp - Metadata that Lumos can use to match the activity record to an application within Lumos.
type ActivityRecordApp struct {
	// The ID of the app as it is identified in the source E.g. the ID that Okta uses to identify the app
	InstanceIdentifier string `json:"instance_identifier"`
}

func (o *ActivityRecordApp) GetInstanceIdentifier() string {
	if o == nil {
		return ""
	}
	return o.InstanceIdentifier
}

type ActivityRecord struct {
	// Metadata that Lumos can use to match the activity record to a software account within Lumos.
	Account ActivityRecordAccount `json:"account"`
	// Metadata about the event being uploaded.
	Event Event `json:"event"`
	// The timestamp of this event, in ISO 8601 format.
	Timestamp time.Time `json:"timestamp"`
	// UUID of the application in Lumos where this activity record was sourced (e.g. the ID of Okta within Lumos found by going to Apps > Find your app in the list > Click '...' > Copy Stable Identifier)
	SourceAppID string `json:"source_app_id"`
	// Metadata that Lumos can use to match the activity record to an application within Lumos.
	App *ActivityRecordApp `json:"app,omitempty"`
}

func (a ActivityRecord) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(a, "", false)
}

func (a *ActivityRecord) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &a, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *ActivityRecord) GetAccount() ActivityRecordAccount {
	if o == nil {
		return ActivityRecordAccount{}
	}
	return o.Account
}

func (o *ActivityRecord) GetEvent() Event {
	if o == nil {
		return Event{}
	}
	return o.Event
}

func (o *ActivityRecord) GetTimestamp() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.Timestamp
}

func (o *ActivityRecord) GetSourceAppID() string {
	if o == nil {
		return ""
	}
	return o.SourceAppID
}

func (o *ActivityRecord) GetApp() *ActivityRecordApp {
	if o == nil {
		return nil
	}
	return o.App
}
