package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	RegisterType(TypeRemoveContactGroups, func() flows.Action { return &RemoveContactGroupsAction{} })
}

// TypeRemoveContactGroups is the type for the remove from groups action
const TypeRemoveContactGroups string = "remove_contact_groups"

// RemoveContactGroupsAction can be used to remove a contact from one or more groups. A [event:contact_groups_changed] event will be created
// for the groups which the contact is removed from. Groups can either be explicitly provided or `all_groups` can be set to true to remove
// the contact from all non-dynamic groups.
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
	universalAction

	Groups    []*assets.GroupReference `json:"groups,omitempty" validate:"dive"`
	AllGroups bool                     `json:"all_groups"`
}

// Type returns the type of this action
func (a *RemoveContactGroupsAction) Type() string { return TypeRemoveContactGroups }

// Validate validates our action is valid and has all the assets it needs
func (a *RemoveContactGroupsAction) Validate(assets flows.SessionAssets) error {
	if a.AllGroups && len(a.Groups) > 0 {
		return fmt.Errorf("can't specify specific groups when all_groups=true")
	}

	// check we have all specified groups
	return a.validateGroups(assets, a.Groups)
}

// Execute runs the action
func (a *RemoveContactGroupsAction) Execute(run flows.FlowRun, step flows.Step) error {
	contact := run.Contact()
	if contact == nil {
		a.logError(run, step, fmt.Errorf("can't execute action in session without a contact"))
		return nil
	}

	var groups []*flows.Group
	var err error

	if a.AllGroups {
		for _, group := range run.Session().Assets().Groups().All() {
			if !group.IsDynamic() {
				groups = append(groups, group)
			}
		}
	} else {
		if groups, err = a.resolveGroups(run, step, a.Groups); err != nil {
			return err
		}
	}

	removed := make([]*flows.Group, 0, len(groups))
	for _, group := range groups {
		// ignore group if contact isn't actually in it
		if contact.Groups().FindByUUID(group.UUID()) == nil {
			continue
		}

		// error if group is dynamic
		if group.IsDynamic() {
			a.logError(run, step, fmt.Errorf("can't manually remove contact from dynamic group '%s' (%s)", group.Name(), group.UUID()))
			continue
		}

		run.Contact().Groups().Remove(group)
		removed = append(removed, group)
	}

	// only generate event if contact's groups change
	if len(removed) > 0 {
		a.log(run, step, events.NewContactGroupsChangedEvent(nil, removed))
	}

	return nil
}
