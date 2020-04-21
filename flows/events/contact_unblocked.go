package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeContactUnblocked, func() flows.Event { return &ContactUnblockedEvent{} })
}

// TypeContactUnblocked is the type of our contact unblocked events
const TypeContactUnblocked string = "contact_unblocked"

// ContactUnblockedEvent events are created when the contact is unblocked
//
//   {
//     "type": "contact_unblocked",
//     "created_on": "2006-01-02T15:04:05Z"
//   }
//
// @event contact_unblocked
type ContactUnblockedEvent struct {
	baseEvent
}

// NewContactUnblocked returns a new contact_unblocked event
func NewContactUnblocked() *ContactUnblockedEvent {
	return &ContactUnblockedEvent{baseEvent: newBaseEvent(TypeContactUnblocked)}
}
