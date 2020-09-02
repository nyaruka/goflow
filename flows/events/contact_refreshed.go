package events

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeContactRefreshed, func() flows.Event { return &ContactRefreshedEvent{} })
}

// TypeContactRefreshed is the type of our contact refreshed event
const TypeContactRefreshed string = "contact_refreshed"

// ContactRefreshedEvent events are generated when the resume has a contact with differences to the current session contact.
//
//   {
//     "type": "contact_refreshed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "contact": {
//       "uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
//       "name": "Bob",
//       "urns": ["tel:+11231234567"]
//     }
//   }
//
// @event contact_refreshed
type ContactRefreshedEvent struct {
	baseEvent

	Contact json.RawMessage `json:"contact"`
}

// NewContactRefreshed creates a new contact changed event
func NewContactRefreshed(contact *flows.Contact) *ContactRefreshedEvent {
	marshalled, _ := jsonx.Marshal(contact)
	return &ContactRefreshedEvent{
		baseEvent: newBaseEvent(TypeContactRefreshed),
		Contact:   marshalled,
	}
}

var _ flows.Event = (*ContactRefreshedEvent)(nil)
