package events

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeContactChanged, func() flows.Event { return &ContactChangedEvent{} })
}

// TypeContactChanged is the type of our set contact event
const TypeContactChanged string = "contact_changed"

// ContactChangedEvent events are sent by the caller to tell the engine to update the session contact.
//
//   {
//     "type": "contact_changed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "contact": {
//       "uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
//       "name": "Bob",
//       "urns": ["tel:+11231234567"]
//     }
//   }
//
// @event contact_changed
type ContactChangedEvent struct {
	BaseEvent

	Contact json.RawMessage `json:"contact"`
}

// NewContactChangedEvent creates a new contact changed event
func NewContactChangedEvent(contact *flows.Contact) *ContactChangedEvent {
	marshalled, _ := json.Marshal(contact)
	return &ContactChangedEvent{
		BaseEvent: NewBaseEvent(TypeContactChanged),
		Contact:   marshalled,
	}
}

var _ flows.Event = (*ContactChangedEvent)(nil)
