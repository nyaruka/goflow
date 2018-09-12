package types

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// json serializable implementation of a location asset
type location struct {
	name     string
	aliases  []string
	children []assets.Location
}

// Name returns the name of the location
func (l *location) Name() string { return l.name }

// Aliases gets the aliases of this location
func (l *location) Aliases() []string { return l.aliases }

// Children gets the children of this location
func (l *location) Children() []assets.Location { return l.children }

type locationEnvelope struct {
	Name     string              `json:"name" validate:"required"`
	Aliases  []string            `json:"aliases,omitempty"`
	Children []*locationEnvelope `json:"children,omitempty"`
}

func locationFromEnvelope(envelope *locationEnvelope) assets.Location {
	location := &location{
		name:     envelope.Name,
		aliases:  envelope.Aliases,
		children: make([]assets.Location, len(envelope.Children)),
	}
	for c := range envelope.Children {
		location.children[c] = locationFromEnvelope(envelope.Children[c])
	}
	return location
}

// ReadLocation reads a location from the given JSON
func ReadLocation(data json.RawMessage) (assets.Location, error) {
	le := &locationEnvelope{}
	if err := utils.UnmarshalAndValidate(data, le); err != nil {
		return nil, err
	}

	return locationFromEnvelope(le), nil
}

// ReadLocations reads locations from the given JSON
func ReadLocations(data json.RawMessage) ([]assets.Location, error) {
	items, err := utils.UnmarshalArray(data)
	if err != nil {
		return nil, err
	}

	locations := make([]assets.Location, len(items))
	for d := range items {
		if locations[d], err = ReadLocation(items[d]); err != nil {
			return nil, err
		}
	}

	return locations, nil
}
