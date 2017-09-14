package events

import "github.com/nyaruka/goflow/flows"
import "encoding/json"

// TypeSetContact is the type of our set contact event
const TypeSetContact string = "set_contact"

// SetContactEvent events are created to set a contact on a session
//
// ```
//   {
//     "type": "set_contact",
//     "created_on": "2006-01-02T15:04:05Z",
//     "contact": {
//       "uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
//       "name": "Bob",
//       "urns": ["tel:+11231234567"]
//     }
//   }
// ```
//
// @event set_contact
type SetContactEvent struct {
	BaseEvent
	Contact json.RawMessage `json:"contact"`
}

// Type returns the type of this event
func (e *SetContactEvent) Type() string { return TypeSetContact }

// Apply applies this event to the given run
func (e *SetContactEvent) Apply(run flows.FlowRun) error {
	contact, err := flows.ReadContact(run.Session(), e.Contact)
	if err != nil {
		return err
	}

	run.SetContact(contact)
	run.Session().SetContact(contact)
	return nil
}
