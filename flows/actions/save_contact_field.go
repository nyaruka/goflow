package actions

import (
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// TypeSaveContactField is the type for our save to contact action
const TypeSaveContactField string = "save_contact_field"

// SaveContactField can be used to save a value to a contact. The value can be a template and will
// be evaluated during the flow. A `save_contact_field` event will be created with the corresponding value.
//
// ```
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "save_contact_field",
//     "field_uuid": "0cb17b2a-3bfe-4a19-8c99-98ab9561045d",
//     "field_name": "Gender",
//     "value": "Male"
//   }
// ```
//
// @action save_contact_field
type SaveContactField struct {
	BaseAction
	FieldUUID flows.FieldUUID `json:"field_uuid"    validate:"required,uuid4"`
	FieldName string          `json:"field_name"    validate:"required"`
	Value     string          `json:"value"`
}

// Type returns the type of this action
func (a *SaveContactField) Type() string { return TypeSaveContactField }

// Validate validates this action
func (a *SaveContactField) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute runs this action
func (a *SaveContactField) Execute(run flows.FlowRun, step flows.Step) error {
	// this is a no-op if we have no contact
	if run.Contact() == nil {
		return nil
	}

	// get our localized value if any
	template := run.GetText(flows.UUID(a.UUID()), "value", a.Value)
	value, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), template)

	// if we received an error, log it
	if err != nil {
		run.AddError(step, a, err)
	}

	run.Contact().Fields().Save(a.FieldUUID, a.FieldName, value)

	// log our event
	if err == nil {
		run.ApplyEvent(step, a, events.NewSaveToContact(a.FieldUUID, a.FieldName, value))
	}

	return nil
}
