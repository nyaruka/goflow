package engine

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/rest"
	"github.com/nyaruka/goflow/flows"
	_ "github.com/nyaruka/goflow/flows/definition"
)

const (
	assetTypeLocationHierarchy rest.AssetType = "location_hierarchy"
)

func init() {
	rest.RegisterType(assetTypeLocationHierarchy, true, func(data json.RawMessage) (interface{}, error) { return flows.ReadLocationHierarchySet(data) })
}

// our implementation of SessionAssets - the high-level API for asset access from the engine
type sessionAssets struct {
	source assets.AssetSource
	legacy rest.LegacyServer

	channels  *flows.ChannelAssets
	fields    *flows.FieldAssets
	flows     *flows.FlowAssets
	groups    *flows.GroupAssets
	labels    *flows.LabelAssets
	resthooks *flows.ResthookAssets
}

var _ flows.SessionAssets = (*sessionAssets)(nil)

// NewSessionAssets creates a new session assets instance with the provided base URLs
func NewSessionAssets(source assets.AssetSource) (flows.SessionAssets, error) {
	channels, err := source.Channels()
	if err != nil {
		return nil, err
	}
	fields, err := source.Fields()
	if err != nil {
		return nil, err
	}
	groups, err := source.Groups()
	if err != nil {
		return nil, err
	}
	labels, err := source.Labels()
	if err != nil {
		return nil, err
	}
	resthooks, err := source.Resthooks()
	if err != nil {
		return nil, err
	}

	return &sessionAssets{
		source:    source,
		legacy:    source.(rest.LegacyServer),
		channels:  flows.NewChannelAssets(channels),
		fields:    flows.NewFieldAssets(fields),
		flows:     flows.NewFlowAssets(source),
		groups:    flows.NewGroupAssets(groups),
		labels:    flows.NewLabelAssets(labels),
		resthooks: flows.NewResthookAssets(resthooks),
	}, nil
}

func (s *sessionAssets) Channels() *flows.ChannelAssets   { return s.channels }
func (s *sessionAssets) Fields() *flows.FieldAssets       { return s.fields }
func (s *sessionAssets) Flows() *flows.FlowAssets         { return s.flows }
func (s *sessionAssets) Groups() *flows.GroupAssets       { return s.groups }
func (s *sessionAssets) Labels() *flows.LabelAssets       { return s.labels }
func (s *sessionAssets) Resthooks() *flows.ResthookAssets { return s.resthooks }

// HasLocations returns whether locations are supported as an asset item type
func (s *sessionAssets) HasLocations() bool {
	return s.source.HasLocations()
}

// GetLocationHierarchy gets the location hierarchy asset for the session
func (s *sessionAssets) GetLocationHierarchySet() (*flows.LocationHierarchySet, error) {
	asset, err := s.legacy.GetAsset(assetTypeLocationHierarchy, "")
	if err != nil {
		return nil, err
	}
	set, isType := asset.(*flows.LocationHierarchySet)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return set, nil
}
