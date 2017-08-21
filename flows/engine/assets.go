package engine

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/utils"
)

type assetContainer struct {
	assetType  flows.AssetType
	asset      flows.Asset
	accessedOn *time.Time
	expiresOn  time.Time
	fetchURL   string
}

type assetManager struct {
	assets map[flows.AssetUUID]assetContainer
}

func NewAssetManager() flows.AssetManager {
	return &assetManager{assets: make(map[flows.AssetUUID]assetContainer)}
}

func (m *assetManager) requestAsset(uuid flows.AssetUUID, assetType string) (flows.Asset, error) {
	container, found := m.assets[uuid]
	if !found {
		return nil, fmt.Errorf("no such asset of type '%s' with UUID '%s'", assetType, uuid)
	}

	if container.asset == nil || time.Now().After(container.expiresOn) {
		// TODO try to (re)fetch from URL, and validate
	}

	accessedOn := time.Now().UTC()
	container.accessedOn = &accessedOn

	return container.asset, nil
}

func (m *assetManager) AddAsset(asset flows.Asset, fetchURL string) {
	m.assets[asset.AssetUUID()] = assetContainer{
		assetType: asset.AssetType(),
		asset:     asset,
		expiresOn: time.Now().Add(5 * time.Minute),
		fetchURL:  fetchURL,
	}
}

func (m *assetManager) AddLazyAsset(assetType flows.AssetType, assetUUID flows.AssetUUID, fetchURL string) {
	m.assets[assetUUID] = assetContainer{assetType: assetType, fetchURL: fetchURL}
}

func (m *assetManager) ClearAssets(asset flows.Asset, expiresOn *time.Time, fetchURL string) {
	m.assets = make(map[flows.AssetUUID]assetContainer)
}

func (m *assetManager) GetFlow(uuid flows.FlowUUID) (flows.Flow, error) {
	asset, err := m.requestAsset(flows.AssetUUID(uuid), "flow")
	if err != nil {
		return nil, err
	}
	flow, isType := asset.(flows.Flow)
	if !isType {
		return nil, fmt.Errorf("unable to find flow with UUID '%s'", uuid)
	}
	return flow, nil
}

func (m *assetManager) GetChannel(uuid flows.ChannelUUID) (flows.Channel, error) {
	asset, err := m.requestAsset(flows.AssetUUID(uuid), "channel")
	if err != nil {
		return nil, err
	}
	channel, isType := asset.(flows.Channel)
	if !isType {
		return nil, fmt.Errorf("unable to find channel with UUID '%s'", uuid)
	}
	return channel, nil
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type assetEnvelope struct {
	Type    flows.AssetType  `json:"type"    validate:"required"`
	UUID    *flows.AssetUUID `json:"uuid"    validate:"omitempty,uuid"`
	Content *json.RawMessage `json:"content" validate:"omitempty"`
	URL     string           `json:"url"     validate:"omitempty,url"`
}

func (m *assetManager) IncludeAssets(data json.RawMessage) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	envelopes := make([]assetEnvelope, len(raw))
	for e := range raw {
		if err := utils.UnmarshalAndValidate(raw[e], &envelopes[e], "asset"); err != nil {
			return err
		}
	}

	nonLazy := make([]flows.Asset, 0)

	for _, envelope := range envelopes {
		var asset flows.Asset

		if envelope.Content != nil {
			var err error
			switch envelope.Type {
			case flows.AssetTypeFlow:
				asset, err = definition.ReadFlow(*envelope.Content)
			case flows.AssetTypeChannel:
				asset, err = flows.ReadChannel(*envelope.Content)
			default:
				err = fmt.Errorf("Invalid asset type: %s", envelope.Type)
			}
			if err != nil {
				return err
			}

			nonLazy = append(nonLazy, asset)
			m.AddAsset(asset, envelope.URL)
		} else {
			m.AddLazyAsset(envelope.Type, *envelope.UUID, envelope.URL)
		}
	}

	// any non-lazy assets can be validated now
	for _, asset := range nonLazy {
		if err := asset.Validate(m); err != nil {
			return err
		}
	}

	return nil
}
