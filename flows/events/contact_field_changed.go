package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeContactFieldChanged, func() flows.Event { return &ContactFieldChangedEvent{} })
}

// TypeContactFieldChanged is the type of our save to contact event
const TypeContactFieldChanged string = "contact_field_changed"

// ContactFieldChangedEvent events are created when a custom field value of the contact has been changed.
//
//   {
//     "type": "contact_field_changed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "field": {"key": "gender", "name": "Gender"},
//     "value": "Male"
//   }
//
// @event contact_field_changed
type ContactFieldChangedEvent struct {
	BaseEvent

	Field *assets.FieldReference `json:"field" validate:"required"`
	Value string                 `json:"value" validate:"required"`
}

// NewContactFieldChangedEvent returns a new save to contact event
func NewContactFieldChangedEvent(field *assets.FieldReference, value string) *ContactFieldChangedEvent {
	return &ContactFieldChangedEvent{
		BaseEvent: NewBaseEvent(),
		Field:     field,
		Value:     value,
	}
}

// Type returns the type of this event
func (e *ContactFieldChangedEvent) Type() string { return TypeContactFieldChanged }
