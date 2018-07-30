package assets

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/nyaruka/goflow/flows"

	"github.com/stretchr/testify/assert"
)

func TestAssetCache(t *testing.T) {
	server := NewMockAssetServer()
	server.MockResponse("http://testserver/assets/label/", json.RawMessage(`[{
		"uuid": "f2a3e00c-e86a-4282-a9e8-bb2275e1b9a4",
		"name": "Spam"
	}]`))
	cache := NewAssetCache(100, 10)

	asset, err := cache.GetAsset(server, assetType("pizza"), "")
	assert.EqualError(t, err, "asset type 'pizza' not supported by asset server")

	asset, err = cache.GetAsset(server, assetTypeLabel, "")
	assert.NoError(t, err)
	assert.Equal(t, server.MockedRequests(), []string{"http://testserver/assets/label/"})

	labelSet, isLabelSet := asset.(*flows.LabelSet)
	assert.True(t, isLabelSet, "expecting label set but got something of type %s", reflect.TypeOf(asset))
	assert.NotNil(t, labelSet.FindByName("Spam"))

	// check that we can refetch without making another server request
	asset, err = cache.GetAsset(server, assetTypeLabel, "")
	assert.NoError(t, err)
	assert.Equal(t, server.MockedRequests(), []string{"http://testserver/assets/label/"})
}

func TestAssetServer(t *testing.T) {
	server := NewMockAssetServer()
	server.MockResponse("http://testserver/assets/group/", json.RawMessage(`[{
		"uuid": "da310302-2340-4cee-b5bb-5ee37a24a122",
		"name": "Survey Audience"
	}]`))
	server.MockResponse("http://testserver/assets/flow/2aad21f6-30b7-42c5-bd7f-1b720c154817/", json.RawMessage(`{
		"uuid": "2aad21f6-30b7-42c5-bd7f-1b720c154817",
		"name": "Registration",
		"nodes": []
	}`))

	url, err := server.getAssetURL(assetType("pizza"), "")
	assert.EqualError(t, err, "asset type 'pizza' not supported by asset server")

	url, err = server.getAssetURL(assetTypeGroup, "")
	assert.NoError(t, err)
	assert.Equal(t, "http://testserver/assets/group/", url)

	url, err = server.getAssetURL(assetTypeFlow, "2aad21f6-30b7-42c5-bd7f-1b720c154817")
	assert.NoError(t, err)
	assert.Equal(t, "http://testserver/assets/flow/2aad21f6-30b7-42c5-bd7f-1b720c154817/", url)

	asset, err := server.fetchAsset(url, assetTypeFlow)
	assert.NoError(t, err)
	assert.Equal(t, []string{"http://testserver/assets/flow/2aad21f6-30b7-42c5-bd7f-1b720c154817/"}, server.mockedRequests)

	flow, isFlow := asset.(flows.Flow)
	assert.True(t, isFlow, "expecting flow but got something of type %s", reflect.TypeOf(asset))
	assert.Equal(t, "Registration", flow.Name())
}

func TestSessionAssets(t *testing.T) {
	server := NewMockAssetServer()
	server.mockResponses["http://testserver/assets/group/"] = json.RawMessage(`[
		{
			"uuid": "2aad21f6-30b7-42c5-bd7f-1b720c154817",
			"name": "Survey Audience"
		}
	]`)
	cache := NewAssetCache(100, 10)
	sessionAssets := NewSessionAssets(cache, server)

	group, err := sessionAssets.GetGroup(flows.GroupUUID("2aad21f6-30b7-42c5-bd7f-1b720c154817"))
	assert.NoError(t, err)
	assert.Equal(t, flows.GroupUUID("2aad21f6-30b7-42c5-bd7f-1b720c154817"), group.UUID())
	assert.Equal(t, "Survey Audience", group.Name())

	// requesting a group actually fetches and caches the entire group set
	assert.Equal(t, server.mockedRequests, []string{"http://testserver/assets/group/"})
}
