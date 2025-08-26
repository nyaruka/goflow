package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeContactFieldChanged, func() flows.Event { return &ContactFieldChanged{} })
}

// TypeContactFieldChanged is the type of our save to contact event
const TypeContactFieldChanged string = "contact_field_changed"

// ContactFieldChanged events are created when a field value of the contact has been changed. An empty value indicates
// that the field value has been cleared.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "contact_field_changed",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "field": {"key": "gender", "name": "Gender"},
//	  "value": {"text": "Male"}
//	}
//
// @event contact_field_changed
type ContactFieldChanged struct {
	BaseEvent

	Field *assets.FieldReference `json:"field" validate:"required"`
	Value *flows.Value           `json:"value"`
}

// NewContactFieldChanged returns a new save to contact event
func NewContactFieldChanged(field *flows.Field, value *flows.Value) *ContactFieldChanged {
	return &ContactFieldChanged{
		BaseEvent: NewBaseEvent(TypeContactFieldChanged),
		Field:     field.Reference(),
		Value:     value,
	}
}
