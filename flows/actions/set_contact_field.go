package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

// TypeSetContactField is the type for the set contact field action
const TypeSetContactField string = "set_contact_field"

// SetContactFieldAction can be used to save a value to a contact. The value can be a template and will
// be evaluated during the flow. A `contact_field_changed` event will be created with the corresponding value.
//
// ```
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "set_contact_field",
//     "field": {"key": "gender", "label": "Gender"},
//     "value": "Male"
//   }
// ```
//
// @action set_contact_field
type SetContactFieldAction struct {
	BaseAction
	Field *flows.FieldReference `json:"field" validate:"required"`
	Value string                `json:"value"`
}

// Type returns the type of this action
func (a *SetContactFieldAction) Type() string { return TypeSetContactField }

// Validate validates our action is valid and has all the assets it needs
func (a *SetContactFieldAction) Validate(assets flows.SessionAssets) error {
	_, err := assets.GetField(a.Field.Key)
	return err
}

// Execute runs this action
func (a *SetContactFieldAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	if run.Contact() == nil {
		log.Add(events.NewFatalErrorEvent(fmt.Errorf("can't execute action in session without a contact")))
		return nil
	}

	// get our localized value if any
	template := run.GetText(utils.UUID(a.UUID()), "value", a.Value)
	value, err := run.EvaluateTemplate(template, false)

	// if we received an error, log it
	if err != nil {
		log.Add(events.NewErrorEvent(err))
		return nil
	}

	log.Add(events.NewContactFieldChangedEvent(a.Field, value))
	return nil
}
