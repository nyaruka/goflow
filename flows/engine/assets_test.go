package engine_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows/engine"

	"github.com/stretchr/testify/assert"
)

func TestSessionAssets(t *testing.T) {
	server := engine.NewMockAssetServer(assets.NewAssetCache(100, 10))
	server.MockResponse("http://testserver/assets/label/", json.RawMessage(`{"results": []}`))
	server.MockResponse("http://testserver/assets/group/", json.RawMessage(`{
		"results": [
			{
				"uuid": "2aad21f6-30b7-42c5-bd7f-1b720c154817",
				"name": "Survey Audience"
			}
		]
	}`))

	sessionAssets, err := engine.NewSessionAssets(engine.NewServerSource(server))
	assert.NoError(t, err)

	group, err := sessionAssets.GetGroup(assets.GroupUUID("2aad21f6-30b7-42c5-bd7f-1b720c154817"))
	assert.NoError(t, err)
	assert.Equal(t, assets.GroupUUID("2aad21f6-30b7-42c5-bd7f-1b720c154817"), group.UUID())
	assert.Equal(t, "Survey Audience", group.Name())

	// requesting a group actually fetches and caches the entire group set
	assert.Equal(t, server.MockedRequests(), []string{"http://testserver/assets/label/", "http://testserver/assets/group/"})
}
