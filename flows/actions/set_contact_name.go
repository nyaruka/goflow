package actions

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/flows"
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

// Type returns the type of this action
func (a *SetContactNameAction) Type() string { return TypeSetContactName }

// Validate validates our action is valid and has all the assets it needs
func (a *SetContactNameAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute runs this action
func (a *SetContactNameAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	if run.Contact() == nil {
		log.Add(a.fatalError(run, fmt.Errorf("can't execute action in session without a contact")))
		return nil
	}

	name, err := a.evaluateLocalizableTemplate(run, "name", a.Name)
	name = strings.TrimSpace(name)

	// if we received an error, log it
	if err != nil {
		log.Add(events.NewErrorEvent(err))
		return nil
	}

	if run.Contact().Name() != name {
		run.Contact().SetName(name)
		log.Add(events.NewContactNameChangedEvent(name))
	}

	a.reevaluateDynamicGroups(run, log)
	return nil
}
