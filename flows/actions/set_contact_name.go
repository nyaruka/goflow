package actions

import (
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/modifiers"
)

func init() {
	registerType(TypeSetContactName, func() flows.Action { return &SetContactNameAction{} })
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
	baseAction
	universalAction

	Name string `json:"name" engine:"evaluated"`
}

// NewSetContactName creates a new set name action
func NewSetContactName(uuid flows.ActionUUID, name string) *SetContactNameAction {
	return &SetContactNameAction{
		baseAction: newBaseAction(TypeSetContactName, uuid),
		Name:       name,
	}
}

// Execute runs this action
func (a *SetContactNameAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	if run.Contact() == nil {
		logEvent(events.NewErrorf("can't execute action in session without a contact"))
		return nil
	}

	name, err := run.EvaluateTemplate(a.Name)
	name = strings.TrimSpace(name)

	// if we received an error, log it
	if err != nil {
		logEvent(events.NewError(err))
		return nil
	}

	a.applyModifier(run, modifiers.NewName(name), logModifier, logEvent)
	return nil
}
