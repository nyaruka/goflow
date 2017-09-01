package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/utils"
)

type assetURL string
type itemUUID string

type assetType string

const (
	assetTypeObject assetType = "object"
	assetTypeSet    assetType = "set"
)

type AssetItemType string

const (
	assetItemTypeChannel AssetItemType = "channel"
	assetItemTypeFlow    AssetItemType = "flow"
	assetItemTypeGroup   AssetItemType = "group"
)

// container for any asset in the cache
type cachedAsset struct {
	asset      interface{}
	addedOn    time.Time
	accessedOn time.Time
}

// AssetCache fetches and caches assets for the engine
type AssetCache struct {
	cache map[assetURL]cachedAsset
	mutex sync.Mutex
}

// NewAssetCache creates a new asset cache
func NewAssetCache() *AssetCache {
	return &AssetCache{cache: make(map[assetURL]cachedAsset)}
}

func (c *AssetCache) putAsset(url assetURL, asset interface{}) {
	c.cache[url] = cachedAsset{asset: asset, addedOn: time.Now().UTC()}
}

func (c *AssetCache) addAsset(url assetURL, asset interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.putAsset(url, asset)
}

// gets an asset from the cache if it's there or from the asset server
func (c *AssetCache) getAsset(url assetURL, aType assetType, itemType AssetItemType) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	cached, found := c.cache[url]
	if !found {
		asset, err := c.fetchAsset(url, aType, itemType)
		if err != nil {
			return nil, err
		}

		c.putAsset(url, asset)
		return asset, nil
	}

	// update the accessed on time
	cached.accessedOn = time.Now().UTC()

	return cached.asset, nil
}

// fetches an asset by its URL and parses it as the provided type
func (c *AssetCache) fetchAsset(url assetURL, aType assetType, itemType AssetItemType) (interface{}, error) {
	response, err := http.Get(string(url))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("asset request returned non-200 response")
	}
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return readAsset(buf, aType, itemType)
}

type sessionAssets struct {
	cache     *AssetCache
	assetURLs map[AssetItemType]string
}

func NewSessionAssets(cache *AssetCache, assetURLs map[AssetItemType]string) flows.SessionAssets {
	return &sessionAssets{cache: cache, assetURLs: assetURLs}
}

func (s *sessionAssets) GetChannel(uuid flows.ChannelUUID) (flows.Channel, error) {
	url := s.getAssetItemURL(assetItemTypeChannel, itemUUID(uuid))
	asset, err := s.cache.getAsset(url, assetTypeObject, assetItemTypeChannel)
	if err != nil {
		return nil, err
	}
	channel, isType := asset.(flows.Channel)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type for UUID '%s'", uuid)
	}
	return channel, nil
}

func (s *sessionAssets) GetFlow(uuid flows.FlowUUID) (flows.Flow, error) {
	url := s.getAssetItemURL(assetItemTypeFlow, itemUUID(uuid))
	asset, err := s.cache.getAsset(url, assetTypeObject, assetItemTypeFlow)
	if err != nil {
		return nil, err
	}
	flow, isType := asset.(flows.Flow)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type for UUID '%s'", uuid)
	}
	return flow, nil
}

func (s *sessionAssets) GetGroup(uuid flows.GroupUUID) (*flows.Group, error) {
	groups, err := s.GetGroupSet()
	if err != nil {
		return nil, err
	}
	group := groups.FindByUUID(uuid)
	if group != nil {
		return nil, fmt.Errorf("no such group with uuid '%s'", uuid)
	}
	return group, nil
}

func (s *sessionAssets) GetGroupSet() (*flows.GroupSet, error) {
	url := s.getAssetSetURL(assetItemTypeGroup)
	asset, err := s.cache.getAsset(url, assetTypeSet, assetItemTypeGroup)
	if err != nil {
		return nil, err
	}
	groups, isType := asset.(*flows.GroupSet)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return groups, nil
}

func (s *sessionAssets) getAssetSetURL(itemType AssetItemType) assetURL {
	return assetURL(s.assetURLs[itemType])
}

func (s *sessionAssets) getAssetItemURL(itemType AssetItemType, uuid itemUUID) assetURL {
	return assetURL(fmt.Sprintf("%s/%s", s.assetURLs[itemType], uuid))
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type assetEnvelope struct {
	URL      assetURL        `json:"url" validate:"required,url"`
	ItemType AssetItemType   `json:"type" validate:"required"`
	Content  json.RawMessage `json:"content" validate:"required"`
	IsSet    bool            `json:"is_set"`
}

// Include loads assets from the given raw JSON into the cache
func (c *AssetCache) Include(data json.RawMessage) error {
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

	for _, envelope := range envelopes {
		var aType assetType
		if envelope.IsSet {
			aType = assetTypeSet
		} else {
			aType = assetTypeObject
		}

		asset, err := readAsset(envelope.Content, aType, envelope.ItemType)
		if err != nil {
			return err
		}
		c.addAsset(envelope.URL, asset)
	}

	return nil
}

// reads an asset from the given raw JSON data
func readAsset(data json.RawMessage, aType assetType, itemType AssetItemType) (interface{}, error) {
	var assetReader func(data json.RawMessage) (interface{}, error)

	if aType == assetTypeObject && itemType == assetItemTypeChannel {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadChannel(data) }
	} else if aType == assetTypeObject && itemType == assetItemTypeFlow {
		assetReader = func(data json.RawMessage) (interface{}, error) { return definition.ReadFlow(data) }
	} else if aType == assetTypeObject && itemType == assetItemTypeGroup {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadGroup(data) }
	} else if aType == assetTypeSet && itemType == assetItemTypeGroup {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadGroupSet(data) }
	} else {
		return nil, fmt.Errorf("unsupported asset type: %s of %s", aType, itemType)
	}

	return assetReader(data)
}
