package core

import (
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// Group adds some functionality to group assets.
type Group struct {
	assets.Group

	parsedQuery *contactql.ContactQuery
}

// NewGroup returns a new group object from the given group asset and - if it's query based - its parsed
// query. Parsing is the caller's responsibility so that this package doesn't depend on query parsing.
func NewGroup(asset assets.Group, query *contactql.ContactQuery) *Group {
	return &Group{Group: asset, parsedQuery: query}
}

// Asset returns the underlying asset
func (g *Group) Asset() assets.Group { return g.Group }

// UsesQuery returns whether this group is query based
func (g *Group) UsesQuery() bool { return g.Query() != "" }

// CheckQueryBasedMembership returns whether the given contact belongs in a query based group
func (g *Group) CheckQueryBasedMembership(env envs.Environment, contact *Contact) bool {
	if !g.UsesQuery() {
		panic("can't check membership on a non-query based group")
	}

	if contact.Status() != ContactStatusActive {
		return false
	}

	return contactql.EvaluateQuery(env, g.parsedQuery, contact)
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
//	uuid:text -> the UUID of the group
//	name:text -> the name of the group
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
func NewGroupList(groupAssets *GroupAssets, refs []*assets.GroupReference, missing assets.MissingCallback) *GroupList {
	groups := make([]*Group, 0, len(refs))

	for _, ref := range refs {
		group := groupAssets.Get(ref.UUID)
		if group == nil {
			missing(ref, nil)
		} else {
			groups = append(groups, group)
		}
	}
	return &GroupList{groups: groups}
}

// Clone returns a clone of this group list
func (l *GroupList) Clone() *GroupList {
	groups := make([]*Group, len(l.groups))
	copy(groups, l.groups)
	return &GroupList{groups: groups}
}

// References returns this group list as a slice of group references
func (l *GroupList) References() []*assets.GroupReference {
	refs := make([]*assets.GroupReference, len(l.groups))
	for i, group := range l.groups {
		refs[i] = group.Reference()
	}
	return refs
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

// Clear clears this group list
func (l *GroupList) Clear() {
	l.groups = []*Group{}
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
func (l *GroupList) ToXValue(env envs.Environment) types.XValue {
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
func NewGroupAssets(groups []*Group) *GroupAssets {
	s := &GroupAssets{
		all:    groups,
		byUUID: make(map[assets.GroupUUID]*Group, len(groups)),
	}
	for _, group := range groups {
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

// GroupReferences converts a slice of groups to a slice of references
func GroupReferences(groups []*Group) []*assets.GroupReference {
	refs := make([]*assets.GroupReference, len(groups))
	for i := range groups {
		refs[i] = groups[i].Reference()
	}
	return refs
}
