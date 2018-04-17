package assets

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/utils"

	"github.com/karlseguin/ccache"
)

type assetType string

const (
	assetTypeChannel           assetType = "channel"
	assetTypeField             assetType = "field"
	assetTypeFlow              assetType = "flow"
	assetTypeGroup             assetType = "group"
	assetTypeLabel             assetType = "label"
	assetTypeLocationHierarchy assetType = "location_hierarchy"
)

// AssetCache fetches and caches assets for the engine
type AssetCache struct {
	cache          *ccache.Cache
	fetchUserAgent string
	fetchMutex     sync.Mutex
}

// NewAssetCache creates a new asset cache
func NewAssetCache(maxSize int64, pruneItems int, fetchUserAgent string) *AssetCache {
	return &AssetCache{
		cache:          ccache.New(ccache.Configure().MaxSize(maxSize).ItemsToPrune(uint32(pruneItems))),
		fetchUserAgent: fetchUserAgent,
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

// getItemAsset gets an item asset from the cache if it's there or from the asset server
func (c *AssetCache) getItemAsset(server AssetServer, itemType assetType, uuid string) (interface{}, error) {
	url, err := server.getItemAssetURL(itemType, uuid)
	if err != nil {
		return nil, err
	}

	return c.getAsset(url, server, itemType, false)
}

// getSetAsset gets an set asset from the cache if it's there or from the asset server
func (c *AssetCache) getSetAsset(server AssetServer, itemType assetType) (interface{}, error) {
	url, err := server.getSetAssetURL(itemType)
	if err != nil {
		return nil, err
	}

	return c.getAsset(url, server, itemType, true)
}

// gets an asset from the cache if it's there or from the asset server
func (c *AssetCache) getAsset(url string, server AssetServer, itemType assetType, isSet bool) (interface{}, error) {
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
	fetched, err := server.fetchAsset(url, itemType, isSet, c.fetchUserAgent)
	if err != nil {
		return nil, err
	}

	c.addAsset(url, fetched)
	return fetched, nil
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type assetEnvelope struct {
	URL      string          `json:"url" validate:"required,url"`
	ItemType assetType       `json:"type" validate:"required"`
	Content  json.RawMessage `json:"content"`
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
		asset, err := readAsset(envelope.Content, envelope.ItemType, envelope.IsSet)
		if err != nil {
			return fmt.Errorf("unable to read asset[url=%s]: %s", envelope.URL, err)
		}
		c.addAsset(envelope.URL, asset)
	}

	return nil
}

// reads an asset from the given raw JSON data
func readAsset(data json.RawMessage, itemType assetType, isSet bool) (interface{}, error) {
	var assetReader func(data json.RawMessage) (interface{}, error)

	if itemType == assetTypeLocationHierarchy && !isSet {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadLocationHierarchy(data) }
	} else if itemType == assetTypeChannel && !isSet {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadChannel(data) }
	} else if itemType == assetTypeChannel && isSet {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadChannelSet(data) }
	} else if itemType == assetTypeField && !isSet {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadField(data) }
	} else if itemType == assetTypeField && isSet {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadFieldSet(data) }
	} else if itemType == assetTypeFlow && !isSet {
		assetReader = func(data json.RawMessage) (interface{}, error) { return definition.ReadFlow(data) }
	} else if itemType == assetTypeGroup && !isSet {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadGroup(data) }
	} else if itemType == assetTypeGroup && isSet {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadGroupSet(data) }
	} else if itemType == assetTypeLabel && !isSet {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadLabel(data) }
	} else if itemType == assetTypeLabel && isSet {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadLabelSet(data) }
	} else {
		return nil, fmt.Errorf("unsupported asset type: %s (set=%s)", itemType, strconv.FormatBool(isSet))
	}

	return assetReader(data)
}
