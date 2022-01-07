package actions

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/modifiers"

	"github.com/pkg/errors"
)

func init() {
	registerType(TypeRemoveContactGroups, func() flows.Action { return &RemoveContactGroupsAction{} })
}

// TypeRemoveContactGroups is the type for the remove from groups action
const TypeRemoveContactGroups string = "remove_contact_groups"

// RemoveContactGroupsAction can be used to remove a contact from one or more groups. A [event:contact_groups_changed] event will be created
// for the groups which the contact is removed from. Groups can either be explicitly provided or `all_groups` can be set to true to remove
// the contact from all non-query based groups.
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
	baseAction
	universalAction

	Groups    []*assets.GroupReference `json:"groups,omitempty" validate:"dive"`
	AllGroups bool                     `json:"all_groups,omitempty"`
}

// NewRemoveContactGroups creates a new remove from groups action
func NewRemoveContactGroups(uuid flows.ActionUUID, groups []*assets.GroupReference, allGroups bool) *RemoveContactGroupsAction {
	return &RemoveContactGroupsAction{
		baseAction: newBaseAction(TypeRemoveContactGroups, uuid),
		Groups:     groups,
		AllGroups:  allGroups,
	}
}

// Validate validates our action is valid
func (a *RemoveContactGroupsAction) Validate() error {
	if a.AllGroups && len(a.Groups) > 0 {
		return errors.Errorf("can't specify specific groups when all_groups=true")
	}
	return nil
}

// Execute runs the action
func (a *RemoveContactGroupsAction) Execute(run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	contact := run.Contact()
	if contact == nil {
		logEvent(events.NewErrorf("can't execute action in session without a contact"))
		return nil
	}

	var groups []*flows.Group

	if a.AllGroups {
		for _, group := range run.Session().Assets().Groups().All() {
			if !group.UsesQuery() {
				groups = append(groups, group)
			}
		}
	} else {
		groups = resolveGroups(run, a.Groups, logEvent)
	}

	a.applyModifier(run, modifiers.NewGroups(groups, modifiers.GroupsRemove), logModifier, logEvent)
	return nil
}
