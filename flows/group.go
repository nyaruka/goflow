package flows

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/nyaruka/goflow/utils"
)

// GroupReference is a reference to group used in a flow action or event
type GroupReference struct {
	UUID GroupUUID `json:"uuid,omitempty" validate:"omitempty,uuid4"`
	Name string    `json:"name"`
}

func NewGroupReference(uuid GroupUUID, name string) *GroupReference {
	return &GroupReference{UUID: uuid, Name: name}
}

// Group represents a grouping of contacts
type Group struct {
	uuid  GroupUUID
	name  string
	query string
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

// IsDynamic returns whether this group is dynamic
func (g *Group) IsDynamic() bool { return g.query != "" }

// Resolve resolves the passed in key to a value
func (g *Group) Resolve(key string) interface{} {
	switch key {
	case "uuid":
		return g.uuid
	case "name":
		return g.name
	}

	return fmt.Errorf("no field '%s' on group", key)
}

// Default returns the default value for this group
func (g *Group) Default() interface{} { return g }

// String satisfies the stringer interface returning the name of the group
func (g *Group) String() string { return g.name }

var _ utils.VariableResolver = (*Group)(nil)

// GroupList defines a contact's list of groups
type GroupList []*Group

// FindGroup returns the group with the passed in UUID or nil if not found
func (l GroupList) FindByUUID(uuid GroupUUID) *Group {
	for i := range l {
		if l[i].uuid == uuid {
			return l[i]
		}
	}
	return nil
}

// Resolve looks up the passed in key for the group list, which must be either "count" or a numerical index
func (l GroupList) Resolve(key string) interface{} {
	if key == "count" {
		return len(l)
	}

	// key must be a numerical index
	i, err := strconv.Atoi(key)
	if err != nil {
		return fmt.Errorf("not a valid integer '%s'", key)
	}
	if i < len(l) {
		return l[i]
	}
	return nil
}

// Default returns the default value for this group, which is our entire list
func (l GroupList) Default() interface{} {
	return l
}

// String stringifies the group list, joining our names with a comma
func (l GroupList) String() string {
	names := make([]string, len(l))
	for i := range l {
		names[i] = l[i].name
	}
	return strings.Join(names, ", ")
}

var _ utils.VariableResolver = (GroupList)(nil)

// GroupSet defines the unordered set of all groups for a session
type GroupSet struct {
	groups       []*Group
	groupsByUUID map[GroupUUID]*Group
}

func NewGroupSet(groups []*Group) *GroupSet {
	s := &GroupSet{groups: groups, groupsByUUID: make(map[GroupUUID]*Group, len(groups))}
	for _, group := range s.groups {
		s.groupsByUUID[group.uuid] = group
	}
	return s
}

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

func ReadGroup(data json.RawMessage) (*Group, error) {
	var ge groupEnvelope
	if err := utils.UnmarshalAndValidate(data, &ge, "group"); err != nil {
		return nil, err
	}

	return NewGroup(ge.UUID, ge.Name, ge.Query), nil
}

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
