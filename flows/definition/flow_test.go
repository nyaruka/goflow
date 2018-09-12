package definition_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testReadingInvalidFlow(t *testing.T, file string, expectedErr string) {
	var err error
	var assetsJSON json.RawMessage
	assetsJSON, err = ioutil.ReadFile(file)
	require.NoError(t, err)

	_, err = definition.ReadFlow(assetsJSON)
	assert.EqualError(t, err, expectedErr)
}

func TestReadFlow(t *testing.T) {
	testReadingInvalidFlow(t,
		"testdata/flow_with_duplicate_node_uuid.json",
		"duplicate node uuid: a58be63b-907d-4a1a-856b-0bb5579d7507",
	)
	testReadingInvalidFlow(t,
		"testdata/flow_with_invalid_default_exit.json",
		"router is invalid on node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: default exit 0680b01f-ba0b-48f4-a688-d2f963130126 is not a valid exit",
	)
	testReadingInvalidFlow(t,
		"testdata/flow_with_invalid_case_exit.json",
		"router is invalid on node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: case exit 37d8813f-1402-4ad2-9cc2-e9054a96525b is not a valid exit",
	)
	testReadingInvalidFlow(t,
		"testdata/flow_with_invalid_case_exit.json",
		"router is invalid on node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: case exit 37d8813f-1402-4ad2-9cc2-e9054a96525b is not a valid exit",
	)
}

func TestFlowValidation(t *testing.T) {
	assetsJSON, err := ioutil.ReadFile("testdata/flow_validation.json")
	assert.NoError(t, err)

	// build our session
	assetCache := assets.NewAssetCache(100, 5)
	err = assetCache.Include(assetsJSON)
	assert.NoError(t, err)

	assets, err := engine.NewSessionAssets(engine.NewMockServerSource(assetCache))
	assert.NoError(t, err)

	session := engine.NewSession(assets, engine.NewDefaultConfig(), test.TestHTTPClient)
	flow, err := session.Assets().GetFlow("76f0a02f-3b75-4b86-9064-e9195e1b3a02")
	assert.NoError(t, err)

	// break the add_input_labels action so references an invalid label
	addLabelAction := flow.Nodes()[0].Actions()[1].(*actions.AddInputLabelsAction)
	addLabelAction.Labels[0].UUID = "xyx"

	// check that validation fails
	err = flow.Validate(session.Assets())
	assert.EqualError(t, err, "validation failed for action[uuid=ad154980-7bf7-4ab8-8728-545fd6378912, type=add_input_labels]: no such label with uuid 'xyx'")

	// fix the add_input_labels action
	addLabelAction.Labels[0].UUID = "3f65d88a-95dc-4140-9451-943e94e06fea"

	// break the add_group action so references an invalid group
	addGroupAction := flow.Nodes()[0].Actions()[2].(*actions.AddContactGroupsAction)
	addGroupAction.Groups[0].UUID = "xyx"

	// check that validation fails
	err = flow.Validate(session.Assets())
	assert.EqualError(t, err, "validation failed for action[uuid=09cd9762-8700-4d14-bbc9-35f75f711873, type=add_contact_groups]: no such group with uuid 'xyx'")

	// fix the add_group action
	addGroupAction.Groups[0].UUID = "2aad21f6-30b7-42c5-bd7f-1b720c154817"

	// break the set_contact_field action so references an invalid field
	saveContactAction := flow.Nodes()[0].Actions()[3].(*actions.SetContactFieldAction)
	saveContactAction.Field.Key = "xyx"

	// check that validation fails
	err = flow.Validate(session.Assets())
	assert.EqualError(t, err, "validation failed for action[uuid=7bd8b3bf-0a3c-4928-bc46-df416e77ddf4, type=set_contact_field]: no such field with key 'xyx'")

	// fix the set_contact_field action
	saveContactAction.Field.Key = "first_name"

	// break the set_contact_channel action so references an invalid channel
	prefChannelAction := flow.Nodes()[0].Actions()[4].(*actions.SetContactChannelAction)
	prefChannelAction.Channel.UUID = "xyx"

	// check that validation fails
	err = flow.Validate(session.Assets())
	assert.EqualError(t, err, "validation failed for action[uuid=3248a064-bc42-4dff-aa0f-93d85de2f600, type=set_contact_channel]: no such channel with uuid 'xyx'")

	// fix the set_contact_channel action
	prefChannelAction.Channel.UUID = "57f1078f-88aa-46f4-a59a-948a5739c03d"
}
