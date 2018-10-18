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
		Channels  []*types.Channel           `json:"channels" validate:"omitempty,dive"`
		Fields    []*types.Field             `json:"fields" validate:"omitempty,dive"`
		Flows     []*types.Flow              `json:"flows" validate:"omitempty,dive"`
		Groups    []*types.Group             `json:"groups" validate:"omitempty,dive"`
		Labels    []*types.Label             `json:"labels" validate:"omitempty,dive"`
		Locations []*utils.LocationHierarchy `json:"locations"`
		Resthooks []*types.Resthook          `json:"resthooks" validate:"omitempty,dive"`
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
	set := make([]assets.Channel, len(s.s.Channels))
	for i := range s.s.Channels {
		set[i] = s.s.Channels[i]
	}
	return set, nil
}

// Fields returns all field assets
func (s *StaticSource) Fields() ([]assets.Field, error) {
	set := make([]assets.Field, len(s.s.Fields))
	for i := range s.s.Fields {
		set[i] = s.s.Fields[i]
	}
	return set, nil
}

// Flow returns the flow asset with the given UUID
func (s *StaticSource) Flow(uuid assets.FlowUUID) (assets.Flow, error) {
	for _, flow := range s.s.Flows {
		if flow.UUID() == uuid {
			return flow, nil
		}
	}
	return nil, fmt.Errorf("no such flow with UUID '%s'", uuid)
}

// Groups returns all group assets
func (s *StaticSource) Groups() ([]assets.Group, error) {
	set := make([]assets.Group, len(s.s.Groups))
	for i := range s.s.Groups {
		set[i] = s.s.Groups[i]
	}
	return set, nil
}

// Labels returns all label assets
func (s *StaticSource) Labels() ([]assets.Label, error) {
	set := make([]assets.Label, len(s.s.Labels))
	for i := range s.s.Labels {
		set[i] = s.s.Labels[i]
	}
	return set, nil
}

// Locations returns all location assets
func (s *StaticSource) Locations() ([]assets.LocationHierarchy, error) {
	set := make([]assets.LocationHierarchy, len(s.s.Locations))
	for i := range s.s.Locations {
		set[i] = s.s.Locations[i]
	}
	return set, nil
}

// Resthooks returns all resthook assets
func (s *StaticSource) Resthooks() ([]assets.Resthook, error) {
	set := make([]assets.Resthook, len(s.s.Resthooks))
	for i := range s.s.Resthooks {
		set[i] = s.s.Resthooks[i]
	}
	return set, nil
}
