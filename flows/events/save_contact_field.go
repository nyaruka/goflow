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
//     "field": {"key": "gender", "label": "Gender"},
//     "value": "Male"
//   }
// ```
//
// @event save_contact_field
type SaveContactFieldEvent struct {
	BaseEvent
	Field *flows.FieldReference `json:"field" validate:"required"`
	Value string                `json:"value" validate:"required"`
}

// NewSaveToContactEvent returns a new save to contact event
func NewSaveToContactEvent(field *flows.FieldReference, value string) *SaveContactFieldEvent {
	return &SaveContactFieldEvent{
		BaseEvent: NewBaseEvent(),
		Field:     field,
		Value:     value,
	}
}

// Type returns the type of this event
func (e *SaveContactFieldEvent) Type() string { return TypeSaveContactField }

// Apply applies this event to the given run
func (e *SaveContactFieldEvent) Apply(run flows.FlowRun) error {
	field, err := run.Session().Assets().GetField(e.Field.Key)
	if err != nil {
		return err
	}

	run.Contact().Fields().Save(run.Environment(), field, e.Value)

	return run.Contact().UpdateDynamicGroups(run.Session())
}
