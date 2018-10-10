// Package static is an implementation of AssetSource which loads assets from a static JSON file.
package static

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static/types"
	"github.com/nyaruka/goflow/utils"
)

// StaticSource is an asset source which loads assets from a static JSON file
type StaticSource struct {
	s struct {
		Channels  json.RawMessage   `json:"channels"`
		Fields    json.RawMessage   `json:"fields"`
		Flows     []json.RawMessage `json:"flows"`
		Groups    json.RawMessage   `json:"groups"`
		Labels    json.RawMessage   `json:"labels"`
		Locations json.RawMessage   `json:"locations"`
		Resthooks json.RawMessage   `json:"resthooks"`
	}
}

// NewStaticSource creates a new static source from the given JSON
func NewStaticSource(data json.RawMessage) (*StaticSource, error) {
	s := &StaticSource{}
	if err := utils.UnmarshalAndValidate(data, &s.s); err != nil {
		return nil, fmt.Errorf("unable to read assets: %s", err)
	}
	return s, nil
}

// LoadStaticSource loads a new static source from the given JSON file
func LoadStaticSource(path string) (*StaticSource, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file '%s': %s", path, err)
	}
	return NewStaticSource(data)
}

var _ assets.AssetSource = (*StaticSource)(nil)

// Channels returns all channel assets
func (s *StaticSource) Channels() ([]assets.Channel, error) {
	return types.ReadChannels(s.s.Channels)
}

// Fields returns all field assets
func (s *StaticSource) Fields() ([]assets.Field, error) {
	return types.ReadFields(s.s.Fields)
}

// Flow returns the flow asset with the given UUID
func (s *StaticSource) Flow(uuid assets.FlowUUID) (assets.Flow, error) {
	for _, rawFlow := range s.s.Flows {
		// TODO inefficient
		flow, err := types.ReadFlow(rawFlow)
		if err != nil {
			return nil, err
		}
		if flow.UUID() == uuid {
			return flow, nil
		}
	}
	return nil, fmt.Errorf("no such flow with UUID: %s", uuid)
}

// Groups returns all group assets
func (s *StaticSource) Groups() ([]assets.Group, error) {
	return types.ReadGroups(s.s.Groups)
}

// Labels returns all label assets
func (s *StaticSource) Labels() ([]assets.Label, error) {
	return types.ReadLabels(s.s.Labels)
}

// Locations returns all location assets
func (s *StaticSource) Locations() ([]assets.LocationHierarchy, error) {
	return types.ReadLocationHierarchies(s.s.Locations)
}

// Resthooks returns all resthook assets
func (s *StaticSource) Resthooks() ([]assets.Resthook, error) {
	return types.ReadResthooks(s.s.Resthooks)
}
