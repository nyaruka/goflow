package assets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/nyaruka/goflow/utils"

	log "github.com/sirupsen/logrus"
)

// AssetServer describes the asset server we'll fetch missing assets from
type AssetServer interface {
	isTypeSupported(assetType) bool
	getAssetURL(assetType, string) (string, error)
	fetchAsset(string, assetType, string) (interface{}, error)
}

type assetServer struct {
	authToken string
	typeURLs  map[assetType]string
}

type assetServerEnvelope struct {
	TypeURLs map[assetType]string `json:"type_urls"`
}

// ReadAssetServer reads an asset server fronm the given JSON
func ReadAssetServer(authToken string, data json.RawMessage) (AssetServer, error) {
	envelope := &assetServerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, envelope, "asset_server"); err != nil {
		return nil, err
	}

	return NewAssetServer(authToken, envelope.TypeURLs), nil
}

// NewAssetServer creates a new asset server
func NewAssetServer(authToken string, typeURLs map[assetType]string) AssetServer {
	return &assetServer{authToken: authToken, typeURLs: typeURLs}
}

// isTypeSupported returns whether the given asset item type is supported
func (s *assetServer) isTypeSupported(itemType assetType) bool {
	_, hasTypeURL := s.typeURLs[itemType]
	return hasTypeURL
}

// getAssetURL gets the URL for a set of the given asset type
func (s *assetServer) getAssetURL(itemType assetType, itemUUID string) (string, error) {
	url, found := s.typeURLs[itemType]
	if !found {
		return "", fmt.Errorf("asset type '%s' not supported by asset server", itemType)
	}

	if itemUUID != "" {
		url = strings.Replace(url, "{uuid}", itemUUID, -1)
	}

	return url, nil
}

// fetches an asset by its URL and parses it as the provided type
func (s *assetServer) fetchAsset(url string, itemType assetType, userAgent string) (interface{}, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// set request headers
	request.Header.Set("User-Agent", userAgent)
	request.Header.Set("Authorization", fmt.Sprintf("Token %s", s.authToken))

	// make the actual request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	log.WithField("asset_type", string(itemType)).WithField("url", url).Debugf("asset requested")

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("asset request returned non-200 response (%d)", response.StatusCode)
	}

	if response.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("asset request returned non-JSON response")
	}

	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return readAsset(buf, itemType)
}

type MockAssetServer struct {
	assetServer

	mockResponses  map[string]json.RawMessage
	mockedRequests []string
}

// NewMockAssetServer creates a new mocked asset server for testing
func NewMockAssetServer() *MockAssetServer {
	return &MockAssetServer{
		assetServer: assetServer{
			typeURLs: map[assetType]string{
				assetTypeChannelSet:        "http://testserver/assets/channel/",
				assetTypeFieldSet:          "http://testserver/assets/field/",
				assetTypeFlow:              "http://testserver/assets/flow/{uuid}/",
				assetTypeGroupSet:          "http://testserver/assets/group/",
				assetTypeLabelSet:          "http://testserver/assets/label/",
				assetTypeLocationHierarchy: "http://testserver/assets/location_hierarchy/",
			},
		},
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

func (s *MockAssetServer) fetchAsset(url string, itemType assetType, userAgent string) (interface{}, error) {
	s.mockedRequests = append(s.mockedRequests, url)

	assetBuf, found := s.mockResponses[url]
	if !found {
		return nil, fmt.Errorf("mock asset server has no mocked response for URL: %s", url)
	}
	return readAsset(assetBuf, itemType)
}

// MarshalJSON marshals this mock asset server into JSON
func (s *MockAssetServer) MarshalJSON() ([]byte, error) {
	envelope := &assetServerEnvelope{}
	envelope.TypeURLs = s.typeURLs
	return json.Marshal(envelope)
}
