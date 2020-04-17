package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeContactBlocked, func() flows.Event { return &ContactBlockedEvent{} })
}

// TypeContactBlocked is the type of our contact blocked events
const TypeContactBlocked string = "contact_blocked"

// ContactBlockedEvent events are created when the contact is blocked
//
//   {
//     "type": "contact_blocked",
//     "created_on": "2006-01-02T15:04:05Z"
//   }
//
// @event contact_blocked
type ContactBlockedEvent struct {
	baseEvent
}

// NewContactBlocked returns a new contact_blocked event
func NewContactBlocked() *ContactBlockedEvent {
	return &ContactBlockedEvent{baseEvent: newBaseEvent(TypeContactBlocked)}
}
