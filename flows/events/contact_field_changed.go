package events

import "github.com/nyaruka/goflow/flows"

// TypeContactFieldChanged is the type of our save to contact event
const TypeContactFieldChanged string = "contact_field_changed"

// ContactFieldChangedEvent events are created when a contact field is updated.
//
// ```
//   {
//     "type": "contact_field_changed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "field": {"key": "gender", "label": "Gender"},
//     "value": "Male"
//   }
// ```
//
// @event contact_field_changed
type ContactFieldChangedEvent struct {
	BaseEvent
	Field *flows.FieldReference `json:"field" validate:"required"`
	Value string                `json:"value" validate:"required"`
}

// NewContactFieldChangedEvent returns a new save to contact event
func NewContactFieldChangedEvent(field *flows.FieldReference, value string) *ContactFieldChangedEvent {
	return &ContactFieldChangedEvent{
		BaseEvent: NewBaseEvent(),
		Field:     field,
		Value:     value,
	}
}

// Type returns the type of this event
func (e *ContactFieldChangedEvent) Type() string { return TypeContactFieldChanged }

// Apply applies this event to the given run
func (e *ContactFieldChangedEvent) Apply(run flows.FlowRun) error {
	field, err := run.Session().Assets().GetField(e.Field.Key)
	if err != nil {
		return err
	}

	run.Contact().Fields().Save(run.Environment(), field, e.Value)

	return run.Contact().UpdateDynamicGroups(run.Session())
}
