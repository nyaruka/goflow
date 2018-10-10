package types

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// Group is a JSON serializable implementation of a group asset
type Group struct {
	UUID_  assets.GroupUUID `json:"uuid" validate:"required,uuid4"`
	Name_  string           `json:"name"`
	Query_ string           `json:"query,omitempty"`
}

// NewGroup creates a new group from the passed in UUID, name and query
func NewGroup(uuid assets.GroupUUID, name string, query string) assets.Group {
	return &Group{UUID_: uuid, Name_: name, Query_: query}
}

// UUID returns the UUID of the group
func (g *Group) UUID() assets.GroupUUID { return g.UUID_ }

// Name returns the name of the group
func (g *Group) Name() string { return g.Name_ }

// Query returns the query of a dynamic group
func (g *Group) Query() string { return g.Query_ }

// ReadGroup reads a group from the given JSON
func ReadGroup(data json.RawMessage) (assets.Group, error) {
	g := &Group{}
	if err := utils.UnmarshalAndValidate(data, g); err != nil {
		return nil, fmt.Errorf("unable to read group: %s", err)
	}
	return g, nil
}

// ReadGroups reads groups from the given JSON
func ReadGroups(data json.RawMessage) ([]assets.Group, error) {
	items, err := utils.UnmarshalArray(data)
	if err != nil {
		return nil, err
	}

	groups := make([]assets.Group, len(items))
	for d := range items {
		if groups[d], err = ReadGroup(items[d]); err != nil {
			return nil, err
		}
	}

	return groups, nil
}
