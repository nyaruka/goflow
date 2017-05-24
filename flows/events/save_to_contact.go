package events

import "github.com/nyaruka/goflow/flows"

// TypeSaveToContact is the type of our save to contact event
const TypeSaveToContact string = "save_to_contact"

// SaveToContactEvent events are created when a contact field is updated.
//
// ```
//   {
//    "step": "8eebd020-1af5-431c-b943-aa670fc74da9",
//    "created_on": "2006-01-02T15:04:05Z",
//    "type": "save_to_contact",
//    "field": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
//    "name": "Gender",
//    "value": "Male"
//   }
// ```
//
// @event save_to_contact
type SaveToContactEvent struct {
	BaseEvent
	Field flows.FieldUUID `json:"field"  validate:"required"`
	Name  string          `json:"name"   validate:"required"`
	Value string          `json:"value"  validate:"required"`
}

// NewSaveToContact returns a new save to contact event
func NewSaveToContact(field flows.FieldUUID, name string, value string) *SaveToContactEvent {
	return &SaveToContactEvent{
		Field: field,
		Name:  name,
		Value: value,
	}
}

// Type returns the type of this event
func (e *SaveToContactEvent) Type() string { return TypeSaveToContact }
