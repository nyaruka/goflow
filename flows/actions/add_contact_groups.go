package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	RegisterType(TypeAddContactGroups, func() flows.Action { return &AddContactGroupsAction{} })
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
	BaseAction
	universalAction

	Groups []*assets.GroupReference `json:"groups" validate:"required,dive"`
}

// NewAddContactGroupsAction creates a new add to groups action
func NewAddContactGroupsAction(uuid flows.ActionUUID, groups []*assets.GroupReference) *AddContactGroupsAction {
	return &AddContactGroupsAction{
		BaseAction: NewBaseAction(TypeAddContactGroups, uuid),
		Groups:     groups,
	}
}

// Type returns the type of this action
func (a *AddContactGroupsAction) Type() string { return TypeAddContactGroups }

// Validate validates our action is valid and has all the assets it needs
func (a *AddContactGroupsAction) Validate(assets flows.SessionAssets, context *flows.ValidationContext) error {
	// check we have all groups
	return a.validateGroups(assets, a.Groups)
}

// Execute adds our contact to the specified groups
func (a *AddContactGroupsAction) Execute(run flows.FlowRun, step flows.Step) error {
	contact := run.Contact()
	if contact == nil {
		a.logError(run, step, fmt.Errorf("can't execute action in session without a contact"))
		return nil
	}

	groups, err := a.resolveGroups(run, step, a.Groups)
	if err != nil {
		return err
	}

	added := make([]*flows.Group, 0, len(groups))
	for _, group := range groups {
		// ignore group if contact is already in it
		if contact.Groups().FindByUUID(group.UUID()) != nil {
			continue
		}

		// error if group is dynamic
		if group.IsDynamic() {
			a.logError(run, step, fmt.Errorf("can't manually add contact to dynamic group '%s' (%s)", group.Name(), group.UUID()))
			continue
		}

		run.Contact().Groups().Add(group)
		added = append(added, group)
	}

	// only generate event if contact's groups change
	if len(added) > 0 {
		a.log(run, step, events.NewContactGroupsChangedEvent(added, nil))
	}

	return nil
}
