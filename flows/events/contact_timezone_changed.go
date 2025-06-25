package events

import (
	"time"

	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeContactTimezoneChanged, func() flows.Event { return &ContactTimezoneChanged{} })
}

// TypeContactTimezoneChanged is the type of our contact timezone changed event
const TypeContactTimezoneChanged string = "contact_timezone_changed"

// ContactTimezoneChanged events are created when the timezone of the contact has been changed.
//
//	{
//	  "type": "contact_timezone_changed",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "timezone": "Africa/Kigali"
//	}
//
// @event contact_timezone_changed
type ContactTimezoneChanged struct {
	BaseEvent

	Timezone string `json:"timezone"`
}

// NewContactTimezoneChanged returns a new contact timezone changed event
func NewContactTimezoneChanged(timezone *time.Location) *ContactTimezoneChanged {
	var tzname string
	if timezone != nil {
		tzname = timezone.String()
	}

	return &ContactTimezoneChanged{
		BaseEvent: NewBaseEvent(TypeContactTimezoneChanged),
		Timezone:  tzname,
	}
}
