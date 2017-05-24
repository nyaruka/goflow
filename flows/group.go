package flows

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/utils"
	validator "gopkg.in/go-playground/validator.v9"
)

func init() {
	utils.Validator.RegisterStructValidation(ValidateGroup, Group{})
}

// Group represents a grouping of contacts. From an engine perspective the only piece that matter is
// the UUID of the group and its name
type Group struct {
	uuid GroupUUID
	name string
}

// ValidateGroup is our global validator for our group struct
func ValidateGroup(sl validator.StructLevel) {
	group := sl.Current().Interface().(Group)
	if len(group.uuid) == 0 {
		sl.ReportError(group.uuid, "uuid", "uuid", "uuid4", "")
	}
	if len(group.name) == 0 {
		sl.ReportError(group.uuid, "name", "name", "required", "")
	}
}

// NewGroup returns a new group object with the passed in uuid and name
func NewGroup(uuid GroupUUID, name string) *Group {
	return &Group{uuid, name}
}

// Name returns the name of the group
func (g *Group) Name() string { return g.name }

// UUID returns the UUID of the group
func (g *Group) UUID() GroupUUID { return g.uuid }

// Resolve resolves the passed in key to a value
func (g *Group) Resolve(key string) interface{} {
	switch key {

	case "name":
		return g.name

	case "uuid":
		return g.uuid
	}

	return fmt.Errorf("no field '%s' on group", key)
}

// Default returns the default value for this group
func (g *Group) Default() interface{} { return g }

// String satisfies the stringer interface returning the name of the group
func (g *Group) String() string { return g.name }

// GroupList defines a list of groups
type GroupList []*Group

// FindGroup returns the group with the passed in UUID or nil if not found
func (l GroupList) FindGroup(uuid GroupUUID) *Group {
	for i := range l {
		if l[i].uuid == uuid {
			return l[i]
		}
	}
	return nil
}

// Resolve looks up the passed in key for the group list, we attempt to find the group with the uuid
// of the passed in key, which can be used for testing whether a contact is part of a group
func (l GroupList) Resolve(key string) interface{} {
	uuid := GroupUUID(key)
	return l.FindGroup(uuid)
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

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type groupEnvelope struct {
	UUID GroupUUID `json:"uuid"`
	Name string    `json:"name"`
}

// UnmarshalJSON unmarshals the group from the passed in json
func (g *Group) UnmarshalJSON(data []byte) error {
	var ge groupEnvelope
	var err error

	err = json.Unmarshal(data, &ge)
	g.uuid = ge.UUID
	g.name = ge.Name

	return err
}

// MarshalJSON marshals the Group into json
func (g *Group) MarshalJSON() ([]byte, error) {
	var ge groupEnvelope

	ge.Name = g.name
	ge.UUID = g.uuid

	return json.Marshal(ge)
}
