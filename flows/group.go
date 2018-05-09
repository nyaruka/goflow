package flows

import (
	"encoding/json"
	"fmt"
	"strings"

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
	uuid        GroupUUID
	name        string
	query       string
	parsedQuery *contactql.ContactQuery
}

// NewGroup returns a new group object with the passed in uuid and name
func NewGroup(uuid GroupUUID, name string, query string) *Group {
	return &Group{uuid: uuid, name: name, query: query}
}

// UUID returns the UUID of the group
func (g *Group) UUID() GroupUUID { return g.uuid }

// Name returns the name of the group
func (g *Group) Name() string { return g.name }

// Query returns the query of a dynamic group
func (g *Group) Query() string { return g.query }

// ParsedQuery returns the parsed query of a dynamic group (cached)
func (g *Group) ParsedQuery() (*contactql.ContactQuery, error) {
	if g.query != "" && g.parsedQuery == nil {
		var err error
		if g.parsedQuery, err = contactql.ParseQuery(g.query); err != nil {
			return nil, err
		}
	}
	return g.parsedQuery, nil
}

// IsDynamic returns whether this group is dynamic
func (g *Group) IsDynamic() bool { return g.query != "" }

// CheckDynamicMembership returns whether the given contact belongs in this dynamic group
func (g *Group) CheckDynamicMembership(session Session, contact *Contact) (bool, error) {
	if !g.IsDynamic() {
		return false, fmt.Errorf("can't check membership on a non-dynamic group")
	}
	parsedQuery, err := g.ParsedQuery()
	if err != nil {
		return false, err
	}

	return contactql.EvaluateQuery(session.Environment(), parsedQuery, contact)
}

// Reference returns a reference to this group
func (g *Group) Reference() *GroupReference { return NewGroupReference(g.uuid, g.name) }

// Resolve resolves the given key when this group is referenced in an expression
func (g *Group) Resolve(key string) types.XValue {
	switch key {
	case "uuid":
		return types.NewXText(string(g.uuid))
	case "name":
		return types.NewXText(g.name)
	}

	return types.NewXResolveError(g, key)
}

// Repr returns the representation of this type
func (g *Group) Repr() string { return "group" }

// Reduce is called when this object needs to be reduced to a primitive
func (g *Group) Reduce() types.XPrimitive { return types.NewXText(g.name) }

// ToXJSON is called when this type is passed to @(json(...))
func (g *Group) ToXJSON() types.XText {
	return types.ResolveKeys(g, "uuid", "name").ToXJSON()
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

// Clone returns a clone of this group list
func (l *GroupList) clone() *GroupList {
	groups := make([]*Group, len(l.groups))
	copy(groups, l.groups)
	return NewGroupList(groups)
}

// FindByUUID returns the group with the passed in UUID or nil if not found
func (l *GroupList) FindByUUID(uuid GroupUUID) *Group {
	for _, group := range l.groups {
		if group.uuid == uuid {
			return group
		}
	}
	return nil
}

// Add adds the given group to this group list
func (l *GroupList) Add(group *Group) bool {
	if l.FindByUUID(group.uuid) == nil {
		l.groups = append(l.groups, group)
		return true
	}
	return false
}

// Remove removes the given group from this group list
func (l *GroupList) Remove(group *Group) bool {
	for g := range l.groups {
		if l.groups[g].uuid == group.uuid {
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

// Repr returns the representation of this type
func (l GroupList) Repr() string { return "groups" }

// Reduce is called when this object needs to be reduced to a primitive
func (l GroupList) Reduce() types.XPrimitive {
	array := types.NewXArray()
	for _, group := range l.groups {
		array.Append(group)
	}
	return array
}

// ToXJSON is called when this type is passed to @(json(...))
func (l GroupList) ToXJSON() types.XText {
	return l.Reduce().ToXJSON()
}

var _ types.XValue = (*GroupList)(nil)
var _ types.XIndexable = (*GroupList)(nil)

// GroupSet defines the unordered set of all groups for a session
type GroupSet struct {
	groups       []*Group
	groupsByUUID map[GroupUUID]*Group
}

// NewGroupSet creates a new group set from the given list of groups
func NewGroupSet(groups []*Group) *GroupSet {
	s := &GroupSet{groups: groups, groupsByUUID: make(map[GroupUUID]*Group, len(groups))}
	for _, group := range s.groups {
		s.groupsByUUID[group.uuid] = group
	}
	return s
}

// All returns all groups in this group set
func (s *GroupSet) All() []*Group {
	return s.groups
}

// FindByUUID finds the group with the given UUID
func (s *GroupSet) FindByUUID(uuid GroupUUID) *Group {
	return s.groupsByUUID[uuid]
}

// FindByName looks for a group with the given name (case-insensitive)
func (s *GroupSet) FindByName(name string) *Group {
	name = strings.ToLower(name)
	for _, group := range s.groups {
		if strings.ToLower(group.name) == name {
			return group
		}
	}
	return nil
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type groupEnvelope struct {
	UUID  GroupUUID `json:"uuid" validate:"required,uuid4"`
	Name  string    `json:"name"`
	Query string    `json:"query,omitempty"`
}

// ReadGroup reads a group from the given JSON
func ReadGroup(data json.RawMessage) (*Group, error) {
	var ge groupEnvelope
	if err := utils.UnmarshalAndValidate(data, &ge, "group"); err != nil {
		return nil, err
	}

	return NewGroup(ge.UUID, ge.Name, ge.Query), nil
}

// ReadGroupSet reads a group set from the given JSON
func ReadGroupSet(data json.RawMessage) (*GroupSet, error) {
	items, err := utils.UnmarshalArray(data)
	if err != nil {
		return nil, err
	}

	groups := make([]*Group, len(items))
	for d := range items {
		if groups[d], err = ReadGroup(items[d]); err != nil {
			return nil, err
		}
	}

	return NewGroupSet(groups), nil
}
