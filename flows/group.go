package flows

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// Group represents a grouping of contacts. It can be static (contacts are added and removed manually through
// [actions](#action:add_contact_groups)) or dynamic (contacts are added automatically by a query). It renders as its name in a
// template, and has the following properties which can be accessed:
//
//  * `uuid` the UUID of the group
//  * `name` the name of the group
//
// Examples:
//
//   @contact.groups -> ["Testers","Males"]
//   @contact.groups.0.uuid -> b7cf0d83-f1c9-411c-96fd-c511a4cfa86d
//   @contact.groups.1.name -> Males
//   @(json(contact.groups.1)) -> {"name":"Males","uuid":"4f1f98fc-27a7-4a69-bbdb-24744ba739a9"}
//
// @context group
type Group struct {
	assets.Group

	parsedQuery *contactql.ContactQuery
}

// NewGroup returns a new group object from the given group asset
func NewGroup(asset assets.Group) *Group {
	return &Group{Group: asset}
}

// Asset returns the underlying asset
func (g *Group) Asset() assets.Group { return g.Group }

// ParsedQuery returns the parsed query of a dynamic group (cached)
func (g *Group) ParsedQuery() (*contactql.ContactQuery, error) {
	if g.Query() != "" && g.parsedQuery == nil {
		var err error
		if g.parsedQuery, err = contactql.ParseQuery(g.Query()); err != nil {
			return nil, err
		}
	}
	return g.parsedQuery, nil
}

// IsDynamic returns whether this group is dynamic
func (g *Group) IsDynamic() bool { return g.Query() != "" }

// CheckDynamicMembership returns whether the given contact belongs in this dynamic group
func (g *Group) CheckDynamicMembership(env utils.Environment, contact *Contact) (bool, error) {
	if !g.IsDynamic() {
		return false, fmt.Errorf("can't check membership on a non-dynamic group")
	}
	parsedQuery, err := g.ParsedQuery()
	if err != nil {
		return false, err
	}

	return contactql.EvaluateQuery(env, parsedQuery, contact)
}

// Reference returns a reference to this group
func (g *Group) Reference() *assets.GroupReference {
	return assets.NewGroupReference(g.UUID(), g.Name())
}

// Resolve resolves the given key when this group is referenced in an expression
func (g *Group) Resolve(env utils.Environment, key string) types.XValue {
	switch key {
	case "uuid":
		return types.NewXText(string(g.UUID()))
	case "name":
		return types.NewXText(g.Name())
	}

	return types.NewXResolveError(g, key)
}

// Describe returns a representation of this type for error messages
func (g *Group) Describe() string { return "group" }

// Reduce is called when this object needs to be reduced to a primitive
func (g *Group) Reduce(env utils.Environment) types.XPrimitive { return types.NewXText(g.Name()) }

// ToXJSON is called when this type is passed to @(json(...))
func (g *Group) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, g, "uuid", "name").ToXJSON(env)
}

var _ types.XValue = (*Group)(nil)
var _ types.XResolvable = (*Group)(nil)

// GroupList defines a contact's list of groups
type GroupList struct {
	groups []*Group
}

// NewGroupList creates a new group list
func NewGroupList(groups []*Group) *GroupList {
	return &GroupList{groups: groups}
}

// NewGroupListFromAssets creates a new group list
func NewGroupListFromAssets(a SessionAssets, groupAssets []assets.Group) (*GroupList, error) {
	groups := make([]*Group, len(groupAssets))
	var err error

	for g, asset := range groupAssets {
		if groups[g], err = a.Groups().Get(asset.UUID()); err != nil {
			return nil, err
		}
	}
	return &GroupList{groups: groups}, nil
}

// Clone returns a clone of this group list
func (l *GroupList) clone() *GroupList {
	groups := make([]*Group, len(l.groups))
	copy(groups, l.groups)
	return NewGroupList(groups)
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
	for g := range l.groups {
		if l.groups[g].UUID() == group.UUID() {
			l.groups = append(l.groups[:g], l.groups[g+1:]...)
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

// Index is called when this object is indexed into in an expression
func (l *GroupList) Index(index int) types.XValue {
	return l.groups[index]
}

// Length is called when the length of this object is requested in an expression
func (l *GroupList) Length() int {
	return len(l.groups)
}

// Describe returns a representation of this type for error messages
func (l GroupList) Describe() string { return "groups" }

// Reduce is called when this object needs to be reduced to a primitive
func (l GroupList) Reduce(env utils.Environment) types.XPrimitive {
	array := types.NewXArray()
	for _, group := range l.groups {
		array.Append(group)
	}
	return array
}

// ToXJSON is called when this type is passed to @(json(...))
func (l GroupList) ToXJSON(env utils.Environment) types.XText {
	return l.Reduce(env).ToXJSON(env)
}

var _ types.XValue = (*GroupList)(nil)
var _ types.XIndexable = (*GroupList)(nil)

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
	for g, asset := range groups {
		group := NewGroup(asset)
		s.all[g] = group
		s.byUUID[group.UUID()] = group
	}
	return s
}

// All returns all the groups
func (s *GroupAssets) All() []*Group {
	return s.all
}

// Get returns the group with the given UUID
func (s *GroupAssets) Get(uuid assets.GroupUUID) (*Group, error) {
	c, found := s.byUUID[uuid]
	if !found {
		return nil, fmt.Errorf("no such group with uuid '%s'", uuid)
	}
	return c, nil
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
