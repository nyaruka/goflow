package actions

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions/modifiers"
	"github.com/nyaruka/goflow/flows/events"

	"github.com/pkg/errors"
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
	AllGroups bool                     `json:"all_groups,omitempty"`
}

// NewRemoveContactGroupsAction creates a new remove from groups action
func NewRemoveContactGroupsAction(uuid flows.ActionUUID, groups []*assets.GroupReference, allGroups bool) *RemoveContactGroupsAction {
	return &RemoveContactGroupsAction{
		BaseAction: NewBaseAction(TypeRemoveContactGroups, uuid),
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
func (a *RemoveContactGroupsAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	contact := run.Contact()
	if contact == nil {
		logEvent(events.NewErrorEventf("can't execute action in session without a contact"))
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
		if groups, err = a.resolveGroups(run, a.Groups, true, logEvent); err != nil {
			return err
		}
	}

	a.applyModifier(run, modifiers.NewGroupsModifier(groups, modifiers.GroupsRemove), logModifier, logEvent)
	return nil
}

// Inspect inspects this object and any children
func (a *RemoveContactGroupsAction) Inspect(inspect func(flows.Inspectable)) {
	inspect(a)
}

// EnumerateDependencies enumerates all dependencies on this object and its children
func (a *RemoveContactGroupsAction) EnumerateDependencies(localization flows.Localization, include func(assets.Reference)) {
	for _, g := range a.Groups {
		if !g.Variable() {
			include(g)
		}
	}
}
