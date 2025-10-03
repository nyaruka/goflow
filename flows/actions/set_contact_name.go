package actions

import (
	"context"
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/modifiers"
)

func init() {
	registerType(TypeSetContactName, func() flows.Action { return &SetContactName{} })
}

// TypeSetContactName is the type for the set contact name action
const TypeSetContactName string = "set_contact_name"

// SetContactName can be used to update the name of the contact. The name is a localizable
// template and white space is trimmed from the final value. An empty string clears the name.
// A [event:contact_name_changed] event will be created with the corresponding value.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "set_contact_name",
//	  "name": "Bob Smith"
//	}
//
// @action set_contact_name
type SetContactName struct {
	baseAction
	universalAction

	Name string `json:"name" validate:"max=1000" engine:"evaluated"`
}

// NewSetContactName creates a new set name action
func NewSetContactName(uuid flows.ActionUUID, name string) *SetContactName {
	return &SetContactName{
		baseAction: newBaseAction(TypeSetContactName, uuid),
		Name:       name,
	}
}

// Execute runs this action
func (a *SetContactName) Execute(ctx context.Context, run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	name, ok := run.EvaluateTemplate(a.Name, logEvent)
	name = strings.TrimSpace(name)

	if !ok {
		return nil
	}

	_, err := a.applyModifier(run, modifiers.NewName(name), logModifier, logEvent)
	return err
}
