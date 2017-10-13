package engine

import (
	"io/ioutil"
	"testing"

	"github.com/nyaruka/goflow/flows/actions"
	"github.com/stretchr/testify/assert"
)

var testAssetURLs = map[AssetItemType]string{
	"channel": "http://testserver/assets/channel",
	"field":   "http://testserver/assets/field",
	"flow":    "http://testserver/assets/flow",
	"group":   "http://testserver/assets/group",
	"label":   "http://testserver/assets/label",
}

func TestFlowValidation(t *testing.T) {
	assetsJSON, err := ioutil.ReadFile("testdata/assets.json")
	assert.NoError(t, err)

	// build our session
	assetCache := NewAssetCache(100, 5)
	err = assetCache.Include(assetsJSON)
	assert.NoError(t, err)

	session := NewSession(assetCache, testAssetURLs)
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

	// check that validation fails
	err = flow.Validate(session.Assets())
	assert.EqualError(t, err, "validation failed for action[uuid=3248a064-bc42-4dff-aa0f-93d85de2f600, type=set_preferred_channel]: Get http://testserver/assets/channel/xyx: dial tcp: lookup testserver: no such host")

	// fix the set_preferred_channel action
	prefChannelAction.Channel.UUID = "57f1078f-88aa-46f4-a59a-948a5739c03d"
}
