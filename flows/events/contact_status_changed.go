package events

import "github.com/nyaruka/goflow/flows"

func init() {
	registerType(TypeContactStatusChanged, func() flows.Event { return &ContactStatusChangedEvent{} })
}

// TypeContactStatusChanged is the type of our contact status changed event
const TypeContactStatusChanged string = "contact_status_changed"

// ContactStatusChangedEvent events are created when the status of the contact has been changed.
//
//   {
//     "type": "contact_timezone_changed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "status": "blocked"
//   }
//
// @event contact_status_changed
type ContactStatusChangedEvent struct {
	baseEvent

	Status flows.ContactStatus `json:"status"`
}

// NewContactStatusChanged returns a new contact_status_changed event
func NewContactStatusChanged(status flows.ContactStatus) *ContactStatusChangedEvent {
	return &ContactStatusChangedEvent{
		baseEvent: newBaseEvent(TypeContactStatusChanged),
		Status:    status,
	}
}
