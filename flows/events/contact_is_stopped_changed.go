package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeContactIsStoppedChanged, func() flows.Event { return &ContactIsStoppedChangedEvent{} })
}

// TypeContactIsStoppedChanged is the type of our contact is_stopped changed event
const TypeContactIsStoppedChanged string = "contact_is_stopped_changed"

// ContactIsStoppedChangedEvent events are created when the is_stopped of the contact has been changed.
//
//   {
//     "type": "contact_is_stopped_changed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "is_stopped": false
//   }
//
// @event contact_is_stopped_changed
type ContactIsStoppedChangedEvent struct {
	baseEvent

	IsStopped bool `json:"is_stopped"`
}

// NewContactIsStoppedChanged returns a new contact is_stopped changed event
func NewContactIsStoppedChanged(isStopped bool) *ContactIsStoppedChangedEvent {
	return &ContactIsStoppedChangedEvent{
		baseEvent: newBaseEvent(TypeContactIsStoppedChanged),
		IsStopped: isStopped,
	}
}
