package rest_test

import (
	"encoding/json"
	"github.com/nyaruka/goflow/assets"
	"reflect"
	"testing"

	"github.com/nyaruka/goflow/assets/rest"

	"github.com/stretchr/testify/assert"
)

func TestAssetCache(t *testing.T) {
	source := rest.NewMockServerSource(rest.NewAssetCache(100, 10))
	source.MockResponse("http://testserver/assets/label/", json.RawMessage(`{
		"results": [
			{"uuid": "0bb7d9b4-67f9-419d-88c8-806423395067", "name": "Test"},
			{"uuid": "4b2e8e3e-bf59-427d-a969-dc0ec3c55773", "name": "Spam"}
		]
	}`))

	// can't get an non-registered asset type
	asset, err := source.GetAsset(rest.AssetType("pizza"), "")
	assert.EqualError(t, err, "asset type 'pizza' not supported by asset server")

	// try to get all labels
	asset, err = source.GetAsset(rest.AssetTypeLabel, "")
	assert.NoError(t, err)
	assert.Equal(t, source.MockedRequests(), []string{"http://testserver/assets/label/"})

	// check we got an asset of the expected type
	labels, isLabelSlice := asset.([]assets.Label)
	assert.True(t, isLabelSlice, "expecting slice of label objects but got something of type %s", reflect.TypeOf(asset))
	assert.Equal(t, 2, len(labels))
	assert.Equal(t, "Test", labels[0].Name())

	// check that we can refetch without making another server request
	asset, err = source.GetAsset(rest.AssetTypeLabel, "")
	assert.NoError(t, err)
	assert.Equal(t, source.MockedRequests(), []string{"http://testserver/assets/label/"})
}
