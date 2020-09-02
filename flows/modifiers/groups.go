package modifiers

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeGroups, readGroupsModifier)
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

	groups       []*flows.Group
	modification GroupsModification
}

// NewGroups creates a new groups modifier
func NewGroups(groups []*flows.Group, modification GroupsModification) *GroupsModifier {
	return &GroupsModifier{
		baseModifier: newBaseModifier(TypeGroups),
		groups:       groups,
		modification: modification,
	}
}

// Apply applies this modification to the given contact
func (m *GroupsModifier) Apply(env envs.Environment, assets flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) {
	if contact.Status() == flows.ContactStatusBlocked || contact.Status() == flows.ContactStatusStopped {
		log(events.NewErrorf("can't add blocked or stopped contacts to groups"))
		return
	}

	diff := make([]*flows.Group, 0, len(m.groups))

	if m.modification == GroupsAdd {
		for _, group := range m.groups {
			if group.UsesQuery() {
				log(events.NewErrorf("can't add contacts to the query based group '%s'", group.Name()))
				continue
			}

			// ignore group if contact is already in it
			if contact.Groups().FindByUUID(group.UUID()) != nil {
				continue
			}

			contact.Groups().Add(group)
			diff = append(diff, group)
		}

		// only generate event if contact's groups change
		if len(diff) > 0 {
			log(events.NewContactGroupsChanged(diff, nil))
		}

	} else if m.modification == GroupsRemove {
		for _, group := range m.groups {
			if group.UsesQuery() {
				log(events.NewErrorf("can't remove contacts from the query based group '%s'", group.Name()))
				continue
			}

			// ignore group if contact isn't actually in it
			if contact.Groups().FindByUUID(group.UUID()) == nil {
				continue
			}

			contact.Groups().Remove(group)
			diff = append(diff, group)
		}

		// only generate event if contact's groups change
		if len(diff) > 0 {
			log(events.NewContactGroupsChanged(nil, diff))
		}
	}
}

var _ flows.Modifier = (*GroupsModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type groupsModifierEnvelope struct {
	utils.TypedEnvelope
	Groups       []*assets.GroupReference `json:"groups" validate:"required,dive"`
	Modification GroupsModification       `json:"modification" validate:"eq=add|eq=remove"`
}

func readGroupsModifier(assets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &groupsModifierEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	groups := make([]*flows.Group, 0, len(e.Groups))
	for _, groupRef := range e.Groups {
		group := assets.Groups().Get(groupRef.UUID)
		if group == nil {
			missing(groupRef, nil)
		} else {
			groups = append(groups, group)
		}
	}

	if len(groups) > 0 {
		return NewGroups(groups, e.Modification), nil
	}

	return nil, ErrNoModifier // nothing left to modify if there are no groups
}

func (m *GroupsModifier) MarshalJSON() ([]byte, error) {
	groupRefs := make([]*assets.GroupReference, len(m.groups))
	for i := range m.groups {
		groupRefs[i] = m.groups[i].Reference()
	}

	return jsonx.Marshal(&groupsModifierEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Groups:        groupRefs,
		Modification:  m.modification,
	})
}
