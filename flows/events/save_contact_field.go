package events

import "github.com/nyaruka/goflow/flows"

// TypeSaveContactField is the type of our save to contact event
const TypeSaveContactField string = "save_contact_field"

// SaveContactFieldEvent events are created when a contact field is updated.
//
// ```
//   {
//    "step_uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//    "type": "save_contact_field",
//    "created_on": "2006-01-02T15:04:05Z",
//    "field_uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
//    "field_name": "Gender",
//    "value": "Male"
//   }
// ```
//
// @event save_contact_field
type SaveContactFieldEvent struct {
	BaseEvent
	FieldUUID flows.FieldUUID `json:"field_uuid"  validate:"required"`
	FieldName string          `json:"field_name"  validate:"required"`
	Value     string          `json:"value"       validate:"required"`
}

// NewSaveToContact returns a new save to contact event
func NewSaveToContact(field flows.FieldUUID, name string, value string) *SaveContactFieldEvent {
	return &SaveContactFieldEvent{
		FieldUUID: field,
		FieldName: name,
		Value:     value,
	}
}

// Type returns the type of this event
func (e *SaveContactFieldEvent) Type() string { return TypeSaveContactField }
