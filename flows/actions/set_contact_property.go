package actions

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

// TypeSetContactProperty is the type for the set contact property action
const TypeSetContactProperty string = "set_contact_property"

// SetContactPropertyAction can be used to update one of the built in fields for a contact of "name" or
// "language". An `contact_property_changed` event will be created with the corresponding values.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "set_contact_property",
//     "property": "language",
//     "value": "eng"
//   }
//
// @action set_contact_property
type SetContactPropertyAction struct {
	BaseAction
	Property string `json:"property" validate:"required,eq=name|eq=language"`
	Value    string `json:"value"`
}

// Type returns the type of this action
func (a *SetContactPropertyAction) Type() string { return TypeSetContactProperty }

// Validate validates our action is valid and has all the assets it needs
func (a *SetContactPropertyAction) Validate(assets flows.SessionAssets) error {
	// check language is valid if specified
	if a.Property == "language" && a.Value != "" {
		if _, err := utils.ParseLanguage(a.Value); err != nil {
			return err
		}
	}
	return nil
}

// Execute runs this action
func (a *SetContactPropertyAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	if run.Contact() == nil {
		log.Add(events.NewFatalErrorEvent(fmt.Errorf("can't execute action in session without a contact")))
		return nil
	}

	// get our localized value if any
	template := run.GetText(utils.UUID(a.UUID()), "value", a.Value)
	value, err := run.EvaluateTemplateAsString(template, false)
	value = strings.TrimSpace(value)

	// if we received an error, log it
	if err != nil {
		log.Add(events.NewErrorEvent(err))
		return nil
	}

	log.Add(events.NewContactPropertyChangedEvent(a.Property, value))
	return nil
}
