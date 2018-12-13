package modifiers

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	RegisterType(TypeGroups, func() Modifier { return &GroupsModifier{} })
}

// TypeGroups is the type of our groups modifier
const TypeGroups string = "groups"

// GroupsModifier modifies the group membership of the contact
type GroupsModifier struct {
	baseModifier

	Groups []*flows.Group
	Add    bool
}

// NewGroupsModifier creates a new groups modifier
func NewGroupsModifier(groups []*flows.Group, add bool) *GroupsModifier {
	return &GroupsModifier{
		baseModifier: newBaseModifier(TypeGroups),
		Groups:       groups,
		Add:          add,
	}
}

func (m *GroupsModifier) Apply(assets flows.SessionAssets, contact *flows.Contact) flows.Event {
	diff := make([]*flows.Group, 0, len(m.Groups))
	if m.Add {
		for _, group := range m.Groups {

			// ignore group if contact is already in it
			if contact.Groups().FindByUUID(group.UUID()) != nil {
				continue
			}

			contact.Groups().Add(group)
			diff = append(diff, group)
		}

		// only generate event if contact's groups change
		if len(diff) > 0 {
			return events.NewContactGroupsChangedEvent(diff, nil)
		}
	} else {
		for _, group := range m.Groups {
			// ignore group if contact isn't actually in it
			if contact.Groups().FindByUUID(group.UUID()) == nil {
				continue
			}

			contact.Groups().Remove(group)
			diff = append(diff, group)
		}

		// only generate event if contact's groups change
		if len(diff) > 0 {
			return events.NewContactGroupsChangedEvent(nil, diff)
		}
	}
	return nil
}

var _ Modifier = (*GroupsModifier)(nil)
