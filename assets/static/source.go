// Package static is an implementation of Source which loads assets from a static JSON file.
package static

import (
	"encoding/json"
	"io/ioutil"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static/types"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

// StaticSource is an asset source which loads assets from a static JSON file
type StaticSource struct {
	s struct {
		Channels    []*types.Channel          `json:"channels" validate:"omitempty,dive"`
		Classifiers []*types.Classifier       `json:"classifiers" validate:"omitempty,dive"`
		Fields      []*types.Field            `json:"fields" validate:"omitempty,dive"`
		Flows       []*types.Flow             `json:"flows" validate:"omitempty,dive"`
		Globals     []*types.Global           `json:"globals" validate:"omitempty,dive"`
		Groups      []*types.Group            `json:"groups" validate:"omitempty,dive"`
		Labels      []*types.Label            `json:"labels" validate:"omitempty,dive"`
		Locations   []*envs.LocationHierarchy `json:"locations"`
		Resthooks   []*types.Resthook         `json:"resthooks" validate:"omitempty,dive"`
		Templates   []*types.Template         `json:"templates" validate:"omitempty,dive"`
		Ticketers   []*types.Ticketer         `json:"ticketers" validate:"omitempty,dive"`
		Users       []*types.User             `json:"users" validate:"omitempty,dive"`
	}
}

// NewEmptySource creates a new empty source with no assets
func NewEmptySource() *StaticSource {
	return &StaticSource{}
}

// NewSource creates a new static source from the given JSON
func NewSource(data json.RawMessage) (*StaticSource, error) {
	s := &StaticSource{}
	if err := utils.UnmarshalAndValidate(data, &s.s); err != nil {
		return nil, errors.Wrap(err, "unable to read assets")
	}
	return s, nil
}

// LoadSource loads a new static source from the given JSON file
func LoadSource(path string) (*StaticSource, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading file '%s'", path)
	}
	return NewSource(data)
}

var _ assets.Source = (*StaticSource)(nil)

// Channels returns all channel assets
func (s *StaticSource) Channels() ([]assets.Channel, error) {
	set := make([]assets.Channel, len(s.s.Channels))
	for i := range s.s.Channels {
		set[i] = s.s.Channels[i]
	}
	return set, nil
}

// Classifiers returns all classifier assets
func (s *StaticSource) Classifiers() ([]assets.Classifier, error) {
	set := make([]assets.Classifier, len(s.s.Classifiers))
	for i := range s.s.Classifiers {
		set[i] = s.s.Classifiers[i]
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
	return nil, errors.Errorf("no such flow with UUID '%s'", uuid)
}

// Globals returns all global assets
func (s *StaticSource) Globals() ([]assets.Global, error) {
	set := make([]assets.Global, len(s.s.Globals))
	for i := range s.s.Globals {
		set[i] = s.s.Globals[i]
	}
	return set, nil
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

// Templates returns all template assets
func (s *StaticSource) Templates() ([]assets.Template, error) {
	set := make([]assets.Template, len(s.s.Templates))
	for i := range s.s.Templates {
		set[i] = s.s.Templates[i]
	}
	return set, nil
}

// Ticketers returns all ticketer assets
func (s *StaticSource) Ticketers() ([]assets.Ticketer, error) {
	set := make([]assets.Ticketer, len(s.s.Ticketers))
	for i := range s.s.Ticketers {
		set[i] = s.s.Ticketers[i]
	}
	return set, nil
}

// Users returns all user assets
func (s *StaticSource) Users() ([]assets.User, error) {
	set := make([]assets.User, len(s.s.Users))
	for i := range s.s.Users {
		set[i] = s.s.Users[i]
	}
	return set, nil
}
