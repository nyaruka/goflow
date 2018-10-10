// Package static is an implementation of AssetSource which loads assets from a static JSON file.
package static

import (
	"fmt"
	"io/ioutil"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// StaticSource is an asset source which loads assets from a static JSON file
type StaticSource struct {
	s struct {
		Channels  []assets.Channel           `json:"channels" validate:"omitempty,dive"`
		Fields    []assets.Field             `json:"fields" validate:"omitempty,dive"`
		Flows     []assets.Flow              `json:"flows" validate:"omitempty,dive"`
		Groups    []assets.Group             `json:"groups" validate:"omitempty,dive"`
		Labels    []assets.Label             `json:"labels" validate:"omitempty,dive"`
		Locations []assets.LocationHierarchy `json:"locations" validate:"omitempty,dive"`
		Resthooks []assets.Resthook          `json:"resthooks" validate:"omitempty,dive"`
	}
}

// NewStaticSource creates a new static loaded from the given JSON file
func NewStaticSource(path string) (*StaticSource, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file '%s': %s", path, err)
	}

	s := &StaticSource{}
	if err := utils.UnmarshalAndValidate(data, s); err != nil {
		return nil, fmt.Errorf("unable to read assets: %s", err)
	}
	return s, nil
}

var _ assets.AssetSource = (*StaticSource)(nil)

// Channels returns all channel assets
func (s *StaticSource) Channels() ([]assets.Channel, error) {
	return s.s.Channels, nil
}

// Fields returns all field assets
func (s *StaticSource) Fields() ([]assets.Field, error) {
	return s.s.Fields, nil
}

// Flow returns the flow asset with the given UUID
func (s *StaticSource) Flow(uuid assets.FlowUUID) (assets.Flow, error) {
	for _, flow := range s.s.Flows {
		if flow.UUID() == uuid {
			return flow, nil
		}
	}
	return nil, fmt.Errorf("no such flow with UUID: %s", uuid)
}

// Groups returns all group assets
func (s *StaticSource) Groups() ([]assets.Group, error) {
	return s.s.Groups, nil
}

// Labels returns all label assets
func (s *StaticSource) Labels() ([]assets.Label, error) {
	return s.s.Labels, nil
}

// Locations returns all location assets
func (s *StaticSource) Locations() ([]assets.LocationHierarchy, error) {
	return s.s.Locations, nil
}

// Resthooks returns all resthook assets
func (s *StaticSource) Resthooks() ([]assets.Resthook, error) {
	return s.s.Resthooks, nil
}
