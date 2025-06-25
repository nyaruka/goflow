package actions

import (
	"context"
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/modifiers"
)

func init() {
	registerType(TypeRemoveContactGroups, func() flows.Action { return &RemoveContactGroups{} })
}

// TypeRemoveContactGroups is the type for the remove from groups action
const TypeRemoveContactGroups string = "remove_contact_groups"

// RemoveContactGroups can be used to remove a contact from one or more groups. A [event:contact_groups_changed] event will be created
// for the groups which the contact is removed from. Groups can either be explicitly provided or `all_groups` can be set to true to remove
// the contact from all non-query based groups.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "remove_contact_groups",
//	  "groups": [{
//	    "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
//	    "name": "Registered Users"
//	  }]
//	}
//
// @action remove_contact_groups
type RemoveContactGroups struct {
	baseAction
	universalAction

	Groups    []*assets.GroupReference `json:"groups,omitempty" validate:"dive"`
	AllGroups bool                     `json:"all_groups,omitempty"`
}

// NewRemoveContactGroups creates a new remove from groups action
func NewRemoveContactGroups(uuid flows.ActionUUID, groups []*assets.GroupReference, allGroups bool) *RemoveContactGroups {
	return &RemoveContactGroups{
		baseAction: newBaseAction(TypeRemoveContactGroups, uuid),
		Groups:     groups,
		AllGroups:  allGroups,
	}
}

// Validate validates our action is valid
func (a *RemoveContactGroups) Validate() error {
	if a.AllGroups && len(a.Groups) > 0 {
		return fmt.Errorf("can't specify specific groups when all_groups=true")
	}
	return nil
}

// Execute runs the action
func (a *RemoveContactGroups) Execute(ctx context.Context, run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
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

func (a *RemoveContactGroups) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	for _, group := range a.Groups {
		dependency(group)
	}
}
