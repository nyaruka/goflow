package events

import "github.com/nyaruka/goflow/flows"

// TypeUpdateContact is the type of our update contact event
const TypeUpdateContact string = "update_contact"

// UpdateContactEvent events are created when a contact's built in field is updated.
//
// ```
//   {
//     "type": "update_contact",
//     "created_on": "2006-01-02T15:04:05Z",
//     "field_name": "Language",
//     "value": "eng"
//   }
// ```
//
// @event update_contact
type UpdateContactEvent struct {
	BaseEvent
	FieldName string `json:"field_name"  validate:"required"`
	Value     string `json:"value"`
}

// NewUpdateContact returns a new save to contact event
func NewUpdateContact(name string, value string) *UpdateContactEvent {
	return &UpdateContactEvent{
		BaseEvent: NewBaseEvent(),
		FieldName: name,
		Value:     value,
	}
}

// Type returns the type of this event
func (e *UpdateContactEvent) Type() string { return TypeUpdateContact }

// Apply applies this event to the given run
func (e *UpdateContactEvent) Apply(run flows.FlowRun) error { return nil }
