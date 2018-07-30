package assets

import (
	"encoding/json"
	"fmt"
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
	assetTypeChannelSet           assetType = "channel_set"
	assetTypeFieldSet             assetType = "field_set"
	assetTypeFlow                 assetType = "flow"
	assetTypeGroupSet             assetType = "group_set"
	assetTypeLabelSet             assetType = "label_set"
	assetTypeLocationHierarchySet assetType = "location_hierarchy_set"
	assetTypeResthookSet          assetType = "resthook_set"
)

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

// GetAsset gets an asset from the cache if it's there or from the asset server
func (c *AssetCache) GetAsset(server AssetServer, itemType assetType, itemUUID string) (interface{}, error) {
	url, err := server.getAssetURL(itemType, itemUUID)
	if err != nil {
		return nil, err
	}

	return c.getAsset(url, server, itemType)
}

// gets an asset from the cache if it's there or from the asset server
func (c *AssetCache) getAsset(url string, server AssetServer, itemType assetType) (interface{}, error) {
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
	fetched, err := server.fetchAsset(url, itemType)
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
		asset, err := readAsset(envelope.Content, envelope.ItemType)
		if err != nil {
			return fmt.Errorf("unable to read asset[url=%s]: %s", envelope.URL, err)
		}
		c.addAsset(envelope.URL, asset)
	}

	return nil
}

// reads an asset from the given raw JSON data
func readAsset(data json.RawMessage, itemType assetType) (interface{}, error) {
	var assetReader func(data json.RawMessage) (interface{}, error)

	if itemType == assetTypeLocationHierarchySet {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadLocationHierarchySet(data) }
	} else if itemType == assetTypeChannelSet {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadChannelSet(data) }
	} else if itemType == assetTypeFieldSet {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadFieldSet(data) }
	} else if itemType == assetTypeFlow {
		assetReader = func(data json.RawMessage) (interface{}, error) { return definition.ReadFlow(data) }
	} else if itemType == assetTypeGroupSet {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadGroupSet(data) }
	} else if itemType == assetTypeLabelSet {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadLabelSet(data) }
	} else if itemType == assetTypeResthookSet {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadResthookSet(data) }
	} else {
		return nil, fmt.Errorf("unsupported asset type: %s", itemType)
	}

	return assetReader(data)
}
