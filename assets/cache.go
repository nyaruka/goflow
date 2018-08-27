package assets

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nyaruka/goflow/utils"

	"github.com/karlseguin/ccache"
)

// AssetType is the unique slug for an asset type
type AssetType string

// AssetReader is a function capable of reading an asset type from JSON
type AssetReader func(data json.RawMessage) (interface{}, error)

type assetTypeConfig struct {
	manageAsSet bool
	reader      AssetReader
}

var typeConfigs = map[AssetType]*assetTypeConfig{}

// RegisterType registers a new asset type for use with this cache
func RegisterType(name AssetType, manageAsSet bool, reader AssetReader) {
	typeConfigs[name] = &assetTypeConfig{manageAsSet: manageAsSet, reader: reader}
}

// anything which the cache can use to fetch missing items
type assetFetcher interface {
	fetchAsset(url string, itemType AssetType) ([]byte, error)
}

// AssetCache fetches and caches assets for the engine
type AssetCache struct {
	cache      *ccache.Cache
	fetchMutex sync.Mutex
}

// NewAssetCache creates a new asset cache
func NewAssetCache(maxSize int64, pruneItems int) *AssetCache {
	return &AssetCache{
		cache: ccache.New(ccache.Configure().MaxSize(maxSize).ItemsToPrune(uint32(pruneItems))),
	}
}

// Shutdown shuts down this asset cache
func (c *AssetCache) Shutdown() {
	c.cache.Stop()
}

func (c *AssetCache) normalizeURL(url string) string {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	return url
}

// adds an asset to the cache identified by the given URL
func (c *AssetCache) addAsset(url string, asset interface{}) {
	c.cache.Set(c.normalizeURL(url), asset, time.Hour*24)
}

// gets an asset from the cache if it's there or from the asset server
func (c *AssetCache) getAsset(url string, fetcher assetFetcher, itemType AssetType) (interface{}, error) {
	item := c.cache.Get(c.normalizeURL(url))

	// asset was in cache, so just return it
	if item != nil {
		return item.Value(), nil
	}

	// multiple threads might get but we don't want to perform multiple fetches
	c.fetchMutex.Lock()
	defer c.fetchMutex.Unlock()

	// check again in case we weren't the first thread to reach the fetch mutex
	item = c.cache.Get(string(url))
	if item != nil {
		return item.Value(), nil
	}

	// actually fetch the asset from it's URL
	data, err := fetcher.fetchAsset(url, itemType)
	if err != nil {
		return nil, fmt.Errorf("error fetching asset %s: %s", url, err)
	}

	cfg := typeConfigs[itemType]
	if cfg == nil {
		return nil, fmt.Errorf("unsupported asset type: %s", itemType)
	}

	a, err := c.readAsset(data, itemType, true)
	if err != nil {
		return nil, err
	}

	c.addAsset(url, a)
	return a, nil
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type assetEnvelope struct {
	URL      string          `json:"url" validate:"required,url"`
	ItemType AssetType       `json:"type" validate:"required"`
	Content  json.RawMessage `json:"content"`
}

// Include loads assets from the given raw JSON into the cache
func (c *AssetCache) Include(data json.RawMessage) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	envelopes := make([]assetEnvelope, len(raw))
	for e := range raw {
		if err := utils.UnmarshalAndValidate(raw[e], &envelopes[e]); err != nil {
			return fmt.Errorf("unable to read asset: %s", err)
		}
	}

	for _, envelope := range envelopes {
		asset, err := c.readAsset(envelope.Content, envelope.ItemType, false)
		if err != nil {
			return fmt.Errorf("unable to read asset[url=%s]: %s", envelope.URL, err)
		}
		c.addAsset(envelope.URL, asset)
	}

	return nil
}

// reads an asset from the given raw JSON data
func (c *AssetCache) readAsset(data json.RawMessage, itemType AssetType, fromRequest bool) (interface{}, error) {
	cfg := typeConfigs[itemType]
	if cfg == nil {
		return nil, fmt.Errorf("unsupported asset type: %s", itemType)
	}

	if cfg.manageAsSet && fromRequest {
		listResponse := &struct {
			Results json.RawMessage `json:"results"`
		}{}
		if err := json.Unmarshal(data, listResponse); err != nil {
			return nil, fmt.Errorf("expected result set")
		}
		data = listResponse.Results
	}

	return cfg.reader(data)
}
