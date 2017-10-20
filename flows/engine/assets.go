package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/karlseguin/ccache"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/utils"
)

type assetURL string
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

// gets an item asset from the cache if it's there or from the asset server
func (c *AssetCache) getItemAsset(server *AssetServer, itemType assetType, uuid string) (interface{}, error) {
	url, err := server.getItemAssetURL(itemType, uuid)
	if err != nil {
		return nil, err
	}

	return c.getAsset(url, server, itemType, false)
}

// gets an set asset from the cache if it's there or from the asset server
func (c *AssetCache) getSetAsset(server *AssetServer, itemType assetType) (interface{}, error) {
	url, err := server.getSetAssetURL(itemType)
	if err != nil {
		return nil, err
	}

	return c.getAsset(url, server, itemType, true)
}

// gets an asset from the cache if it's there or from the asset server
func (c *AssetCache) getAsset(url assetURL, server *AssetServer, itemType assetType, isSet bool) (interface{}, error) {
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
	fetched, err := server.fetchAsset(url, itemType, isSet)
	if err != nil {
		return nil, err
	}

	c.addAsset(url, fetched)
	return fetched, nil
}

// AssetServer describes the asset server we'll fetch missing assets from
type AssetServer struct {
	RequestHeaders map[string]string    `json:"request_headers"`
	TypeURLs       map[assetType]string `json:"type_urls"`
}

// NewAssetServer creates a new asset server
func NewAssetServer(requestHeaders map[string]string, typeURLs map[assetType]string) *AssetServer {
	return &AssetServer{RequestHeaders: requestHeaders, TypeURLs: typeURLs}
}

// NewTestAssetServer creates a new asset server for testing
func NewTestAssetServer() *AssetServer {
	return &AssetServer{
		RequestHeaders: map[string]string{
			"User-Agent": "FlowServerTest/1.0",
		},
		TypeURLs: map[assetType]string{
			assetTypeChannel:           "http://testserver/assets/channel",
			assetTypeField:             "http://testserver/assets/field",
			assetTypeFlow:              "http://testserver/assets/flow",
			assetTypeGroup:             "http://testserver/assets/group",
			assetTypeLabel:             "http://testserver/assets/label",
			assetTypeLocationHierarchy: "http://testserver/assets/location_hierarchy",
		},
	}
}

// isTypeSupported returns whether the given asset item type is supported
func (s *AssetServer) isTypeSupported(itemType assetType) bool {
	_, hasTypeURL := s.TypeURLs[itemType]
	return hasTypeURL
}

// getSetAssetURL gets the URL for a set of the given asset type
func (s *AssetServer) getSetAssetURL(itemType assetType) (assetURL, error) {
	typeURL, typeFound := s.TypeURLs[itemType]
	if !typeFound {
		return "", fmt.Errorf("asset type '%s' not supported by asset server", itemType)
	}

	return assetURL(typeURL), nil
}

// getItemAssetURL gets the URL for an item of the given asset type
func (s *AssetServer) getItemAssetURL(itemType assetType, uuid string) (assetURL, error) {
	setURL, err := s.getSetAssetURL(itemType)
	if err != nil {
		return "", err
	}

	return assetURL(fmt.Sprintf("%s/%s", setURL, uuid)), nil
}

// fetches an asset by its URL and parses it as the provided type
func (s *AssetServer) fetchAsset(url assetURL, itemType assetType, isSet bool) (interface{}, error) {
	request, err := http.NewRequest("GET", string(url), nil)
	if err != nil {
		return nil, err
	}

	// set request headers
	for header, value := range s.RequestHeaders {
		request.Header.Set(header, value)
	}

	// make the actual request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("asset request returned non-200 response")
	}

	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return readAsset(buf, itemType, isSet)
}

// a higher level wrapper for the cache
type sessionAssets struct {
	cache  *AssetCache
	server *AssetServer
}

// NewSessionAssets creates a new session assets instance with the provided base URLs
func NewSessionAssets(cache *AssetCache, server *AssetServer) flows.SessionAssets {
	return &sessionAssets{cache: cache, server: server}
}

// HasLocations returns whether locations are supported as an asset item type
func (s *sessionAssets) HasLocations() bool {
	return s.server.isTypeSupported(assetTypeLocationHierarchy)
}

// GetLocationHierarchy gets the location hierarchy asset for the session
func (s *sessionAssets) GetLocationHierarchy() (*utils.LocationHierarchy, error) {
	asset, err := s.cache.getSetAsset(s.server, assetTypeLocationHierarchy)
	if err != nil {
		return nil, err
	}
	hierarchy, isType := asset.(*utils.LocationHierarchy)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return hierarchy, nil
}

// GetChannel gets a channel asset for the session
func (s *sessionAssets) GetChannel(uuid flows.ChannelUUID) (flows.Channel, error) {
	asset, err := s.cache.getItemAsset(s.server, assetTypeChannel, string(uuid))
	if err != nil {
		return nil, err
	}
	channel, isType := asset.(flows.Channel)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type for UUID '%s'", uuid)
	}
	return channel, nil
}

// GetField gets a contact field asset for the session
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

// GetFieldSet gets the set of all fields asset for the session
func (s *sessionAssets) GetFieldSet() (*flows.FieldSet, error) {
	asset, err := s.cache.getSetAsset(s.server, assetTypeField)
	if err != nil {
		return nil, err
	}
	fields, isType := asset.(*flows.FieldSet)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return fields, nil
}

// GetFlow gets a flow asset for the session
func (s *sessionAssets) GetFlow(uuid flows.FlowUUID) (flows.Flow, error) {
	asset, err := s.cache.getItemAsset(s.server, assetTypeFlow, string(uuid))
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

// GetGroupSet gets the set of all groups asset for the session
func (s *sessionAssets) GetGroupSet() (*flows.GroupSet, error) {
	asset, err := s.cache.getSetAsset(s.server, assetTypeGroup)
	if err != nil {
		return nil, err
	}
	groups, isType := asset.(*flows.GroupSet)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return groups, nil
}

// GetLabel gets a message label asset for the session
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
	asset, err := s.cache.getSetAsset(s.server, assetTypeLabel)
	if err != nil {
		return nil, err
	}
	labels, isType := asset.(*flows.LabelSet)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return labels, nil
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type assetEnvelope struct {
	URL      assetURL        `json:"url" validate:"required,url"`
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
		assetReader = func(data json.RawMessage) (interface{}, error) { return utils.ReadLocationHierarchy(data) }
	} else if itemType == assetTypeChannel && !isSet {
		assetReader = func(data json.RawMessage) (interface{}, error) { return flows.ReadChannel(data) }
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
