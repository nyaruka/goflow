package actions

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

// TypeRemoveFromGroup is our type for our remove from group action
const TypeRemoveFromGroup string = "remove_from_group"

// RemoveFromGroupAction can be used to remove a contact from one or more groups. A `remove_from_group` event will be created
// for each group which the contact is removed from. If no groups are specified, then the contact will be removed from
// all groups.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "remove_from_group",
//     "groups": [{
//       "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
//       "name": "Registered Users"
//     }]
//   }
//
// @action remove_from_group
type RemoveFromGroupAction struct {
	BaseAction
	Groups []*flows.Group `json:"groups"     validate:"dive"`
}

// Type returns the type of this action
func (a *RemoveFromGroupAction) Type() string { return TypeRemoveFromGroup }

// Validate validates the fields on this action
func (a *RemoveFromGroupAction) Validate() error {
	return utils.ValidateAll(a)
}

// Execute runs the action
func (a *RemoveFromGroupAction) Execute(run flows.FlowRun, step flows.Step) error {
	contact := run.Contact()
	if contact != nil {
		groups := make([]flows.GroupUUID, 0)

		// no groups in our action means remove all
		if len(a.Groups) == 0 {
			for _, group := range contact.Groups() {
				contact.RemoveGroup(group.UUID())
				groups = append(groups, group.UUID())
			}

		} else {
			for _, group := range a.Groups {
				removed := contact.RemoveGroup(group.UUID())
				if removed {
					groups = append(groups, group.UUID())
				}
			}
		}

		if len(groups) > 0 {
			run.AddEvent(step, events.NewRemoveFromGroup(groups))
		}
	}

	return nil
}
