// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/utils"
	"time"
)

type Targets struct {
}

type Actor struct {
}

type EventMetadata struct {
}

// ActivityLog - API version of SIEMEvent
type ActivityLog struct {
	EventHash             string        `json:"event_hash"`
	EventType             string        `json:"event_type"`
	EventTypeUserFriendly string        `json:"event_type_user_friendly"`
	Outcome               string        `json:"outcome"`
	Targets               []Targets     `json:"targets"`
	Actor                 Actor         `json:"actor"`
	EventBeganAt          time.Time     `json:"event_began_at"`
	EventMetadata         EventMetadata `json:"event_metadata"`
}

func (a ActivityLog) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(a, "", false)
}

func (a *ActivityLog) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &a, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *ActivityLog) GetEventHash() string {
	if o == nil {
		return ""
	}
	return o.EventHash
}

func (o *ActivityLog) GetEventType() string {
	if o == nil {
		return ""
	}
	return o.EventType
}

func (o *ActivityLog) GetEventTypeUserFriendly() string {
	if o == nil {
		return ""
	}
	return o.EventTypeUserFriendly
}

func (o *ActivityLog) GetOutcome() string {
	if o == nil {
		return ""
	}
	return o.Outcome
}

func (o *ActivityLog) GetTargets() []Targets {
	if o == nil {
		return []Targets{}
	}
	return o.Targets
}

func (o *ActivityLog) GetActor() Actor {
	if o == nil {
		return Actor{}
	}
	return o.Actor
}

func (o *ActivityLog) GetEventBeganAt() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.EventBeganAt
}

func (o *ActivityLog) GetEventMetadata() EventMetadata {
	if o == nil {
		return EventMetadata{}
	}
	return o.EventMetadata
}
