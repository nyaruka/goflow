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
	server.MockResponse("http://testserver/assets/label/f2a3e00c-e86a-4282-a9e8-bb2275e1b9a4/", json.RawMessage(`{
		"uuid": "f2a3e00c-e86a-4282-a9e8-bb2275e1b9a4",
		"name": "Spam"
	}`))
	cache := NewAssetCache(100, 10, "testing/1.0")

	asset, err := cache.getSetAsset(server, assetType("pizza"))
	assert.EqualError(t, err, "asset type 'pizza' not supported by asset server")

	asset, err = cache.getItemAsset(server, assetTypeLabel, "f2a3e00c-e86a-4282-a9e8-bb2275e1b9a4")
	assert.NoError(t, err)
	assert.Equal(t, server.MockedRequests(), []string{"http://testserver/assets/label/f2a3e00c-e86a-4282-a9e8-bb2275e1b9a4/"})

	label, isLabel := asset.(*flows.Label)
	assert.True(t, isLabel, "expecting label but got something of type %s", reflect.TypeOf(asset))
	assert.Equal(t, flows.LabelUUID("f2a3e00c-e86a-4282-a9e8-bb2275e1b9a4"), label.UUID())
	assert.Equal(t, "Spam", label.Name())

	// check that we can refetch without making another server request
	asset, err = cache.getItemAsset(server, assetTypeLabel, "f2a3e00c-e86a-4282-a9e8-bb2275e1b9a4")
	assert.NoError(t, err)
	assert.Equal(t, server.MockedRequests(), []string{"http://testserver/assets/label/f2a3e00c-e86a-4282-a9e8-bb2275e1b9a4/"})
}

func TestAssetServer(t *testing.T) {
	server := NewMockAssetServer()
	server.MockResponse("http://testserver/assets/group/2aad21f6-30b7-42c5-bd7f-1b720c154817/", json.RawMessage(`{
		"uuid": "2aad21f6-30b7-42c5-bd7f-1b720c154817",
		"name": "Survey Audience"
	}`))

	url, err := server.getSetAssetURL(assetType("pizza"))
	assert.EqualError(t, err, "asset type 'pizza' not supported by asset server")

	url, err = server.getSetAssetURL(assetTypeGroup)
	assert.NoError(t, err)
	assert.Equal(t, "http://testserver/assets/group/", url)

	url, err = server.getItemAssetURL(assetTypeGroup, "2aad21f6-30b7-42c5-bd7f-1b720c154817")
	assert.NoError(t, err)
	assert.Equal(t, "http://testserver/assets/group/2aad21f6-30b7-42c5-bd7f-1b720c154817/", url)

	asset, err := server.fetchAsset(url, assetTypeGroup, false, "testing/1.0")
	assert.NoError(t, err)
	assert.Equal(t, server.mockedRequests, []string{"http://testserver/assets/group/2aad21f6-30b7-42c5-bd7f-1b720c154817/"})

	group, isGroup := asset.(*flows.Group)
	assert.True(t, isGroup, "expecting group but got something of type %s", reflect.TypeOf(asset))
	assert.Equal(t, flows.GroupUUID("2aad21f6-30b7-42c5-bd7f-1b720c154817"), group.UUID())
	assert.Equal(t, "Survey Audience", group.Name())
}

func TestSessionAssets(t *testing.T) {
	server := NewMockAssetServer()
	server.mockResponses["http://testserver/assets/group/"] = json.RawMessage(`[
		{
			"uuid": "2aad21f6-30b7-42c5-bd7f-1b720c154817",
			"name": "Survey Audience"
		}
	]`)
	cache := NewAssetCache(100, 10, "testing/1.0")
	sessionAssets := NewSessionAssets(cache, server)

	group, err := sessionAssets.GetGroup(flows.GroupUUID("2aad21f6-30b7-42c5-bd7f-1b720c154817"))
	assert.NoError(t, err)
	assert.Equal(t, flows.GroupUUID("2aad21f6-30b7-42c5-bd7f-1b720c154817"), group.UUID())
	assert.Equal(t, "Survey Audience", group.Name())

	// requesting a group actually fetches and caches the entire group set
	assert.Equal(t, server.mockedRequests, []string{"http://testserver/assets/group/"})
}
