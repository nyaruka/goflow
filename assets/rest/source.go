// Package rest is an implementation of AssetSource which fetches assets from a REST server. It maintains
// a cache which can also be preloaded with assets.
package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"

	log "github.com/sirupsen/logrus"
)

const (
	assetTypeChannel           AssetType = "channel"
	assetTypeField             AssetType = "field"
	assetTypeFlow              AssetType = "flow"
	assetTypeGroup             AssetType = "group"
	assetTypeLabel             AssetType = "label"
	assetTypeLocationHierarchy AssetType = "location_hierarchy"
	assetTypeResthook          AssetType = "resthook"
)

// ServerSource is an asset source which fetches assets from a server and caches them
type ServerSource struct {
	authToken  string
	typeURLs   map[AssetType]string
	httpClient *utils.HTTPClient
	cache      *AssetCache

	fetcher assetFetcher
}

var _ assets.AssetSource = (*ServerSource)(nil)
var _ assetFetcher = (*ServerSource)(nil)

// NewServerSource creates a new server asset source
func NewServerSource(authToken string, typeURLs map[AssetType]string, httpClient *utils.HTTPClient, cache *AssetCache) *ServerSource {
	// TODO validate typeURLs are for registered types?

	s := &ServerSource{authToken: authToken, typeURLs: typeURLs, httpClient: httpClient, cache: cache}
	s.fetcher = s
	return s
}

type serverSourceEnvelope struct {
	TypeURLs map[AssetType]string `json:"type_urls"`
}

// ReadServerSource reads a server asset source fronm the given JSON
func ReadServerSource(authToken string, httpClient *utils.HTTPClient, cache *AssetCache, data json.RawMessage) (*ServerSource, error) {
	envelope := &serverSourceEnvelope{}
	if err := utils.UnmarshalAndValidate(data, envelope); err != nil {
		return nil, fmt.Errorf("unable to read asset server: %s", err)
	}

	return NewServerSource(authToken, envelope.TypeURLs, httpClient, cache), nil
}

// Channels returns all channel assets
func (s *ServerSource) Channels() ([]assets.Channel, error) {
	if _, supported := s.typeURLs[assetTypeResthook]; !supported {
		return nil, nil
	}
	asset, err := s.getAsset(assetTypeChannel, "")
	if err != nil {
		return nil, err
	}
	set, isType := asset.([]assets.Channel)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return set, nil
}

// Fields returns all field assets
func (s *ServerSource) Fields() ([]assets.Field, error) {
	if _, supported := s.typeURLs[assetTypeResthook]; !supported {
		return nil, nil
	}
	asset, err := s.getAsset(assetTypeField, "")
	if err != nil {
		return nil, err
	}
	set, isType := asset.([]assets.Field)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return set, nil
}

// Flow returns the flow asset with the given UUID
func (s *ServerSource) Flow(uuid assets.FlowUUID) (assets.Flow, error) {
	if _, supported := s.typeURLs[assetTypeResthook]; !supported {
		return nil, nil
	}
	asset, err := s.getAsset(assetTypeFlow, string(uuid))
	if err != nil {
		return nil, err
	}
	flow, isType := asset.(assets.Flow)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type for UUID '%s'", uuid)
	}
	return flow, nil
}

// Groups returns all group assets
func (s *ServerSource) Groups() ([]assets.Group, error) {
	if _, supported := s.typeURLs[assetTypeResthook]; !supported {
		return nil, nil
	}
	asset, err := s.getAsset(assetTypeGroup, "")
	if err != nil {
		return nil, err
	}
	set, isType := asset.([]assets.Group)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return set, nil
}

// Labels returns all label assets
func (s *ServerSource) Labels() ([]assets.Label, error) {
	if _, supported := s.typeURLs[assetTypeResthook]; !supported {
		return nil, nil
	}
	asset, err := s.getAsset(assetTypeLabel, "")
	if err != nil {
		return nil, err
	}
	set, isType := asset.([]assets.Label)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return set, nil
}

// Locations returns all location assets
func (s *ServerSource) Locations() ([]*utils.LocationHierarchy, error) {
	if _, supported := s.typeURLs[assetTypeResthook]; !supported {
		return nil, nil
	}
	asset, err := s.getAsset(assetTypeLocationHierarchy, "")
	if err != nil {
		return nil, err
	}
	set, isType := asset.([]*utils.LocationHierarchy)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return set, nil
}

// Resthooks returns all resthook assets
func (s *ServerSource) Resthooks() ([]assets.Resthook, error) {
	if _, supported := s.typeURLs[assetTypeResthook]; !supported {
		return nil, nil
	}
	asset, err := s.getAsset(assetTypeResthook, "")
	if err != nil {
		return nil, err
	}
	set, isType := asset.([]assets.Resthook)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return set, nil
}

func (s *ServerSource) getAsset(itemType AssetType, itemUUID string) (interface{}, error) {
	url, err := s.getAssetURL(itemType, itemUUID)
	if err != nil {
		return nil, err
	}

	return s.cache.getAsset(url, s.fetcher, itemType)
}

// getAssetURL gets the URL for a set of the given asset type
func (s *ServerSource) getAssetURL(itemType AssetType, itemUUID string) (string, error) {
	url, found := s.typeURLs[itemType]
	if !found {
		return "", fmt.Errorf("asset type '%s' not supported by asset server", itemType)
	}

	if itemUUID != "" {
		url = fmt.Sprintf("%s%s/", url, itemUUID)
	}

	return url, nil
}

// fetches an asset by its URL and parses it as the provided type
func (s *ServerSource) fetchAsset(url string, itemType AssetType) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// set request headers
	request.Header.Set("Authorization", fmt.Sprintf("Token %s", s.authToken))

	// make the actual request
	response, err := s.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	log.WithField("asset_type", string(itemType)).WithField("url", url).Debugf("asset requested")

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("request returned non-200 response (%d)", response.StatusCode)
	}

	if response.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("request returned non-JSON response")
	}

	return ioutil.ReadAll(response.Body)
}
