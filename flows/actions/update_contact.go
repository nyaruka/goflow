package actions

import (
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
	"strings"
)

// TypeUpdateContact is the type for our update contact action
const TypeUpdateContact string = "update_contact"

// UpdateContactAction can be used to update one of the built in fields for a contact of "name" or
// "language". An `contact_changed` event will be created with the corresponding values.
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
	FieldName string `json:"field_name"    validate:"required,eq=name|eq=language"`
	Value     string `json:"value"`
}

// Type returns the type of this action
func (a *UpdateContactAction) Type() string { return TypeUpdateContact }

// Validate validates our action is valid and has all the assets it needs
func (a *UpdateContactAction) Validate(assets flows.SessionAssets) error {
	// check language is valid if specified
	if a.FieldName == "language" && a.Value != "" {
		if _, err := utils.ParseLanguage(a.Value); err != nil {
			return err
		}
	}
	return nil
}

// Execute runs this action
func (a *UpdateContactAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	// this is a no-op if we have no contact
	if run.Contact() == nil {
		return nil
	}

	// get our localized value if any
	template := run.GetText(flows.UUID(a.UUID()), "value", a.Value)
	value, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), template, false)
	value = strings.TrimSpace(value)

	// if we received an error, log it
	if err != nil {
		log.Add(events.NewErrorEvent(err))
		return nil
	}

	log.Add(events.NewContactChangedEvent(a.FieldName, value))
	return nil
}
