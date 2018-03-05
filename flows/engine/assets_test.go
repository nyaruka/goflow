package engine

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/stretchr/testify/assert"
)

func TestAssetCache(t *testing.T) {
	server := NewMockAssetServer().(*mockAssetServer)
	server.mockResponses["http://testserver/assets/label/f2a3e00c-e86a-4282-a9e8-bb2275e1b9a4/"] = json.RawMessage(`{
		"uuid": "f2a3e00c-e86a-4282-a9e8-bb2275e1b9a4",
		"name": "Spam"
	}`)
	cache := NewAssetCache(100, 10, "testing/1.0")

	asset, err := cache.getSetAsset(server, assetType("pizza"))
	assert.EqualError(t, err, "asset type 'pizza' not supported by asset server")

	asset, err = cache.getItemAsset(server, assetTypeLabel, "f2a3e00c-e86a-4282-a9e8-bb2275e1b9a4")
	assert.NoError(t, err)
	assert.Equal(t, server.mockedRequests, []string{"http://testserver/assets/label/f2a3e00c-e86a-4282-a9e8-bb2275e1b9a4/"})

	label, isLabel := asset.(*flows.Label)
	assert.True(t, isLabel, "expecting label but got something of type %s", reflect.TypeOf(asset))
	assert.Equal(t, flows.LabelUUID("f2a3e00c-e86a-4282-a9e8-bb2275e1b9a4"), label.UUID())
	assert.Equal(t, "Spam", label.Name())

	// check that we can refetch without making another server request
	asset, err = cache.getItemAsset(server, assetTypeLabel, "f2a3e00c-e86a-4282-a9e8-bb2275e1b9a4")
	assert.NoError(t, err)
	assert.Equal(t, server.mockedRequests, []string{"http://testserver/assets/label/f2a3e00c-e86a-4282-a9e8-bb2275e1b9a4/"})
}

func TestAssetServer(t *testing.T) {
	server := NewMockAssetServer().(*mockAssetServer)
	server.mockResponses["http://testserver/assets/group/2aad21f6-30b7-42c5-bd7f-1b720c154817/"] = json.RawMessage(`{
		"uuid": "2aad21f6-30b7-42c5-bd7f-1b720c154817",
		"name": "Survey Audience"
	}`)

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
	server := NewMockAssetServer().(*mockAssetServer)
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

func TestFlowValidation(t *testing.T) {
	assetsJSON, err := ioutil.ReadFile("testdata/assets.json")
	assert.NoError(t, err)

	// build our session
	assetCache := NewAssetCache(100, 5, "testing/1.0")
	err = assetCache.Include(assetsJSON)
	assert.NoError(t, err)

	session := NewSession(assetCache, NewMockAssetServer())
	flow, err := session.Assets().GetFlow("76f0a02f-3b75-4b86-9064-e9195e1b3a02")
	assert.NoError(t, err)

	// break the add_label action so references an invalid label
	addLabelAction := flow.Nodes()[0].Actions()[0].(*actions.AddLabelAction)
	addLabelAction.Labels[0].UUID = "xyx"

	// check that validation fails
	err = flow.Validate(session.Assets())
	assert.EqualError(t, err, "validation failed for action[uuid=ad154980-7bf7-4ab8-8728-545fd6378912, type=add_label]: no such label with uuid 'xyx'")

	// fix the add_label action
	addLabelAction.Labels[0].UUID = "3f65d88a-95dc-4140-9451-943e94e06fea"

	// break the add_group action so references an invalid group
	addGroupAction := flow.Nodes()[0].Actions()[1].(*actions.AddToGroupAction)
	addGroupAction.Groups[0].UUID = "xyx"

	// check that validation fails
	err = flow.Validate(session.Assets())
	assert.EqualError(t, err, "validation failed for action[uuid=ad154980-7bf7-4ab8-8728-545fd6378912, type=add_to_group]: no such group with uuid 'xyx'")

	// fix the add_group action
	addGroupAction.Groups[0].UUID = "2aad21f6-30b7-42c5-bd7f-1b720c154817"

	// break the save_contact_field action so references an invalid field
	saveContactAction := flow.Nodes()[0].Actions()[2].(*actions.SaveContactField)
	saveContactAction.Field.Key = "xyx"

	// check that validation fails
	err = flow.Validate(session.Assets())
	assert.EqualError(t, err, "validation failed for action[uuid=7bd8b3bf-0a3c-4928-bc46-df416e77ddf4, type=save_contact_field]: no such field with key 'xyx'")

	// fix the save_contact_field action
	saveContactAction.Field.Key = "first_name"

	// break the set_preferred_channel action so references an invalid channel
	prefChannelAction := flow.Nodes()[0].Actions()[3].(*actions.PreferredChannelAction)
	prefChannelAction.Channel.UUID = "xyx"

	// can't simulate a missing channel without the asset store trying to fetch it... but we can stuff the wrong thing in the store
	session.Assets().(*sessionAssets).cache.addAsset("http://testserver/assets/channel/xyx", flow)

	// check that validation fails
	err = flow.Validate(session.Assets())
	assert.EqualError(t, err, "validation failed for action[uuid=3248a064-bc42-4dff-aa0f-93d85de2f600, type=set_preferred_channel]: asset cache contains asset with wrong type for UUID 'xyx'")

	// fix the set_preferred_channel action
	prefChannelAction.Channel.UUID = "57f1078f-88aa-46f4-a59a-948a5739c03d"
}
