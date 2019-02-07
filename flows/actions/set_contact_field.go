package actions

import (
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions/modifiers"
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
//     "value": "Female"
//   }
//
// @action set_contact_field
type SetContactFieldAction struct {
	BaseAction
	universalAction

	Field *assets.FieldReference `json:"field" validate:"required"`
	Value string                 `json:"value" engine:"evaluate"`
}

// NewSetContactFieldAction creates a new set channel action
func NewSetContactFieldAction(uuid flows.ActionUUID, field *assets.FieldReference, value string) *SetContactFieldAction {
	return &SetContactFieldAction{
		BaseAction: NewBaseAction(TypeSetContactField, uuid),
		Field:      field,
		Value:      value,
	}
}

// Validate validates our action is valid and has all the assets it needs
func (a *SetContactFieldAction) Validate(assets flows.SessionAssets, context *flows.ValidationContext) error {
	_, err := assets.Fields().Get(a.Field.Key)
	return err
}

// Execute runs this action
func (a *SetContactFieldAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	if run.Contact() == nil {
		logEvent(events.NewErrorEventf("can't execute action in session without a contact"))
		return nil
	}

	rawValue, err := a.evaluateLocalizableTemplate(run, "value", a.Value)
	rawValue = strings.TrimSpace(rawValue)

	// if we received an error, log it
	if err != nil {
		logEvent(events.NewErrorEvent(err))
		return nil
	}

	fields := run.Session().Assets().Fields()

	field, err := fields.Get(a.Field.Key)
	if err != nil {
		return err
	}

	newValue := run.Contact().Fields().Parse(run.Environment(), fields, field, rawValue)

	a.applyModifier(run, modifiers.NewFieldModifier(field, newValue), logModifier, logEvent)
	return nil
}
