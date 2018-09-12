package assets_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/nyaruka/goflow/assets"

	"github.com/stretchr/testify/assert"
)

const (
	fooAssetType assets.AssetType = "foo"
)

type fooAsset struct {
	Value int `json:"value"`
}

func readFooAssets(data json.RawMessage) ([]*fooAsset, error) {
	var foos []*fooAsset
	if err := json.Unmarshal(data, &foos); err != nil {
		return nil, err
	}
	return foos, nil
}

func TestAssetCache(t *testing.T) {
	assets.RegisterType(fooAssetType, true, func(data json.RawMessage) (interface{}, error) { return readFooAssets(data) })

	server := assets.NewMockServerSource(map[assets.AssetType]string{
		fooAssetType: "http://testserver/assets/foo/",
	}, assets.NewAssetCache(100, 10))
	server.MockResponse("http://testserver/assets/foo/", json.RawMessage(`{
		"results": [
			{"value": 123},
			{"value": 234}
		]
	}`))

	// can't get an non-registered asset type
	asset, err := server.GetAsset(assets.AssetType("pizza"), "")
	assert.EqualError(t, err, "asset type 'pizza' not supported by asset server")

	// try to get all foos
	asset, err = server.GetAsset(fooAssetType, "")
	assert.NoError(t, err)
	assert.Equal(t, server.MockedRequests(), []string{"http://testserver/assets/foo/"})

	// check we got an asset of the expected type
	foos, isFooSlice := asset.([]*fooAsset)
	assert.True(t, isFooSlice, "expecting slice of foo objects but got something of type %s", reflect.TypeOf(asset))
	assert.Equal(t, 2, len(foos))
	assert.Equal(t, 123, foos[0].Value)

	// check that we can refetch without making another server request
	asset, err = server.GetAsset(fooAssetType, "")
	assert.NoError(t, err)
	assert.Equal(t, server.MockedRequests(), []string{"http://testserver/assets/foo/"})
}
