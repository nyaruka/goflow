package engine

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/simple"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
)

const (
	assetTypeChannel           assets.AssetType = "channel"
	assetTypeField             assets.AssetType = "field"
	assetTypeFlow              assets.AssetType = "flow"
	assetTypeGroup             assets.AssetType = "group"
	assetTypeLabel             assets.AssetType = "label"
	assetTypeLocationHierarchy assets.AssetType = "location_hierarchy"
	assetTypeResthook          assets.AssetType = "resthook"
)

func init() {
	assets.RegisterType(assetTypeChannel, true, func(data json.RawMessage) (interface{}, error) { return flows.ReadChannelSet(data) })
	assets.RegisterType(assetTypeField, true, func(data json.RawMessage) (interface{}, error) { return flows.ReadFieldSet(data) })
	assets.RegisterType(assetTypeFlow, false, func(data json.RawMessage) (interface{}, error) { return definition.ReadFlow(data) })
	assets.RegisterType(assetTypeGroup, true, func(data json.RawMessage) (interface{}, error) { return simple.ReadGroups(data) })
	assets.RegisterType(assetTypeLabel, true, func(data json.RawMessage) (interface{}, error) { return simple.ReadLabels(data) })
	assets.RegisterType(assetTypeLocationHierarchy, true, func(data json.RawMessage) (interface{}, error) { return flows.ReadLocationHierarchySet(data) })
	assets.RegisterType(assetTypeResthook, true, func(data json.RawMessage) (interface{}, error) { return flows.ReadResthookSet(data) })
}

type ServerSource struct {
	server assets.AssetServer
}

func NewServerSource(server assets.AssetServer) assets.AssetSource {
	return &ServerSource{server: server}
}

func (s *ServerSource) Server() assets.AssetServer {
	return s.server
}

func (s *ServerSource) Labels() ([]assets.Label, error) {
	asset, err := s.server.GetAsset(assetTypeLabel, "")
	if err != nil {
		return nil, err
	}
	set, isType := asset.([]assets.Label)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return set, nil
}

func (s *ServerSource) Groups() ([]assets.Group, error) {
	asset, err := s.server.GetAsset(assetTypeGroup, "")
	if err != nil {
		return nil, err
	}
	set, isType := asset.([]assets.Group)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return set, nil
}

// our implementation of SessionAssets - the high-level API for asset access from the engine
type sessionAssets struct {
	source assets.AssetSource
	server assets.AssetServer

	groups       []*flows.Group
	groupsByUUID map[assets.GroupUUID]*flows.Group
	labels       []*flows.Label
	labelsByUUID map[assets.LabelUUID]*flows.Label
}

var _ flows.SessionAssets = (*sessionAssets)(nil)

// NewSessionAssets creates a new session assets instance with the provided base URLs
func NewSessionAssets(source assets.AssetSource) (flows.SessionAssets, error) {
	rawLabels, err := source.Labels()
	if err != nil {
		return nil, err
	}
	labels := make([]*flows.Label, len(rawLabels))
	labelsByUUID := make(map[assets.LabelUUID]*flows.Label, len(rawLabels))
	for l, rawLabel := range rawLabels {
		label := flows.NewLabel(rawLabel)
		labels[l] = label
		labelsByUUID[label.UUID()] = label
	}

	rawGroups, err := source.Groups()
	if err != nil {
		return nil, err
	}
	groups := make([]*flows.Group, len(rawGroups))
	groupsByUUID := make(map[assets.GroupUUID]*flows.Group, len(rawGroups))
	for g, rawGroup := range rawGroups {
		group := flows.NewGroup(rawGroup)
		groups[g] = group
		groupsByUUID[group.UUID()] = group
	}

	return &sessionAssets{
		server:       source.(*ServerSource).Server(),
		groups:       groups,
		groupsByUUID: groupsByUUID,
		labels:       labels,
		labelsByUUID: labelsByUUID,
	}, nil
}

// HasLocations returns whether locations are supported as an asset item type
func (s *sessionAssets) HasLocations() bool {
	return s.server.IsTypeSupported(assetTypeLocationHierarchy)
}

// GetLocationHierarchy gets the location hierarchy asset for the session
func (s *sessionAssets) GetLocationHierarchySet() (*flows.LocationHierarchySet, error) {
	asset, err := s.server.GetAsset(assetTypeLocationHierarchy, "")
	if err != nil {
		return nil, err
	}
	set, isType := asset.(*flows.LocationHierarchySet)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return set, nil
}

// GetChannel gets a channel asset for the session
func (s *sessionAssets) GetChannel(uuid flows.ChannelUUID) (flows.Channel, error) {
	set, err := s.GetChannelSet()
	if err != nil {
		return nil, err
	}
	channel := set.FindByUUID(uuid)
	if channel == nil {
		return nil, fmt.Errorf("no such channel with uuid '%s'", uuid)
	}
	return channel, nil
}

// GetChannelSet gets the set of all channels asset for the session
func (s *sessionAssets) GetChannelSet() (*flows.ChannelSet, error) {
	asset, err := s.server.GetAsset(assetTypeChannel, "")
	if err != nil {
		return nil, err
	}
	set, isType := asset.(*flows.ChannelSet)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return set, nil
}

// GetField gets a contact field asset for the session
func (s *sessionAssets) GetField(key string) (*flows.Field, error) {
	set, err := s.GetFieldSet()
	if err != nil {
		return nil, err
	}
	field := set.FindByKey(key)
	if field == nil {
		return nil, fmt.Errorf("no such field with key '%s'", key)
	}
	return field, nil
}

// GetFieldSet gets the set of all fields asset for the session
func (s *sessionAssets) GetFieldSet() (*flows.FieldSet, error) {
	asset, err := s.server.GetAsset(assetTypeField, "")
	if err != nil {
		return nil, err
	}
	set, isType := asset.(*flows.FieldSet)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return set, nil
}

// GetFlow gets a flow asset for the session
func (s *sessionAssets) GetFlow(uuid flows.FlowUUID) (flows.Flow, error) {
	asset, err := s.server.GetAsset(assetTypeFlow, string(uuid))
	if err != nil {
		return nil, err
	}
	flow, isType := asset.(flows.Flow)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type for UUID '%s'", uuid)
	}
	return flow, nil
}

// GetGroup gets the group with the given UUID
func (s *sessionAssets) GetGroup(uuid assets.GroupUUID) (*flows.Group, error) {
	group, found := s.groupsByUUID[uuid]
	if !found {
		return nil, fmt.Errorf("no such group with uuid '%s'", uuid)
	}
	return group, nil
}

// FindGroupByName gets the group with the given name if its exists
func (s *sessionAssets) FindGroupByName(name string) *flows.Group {
	name = strings.ToLower(name)
	for _, group := range s.groups {
		if strings.ToLower(group.Name()) == name {
			return group
		}
	}
	return nil
}

// GetAllGroups gets all groups
func (s *sessionAssets) GetAllGroups() []*flows.Group {
	return s.groups
}

// GetLabel gets the label with the given UUID
func (s *sessionAssets) GetLabel(uuid assets.LabelUUID) (*flows.Label, error) {
	label, found := s.labelsByUUID[uuid]
	if !found {
		return nil, fmt.Errorf("no such label with uuid '%s'", uuid)
	}
	return label, nil
}

// FindLabelByName gets the label with the given name if its exists
func (s *sessionAssets) FindLabelByName(name string) *flows.Label {
	name = strings.ToLower(name)
	for _, label := range s.labels {
		if strings.ToLower(label.Name()) == name {
			return label
		}
	}
	return nil
}

func (s *sessionAssets) GetResthookSet() (*flows.ResthookSet, error) {
	asset, err := s.server.GetAsset(assetTypeResthook, "")
	if err != nil {
		return nil, err
	}
	set, isType := asset.(*flows.ResthookSet)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return set, nil
}

// NewMockAssetServer creates a new mocked asset server with URLs for all flow engine types already configured
func NewMockAssetServer(cache *assets.AssetCache) *assets.MockAssetServer {
	return assets.NewMockAssetServer(map[assets.AssetType]string{
		assetTypeChannel:           "http://testserver/assets/channel/",
		assetTypeField:             "http://testserver/assets/field/",
		assetTypeFlow:              "http://testserver/assets/flow/",
		assetTypeGroup:             "http://testserver/assets/group/",
		assetTypeLabel:             "http://testserver/assets/label/",
		assetTypeLocationHierarchy: "http://testserver/assets/location_hierarchy/",
		assetTypeResthook:          "http://testserver/assets/resthook/",
	}, cache)
}
