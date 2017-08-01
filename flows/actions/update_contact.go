package actions

import (
	"strings"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

// TypeUpdateContact is the type for our update contact action
const TypeUpdateContact string = "update_contact"

// UpdateContactAction can be used to update one of the built in fields for a contact of "name" or
// "language". An `update_contact` event will be created with the corresponding values.
//
// ```
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "update_contact",
//     "field_name": "language",
//     "value": "eng"
//   }
// ```
//
// @action update_contact
type UpdateContactAction struct {
	BaseAction
	FieldName string `json:"field_name"    validate:"required,eq=language|eq=name"`
	Value     string `json:"value"         validate:"required"`
}

// Type returns the type of this action
func (a *UpdateContactAction) Type() string { return TypeUpdateContact }

// Validate validates this action
func (a *UpdateContactAction) Validate(assets flows.Assets) error {
	return nil
}

// Execute runs this action
func (a *UpdateContactAction) Execute(run flows.FlowRun, step flows.Step) error {
	// this is a no-op if we have no contact
	if run.Contact() == nil {
		return nil
	}

	// get our localized value if any
	template := run.GetText(flows.UUID(a.UUID()), "value", a.Value)
	value, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), template)

	// if we received an error, log it
	if err != nil {
		run.AddError(step, err)
	}

	// if this is either name or language, we save directly to the contact
	if strings.ToLower(a.FieldName) == "name" {
		run.Contact().SetName(value)
	} else if strings.ToLower(a.FieldName) == "language" {
		// try to parse our language
		lang := utils.NilLanguage
		lang, err = utils.ParseLanguage(value)

		// if this doesn't look valid, log an error and don't set our language
		if err != nil {
			run.AddError(step, err)
		} else {
			run.Contact().SetLanguage(lang)
			value = string(lang)
		}
	}

	// log our event
	if err == nil {
		run.ApplyEvent(step, a, events.NewUpdateContact(strings.ToLower(a.FieldName), value))
	}

	return nil
}
