package modifiers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeGroups, readGroupsModifier)
}

// TypeGroups is the type of our groups modifier
const TypeGroups string = "groups"

// GroupsModification is the type of modification to make
type GroupsModification string

// the supported types of modification
const (
	GroupsAdd    GroupsModification = "add"
	GroupsRemove GroupsModification = "remove"
)

// GroupsModifier modifies the group membership of the contact
type GroupsModifier struct {
	baseModifier

	Groups       []*flows.Group
	Modification GroupsModification
}

// NewGroupsModifier creates a new groups modifier
func NewGroupsModifier(groups []*flows.Group, modification GroupsModification) *GroupsModifier {
	return &GroupsModifier{
		baseModifier: newBaseModifier(TypeGroups),
		Groups:       groups,
		Modification: modification,
	}
}

// Apply applies this modification to the given contact
func (m *GroupsModifier) Apply(env utils.Environment, assets flows.SessionAssets, contact *flows.Contact, log func(flows.Event)) {
	diff := make([]*flows.Group, 0, len(m.Groups))
	if m.Modification == GroupsAdd {
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
			log(events.NewContactGroupsChangedEvent(diff, nil))
		}
	} else if m.Modification == GroupsRemove {
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
			log(events.NewContactGroupsChangedEvent(nil, diff))
		}
	}
}

var _ Modifier = (*GroupsModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readGroupsModifier(assets flows.SessionAssets, data json.RawMessage) (Modifier, error) {
	// TODO
	return nil, nil
}
