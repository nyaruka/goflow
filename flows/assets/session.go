package assets

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
)

// our implementation of SessionAssets - the high-level API for asset access from the engine
type sessionAssets struct {
	cache  *AssetCache
	server AssetServer
}

var _ flows.SessionAssets = (*sessionAssets)(nil)

// NewSessionAssets creates a new session assets instance with the provided base URLs
func NewSessionAssets(cache *AssetCache, server AssetServer) flows.SessionAssets {
	return &sessionAssets{cache: cache, server: server}
}

// HasLocations returns whether locations are supported as an asset item type
func (s *sessionAssets) HasLocations() bool {
	return s.server.isTypeSupported(assetTypeLocationHierarchySet)
}

// GetLocationHierarchy gets the location hierarchy asset for the session
func (s *sessionAssets) GetLocationHierarchySet() (*flows.LocationHierarchySet, error) {
	asset, err := s.cache.GetAsset(s.server, assetTypeLocationHierarchySet, "")
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
	asset, err := s.cache.GetAsset(s.server, assetTypeChannelSet, "")
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
	asset, err := s.cache.GetAsset(s.server, assetTypeFieldSet, "")
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
	asset, err := s.cache.GetAsset(s.server, assetTypeFlow, string(uuid))
	if err != nil {
		return nil, err
	}
	flow, isType := asset.(flows.Flow)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type for UUID '%s'", uuid)
	}
	return flow, nil
}

// GetGroup gets a contact group asset for the session
func (s *sessionAssets) GetGroup(uuid flows.GroupUUID) (*flows.Group, error) {
	set, err := s.GetGroupSet()
	if err != nil {
		return nil, err
	}
	group := set.FindByUUID(uuid)
	if group == nil {
		return nil, fmt.Errorf("no such group with uuid '%s'", uuid)
	}
	return group, nil
}

// GetGroupSet gets the set of all groups asset for the session
func (s *sessionAssets) GetGroupSet() (*flows.GroupSet, error) {
	asset, err := s.cache.GetAsset(s.server, assetTypeGroupSet, "")
	if err != nil {
		return nil, err
	}
	set, isType := asset.(*flows.GroupSet)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return set, nil
}

// GetLabel gets a message label asset for the session
func (s *sessionAssets) GetLabel(uuid flows.LabelUUID) (*flows.Label, error) {
	set, err := s.GetLabelSet()
	if err != nil {
		return nil, err
	}
	label := set.FindByUUID(uuid)
	if label == nil {
		return nil, fmt.Errorf("no such label with uuid '%s'", uuid)
	}
	return label, nil
}

func (s *sessionAssets) GetLabelSet() (*flows.LabelSet, error) {
	asset, err := s.cache.GetAsset(s.server, assetTypeLabelSet, "")
	if err != nil {
		return nil, err
	}
	set, isType := asset.(*flows.LabelSet)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return set, nil
}

func (s *sessionAssets) GetResthookSet() (*flows.ResthookSet, error) {
	asset, err := s.cache.GetAsset(s.server, assetTypeResthookSet, "")
	if err != nil {
		return nil, err
	}
	set, isType := asset.(*flows.ResthookSet)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return set, nil
}
