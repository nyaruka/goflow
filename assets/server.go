package assets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nyaruka/goflow/utils"

	log "github.com/sirupsen/logrus"
)

// AssetServer describes the asset server we'll fetch missing assets from
type AssetServer interface {
	IsTypeSupported(AssetType) bool
	GetAsset(AssetType, string) (interface{}, error)
}

type assetServer struct {
	authToken  string
	typeURLs   map[AssetType]string
	httpClient *utils.HTTPClient
	cache      *AssetCache
}

var _ assetFetcher = (*assetServer)(nil)

type assetServerEnvelope struct {
	TypeURLs map[AssetType]string `json:"type_urls"`
}

// ReadAssetServer reads an asset server fronm the given JSON
func ReadAssetServer(authToken string, httpClient *utils.HTTPClient, cache *AssetCache, data json.RawMessage) (AssetServer, error) {
	envelope := &assetServerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, envelope); err != nil {
		return nil, fmt.Errorf("unable to read asset server: %s", err)
	}

	return NewAssetServer(authToken, envelope.TypeURLs, httpClient, cache), nil
}

// NewAssetServer creates a new asset server
func NewAssetServer(authToken string, typeURLs map[AssetType]string, httpClient *utils.HTTPClient, cache *AssetCache) AssetServer {
	// TODO validate typeURLs are for registered types?

	return &assetServer{authToken: authToken, typeURLs: typeURLs, httpClient: httpClient, cache: cache}
}

// IsTypeSupported returns whether the given asset item type is supported
func (s *assetServer) IsTypeSupported(itemType AssetType) bool {
	_, hasTypeURL := s.typeURLs[itemType]
	return hasTypeURL
}

// GetAsset returns either a single item or the set of items
func (s *assetServer) GetAsset(itemType AssetType, itemUUID string) (interface{}, error) {
	url, err := s.getAssetURL(itemType, itemUUID)
	if err != nil {
		return nil, err
	}

	return s.cache.getAsset(url, s, itemType)
}

// getAssetURL gets the URL for a set of the given asset type
func (s *assetServer) getAssetURL(itemType AssetType, itemUUID string) (string, error) {
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
func (s *assetServer) fetchAsset(url string, itemType AssetType) ([]byte, error) {
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
	assetServer

	mockResponses  map[string]json.RawMessage
	mockedRequests []string
}

// NewMockAssetServer creates a new mocked asset server for testing
func NewMockAssetServer(typeURLs map[AssetType]string, cache *AssetCache) *MockAssetServer {
	return &MockAssetServer{
		assetServer:    assetServer{typeURLs: typeURLs, cache: cache},
		mockResponses:  map[string]json.RawMessage{},
		mockedRequests: []string{},
	}
}

func (s *MockAssetServer) MockResponse(url string, response json.RawMessage) {
	s.mockResponses[url] = response
}

func (s *MockAssetServer) MockedRequests() []string {
	return s.mockedRequests
}

// GetAsset returns either a single item or the set of items
func (s *MockAssetServer) GetAsset(itemType AssetType, itemUUID string) (interface{}, error) {
	url, err := s.getAssetURL(itemType, itemUUID)
	if err != nil {
		return nil, err
	}

	return s.cache.getAsset(url, s, itemType)
}

func (s *MockAssetServer) fetchAsset(url string, itemType AssetType) ([]byte, error) {
	s.mockedRequests = append(s.mockedRequests, url)

	assetBuf, found := s.mockResponses[url]
	if !found {
		return nil, fmt.Errorf("mock asset server has no mocked response for URL: %s", url)
	}
	return assetBuf, nil
}

// MarshalJSON marshals this mock asset server into JSON
func (s *MockAssetServer) MarshalJSON() ([]byte, error) {
	envelope := &assetServerEnvelope{}
	envelope.TypeURLs = s.typeURLs
	return json.Marshal(envelope)
}
