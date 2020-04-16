package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeContactStopped, func() flows.Event { return &ContactStoppedEvent{} })
}

// TypeContactStopped is the type of our contact_stopped events
const TypeContactStopped string = "contact_stopped"

// ContactStoppedEvent events are created when the contact is stopped
//
//   {
//     "type": "contact_stopped",
//     "created_on": "2006-01-02T15:04:05Z"
//   }
//
// @event contact_stopped
type ContactStoppedEvent struct {
	baseEvent
}

// NewContactStopped creates a new contact_stopped event
func NewContactStopped() *ContactStoppedEvent {
	return &ContactStoppedEvent{baseEvent: newBaseEvent(TypeContactStopped)}
}
