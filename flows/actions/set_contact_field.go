package actions

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	RegisterType(TypeSetContactField, func() flows.Action { return &SetContactFieldAction{} })
}

// TypeSetContactField is the type for the set contact field action
const TypeSetContactField string = "set_contact_field"

// SetContactFieldAction can be used to update a field value on the contact. The value is a localizable
// template and white space is trimmed from the final value. An empty string clears the value.
// A [event:contact_field_changed] event will be created with the corresponding value.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "set_contact_field",
//     "field": {"key": "gender", "name": "Gender"},
//     "value": "Male"
//   }
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

	value, err := a.evaluateLocalizableTemplate(run, "value", a.Value)
	value = strings.TrimSpace(value)

	// if we received an error, log it
	if err != nil {
		log.Add(events.NewErrorEvent(err))
		return nil
	}

	log.Add(events.NewContactFieldChangedEvent(a.Field, value))
	return nil
}
