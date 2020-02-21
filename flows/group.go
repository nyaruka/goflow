package flows

import (
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"

	"github.com/pkg/errors"
)

// Group represents a grouping of contacts. It can be static (contacts are added and removed manually through
// [actions](#action:add_contact_groups)) or dynamic (contacts are added automatically by a query).
type Group struct {
	assets.Group

	cachedQuery *contactql.ContactQuery
}

// NewGroup returns a new group object from the given group asset
func NewGroup(asset assets.Group) *Group {
	return &Group{Group: asset}
}

// Asset returns the underlying asset
func (g *Group) Asset() assets.Group { return g.Group }

// the parsed query of a dynamic group (cached)
func (g *Group) parsedQuery(env envs.Environment, fields *FieldAssets) (*contactql.ContactQuery, error) {
	if g.Query() != "" && g.cachedQuery == nil {
		var err error
		if g.cachedQuery, err = contactql.ParseQuery(g.Query(), env.RedactionPolicy(), env.DefaultCountry(), fields.Resolve); err != nil {
			return nil, err
		}
	}
	return g.cachedQuery, nil
}

// IsDynamic returns whether this group is dynamic
func (g *Group) IsDynamic() bool { return g.Query() != "" }

// CheckDynamicMembership returns whether the given contact belongs in this dynamic group
func (g *Group) CheckDynamicMembership(env envs.Environment, contact *Contact, fields *FieldAssets) (bool, error) {
	if !g.IsDynamic() {
		panic("can't check membership on a non-dynamic group")
	}
	parsedQuery, err := g.parsedQuery(env, fields)
	if err != nil {
		return false, err
	}

	return contactql.EvaluateQuery(env, parsedQuery, contact)
}

// Reference returns a reference to this group
func (g *Group) Reference() *assets.GroupReference {
	if g == nil {
		return nil
	}
	return assets.NewGroupReference(g.UUID(), g.Name())
}

// ToXValue returns a representation of this object for use in expressions
//
//   uuid:text -> the UUID of the group
//   name:text -> the name of the group
//
// @context group
func (g *Group) ToXValue(env envs.Environment) types.XValue {
	return types.NewXObject(map[string]types.XValue{
		"uuid": types.NewXText(string(g.UUID())),
		"name": types.NewXText(g.Name()),
	})
}

// GroupList defines a contact's list of groups
type GroupList struct {
	groups []*Group
}

// NewGroupList creates a new group list
func NewGroupList(a SessionAssets, refs []*assets.GroupReference, missing assets.MissingCallback) *GroupList {
	groups := make([]*Group, 0, len(refs))

	for _, ref := range refs {
		group := a.Groups().Get(ref.UUID)
		if group == nil {
			missing(ref, nil)
		} else {
			groups = append(groups, group)
		}
	}
	return &GroupList{groups: groups}
}

// Clone returns a clone of this group list
func (l *GroupList) clone() *GroupList {
	groups := make([]*Group, len(l.groups))
	copy(groups, l.groups)
	return &GroupList{groups: groups}
}

// FindByUUID returns the group with the passed in UUID or nil if not found
func (l *GroupList) FindByUUID(uuid assets.GroupUUID) *Group {
	for _, group := range l.groups {
		if group.UUID() == uuid {
			return group
		}
	}
	return nil
}

// Add adds the given group to this group list
func (l *GroupList) Add(group *Group) bool {
	if l.FindByUUID(group.UUID()) == nil {
		l.groups = append(l.groups, group)
		return true
	}
	return false
}

// Remove removes the given group from this group list
func (l *GroupList) Remove(group *Group) bool {
	for i := range l.groups {
		if l.groups[i].UUID() == group.UUID() {
			l.groups = append(l.groups[:i], l.groups[i+1:]...)
			return true
		}
	}
	return false
}

// All returns all groups in this group list
func (l *GroupList) All() []*Group {
	return l.groups
}

// Count returns the number of groups in this group list
func (l *GroupList) Count() int {
	return len(l.groups)
}

// ToXValue returns a representation of this object for use in expressions
func (l GroupList) ToXValue(env envs.Environment) types.XValue {
	array := make([]types.XValue, len(l.groups))
	for i, group := range l.groups {
		array[i] = group.ToXValue(env)
	}
	return types.NewXArray(array...)
}

// GroupAssets provides access to all group assets
type GroupAssets struct {
	all    []*Group
	byUUID map[assets.GroupUUID]*Group
}

// NewGroupAssets creates a new set of group assets
func NewGroupAssets(groups []assets.Group) *GroupAssets {
	s := &GroupAssets{
		all:    make([]*Group, len(groups)),
		byUUID: make(map[assets.GroupUUID]*Group, len(groups)),
	}
	for i, asset := range groups {
		group := NewGroup(asset)
		s.all[i] = group
		s.byUUID[group.UUID()] = group
	}
	return s
}

// All returns all the groups
func (s *GroupAssets) All() []*Group {
	return s.all
}

// Get returns the group with the given UUID
func (s *GroupAssets) Get(uuid assets.GroupUUID) *Group {
	return s.byUUID[uuid]
}

// FindByName looks for a group with the given name (case-insensitive)
func (s *GroupAssets) FindByName(name string) *Group {
	name = strings.ToLower(name)
	for _, group := range s.all {
		if strings.ToLower(group.Name()) == name {
			return group
		}
	}
	return nil
}
