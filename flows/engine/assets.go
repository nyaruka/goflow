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

// AssetAllItems is used in place of an asset identifier to signify that an asset is the set of all items of a type
var AssetAllItems itemUUID = ""

// AssetItemType is the asset item type, e.g. flow
type AssetItemType string

const (
	AssetItemTypeChannel AssetItemType = "channel"
	AssetItemTypeFlow    AssetItemType = "flow"
	AssetItemTypeGroup   AssetItemType = "group"
)

type cachedAsset struct {
	asset      interface{}
	accessedOn time.Time
}

// AssetCache fetches and caches assets for the engine
type AssetCache struct {
	cache map[assetURL]cachedAsset
	mutex sync.Mutex
}

func NewAssetCache() *AssetCache {
	return &AssetCache{cache: make(map[assetURL]cachedAsset)}
}

func (c *AssetCache) putAsset(url assetURL, asset interface{}) {
	c.cache[url] = cachedAsset{asset: asset, accessedOn: time.Now().UTC()}
}

func (c *AssetCache) addAsset(url assetURL, asset interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.putAsset(url, asset)
}

func (c *AssetCache) getAsset(url assetURL, itemType AssetItemType, allOfType bool) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	cached, found := c.cache[url]
	if !found {
		asset, err := c.fetchAsset(url, itemType, allOfType)
		if err != nil {
			return nil, err
		}

		c.putAsset(url, asset)
		return asset, nil
	}

	cached.accessedOn = time.Now().UTC()

	return cached.asset, nil
}

func (c *AssetCache) fetchAsset(url assetURL, itemType AssetItemType, allOfType bool) (interface{}, error) {
	response, err := http.Get(string(url))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("asset request returned non-200 response")
	}

	buf, _ := ioutil.ReadAll(response.Body)
	return readAsset(buf, itemType, allOfType)
}

type sessionAssets struct {
	cache         *AssetCache
	serverBaseURL string
}

func NewSessionAssets(cache *AssetCache, serverBaseURL string) flows.SessionAssets {
	return &sessionAssets{cache: cache, serverBaseURL: serverBaseURL}
}

func (s *sessionAssets) ServerBaseURL() string {
	return s.serverBaseURL
}

func (s *sessionAssets) GetChannel(uuid flows.ChannelUUID) (flows.Channel, error) {
	url := s.getURLForAsset(AssetItemTypeChannel, itemUUID(uuid))
	asset, err := s.cache.getAsset(url, AssetItemTypeChannel, false)
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
	url := s.getURLForAsset(AssetItemTypeFlow, itemUUID(uuid))
	asset, err := s.cache.getAsset(url, AssetItemTypeFlow, false)
	if err != nil {
		return nil, err
	}
	flow, isType := asset.(flows.Flow)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type for UUID '%s'", uuid)
	}
	return flow, nil
}

func (s *sessionAssets) GetGroups() ([]flows.Group, error) {
	url := s.getURLForAsset(AssetItemTypeGroup, AssetAllItems)
	asset, err := s.cache.getAsset(url, AssetItemTypeGroup, true)
	if err != nil {
		return nil, err
	}
	groups, isType := asset.([]flows.Group)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return groups, nil
}

func (s *sessionAssets) getURLForAsset(itemType AssetItemType, uuid itemUUID) assetURL {
	if uuid == AssetAllItems {
		return assetURL(fmt.Sprintf("%s/%s", s.serverBaseURL, itemType))
	}
	return assetURL(fmt.Sprintf("%s/%s/%s", s.serverBaseURL, itemType, uuid))
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type assetEnvelope struct {
	URL     assetURL        `json:"url"     validate:"required,url"`
	Type    AssetItemType   `json:"type"    validate:"required"`
	Content json.RawMessage `json:"content" validate:"required"`
	IsSet   bool            `json:"is_set"`
}

// Include loads assets from the given raw JSON into the cache
func (s *AssetCache) Include(data json.RawMessage) error {
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
		asset, err := readAsset(envelope.Content, envelope.Type, envelope.IsSet)
		if err != nil {
			return err
		}
		s.addAsset(envelope.URL, asset)
	}

	return nil
}

// reads an asset from the given raw JSON data
func readAsset(data json.RawMessage, itemType AssetItemType, isSet bool) (interface{}, error) {
	var itemReader func(data json.RawMessage) (interface{}, error)

	switch itemType {
	case AssetItemTypeFlow:
		itemReader = func(data json.RawMessage) (interface{}, error) { return definition.ReadFlow(data) }
	case AssetItemTypeChannel:
		itemReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadChannel(data) }
	default:
		return nil, fmt.Errorf("unknown asset type: %s", itemType)
	}

	if isSet {
		var envelopes []json.RawMessage
		if err := json.Unmarshal(data, &envelopes); err != nil {
			return nil, err
		}

		assets := make([]interface{}, len(envelopes))
		var err error
		for e := range envelopes {
			if assets[e], err = itemReader(data); err != nil {
				return nil, err
			}
		}

		return assets, nil
	}

	// asset is a single item
	return itemReader(data)
}
