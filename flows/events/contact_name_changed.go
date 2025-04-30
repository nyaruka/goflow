package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeContactNameChanged, func() flows.Event { return &ContactNameChangedEvent{} })
}

// TypeContactNameChanged is the type of our contact name changed event
const TypeContactNameChanged string = "contact_name_changed"

// ContactNameChangedEvent events are created when the name of the contact has been changed.
//
//	{
//	  "uuid": "019688A6-41d2-7366-958a-630e35c62431",
//	  "type": "contact_name_changed",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "name": "Bob Smith"
//	}
//
// @event contact_name_changed
type ContactNameChangedEvent struct {
	BaseEvent

	Name string `json:"name"`
}

// NewContactNameChanged returns a new contact name changed event
func NewContactNameChanged(name string) *ContactNameChangedEvent {
	return &ContactNameChangedEvent{
		BaseEvent: NewBaseEvent(TypeContactNameChanged),
		Name:      name,
	}
}
