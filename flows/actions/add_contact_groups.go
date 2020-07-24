package actions

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/modifiers"
)

func init() {
	registerType(TypeAddContactGroups, func() flows.Action { return &AddContactGroupsAction{} })
}

// TypeAddContactGroups is our type for the add to groups action
const TypeAddContactGroups string = "add_contact_groups"

// AddContactGroupsAction can be used to add a contact to one or more groups. A [event:contact_groups_changed] event will be created
// for the groups which the contact has been added to.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "add_contact_groups",
//     "groups": [{
//       "uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a",
//       "name": "Customers"
//     }]
//   }
//
// @action add_contact_groups
type AddContactGroupsAction struct {
	baseAction
	universalAction

	Groups []*assets.GroupReference `json:"groups" validate:"required,dive"`
}

// NewAddContactGroups creates a new add to groups action
func NewAddContactGroups(uuid flows.ActionUUID, groups []*assets.GroupReference) *AddContactGroupsAction {
	return &AddContactGroupsAction{
		baseAction: newBaseAction(TypeAddContactGroups, uuid),
		Groups:     groups,
	}
}

// Execute adds our contact to the specified groups
func (a *AddContactGroupsAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	contact := run.Contact()
	if contact == nil {
		logEvent(events.NewErrorf("can't execute action in session without a contact"))
		return nil
	}

	groups, err := resolveGroups(run, a.Groups, logEvent)
	if err != nil {
		return err
	}

	a.applyModifier(run, modifiers.NewGroups(groups, modifiers.GroupsAdd), logModifier, logEvent)
	return nil
}
