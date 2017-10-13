package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// TypeRemoveFromGroup is our type for our remove from group action
const TypeRemoveFromGroup string = "remove_from_group"

// RemoveFromGroupAction can be used to remove a contact from one or more groups. A `remove_from_group` event will be created
// for each group which the contact is removed from. If no groups are specified, then the contact will be removed from
// all groups.
//
// ```
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "remove_from_group",
//     "groups": [{
//       "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
//       "name": "Registered Users"
//     }]
//   }
// ```
//
// @action remove_from_group
type RemoveFromGroupAction struct {
	BaseAction
	Groups []*flows.GroupReference `json:"groups" validate:"required,min=1,dive"`
}

// Type returns the type of this action
func (a *RemoveFromGroupAction) Type() string { return TypeRemoveFromGroup }

// Validate validates the fields on this action
func (a *RemoveFromGroupAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute runs the action
func (a *RemoveFromGroupAction) Execute(run flows.FlowRun, step flows.Step, log flows.ActionLog) error {
	// only generate event if contact's groups change
	contact := run.Contact()
	if contact == nil {
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

	if len(groupRefs) > 0 {
		log.Add(events.NewRemoveFromGroupEvent(groupRefs))
	}

	return nil
}
