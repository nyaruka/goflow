package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeContactIsBlockedChanged, func() flows.Event { return &ContactIsBlockedChangedEvent{} })
}

// TypeContactIsBlockedChanged is the type of our contact is_blocked changed event
const TypeContactIsBlockedChanged string = "contact_is_blocked_changed"

// ContactIsBlockedChangedEvent events are created when the is_blocked of the contact has been changed.
//
//   {
//     "type": "contact_is_blocked_changed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "is_blocked": false
//   }
//
// @event contact_is_blocked_changed
type ContactIsBlockedChangedEvent struct {
	baseEvent

	IsBlocked bool `json:"is_blocked"`
}

// NewContactIsBlockedChanged returns a new contact is_blocked changed event
func NewContactIsBlockedChanged(isBlocked bool) *ContactIsBlockedChangedEvent {
	return &ContactIsBlockedChangedEvent{
		baseEvent: newBaseEvent(TypeContactIsBlockedChanged),
		IsBlocked: isBlocked,
	}
}
