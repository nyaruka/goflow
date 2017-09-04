package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// TypeAddToGroup is our type for the add to group action
const TypeAddToGroup string = "add_to_group"

// AddToGroupAction can be used to add a contact to one or more groups. An `add_to_group` event will be created
// for each group which the contact is added to.
//
// ```
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "add_to_group",
//     "groups": [{
//       "uuid": "2aad21f6-30b7-42c5-bd7f-1b720c154817",
//       "name": "Survey Audience"
//     }]
//   }
// ```
//
// @action add_to_group
type AddToGroupAction struct {
	BaseAction
	Groups []*flows.GroupReference `json:"groups" validate:"required,min=1"`
}

// Type returns the type of this action
func (a *AddToGroupAction) Type() string { return TypeAddToGroup }

// Validate validates that this action is valid
func (a *AddToGroupAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute adds our contact to the specified groups
func (a *AddToGroupAction) Execute(run flows.FlowRun, step flows.Step) error {
	// only generate event if contact's groups change
	contact := run.Contact()
	if contact != nil {
		groupUUIDs := make([]flows.GroupUUID, 0, len(a.Groups))
		for _, group := range a.Groups {
			if group.UUID != "" && contact.Groups().FindByUUID(group.UUID) == nil {
				// group is a fixed group with a UUID, and contact doesn't already belong to it
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
					} else if contact.Groups().FindByUUID(addGroup.UUID()) == nil {
						groupUUIDs = append(groupUUIDs, addGroup.UUID())
					}
				}
			}
		}
		if len(groupUUIDs) > 0 {
			run.ApplyEvent(step, a, events.NewAddToGroupEvent(groupUUIDs))
		}
	}

	return nil
}
