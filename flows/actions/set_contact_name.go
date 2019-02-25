package actions

import (
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions/modifiers"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	RegisterType(TypeSetContactName, func() flows.Action { return &SetContactNameAction{} })
}

// TypeSetContactName is the type for the set contact name action
const TypeSetContactName string = "set_contact_name"

// SetContactNameAction can be used to update the name of the contact. The name is a localizable
// template and white space is trimmed from the final value. An empty string clears the name.
// A [event:contact_name_changed] event will be created with the corresponding value.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "set_contact_name",
//     "name": "Bob Smith"
//   }
//
// @action set_contact_name
type SetContactNameAction struct {
	BaseAction
	universalAction

	Name string `json:"name"`
}

// NewSetContactNameAction creates a new set name action
func NewSetContactNameAction(uuid flows.ActionUUID, name string) *SetContactNameAction {
	return &SetContactNameAction{
		BaseAction: NewBaseAction(TypeSetContactName, uuid),
		Name:       name,
	}
}

// Validate validates our action is valid and has all the assets it needs
func (a *SetContactNameAction) Validate(assets flows.SessionAssets, context *flows.ValidationContext) error {
	return nil
}

// Execute runs this action
func (a *SetContactNameAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	if run.Contact() == nil {
		logEvent(events.NewErrorEventf("can't execute action in session without a contact"))
		return nil
	}

	name, err := run.EvaluateTemplate(a.Name)
	name = strings.TrimSpace(name)

	// if we received an error, log it
	if err != nil {
		logEvent(events.NewErrorEvent(err))
		return nil
	}

	a.applyModifier(run, modifiers.NewNameModifier(name), logModifier, logEvent)
	return nil
}
