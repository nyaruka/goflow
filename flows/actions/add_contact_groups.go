package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// TypeAddContactGroups is our type for the add to groups action
const TypeAddContactGroups string = "add_contact_groups"

// AddContactGroupsAction can be used to add a contact to one or more groups. An `contact_groups_added` event will be created
// for the groups which the contact has been added to.
//
// ```
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "add_contact_groups",
//     "groups": [{
//       "uuid": "2aad21f6-30b7-42c5-bd7f-1b720c154817",
//       "name": "Survey Audience"
//     }]
//   }
// ```
//
// @action add_contact_groups
type AddContactGroupsAction struct {
	BaseAction
	Groups []*flows.GroupReference `json:"groups" validate:"required,min=1,dive"`
}

// Type returns the type of this action
func (a *AddContactGroupsAction) Type() string { return TypeAddContactGroups }

// Validate validates our action is valid and has all the assets it needs
func (a *AddContactGroupsAction) Validate(assets flows.SessionAssets) error {
	// check we have all groups
	return a.validateGroups(assets, a.Groups)
}

// Execute adds our contact to the specified groups
func (a *AddContactGroupsAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
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
		// ignore group if contact is already in it
		if contact.Groups().FindByUUID(group.UUID()) != nil {
			continue
		}

		// error if group is dynamic
		if group.IsDynamic() {
			log.Add(events.NewErrorEvent(fmt.Errorf("can't manually add contact to dynamic group '%s' (%s)", group.Name(), group.UUID())))
			continue
		}

		groupRefs = append(groupRefs, group.Reference())
	}

	// only generate event if contact's groups change
	if len(groupRefs) > 0 {
		log.Add(events.NewContactGroupsAddedEvent(groupRefs))
	}

	return nil
}
