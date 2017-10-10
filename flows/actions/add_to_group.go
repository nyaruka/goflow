package actions

import (
	"fmt"

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
	Groups []*flows.GroupReference `json:"groups" validate:"required,min=1,dive"`
}

// Type returns the type of this action
func (a *AddToGroupAction) Type() string { return TypeAddToGroup }

// Validate validates that this action is valid
func (a *AddToGroupAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute adds our contact to the specified groups
func (a *AddToGroupAction) Execute(run flows.FlowRun, step flows.Step) ([]flows.Event, error) {
	// only generate event if contact's groups change
	contact := run.Contact()
	if contact == nil {
		return nil, nil
	}

	log := make([]flows.Event, 0)

	groups, err := resolveGroups(run, step, a, a.Groups, log)
	if err != nil {
		return log, err
	}

	groupRefs := make([]*flows.GroupReference, 0, len(groups))
	for _, group := range groups {
		// ignore group if contact is already in it
		if contact.Groups().FindByUUID(group.UUID()) != nil {
			continue
		}

		// error if group is dynamic
		if group.IsDynamic() {
			log = append(log, events.NewErrorEvent(fmt.Errorf("can't manually add contact to dynamic group '%s' (%s)", group.Name(), group.UUID())))
			continue
		}

		groupRefs = append(groupRefs, group.Reference())
	}

	if len(groupRefs) > 0 {
		log = append(log, events.NewAddToGroupEvent(groupRefs))
	}

	return log, nil
}
