package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/karlseguin/ccache"

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
	assetItemTypeField   AssetItemType = "field"
	assetItemTypeFlow    AssetItemType = "flow"
	assetItemTypeGroup   AssetItemType = "group"
	assetItemTypeLabel   AssetItemType = "label"
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

// adds an asset to the cache identified by the given URL
func (c *AssetCache) addAsset(url assetURL, asset interface{}) {
	c.cache.Set(string(url), asset, time.Hour*24)
}

// gets an asset from the cache if it's there or from the asset server
func (c *AssetCache) getAsset(url assetURL, aType assetType, itemType AssetItemType) (interface{}, error) {
	item := c.cache.Get(string(url))

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
	fetched, err := c.fetchAsset(url, aType, itemType)
	if err != nil {
		return nil, err
	}

	c.addAsset(url, fetched)
	return fetched, nil
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

// a higher level wrapper for the cache
type sessionAssets struct {
	cache    *AssetCache
	typeURLs map[AssetItemType]string
}

// NewSessionAssets creates a new session assets instance with the provided base URLs
func NewSessionAssets(cache *AssetCache, typeURLs map[AssetItemType]string) flows.SessionAssets {
	return &sessionAssets{cache: cache, typeURLs: typeURLs}
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

func (s *sessionAssets) GetField(key flows.FieldKey) (*flows.Field, error) {
	fields, err := s.GetFieldSet()
	if err != nil {
		return nil, err
	}
	field := fields.FindByKey(key)
	if field == nil {
		return nil, fmt.Errorf("no such field with key '%s'", key)
	}
	return field, nil
}

func (s *sessionAssets) GetFieldSet() (*flows.FieldSet, error) {
	url := s.getAssetSetURL(assetItemTypeField)
	asset, err := s.cache.getAsset(url, assetTypeSet, assetItemTypeField)
	if err != nil {
		return nil, err
	}
	fields, isType := asset.(*flows.FieldSet)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return fields, nil
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
	if group == nil {
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

func (s *sessionAssets) GetLabel(uuid flows.LabelUUID) (*flows.Label, error) {
	labels, err := s.GetLabelSet()
	if err != nil {
		return nil, err
	}
	label := labels.FindByUUID(uuid)
	if label == nil {
		return nil, fmt.Errorf("no such label with uuid '%s'", uuid)
	}
	return label, nil
}

func (s *sessionAssets) GetLabelSet() (*flows.LabelSet, error) {
	url := s.getAssetSetURL(assetItemTypeLabel)
	asset, err := s.cache.getAsset(url, assetTypeSet, assetItemTypeLabel)
	if err != nil {
		return nil, err
	}
	labels, isType := asset.(*flows.LabelSet)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return labels, nil
}

func (s *sessionAssets) getAssetSetURL(itemType AssetItemType) assetURL {
	return assetURL(s.typeURLs[itemType])
}

func (s *sessionAssets) getAssetItemURL(itemType AssetItemType, uuid itemUUID) assetURL {
	return assetURL(fmt.Sprintf("%s/%s", s.typeURLs[itemType], uuid))
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
	} else if aType == assetTypeObject && itemType == assetItemTypeField {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadField(data) }
	} else if aType == assetTypeSet && itemType == assetItemTypeField {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadFieldSet(data) }
	} else if aType == assetTypeObject && itemType == assetItemTypeFlow {
		assetReader = func(data json.RawMessage) (interface{}, error) { return definition.ReadFlow(data) }
	} else if aType == assetTypeObject && itemType == assetItemTypeGroup {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadGroup(data) }
	} else if aType == assetTypeSet && itemType == assetItemTypeGroup {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadGroupSet(data) }
	} else if aType == assetTypeObject && itemType == assetItemTypeLabel {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadLabel(data) }
	} else if aType == assetTypeSet && itemType == assetItemTypeLabel {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadLabelSet(data) }
	} else {
		return nil, fmt.Errorf("unsupported asset type: %s of %s", aType, itemType)
	}

	return assetReader(data)
}
