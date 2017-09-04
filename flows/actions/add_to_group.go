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
	if contact == nil {
		return nil
	}

	groups, err := a.resolveGroups(run, step, a.Groups)
	if err != nil {
		return err
	}

	groupUUIDs := make([]flows.GroupUUID, 0, len(groups))
	for _, group := range groups {
		// ignore group if contact is already in it
		if contact.Groups().FindByUUID(group.UUID()) != nil {
			continue
		}

		// error if group is dynamic
		if group.IsDynamic() {
			run.AddError(step, a, fmt.Errorf("can't manually add contact to dynamic group '%s' (%s)", group.Name(), group.UUID()))
			continue
		}

		groupUUIDs = append(groupUUIDs, group.UUID())
	}

	if len(groupUUIDs) > 0 {
		run.ApplyEvent(step, a, events.NewAddToGroupEvent(groupUUIDs))
	}

	return nil
}

func (a *AddToGroupAction) resolveGroups(run flows.FlowRun, step flows.Step, references []*flows.GroupReference) ([]*flows.Group, error) {
	groupSet, err := run.Session().Assets().GetGroupSet()
	if err != nil {
		return nil, err
	}

	groups := make([]*flows.Group, 0, len(references))

	for _, ref := range references {
		var group *flows.Group

		if ref.UUID != "" {
			// group is a fixed group with a UUID
			group = groupSet.FindByUUID(ref.UUID)
			if group == nil {
				return nil, fmt.Errorf("no such group with UUID '%s'", ref.UUID)
			}
		} else {
			// group is an expression that evaluates to an existing group's name
			// evaluate the expression to get the group name
			evaluatedGroupName, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), ref.Name)
			if err != nil {
				run.AddError(step, a, err)
			} else {
				// look up the set of all groups to see if such a group exists
				group = groupSet.FindByName(evaluatedGroupName)
				if group == nil {
					run.AddError(step, a, fmt.Errorf("no such group with name '%s'", evaluatedGroupName))
				}
			}
		}

		if group != nil {
			groups = append(groups, group)
		}
	}

	return groups, nil
}
