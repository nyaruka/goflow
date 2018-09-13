package events

import (
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeContactFieldChanged, func() flows.Event { return &ContactFieldChangedEvent{} })
}

// TypeContactFieldChanged is the type of our save to contact event
const TypeContactFieldChanged string = "contact_field_changed"

// ContactFieldChangedEvent events are created when a contact field is updated.
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
	callerOrEngineEvent

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

// Validate validates our event is valid and has all the assets it needs
func (e *ContactFieldChangedEvent) Validate(assets flows.SessionAssets) error {
	_, err := assets.Fields().Get(e.Field.Key)
	return err
}

// Apply applies this event to the given run
func (e *ContactFieldChangedEvent) Apply(run flows.FlowRun) error {
	if run.Contact() == nil {
		return fmt.Errorf("can't apply event in session without a contact")
	}

	fields := run.Session().Assets().Fields()

	if err := run.Contact().SetFieldValue(run.Environment(), fields, e.Field.Key, e.Value); err != nil {
		return err
	}

	return run.Contact().ReevaluateDynamicGroups(run.Session())
}
