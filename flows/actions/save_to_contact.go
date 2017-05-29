package actions

import (
	"time"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

// TypeSaveToContact is the type for our save to contact action
const TypeSaveToContact string = "save_to_contact"

// SaveToContactAction can be used to save a value to a contact. The value can be a template and will
// be evaluated during the flow. A `save_to_contact` event will be created with the corresponding value.
//
// Two fields are treated specially, "name" and "language", which can be used as the "field" parameter
// and which will set the special contact fields of the same name.
//
// ```
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "save_to_contact",
//     "field": "0cb17b2a-3bfe-4a19-8c99-98ab9561045d",
//     "name": "Gender",
//     "value": "Male"
//   }
// ```
//
// @action save_to_contact
type SaveToContactAction struct {
	BaseAction
	Field flows.FieldUUID `json:"field"    validate:"required"`
	Name  string          `json:"name"     validate:"required"`
	Value string          `json:"value"    validate:"required"`
}

// Type returns the type of this action
func (a *SaveToContactAction) Type() string { return TypeSaveToContact }

// Validate validates this action
func (a *SaveToContactAction) Validate() error {
	return utils.ValidateAll(a)
}

// Execute runs this action
func (a *SaveToContactAction) Execute(run flows.FlowRun, step flows.Step) error {
	// this is a no-op if we have no contact
	if run.Contact() == nil {
		return nil
	}

	// get our localized value if any
	template := run.GetText(flows.UUID(a.UUID), "value", a.Value)
	value, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), template)

	// if we received an error, log it
	if err != nil {
		run.AddError(step, err)
	}

	// if this is either name or language, we save directly to the contact
	if a.Field == "name" {
		run.Contact().SetName(value)
	} else if a.Field == "language" {
		// try to parse our language
		lang := utils.NilLanguage
		lang, err = utils.ParseLanguage(value)

		// if this doesn't look valid, log an error and don't set our language
		if err != nil {
			run.AddError(step, err)
		} else {
			run.Contact().SetLanguage(lang)
			run.SetLanguage(lang)
			value = string(lang)
		}
	} else {
		// save to our field dictionary
		run.Contact().Fields().Save(a.Field, a.Name, value, time.Now().In(time.UTC))
	}

	// log our event
	if err == nil {
		run.AddEvent(step, events.NewSaveToContact(a.Field, a.Name, value))
	}

	return nil
}
