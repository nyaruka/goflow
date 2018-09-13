package rest

import (
	"encoding/json"
)

// MockServerSource is a version of the server source which allows requests to be mocked for testing
type MockServerSource struct {
	ServerSource

	mockResponses  map[string]json.RawMessage
	mockedRequests []string
}

// NewMockServerSource creates a new mocked asset server for testing
func NewMockServerSource(cache *AssetCache) *MockServerSource {
	s := &MockServerSource{
		ServerSource: ServerSource{typeURLs: map[AssetType]string{
			assetTypeChannel:           "http://testserver/assets/channel/",
			assetTypeField:             "http://testserver/assets/field/",
			assetTypeFlow:              "http://testserver/assets/flow/",
			assetTypeGroup:             "http://testserver/assets/group/",
			assetTypeLabel:             "http://testserver/assets/label/",
			assetTypeLocationHierarchy: "http://testserver/assets/location_hierarchy/",
			assetTypeResthook:          "http://testserver/assets/resthook/",
		}, cache: cache},
		mockResponses:  map[string]json.RawMessage{},
		mockedRequests: []string{},
	}
	s.ServerSource.fetcher = s
	return s
}

// MockResponse creates a new mocked response for the given URL
func (s *MockServerSource) MockResponse(url string, response json.RawMessage) {
	s.mockResponses[url] = response
}

// MockedRequests returns all mocked requests made so far
func (s *MockServerSource) MockedRequests() []string {
	return s.mockedRequests
}

func (s *MockServerSource) fetchAsset(url string, itemType AssetType) ([]byte, error) {
	s.mockedRequests = append(s.mockedRequests, url)

	assetBuf, found := s.mockResponses[url]
	if !found {
		return []byte(`{"results":[]}`), nil
	}
	return assetBuf, nil
}

// MarshalJSON marshals this mock asset server into JSON
func (s *MockServerSource) MarshalJSON() ([]byte, error) {
	envelope := &serverSourceEnvelope{TypeURLs: s.typeURLs}
	return json.Marshal(envelope)
}
