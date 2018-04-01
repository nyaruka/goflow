package events

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
)

// TypeContactFieldChanged is the type of our save to contact event
const TypeContactFieldChanged string = "contact_field_changed"

// ContactFieldChangedEvent events are created when a contact field is updated.
//
// ```
//   {
//     "type": "contact_field_changed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "field": {"key": "gender", "name": "Gender"},
//     "value": "Male"
//   }
// ```
//
// @event contact_field_changed
type ContactFieldChangedEvent struct {
	baseEvent
	callerOrEngineEvent

	Field *flows.FieldReference `json:"field" validate:"required"`
	Value string                `json:"value" validate:"required"`
}

// NewContactFieldChangedEvent returns a new save to contact event
func NewContactFieldChangedEvent(field *flows.FieldReference, value string) *ContactFieldChangedEvent {
	return &ContactFieldChangedEvent{
		baseEvent: newBaseEvent(),
		Field:     field,
		Value:     value,
	}
}

// Type returns the type of this event
func (e *ContactFieldChangedEvent) Type() string { return TypeContactFieldChanged }

// Validate validates our event is valid and has all the assets it needs
func (e *ContactFieldChangedEvent) Validate(assets flows.SessionAssets) error {
	_, err := assets.GetField(e.Field.Key)
	return err
}

// Apply applies this event to the given run
func (e *ContactFieldChangedEvent) Apply(run flows.FlowRun) error {
	if run.Contact() == nil {
		return fmt.Errorf("can't apply event in session without a contact")
	}

	field, err := run.Session().Assets().GetField(e.Field.Key)
	if err != nil {
		return err
	}

	run.Contact().SetFieldValue(run.Environment(), field, e.Value)

	return run.Contact().UpdateDynamicGroups(run.Session())
}
