package engine_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/rest"
	"github.com/nyaruka/goflow/flows/engine"

	"github.com/stretchr/testify/assert"
)

func TestSessionAssets(t *testing.T) {
	server := rest.NewMockServerSource(rest.NewAssetCache(100, 10))
	server.MockResponse("http://testserver/assets/channel/", json.RawMessage(`{"results": []}`))
	server.MockResponse("http://testserver/assets/field/", json.RawMessage(`{"results": []}`))
	server.MockResponse("http://testserver/assets/group/", json.RawMessage(`{
		"results": [
			{
				"uuid": "2aad21f6-30b7-42c5-bd7f-1b720c154817",
				"name": "Survey Audience"
			}
		]
	}`))
	server.MockResponse("http://testserver/assets/label/", json.RawMessage(`{"results": []}`))
	server.MockResponse("http://testserver/assets/resthook/", json.RawMessage(`{"results": []}`))

	sessionAssets, err := engine.NewSessionAssets(server)
	assert.NoError(t, err)

	group, err := sessionAssets.Groups().Get(assets.GroupUUID("2aad21f6-30b7-42c5-bd7f-1b720c154817"))
	assert.NoError(t, err)
	assert.Equal(t, assets.GroupUUID("2aad21f6-30b7-42c5-bd7f-1b720c154817"), group.UUID())
	assert.Equal(t, "Survey Audience", group.Name())

	// requesting a group actually fetches and caches the entire group set
	assert.Equal(t, server.MockedRequests(), []string{
		"http://testserver/assets/channel/",
		"http://testserver/assets/field/",
		"http://testserver/assets/group/",
		"http://testserver/assets/label/",
		"http://testserver/assets/resthook/",
	})
}
