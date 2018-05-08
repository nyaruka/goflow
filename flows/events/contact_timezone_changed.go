package events

import (
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
)

// TypeContactTimezoneChanged is the type of our contact timezone changed event
const TypeContactTimezoneChanged string = "contact_timezone_changed"

// ContactTimezoneChangedEvent events are created when a timezone of a contact has been changed
//
//   {
//     "type": "contact_timezone_changed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "timezone": "Africa/Kigali"
//   }
//
// @event contact_timezone_changed
type ContactTimezoneChangedEvent struct {
	baseEvent
	callerOrEngineEvent

	Timezone string `json:"timezone"`
}

// NewContactTimezoneChangedEvent returns a new contact timezone changed event
func NewContactTimezoneChangedEvent(timezone string) *ContactTimezoneChangedEvent {
	return &ContactTimezoneChangedEvent{
		baseEvent: newBaseEvent(),
		Timezone:  timezone,
	}
}

// Type returns the type of this event
func (e *ContactTimezoneChangedEvent) Type() string { return TypeContactTimezoneChanged }

// Validate validates our event is valid and has all the assets it needs
func (e *ContactTimezoneChangedEvent) Validate(assets flows.SessionAssets) error {
	return nil
}

// Apply applies this event to the given run
func (e *ContactTimezoneChangedEvent) Apply(run flows.FlowRun) error {
	if run.Contact() == nil {
		return fmt.Errorf("can't apply event in session without a contact")
	}

	if e.Timezone != "" {
		timezone, err := time.LoadLocation(e.Timezone)
		if err != nil {
			return err
		}
		run.Contact().SetTimezone(timezone)
	} else {
		run.Contact().SetTimezone(nil)
	}

	return nil
}
