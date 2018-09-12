package assets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

type LegacyServer interface {
	GetAsset(AssetType, string) (interface{}, error)
}

type AssetServer struct {
	authToken  string
	typeURLs   map[AssetType]string
	httpClient *utils.HTTPClient
	cache      *AssetCache

	fetcher assetFetcher
}

var _ AssetSource = (*AssetServer)(nil)
var _ assetFetcher = (*AssetServer)(nil)

// NewAssetServer creates a new asset server
func NewAssetServer(authToken string, typeURLs map[AssetType]string, httpClient *utils.HTTPClient, cache *AssetCache) *AssetServer {
	// TODO validate typeURLs are for registered types?

	s := &AssetServer{authToken: authToken, typeURLs: typeURLs, httpClient: httpClient, cache: cache}
	s.fetcher = s
	return s
}

type assetServerEnvelope struct {
	TypeURLs map[AssetType]string `json:"type_urls"`
}

// ReadAssetServer reads an asset server fronm the given JSON
func ReadAssetServer(authToken string, httpClient *utils.HTTPClient, cache *AssetCache, data json.RawMessage) (*AssetServer, error) {
	envelope := &assetServerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, envelope); err != nil {
		return nil, fmt.Errorf("unable to read asset server: %s", err)
	}

	return NewAssetServer(authToken, envelope.TypeURLs, httpClient, cache), nil
}

func (s *AssetServer) Channels() ([]Channel, error) {
	asset, err := s.GetAsset(assetTypeChannel, "")
	if err != nil {
		return nil, err
	}
	set, isType := asset.([]Channel)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return set, nil
}

func (s *AssetServer) Groups() ([]Group, error) {
	asset, err := s.GetAsset(assetTypeGroup, "")
	if err != nil {
		return nil, err
	}
	set, isType := asset.([]Group)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return set, nil
}

func (s *AssetServer) Labels() ([]Label, error) {
	asset, err := s.GetAsset(assetTypeLabel, "")
	if err != nil {
		return nil, err
	}
	set, isType := asset.([]Label)
	if !isType {
		return nil, fmt.Errorf("asset cache contains asset with wrong type")
	}
	return set, nil
}

func (s *AssetServer) HasLocations() bool {
	_, hasTypeURL := s.typeURLs["location_hierarchy"]
	return hasTypeURL
}

func (s *AssetServer) GetAsset(itemType AssetType, itemUUID string) (interface{}, error) {
	url, err := s.getAssetURL(itemType, itemUUID)
	if err != nil {
		return nil, err
	}

	return s.cache.getAsset(url, s.fetcher, itemType)
}

// getAssetURL gets the URL for a set of the given asset type
func (s *AssetServer) getAssetURL(itemType AssetType, itemUUID string) (string, error) {
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
func (s *AssetServer) fetchAsset(url string, itemType AssetType) ([]byte, error) {
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

type MockAssetServer struct {
	AssetServer

	mockResponses  map[string]json.RawMessage
	mockedRequests []string
}

// NewMockAssetServer creates a new mocked asset server for testing
func NewMockAssetServer(typeURLs map[AssetType]string, cache *AssetCache) *MockAssetServer {
	s := &MockAssetServer{
		AssetServer:    AssetServer{typeURLs: typeURLs, cache: cache},
		mockResponses:  map[string]json.RawMessage{},
		mockedRequests: []string{},
	}
	s.AssetServer.fetcher = s
	return s
}

func (s *MockAssetServer) MockResponse(url string, response json.RawMessage) {
	s.mockResponses[url] = response
}

func (s *MockAssetServer) MockedRequests() []string {
	return s.mockedRequests
}

func (s *MockAssetServer) fetchAsset(url string, itemType AssetType) ([]byte, error) {
	s.mockedRequests = append(s.mockedRequests, url)

	assetBuf, found := s.mockResponses[url]
	if !found {
		return []byte(`{"results":[]}`), nil
	}
	return assetBuf, nil
}

// MarshalJSON marshals this mock asset server into JSON
func (s *MockAssetServer) MarshalJSON() ([]byte, error) {
	envelope := &assetServerEnvelope{TypeURLs: s.typeURLs}
	return json.Marshal(envelope)
}

var _ LegacyServer = (*MockAssetServer)(nil)
var _ assetFetcher = (*MockAssetServer)(nil)
