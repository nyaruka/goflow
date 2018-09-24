package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeContactTimezoneChanged, func() flows.Event { return &ContactTimezoneChangedEvent{} })
}

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
	BaseEvent
	callerOrEngineEvent

	Timezone string `json:"timezone"`
}

// NewContactTimezoneChangedEvent returns a new contact timezone changed event
func NewContactTimezoneChangedEvent(timezone string) *ContactTimezoneChangedEvent {
	return &ContactTimezoneChangedEvent{
		BaseEvent: NewBaseEvent(),
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
	return nil
}
