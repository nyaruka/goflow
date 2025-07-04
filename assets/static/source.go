// Package static is an implementation of Source which loads assets from a static JSON file.
package static

import (
	"fmt"
	"os"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
)

// StaticSource is an asset source which loads assets from a static JSON file
type StaticSource struct {
	s struct {
		Campaigns   []*Campaign               `json:"campaigns" validate:"omitempty,dive"`
		Channels    []*Channel                `json:"channels" validate:"omitempty,dive"`
		Classifiers []*Classifier             `json:"classifiers" validate:"omitempty,dive"`
		Fields      []*Field                  `json:"fields" validate:"omitempty,dive"`
		Flows       []*Flow                   `json:"flows" validate:"omitempty,dive"`
		Globals     []*Global                 `json:"globals" validate:"omitempty,dive"`
		Groups      []*Group                  `json:"groups" validate:"omitempty,dive"`
		Labels      []*Label                  `json:"labels" validate:"omitempty,dive"`
		LLMs        []*LLM                    `json:"llms" validate:"omitempty,dive"`
		Locations   []*envs.LocationHierarchy `json:"locations"`
		OptIns      []*OptIn                  `json:"optins" validate:"omitempty,dive"`
		Resthooks   []*Resthook               `json:"resthooks" validate:"omitempty,dive"`
		Templates   []*Template               `json:"templates" validate:"omitempty,dive"`
		Topics      []*Topic                  `json:"topics" validate:"omitempty,dive"`
		Users       []*User                   `json:"users" validate:"omitempty,dive"`
	}
}

// NewEmptySource creates a new empty source with no assets
func NewEmptySource() *StaticSource {
	return &StaticSource{}
}

// NewSource creates a new static source from the given JSON
func NewSource(data []byte) (*StaticSource, error) {
	s := &StaticSource{}
	if err := utils.UnmarshalAndValidate(data, &s.s); err != nil {
		return nil, fmt.Errorf("unable to read assets: %w", err)
	}
	return s, nil
}

// LoadSource loads a new static source from the given JSON file
func LoadSource(path string) (*StaticSource, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file '%s': %w", path, err)
	}
	return NewSource(data)
}

// Campaigns returns all campaign assets
func (s *StaticSource) Campaigns() ([]assets.Campaign, error) {
	set := make([]assets.Campaign, len(s.s.Campaigns))
	for i := range s.s.Campaigns {
		set[i] = s.s.Campaigns[i]
	}
	return set, nil
}

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
func (s *StaticSource) FlowByUUID(uuid assets.FlowUUID) (assets.Flow, error) {
	for _, flow := range s.s.Flows {
		if flow.UUID() == uuid {
			return flow, nil
		}
	}
	return nil, fmt.Errorf("no such flow with UUID '%s'", uuid)
}

// Flow returns the flow asset with the given UUID
func (s *StaticSource) FlowByName(name string) (assets.Flow, error) {
	for _, flow := range s.s.Flows {
		if strings.EqualFold(flow.Name(), name) {
			return flow, nil
		}
	}
	return nil, fmt.Errorf("no such flow with name '%s'", name)
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

// LLMs returns all LLM assets
func (s *StaticSource) LLMs() ([]assets.LLM, error) {
	set := make([]assets.LLM, len(s.s.LLMs))
	for i := range s.s.LLMs {
		set[i] = s.s.LLMs[i]
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

// OptIns returns all optin assets
func (s *StaticSource) OptIns() ([]assets.OptIn, error) {
	set := make([]assets.OptIn, len(s.s.OptIns))
	for i := range s.s.OptIns {
		set[i] = s.s.OptIns[i]
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

// Topics returns all topic assets
func (s *StaticSource) Topics() ([]assets.Topic, error) {
	set := make([]assets.Topic, len(s.s.Topics))
	for i := range s.s.Topics {
		set[i] = s.s.Topics[i]
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

var _ assets.Source = (*StaticSource)(nil)
