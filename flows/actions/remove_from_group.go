package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/excellent"
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
	Groups []*flows.GroupReference `json:"groups" validate:"dive"`
}

// Type returns the type of this action
func (a *RemoveFromGroupAction) Type() string { return TypeRemoveFromGroup }

// Validate validates the fields on this action
func (a *RemoveFromGroupAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute runs the action
func (a *RemoveFromGroupAction) Execute(run flows.FlowRun, step flows.Step) error {
	// only generate event if contact's groups change
	contact := run.Contact()
	if contact != nil {
		groupUUIDs := make([]flows.GroupUUID, 0)

		// no groups in our action means remove all
		if len(a.Groups) == 0 {
			for _, group := range contact.Groups() {
				groupUUIDs = append(groupUUIDs, group.UUID())
			}
		} else {
			for _, group := range a.Groups {
				if group.UUID != "" && contact.Groups().FindByUUID(group.UUID) != nil {
					// group is a fixed group with a UUID, and contact does belong to it
					groupUUIDs = append(groupUUIDs, group.UUID)
				} else {
					// group is an expression that evaluates to an existing group's name
					allGroups, err := run.Session().Assets().GetGroupSet()
					if err != nil {
						return err
					}

					// evaluate the expression to get the group name
					evaluatedGroupName, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), group.Name)
					if err != nil {
						run.AddError(step, a, err)
					} else {
						// look up the set of all groups to see if such a group exists
						addGroup := allGroups.FindByName(evaluatedGroupName)
						if addGroup == nil {
							run.AddError(step, a, fmt.Errorf("no such group with name '%s'", evaluatedGroupName))
						} else if contact.Groups().FindByUUID(addGroup.UUID()) != nil {
							groupUUIDs = append(groupUUIDs, addGroup.UUID())
						}
					}
				}
			}
		}

		if len(groupUUIDs) > 0 {
			run.ApplyEvent(step, a, events.NewRemoveFromGroupEvent(groupUUIDs))
		}
	}

	return nil
}
