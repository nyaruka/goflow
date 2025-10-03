package actions

import (
	"context"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/modifiers"
)

func init() {
	registerType(TypeAddContactGroups, func() flows.Action { return &AddContactGroups{} })
}

// TypeAddContactGroups is our type for the add to groups action
const TypeAddContactGroups string = "add_contact_groups"

// AddContactGroups can be used to add a contact to one or more groups. A [event:contact_groups_changed] event will be created
// for the groups which the contact has been added to.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "add_contact_groups",
//	  "groups": [{
//	    "uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a",
//	    "name": "Customers"
//	  }]
//	}
//
// @action add_contact_groups
type AddContactGroups struct {
	baseAction
	universalAction

	Groups []*assets.GroupReference `json:"groups" validate:"required,max=100,dive"`
}

// NewAddContactGroups creates a new add to groups action
func NewAddContactGroups(uuid flows.ActionUUID, groups []*assets.GroupReference) *AddContactGroups {
	return &AddContactGroups{
		baseAction: newBaseAction(TypeAddContactGroups, uuid),
		Groups:     groups,
	}
}

// Execute adds our contact to the specified groups
func (a *AddContactGroups) Execute(ctx context.Context, run flows.Run, step flows.Step, logEvent flows.EventCallback) error {
	groups := resolveGroups(run, a.Groups, logEvent)

	_, err := a.applyModifier(run, modifiers.NewGroups(groups, modifiers.GroupsAdd), logEvent)
	return err
}

func (a *AddContactGroups) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	for _, group := range a.Groups {
		dependency(group)
	}
}
