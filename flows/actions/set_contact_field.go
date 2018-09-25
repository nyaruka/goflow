package actions

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/assets"
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
	universalAction

	Field *assets.FieldReference `json:"field" validate:"required"`
	Value string                 `json:"value"`
}

// Type returns the type of this action
func (a *SetContactFieldAction) Type() string { return TypeSetContactField }

// Validate validates our action is valid and has all the assets it needs
func (a *SetContactFieldAction) Validate(assets flows.SessionAssets) error {
	_, err := assets.Fields().Get(a.Field.Key)
	return err
}

// Execute runs this action
func (a *SetContactFieldAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	if run.Contact() == nil {
		a.logError(fmt.Errorf("can't execute action in session without a contact"), log)
		return nil
	}

	rawValue, err := a.evaluateLocalizableTemplate(run, "value", a.Value)
	rawValue = strings.TrimSpace(rawValue)

	// if we received an error, log it
	if err != nil {
		a.logError(err, log)
		return nil
	}

	fields := run.Session().Assets().Fields()
	value, err := run.Contact().SetFieldValue(run.Environment(), fields, a.Field.Key, rawValue)
	if err != nil {
		return err
	}

	a.log(events.NewContactFieldChangedEvent(a.Field, value), log)

	a.reevaluateDynamicGroups(run, log)
	return nil
}
