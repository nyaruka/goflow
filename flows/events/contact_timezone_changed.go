package events

import (
	"time"

	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeContactTimezoneChanged, func() flows.Event { return &ContactTimezoneChangedEvent{} })
}

// TypeContactTimezoneChanged is the type of our contact timezone changed event
const TypeContactTimezoneChanged string = "contact_timezone_changed"

// ContactTimezoneChangedEvent events are created when the timezone of the contact has been changed.
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

	Timezone string `json:"timezone"`
}

// NewContactTimezoneChanged returns a new contact timezone changed event
func NewContactTimezoneChanged(timezone *time.Location) *ContactTimezoneChangedEvent {
	var tzname string
	if timezone != nil {
		tzname = timezone.String()
	}

	return &ContactTimezoneChangedEvent{
		BaseEvent: NewBaseEvent(TypeContactTimezoneChanged),
		Timezone:  tzname,
	}
}
