package actions_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var actionTests = []struct {
	action flows.Action
	json   string
}{
	{
		actions.NewAddContactGroupsAction(
			flows.ActionUUID("ad154980-7bf7-4ab8-8728-545fd6378912"),
			[]*assets.GroupReference{
				assets.NewGroupReference(assets.GroupUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Testers"),
				assets.NewVariableGroupReference("@(format_location(contact.fields.state)) Members"),
			},
		),
		`{
			"type": "add_contact_groups",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"groups": [
				{
					"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
					"name": "Testers"
				},
				{
					"name_match": "@(format_location(contact.fields.state)) Members"
				}
			]
		}`,
	},
	{
		actions.NewAddContactURNAction(
			flows.ActionUUID("ad154980-7bf7-4ab8-8728-545fd6378912"),
			"tel",
			"+234532626677",
		),
		`{
			"type": "add_contact_urn",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"scheme": "tel",
			"path": "+234532626677"
		}`,
	},
	{
		actions.NewAddInputLabelsAction(
			flows.ActionUUID("ad154980-7bf7-4ab8-8728-545fd6378912"),
			[]*assets.LabelReference{
				assets.NewLabelReference(assets.LabelUUID("3f65d88a-95dc-4140-9451-943e94e06fea"), "Spam"),
				assets.NewVariableLabelReference("@(format_location(contact.fields.state)) Messages"),
			},
		),
		`{
			"type": "add_input_labels",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"labels": [
				{
					"uuid": "3f65d88a-95dc-4140-9451-943e94e06fea",
					"name": "Spam"
				},
				{
					"name_match": "@(format_location(contact.fields.state)) Messages"
				}
			]
		}`,
	},
	{
		actions.NewCallResthookAction(
			flows.ActionUUID("ad154980-7bf7-4ab8-8728-545fd6378912"),
			"new-registration",
			"My Result",
		),
		`{
			"type": "call_resthook",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"resthook": "new-registration",
			"result_name": "My Result"
		}`,
	},
}

func TestActions(t *testing.T) {
	session, err := test.CreateTestSession("", nil)
	require.NoError(t, err)

	for _, tc := range actionTests {
		// test validating the action
		err := tc.action.Validate(session.Assets())
		assert.NoError(t, err)

		actualJSON, err := json.Marshal(tc.action)
		assert.NoError(t, err)

		test.AssertEqualJSON(t, json.RawMessage(tc.json), actualJSON, "new action produced unexpected JSON")
	}
}
