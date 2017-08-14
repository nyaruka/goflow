package events

import "github.com/nyaruka/goflow/flows"

// TypeSaveContactField is the type of our save to contact event
const TypeSaveContactField string = "save_contact_field"

// SaveContactFieldEvent events are created when a contact field is updated.
//
// ```
//   {
//     "type": "save_contact_field",
//     "created_on": "2006-01-02T15:04:05Z",
//     "field_uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
//     "field_name": "Gender",
//     "value": "Male"
//   }
// ```
//
// @event save_contact_field
type SaveContactFieldEvent struct {
	BaseEvent
	FieldUUID flows.FieldUUID `json:"field_uuid"  validate:"required"`
	FieldName string          `json:"field_name"  validate:"required"`
	Value     string          `json:"value"`
}

// NewSaveToContact returns a new save to contact event
func NewSaveToContact(field flows.FieldUUID, name string, value string) *SaveContactFieldEvent {
	return &SaveContactFieldEvent{
		BaseEvent: NewBaseEvent(),
		FieldUUID: field,
		FieldName: name,
		Value:     value,
	}
}

// Type returns the type of this event
func (e *SaveContactFieldEvent) Type() string { return TypeSaveContactField }

// Apply applies this event to the given run
func (e *SaveContactFieldEvent) Apply(run flows.FlowRun) error { return nil }
