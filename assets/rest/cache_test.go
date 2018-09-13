package rest_test

import (
	"encoding/json"
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

	// try to get all labels
	labels, err := source.Labels()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(labels))
	assert.Equal(t, "Test", labels[0].Name())
	assert.Equal(t, source.MockedRequests(), []string{"http://testserver/assets/label/"})

	// check that we can refetch without making another server request
	labels, err = source.Labels()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(labels))
	assert.Equal(t, source.MockedRequests(), []string{"http://testserver/assets/label/"})
}
