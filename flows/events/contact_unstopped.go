package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeContactUnstopped, func() flows.Event { return &ContactUnstoppedEvent{} })
}

// TypeContactUnstopped is the type of our contact_unstopped events
const TypeContactUnstopped string = "contact_unstopped"

// ContactUnstoppedEvent events are created when the contact is stopped
//
//   {
//     "type": "contact_unstopped",
//     "created_on": "2006-01-02T15:04:05Z"
//   }
//
// @event contact_unstopped
type ContactUnstoppedEvent struct {
	baseEvent
}

// NewContactUnstopped creates a new contact_unstopped event
func NewContactUnstopped() *ContactUnstoppedEvent {
	return &ContactUnstoppedEvent{baseEvent: newBaseEvent(TypeContactUnstopped)}
}
