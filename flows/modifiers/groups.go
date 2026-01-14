package modifiers

import (
	"context"
	"fmt"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeGroups, readGroups)
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

// Groups modifies the group membership of the contact
type Groups struct {
	baseModifier

	groups       []*flows.Group
	modification GroupsModification
}

// NewGroups creates a new groups modifier
func NewGroups(groups []*flows.Group, modification GroupsModification) *Groups {
	return &Groups{
		baseModifier: newBaseModifier(TypeGroups),
		groups:       groups,
		modification: modification,
	}
}

// Apply applies this modification to the given contact
func (m *Groups) Apply(ctx context.Context, eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventLogger) (bool, error) {
	if contact.Status() == flows.ContactStatusBlocked || contact.Status() == flows.ContactStatusStopped {
		log(events.NewError("Can't add blocked or stopped contacts to groups", ""))
		return false, nil
	}

	diff := make([]*flows.Group, 0, len(m.groups))

	if m.modification == GroupsAdd {
		for _, group := range m.groups {
			if group.UsesQuery() {
				log(events.NewError(fmt.Sprintf("Can't add contacts to the query based group '%s'", group.Name()), ""))
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
			return true, nil
		}

	} else if m.modification == GroupsRemove {
		for _, group := range m.groups {
			if group.UsesQuery() {
				log(events.NewError(fmt.Sprintf("Can't remove contacts from the query based group '%s'", group.Name()), ""))
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
			return true, nil
		}
	}

	return false, nil
}

var _ flows.Modifier = (*Groups)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type groupsEnvelope struct {
	utils.TypedEnvelope

	Groups       []*assets.GroupReference `json:"groups" validate:"required,dive"`
	Modification GroupsModification       `json:"modification" validate:"eq=add|eq=remove"`
}

func readGroups(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &groupsEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	groups := make([]*flows.Group, 0, len(e.Groups))
	for _, groupRef := range e.Groups {
		group := sa.Groups().Get(groupRef.UUID)
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

func (m *Groups) MarshalJSON() ([]byte, error) {
	groupRefs := make([]*assets.GroupReference, len(m.groups))
	for i := range m.groups {
		groupRefs[i] = m.groups[i].Reference()
	}

	return jsonx.Marshal(&groupsEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Groups:        groupRefs,
		Modification:  m.modification,
	})
}
