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

type cachedAsset struct {
	asset     interface{}
	expiresOn time.Time
}

// AssetCache fetches and caches assets for the engine
type AssetCache struct {
	cache map[string]cachedAsset
	mutex sync.Mutex
}

type assetStore struct {
	cache         *AssetCache
	serverBaseURL string
}

// AssetAllOfType is used in place of an asset identifier to signify that an asset is the set of all of a type
var AssetAllOfType = ""

// AssetType is the asset type, e.g. flow
type AssetType string

const (
	AssetTypeChannel AssetType = "channel"
	AssetTypeFlow    AssetType = "flow"
	AssetTypeGroup   AssetType = "group"
)

func NewAssetCache() *AssetCache {
	return &AssetCache{cache: make(map[string]cachedAsset)}
}

func (c *AssetCache) putAsset(url string, asset interface{}) {
	c.cache[url] = cachedAsset{asset: asset, expiresOn: time.Now().Add(5 * time.Minute)}
}

func (c *AssetCache) addAsset(url string, asset interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.putAsset(url, asset)
}

func (c *AssetCache) getAsset(url string, assetType AssetType) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	cached, found := c.cache[url]
	if !found || time.Now().After(cached.expiresOn) {
		asset, err := c.fetchAsset(url, assetType)
		if err != nil {
			return nil, err
		}

		c.putAsset(url, asset)
		return asset, nil
	}

	return cached.asset, nil
}

func (c *AssetCache) fetchAsset(url string, assetType AssetType) (interface{}, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("asset request returned non-200 response")
	}

	buf, _ := ioutil.ReadAll(response.Body)
	return readAsset(buf, assetType)
}

func NewAssetStore(cache *AssetCache, serverBaseURL string) flows.AssetStore {
	return &assetStore{cache: cache, serverBaseURL: serverBaseURL}
}

func (s *assetStore) ServerBaseURL() string {
	return s.serverBaseURL
}

func (s *assetStore) GetChannel(uuid flows.ChannelUUID) (flows.Channel, error) {
	url := s.getURLForAsset(AssetTypeChannel, string(uuid))
	asset, err := s.cache.getAsset(url, AssetTypeChannel)
	if err != nil {
		return nil, err
	}
	channel, isType := asset.(flows.Channel)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type for UUID '%s'", uuid)
	}
	return channel, nil
}

func (s *assetStore) GetFlow(uuid flows.FlowUUID) (flows.Flow, error) {
	url := s.getURLForAsset(AssetTypeFlow, string(uuid))
	asset, err := s.cache.getAsset(url, AssetTypeFlow)
	if err != nil {
		return nil, err
	}
	flow, isType := asset.(flows.Flow)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type for UUID '%s'", uuid)
	}
	return flow, nil
}

func (s *assetStore) GetGroups() ([]flows.Group, error) {
	url := s.getURLForAsset(AssetTypeGroup, AssetAllOfType)
	asset, err := s.cache.getAsset(url, AssetTypeGroup)
	if err != nil {
		return nil, err
	}
	groups, isType := asset.([]flows.Group)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return groups, nil
}

func (s *assetStore) getURLForAsset(assetType AssetType, assetUUID string) string {
	if assetUUID == AssetAllOfType {
		return fmt.Sprintf("%s/%s", s.serverBaseURL, assetType)
	}
	return fmt.Sprintf("%s/%s/%s", s.serverBaseURL, assetType, assetUUID)
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type assetEnvelope struct {
	URL     string          `json:"url"     validate:"required,url"`
	Type    AssetType       `json:"type"    validate:"required"`
	Content json.RawMessage `json:"content" validate:"required"`
}

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
		asset, err := readAsset(envelope.Content, envelope.Type)
		if err != nil {
			return err
		}

		c.addAsset(envelope.URL, asset)
	}

	// any non-lazy assets can be validated now
	//for _, asset := range nonLazy {
	//	if err := asset.Validate(s); err != nil {
	//		return utils.NewValidationErrors(err.Error())
	//	}
	//}

	return nil
}

func readAsset(data json.RawMessage, assetType AssetType) (interface{}, error) {
	var asset interface{}
	var err error
	switch assetType {
	case AssetTypeFlow:
		asset, err = definition.ReadFlow(data)
	case AssetTypeChannel:
		asset, err = flows.ReadChannel(data)
	default:
		err = fmt.Errorf("unknown asset type: %s", assetType)
	}

	if err != nil {
		return nil, err
	}
	return asset, nil
}
