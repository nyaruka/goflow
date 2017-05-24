package actions

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

// TypeAddToGroup is our type for the add to group action
const TypeAddToGroup string = "add_to_group"

// AddToGroupAction can be used to add a contact to one or more groups. An `add_to_group` event will be created
// for each group which the contact is added to.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "add_to_group",
//     "groups": [{
//       "uuid": "2aad21f6-30b7-42c5-bd7f-1b720c154817",
//       "name": "Survey Audience"
//     }]
//   }
//
// @action add_to_group
type AddToGroupAction struct {
	BaseAction
	Groups []*flows.Group `json:"groups"    validate:"required,min=1"`
}

// Type returns the type of this action
func (a *AddToGroupAction) Type() string { return TypeAddToGroup }

// Validate validates that this action is valid
func (a *AddToGroupAction) Validate() error {
	return utils.ValidateAll(a)
}

// Execute adds our contact to the specified groups
func (a *AddToGroupAction) Execute(run flows.FlowRun, step flows.Step) error {
	contact := run.Contact()
	if contact != nil {
		groups := make([]flows.GroupUUID, 0, len(a.Groups))
		for _, group := range a.Groups {
			if contact.Groups().FindGroup(group.UUID()) == nil {
				groups = append(groups, group.UUID())
			}

		}
		if len(groups) > 0 {
			run.AddEvent(step, events.NewGroupEvent(groups))
		}
	}

	return nil
}
