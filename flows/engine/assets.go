package engine

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/server"
	"github.com/nyaruka/goflow/assets/server/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
)

const (
	assetTypeChannel           server.AssetType = "channel"
	assetTypeField             server.AssetType = "field"
	assetTypeFlow              server.AssetType = "flow"
	assetTypeGroup             server.AssetType = "group"
	assetTypeLabel             server.AssetType = "label"
	assetTypeLocationHierarchy server.AssetType = "location_hierarchy"
	assetTypeResthook          server.AssetType = "resthook"
)

func init() {
	server.RegisterType(assetTypeChannel, true, func(data json.RawMessage) (interface{}, error) { return types.ReadChannels(data) })
	server.RegisterType(assetTypeField, true, func(data json.RawMessage) (interface{}, error) { return types.ReadFields(data) })
	server.RegisterType(assetTypeGroup, true, func(data json.RawMessage) (interface{}, error) { return types.ReadGroups(data) })
	server.RegisterType(assetTypeLabel, true, func(data json.RawMessage) (interface{}, error) { return types.ReadLabels(data) })
	server.RegisterType(assetTypeResthook, true, func(data json.RawMessage) (interface{}, error) { return types.ReadResthooks(data) })

	server.RegisterType(assetTypeFlow, false, func(data json.RawMessage) (interface{}, error) { return definition.ReadFlow(data) })
	server.RegisterType(assetTypeLocationHierarchy, true, func(data json.RawMessage) (interface{}, error) { return flows.ReadLocationHierarchySet(data) })
}

// our implementation of SessionAssets - the high-level API for asset access from the engine
type sessionAssets struct {
	source assets.AssetSource
	legacy server.LegacyServer

	channels  *flows.ChannelAssets
	fields    *flows.FieldAssets
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
		legacy:    source.(server.LegacyServer),
		channels:  flows.NewChannelAssets(channels),
		fields:    flows.NewFieldAssets(fields),
		groups:    flows.NewGroupAssets(groups),
		labels:    flows.NewLabelAssets(labels),
		resthooks: flows.NewResthookAssets(resthooks),
	}, nil
}

func (s *sessionAssets) Channels() *flows.ChannelAssets   { return s.channels }
func (s *sessionAssets) Fields() *flows.FieldAssets       { return s.fields }
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

// GetFlow gets a flow asset for the session
func (s *sessionAssets) GetFlow(uuid flows.FlowUUID) (flows.Flow, error) {
	asset, err := s.legacy.GetAsset(assetTypeFlow, string(uuid))
	if err != nil {
		return nil, err
	}
	flow, isType := asset.(flows.Flow)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type for UUID '%s'", uuid)
	}
	return flow, nil
}

// NewMockServerSource creates a new mocked asset server with URLs for all flow engine types already configured
func NewMockServerSource(cache *server.AssetCache) *server.MockServerSource {
	return server.NewMockServerSource(map[server.AssetType]string{
		assetTypeChannel:           "http://testserver/assets/channel/",
		assetTypeField:             "http://testserver/assets/field/",
		assetTypeFlow:              "http://testserver/assets/flow/",
		assetTypeGroup:             "http://testserver/assets/group/",
		assetTypeLabel:             "http://testserver/assets/label/",
		assetTypeLocationHierarchy: "http://testserver/assets/location_hierarchy/",
		assetTypeResthook:          "http://testserver/assets/resthook/",
	}, cache)
}
