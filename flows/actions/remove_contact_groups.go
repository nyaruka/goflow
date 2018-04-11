package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// TypeRemoveContactGroups is the type for the remove from groups action
const TypeRemoveContactGroups string = "remove_contact_groups"

// RemoveContactGroupsAction can be used to remove a contact from one or more groups. A `contact_groups_removed` event will be created
// for the groups which the contact is removed from. If no groups are specified, then the contact will be removed from
// all groups.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "remove_contact_groups",
//     "groups": [{
//       "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
//       "name": "Registered Users"
//     }]
//   }
//
// @action remove_contact_groups
type RemoveContactGroupsAction struct {
	BaseAction
	Groups []*flows.GroupReference `json:"groups" validate:"required,min=1,dive"`
}

// Type returns the type of this action
func (a *RemoveContactGroupsAction) Type() string { return TypeRemoveContactGroups }

// Validate validates our action is valid and has all the assets it needs
func (a *RemoveContactGroupsAction) Validate(assets flows.SessionAssets) error {
	// check we have all groups
	return a.validateGroups(assets, a.Groups)
}

// Execute runs the action
func (a *RemoveContactGroupsAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	contact := run.Contact()
	if contact == nil {
		log.Add(events.NewFatalErrorEvent(fmt.Errorf("can't execute action in session without a contact")))
		return nil
	}

	groups, err := a.resolveGroups(run, step, a.Groups, log)
	if err != nil {
		return err
	}

	groupRefs := make([]*flows.GroupReference, 0, len(groups))
	for _, group := range groups {
		// ignore group if contact isn't actually in it
		if contact.Groups().FindByUUID(group.UUID()) == nil {
			continue
		}

		// error if group is dynamic
		if group.IsDynamic() {
			log.Add(events.NewErrorEvent(fmt.Errorf("can't manually remove contact from dynamic group '%s' (%s)", group.Name(), group.UUID())))
			continue
		}

		groupRefs = append(groupRefs, group.Reference())
	}

	// only generate event if contact's groups change
	if len(groupRefs) > 0 {
		log.Add(events.NewContactGroupsRemovedEvent(groupRefs))
	}

	return nil
}
