package definition_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/routers"
	"github.com/nyaruka/goflow/flows/waits"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var invalidFlows = []struct {
	path        string
	expectedErr string
}{
	{
		"flow_with_duplicate_node_uuid.json",
		"node UUID a58be63b-907d-4a1a-856b-0bb5579d7507 isn't unique",
	},
	{
		"flow_with_invalid_default_exit.json",
		"validation failed for node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: validation failed for router: default exit 0680b01f-ba0b-48f4-a688-d2f963130126 is not a valid exit",
	},
	{
		"flow_with_invalid_case_exit.json",
		"validation failed for node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: validation failed for router: case exit 37d8813f-1402-4ad2-9cc2-e9054a96525b is not a valid exit",
	},
	{
		"flow_with_invalid_case_exit.json",
		"validation failed for node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: validation failed for router: case exit 37d8813f-1402-4ad2-9cc2-e9054a96525b is not a valid exit",
	},
	{
		"flow_with_invalid_label_asset_ref.json",
		"validation failed for node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: validation failed for action[uuid=ad154980-7bf7-4ab8-8728-545fd6378912, type=add_input_labels]: no such label with UUID 'ab0c8941-64a9-4f48-8949-907d0565e9ad'",
	},
	{
		"flow_with_invalid_group_asset_ref.json",
		"validation failed for node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: validation failed for action[uuid=09cd9762-8700-4d14-bbc9-35f75f711873, type=add_contact_groups]: no such group with UUID 'b27a413d-d737-4a3b-ab42-8a181b52c908'",
	},
	{
		"flow_with_invalid_channel_asset_ref.json",
		"validation failed for node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: validation failed for action[uuid=3248a064-bc42-4dff-aa0f-93d85de2f600, type=set_contact_channel]: no such channel with UUID '038276e5-9223-4143-992b-ef9d7b907030'",
	},
	{
		"flow_with_invalid_field_asset_ref.json",
		"validation failed for node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: validation failed for action[uuid=7bd8b3bf-0a3c-4928-bc46-df416e77ddf4, type=set_contact_field]: no such field with key 'xyz'",
	},
}

func TestFlowValidation(t *testing.T) {
	session, err := test.CreateTestSession("", nil)
	require.NoError(t, err)

	for _, tc := range invalidFlows {
		assetsJSON, err := ioutil.ReadFile("testdata/" + tc.path)
		require.NoError(t, err)

		flow, err := definition.ReadFlow(assetsJSON)
		require.NoError(t, err)

		err = flow.Validate(session.Assets(), flows.NewValidationContext())
		assert.EqualError(t, err, tc.expectedErr)
	}
}

var flowDef = `{
    "uuid": "8ca44c09-791d-453a-9799-a70dd3303306",
    "name": "Test Flow",
    "language": "eng",
    "type": "messaging",
    "revision": 123,
    "expire_after_minutes": 30,
    "localization": null,
    "nodes": [
        {
            "uuid": "a58be63b-907d-4a1a-856b-0bb5579d7507",
            "actions": [
				{
					"type": "send_msg",
					"uuid": "76112ef2-790e-4b5b-84cb-e910f191a335",
					"text": "Do you like beer?"
				}
			],
			"wait": {
				"type": "msg"
			},
			"router": {
				"cases": [
					{
						"uuid": "9f593e22-7886-4c08-a52f-0e8780504d75",
						"type": "has_any_word",
						"arguments": [
							"yes",
							"yeah"
						],
						"exit_uuid": "97b9451c-2856-475b-af38-32af68100897"
					}
				],
				"default_exit_uuid": "8fd08f1c-8f4e-42c1-af6c-df2db2e0eda6",
				"operand": "@input",
				"result_name": "Response 1",
				"type": "switch"
			},
            "exits": [
                {
					"uuid": "97b9451c-2856-475b-af38-32af68100897",
					"destination_node_uuid": "baaf9085-1198-4b41-9a1c-cc51c6dbec99",
					"name": "Yes"
                },
				{
					"uuid": "8fd08f1c-8f4e-42c1-af6c-df2db2e0eda6",
					"destination_node_uuid": "baaf9085-1198-4b41-9a1c-cc51c6dbec99",
					"name": "No"
				}
            ]
		},
		{
            "uuid": "baaf9085-1198-4b41-9a1c-cc51c6dbec99",
            "actions": [
                {
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
                }
            ],
            "exits": [
                {
                    "uuid": "3e077111-7b62-4407-b8a4-4fddaf0d2f24"
                }
            ]
        }
    ]
}`

func TestNewFlow(t *testing.T) {
	session, err := test.CreateTestSession("", nil)
	require.NoError(t, err)

	flow := definition.NewFlow(
		assets.FlowUUID("8ca44c09-791d-453a-9799-a70dd3303306"),
		"Test Flow",           // name
		utils.Language("eng"), // base language
		flows.FlowTypeMessaging,
		123, // revision
		30,  // expires after minutes
		nil, // localization
		[]flows.Node{
			definition.NewNode(
				flows.NodeUUID("a58be63b-907d-4a1a-856b-0bb5579d7507"),
				[]flows.Action{
					actions.NewSendMsgAction(
						flows.ActionUUID("76112ef2-790e-4b5b-84cb-e910f191a335"),
						"Do you like beer?",
						nil,
						nil,
						false,
					),
				},
				waits.NewMsgWait(nil),
				routers.NewSwitchRouter(
					flows.ExitUUID("8fd08f1c-8f4e-42c1-af6c-df2db2e0eda6"),
					"@input",
					[]*routers.Case{
						routers.NewCase(utils.UUID("9f593e22-7886-4c08-a52f-0e8780504d75"), "has_any_word", []string{"yes", "yeah"}, false, flows.ExitUUID("97b9451c-2856-475b-af38-32af68100897")),
					},
					"Response 1",
				),
				[]flows.Exit{
					definition.NewExit(
						flows.ExitUUID("97b9451c-2856-475b-af38-32af68100897"),
						flows.NodeUUID("baaf9085-1198-4b41-9a1c-cc51c6dbec99"),
						"Yes",
					),
					definition.NewExit(
						flows.ExitUUID("8fd08f1c-8f4e-42c1-af6c-df2db2e0eda6"),
						flows.NodeUUID("baaf9085-1198-4b41-9a1c-cc51c6dbec99"),
						"No",
					),
				},
			),
			definition.NewNode(
				flows.NodeUUID("baaf9085-1198-4b41-9a1c-cc51c6dbec99"),
				[]flows.Action{
					actions.NewAddInputLabelsAction(
						flows.ActionUUID("ad154980-7bf7-4ab8-8728-545fd6378912"),
						[]*assets.LabelReference{
							assets.NewLabelReference(assets.LabelUUID("3f65d88a-95dc-4140-9451-943e94e06fea"), "Spam"),
							assets.NewVariableLabelReference("@(format_location(contact.fields.state)) Messages"),
						},
					),
				},
				nil, // no wait
				nil, // no router
				[]flows.Exit{
					definition.NewExit(flows.ExitUUID("3e077111-7b62-4407-b8a4-4fddaf0d2f24"), "", ""),
				},
			),
		},
		nil, // no UI
	)

	// should validate ok
	err = flow.Validate(session.Assets(), flows.NewValidationContext())
	assert.NoError(t, err)

	marshaled, err := json.Marshal(flow)
	assert.NoError(t, err)

	test.AssertEqualJSON(t, []byte(flowDef), marshaled, "flow definition mismatch")
}
