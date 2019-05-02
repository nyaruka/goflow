package definition_test

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/routers"
	"github.com/nyaruka/goflow/flows/routers/waits"
	"github.com/nyaruka/goflow/flows/routers/waits/hints"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlowReadingAndValidation(t *testing.T) {
	invalidFlows := []struct {
		path        string
		expectedErr string
		duringRead  bool
	}{
		{
			"flow_with_exitless_node.json",
			"unable to read node: field 'exits' must have a minimum of 1 items",
			true,
		},
		{
			"flow_with_exitless_category.json",
			"unable to read router: unable to read category: field 'exit_uuid' is required",
			true,
		},
		{
			"flow_with_duplicate_node_uuid.json",
			"node UUID a58be63b-907d-4a1a-856b-0bb5579d7507 isn't unique",
			false,
		},
		{
			"flow_with_invalid_timeout_category.json",
			"validation failed for node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: validation failed for router: timeout category 13fea3d4-b925-495b-b593-1c9e905e700d is not a valid category",
			false,
		},
		{
			"flow_with_invalid_default_exit.json",
			"validation failed for node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: validation failed for router: default category 37d8813f-1402-4ad2-9cc2-e9054a96525b is not a valid category",
			false,
		},
		{
			"flow_with_invalid_case_category.json",
			"validation failed for node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: validation failed for router: case category 37d8813f-1402-4ad2-9cc2-e9054a96525b is not a valid category",
			false,
		},
		{
			"flow_with_invalid_exit_dest.json",
			"validation failed for node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: destination 714f1409-486e-4e8e-bb08-23e2943ef9f6 of exit[uuid=37d8813f-1402-4ad2-9cc2-e9054a96525b] isn't a known node",
			false,
		},
		{
			"flow_with_missing_asset.json",
			"missing dependencies: group[uuid=7be2f40b-38a0-4b06-9e6d-522dca592cc8,name=Registered]",
			false,
		},
	}

	session, _, err := test.CreateTestSession("", nil)
	require.NoError(t, err)

	for _, tc := range invalidFlows {
		assetsJSON, err := ioutil.ReadFile("testdata/" + tc.path)
		require.NoError(t, err)

		flow, err := definition.ReadFlow(assetsJSON)

		if tc.duringRead {
			assert.EqualError(t, err, tc.expectedErr)
		} else {
			require.NoError(t, err)

			err = flow.Validate(session.Assets())
			assert.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestNewFlow(t *testing.T) {
	var flowDef = `
{
    "uuid": "8ca44c09-791d-453a-9799-a70dd3303306", 
    "name": "Test Flow",
    "spec_version": "13.0.0",
    "language": "eng",
    "type": "messaging",
    "revision": 123,
    "expire_after_minutes": 30,
    "localization": {},
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
            "router": {
				"type": "switch",
				"wait": {
					"type": "msg",
					"hint": {
						"type": "image"
					}
				},
				"operand": "@input.text",
                "cases": [
                    {
                        "uuid": "9f593e22-7886-4c08-a52f-0e8780504d75",
                        "type": "has_any_word",
                        "arguments": [
                            "yes",
                            "yeah"
                        ],
                        "category_uuid": "97b9451c-2856-475b-af38-32af68100897"
                    }
                ],
                "default_category_uuid": "8fd08f1c-8f4e-42c1-af6c-df2db2e0eda6",
                "result_name": "Response 1",
				"categories": [
					{
						"uuid": "97b9451c-2856-475b-af38-32af68100897",
						"name": "Yes",
						"exit_uuid": "023a5c10-d74a-4fad-9560-990caead8170"
					},
					{
						"uuid": "8fd08f1c-8f4e-42c1-af6c-df2db2e0eda6",
						"name": "No",
						"exit_uuid": "8943c032-2a91-456c-8080-2a249f1b420c"
					}
				]
            },
            "exits": [
                {
                    "uuid": "023a5c10-d74a-4fad-9560-990caead8170",
                    "destination_uuid": "baaf9085-1198-4b41-9a1c-cc51c6dbec99"
                },
                {
                    "uuid": "8943c032-2a91-456c-8080-2a249f1b420c",
                    "destination_uuid": "baaf9085-1198-4b41-9a1c-cc51c6dbec99"
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
                            "name_match": "@(format_location(contact.fields.gender)) Messages"
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

	session, _, err := test.CreateTestSession("", nil)
	require.NoError(t, err)

	flow := definition.NewFlow(
		assets.FlowUUID("8ca44c09-791d-453a-9799-a70dd3303306"),
		"Test Flow",           // name
		utils.Language("eng"), // base language
		flows.FlowTypeMessaging,
		123, // revision
		30,  // expires after minutes
		definition.NewLocalization(),
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
				routers.NewSwitchRouter(
					waits.NewMsgWait(nil, hints.NewImageHint()),
					"Response 1",
					[]*routers.Category{
						routers.NewCategory(
							flows.CategoryUUID("97b9451c-2856-475b-af38-32af68100897"),
							"Yes",
							flows.ExitUUID("023a5c10-d74a-4fad-9560-990caead8170"),
						),
						routers.NewCategory(
							flows.CategoryUUID("8fd08f1c-8f4e-42c1-af6c-df2db2e0eda6"),
							"No",
							flows.ExitUUID("8943c032-2a91-456c-8080-2a249f1b420c"),
						),
					},
					"@input.text",
					[]*routers.Case{
						routers.NewCase(utils.UUID("9f593e22-7886-4c08-a52f-0e8780504d75"), "has_any_word", []string{"yes", "yeah"}, flows.CategoryUUID("97b9451c-2856-475b-af38-32af68100897")),
					},
					flows.CategoryUUID("8fd08f1c-8f4e-42c1-af6c-df2db2e0eda6"),
				),
				[]flows.Exit{
					definition.NewExit(
						flows.ExitUUID("023a5c10-d74a-4fad-9560-990caead8170"),
						flows.NodeUUID("baaf9085-1198-4b41-9a1c-cc51c6dbec99"),
					),
					definition.NewExit(
						flows.ExitUUID("8943c032-2a91-456c-8080-2a249f1b420c"),
						flows.NodeUUID("baaf9085-1198-4b41-9a1c-cc51c6dbec99"),
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
							assets.NewVariableLabelReference("@(format_location(contact.fields.gender)) Messages"),
						},
					),
				},
				nil, // no router
				[]flows.Exit{
					definition.NewExit(flows.ExitUUID("3e077111-7b62-4407-b8a4-4fddaf0d2f24"), ""),
				},
			),
		},
		nil, // no UI
	)

	marshaled, err := json.Marshal(flow)
	assert.NoError(t, err)

	test.AssertEqualJSON(t, []byte(flowDef), marshaled, "flow definition mismatch")

	// should validate ok
	err = flow.Validate(session.Assets())
	assert.NoError(t, err)

	// check in expressions
	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("Test Flow"),
		"name":        types.NewXText("Test Flow"),
		"revision":    types.NewXNumberFromInt(123),
		"uuid":        types.NewXText("8ca44c09-791d-453a-9799-a70dd3303306"),
	}), flows.Context(session.Environment(), flow))

	// add expected dependencies and result names to our expected JSON
	flowRaw, err := utils.JSONDecodeGeneric([]byte(flowDef))
	require.NoError(t, err)
	flowAsMap := flowRaw.(map[string]interface{})
	flowAsMap[`_dependencies`] = map[string]interface{}{
		"fields": []interface{}{
			map[string]string{"key": "gender", "name": "Gender"},
		},
		"labels": []interface{}{
			map[string]string{"uuid": "3f65d88a-95dc-4140-9451-943e94e06fea", "name": "Spam"},
		},
	}
	flowAsMap[`_results`] = []map[string]interface{}{
		{
			"key":        "response_1",
			"name":       "Response 1",
			"categories": []string{"Yes", "No"},
		},
	}
	flowAsMap[`_waiting_exits`] = []string{
		"023a5c10-d74a-4fad-9560-990caead8170",
		"8943c032-2a91-456c-8080-2a249f1b420c",
	}

	// now when we marshal to JSON, those should be included
	newFlowDef, err := json.Marshal(flowAsMap)
	require.NoError(t, err)

	marshaled, err = json.Marshal(flow)
	assert.NoError(t, err)

	test.AssertEqualJSON(t, []byte(newFlowDef), marshaled, "flow definition mismatch")
}

func TestValidateEmptyFlow(t *testing.T) {
	flow, err := test.LoadFlowFromAssets("../../test/testdata/runner/empty.json", "76f0a02f-3b75-4b86-9064-e9195e1b3a02")
	require.NoError(t, err)

	err = flow.Validate(nil)
	assert.NoError(t, err)

	marshaled, err := json.Marshal(flow)
	require.NoError(t, err)

	test.AssertEqualJSON(t, []byte(`{
    "uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
    "name": "Empty Flow",
    "revision": 0,
    "spec_version": "13.0.0",
    "type": "messaging",
    "expire_after_minutes": 0,
    "language": "eng",
    "localization": {},
    "nodes": [],
    "_dependencies": {},
	"_results": [],
	"_waiting_exits": []
  }`), marshaled, "flow definition mismatch")
}

func assertFlowSection(t *testing.T, definition []byte, key string, data []byte) {
	flowAsMap, err := utils.JSONDecodeGeneric(definition)
	require.NoError(t, err)

	sectionJSON, _ := json.Marshal(flowAsMap.(map[string]interface{})[key])

	test.AssertEqualJSON(t, data, sectionJSON, "flow JSON mismatch")
}

func TestValidateFlow(t *testing.T) {
	sa, err := test.LoadSessionAssets("../../test/testdata/runner/brochure.json")
	require.NoError(t, err)

	flow, err := sa.Flows().Get(assets.FlowUUID("25a2d8b2-ae7c-4fed-964a-506fb8c3f0c0"))
	require.NoError(t, err)

	// validate with session assets
	err = flow.Validate(sa)
	assert.NoError(t, err)

	marshaled, err := json.Marshal(flow)
	require.NoError(t, err)

	// name of group will have been corrected
	assertFlowSection(t, marshaled, "_dependencies", []byte(`{
		"groups": [
			{
				"name": "Registered Users",
				"uuid": "7be2f40b-38a0-4b06-9e6d-522dca592cc8"
			}
		]
	}`))
	assertFlowSection(t, marshaled, "_results", []byte(`[
		{
			"key": "name",
			"name": "Name",
			"categories": ["Not Empty", "Other"]
		}
	]`))
	assertFlowSection(t, marshaled, "_waiting_exits", []byte(`[
		"fc2fcd23-7c4a-44bd-a8c6-6c88e6ed09f8",
        "43accf99-4940-44f7-926b-a8b35d9403d6"
	]`))

	// validate without session assets
	sa, _ = test.LoadSessionAssets("../../test/testdata/runner/brochure.json")
	flow, _ = sa.Flows().Get(assets.FlowUUID("25a2d8b2-ae7c-4fed-964a-506fb8c3f0c0"))
	err = flow.Validate(nil)
	assert.NoError(t, err)

	marshaled, err = json.Marshal(flow)
	require.NoError(t, err)

	// name of group won't have been corrected
	assertFlowSection(t, marshaled, "_dependencies", []byte(`{
    "groups": [
      {
        "name": "Registered",
        "uuid": "7be2f40b-38a0-4b06-9e6d-522dca592cc8"
      }
    ]
  }`))
}

func TestReadFlow(t *testing.T) {
	// try reading something without a flow header
	_, err := definition.ReadFlow([]byte(`{"nodes":[]}`))
	assert.EqualError(t, err, "unable to read flow header: field 'uuid' is required, field 'spec_version' is required")

	// try reading a definition with a newer major version
	_, err = definition.ReadFlow([]byte(`{
    "uuid": "8ca44c09-791d-453a-9799-a70dd3303306", 
    "name": "Test Flow",
    "spec_version": "2000.0",
    "language": "eng",
    "type": "messaging",
    "revision": 123,
    "expire_after_minutes": 30,
    "nodes": []
  }`))
	assert.EqualError(t, err, "spec version 2000.0.0 is newer than this library (13.0.0)")

	// try reading a definition with a newer minor version
	_, err = definition.ReadFlow([]byte(`{
    "uuid": "8ca44c09-791d-453a-9799-a70dd3303306", 
    "name": "Test Flow",
    "spec_version": "13.9999",
    "language": "eng",
    "type": "messaging",
    "revision": 123,
    "expire_after_minutes": 30,
    "nodes": []
  }`))
	assert.NoError(t, err)

	// try reading a definition without a type (a required field in this major version)
	_, err = definition.ReadFlow([]byte(`{
    "uuid": "8ca44c09-791d-453a-9799-a70dd3303306", 
    "name": "Test Flow",
    "spec_version": "13.0",
    "language": "eng",
    "revision": 123,
    "expire_after_minutes": 30,
    "nodes": []
  }`))
	assert.EqualError(t, err, "field 'type' is required")

	// try reading a definition with UI
	flow, err := definition.ReadFlow([]byte(`{
		"uuid": "8ca44c09-791d-453a-9799-a70dd3303306", 
		"name": "Test Flow",
		"spec_version": "13.0",
		"language": "eng",
		"type": "messaging",
		"revision": 123,
		"expire_after_minutes": 30,
		"nodes": [
			{
				"uuid": "b1c5f247-565d-4a7a-8763-c59abbed0a57",
				"exits": [
					{
						"uuid": "9c2412f7-4e8f-44f1-9b4f-0e8f1a274261"
					}
				]
			}
		],
		"_ui": { 
			"nodes": {
				"b1c5f247-565d-4a7a-8763-c59abbed0a57": {
                    "type": "execute_actions"
				}
			},
            "stickies": {}
		}
	  }`))
	assert.NoError(t, err)
	test.AssertEqualJSON(t, []byte(`{
		"nodes": {
			"b1c5f247-565d-4a7a-8763-c59abbed0a57": {
				"type": "execute_actions"
			}
		},
		"stickies": {}
	}`), flow.UI(), "ui mismatch for read flow")
}

func TestExtractAndRewriteTemplates(t *testing.T) {
	testCases := []struct {
		path      string
		uuid      string
		templates []string
	}{
		{
			"../../test/testdata/runner/two_questions.json",
			"615b8a0f-588c-4d20-a05f-363b0b4ce6f4",
			[]string{
				`Hi @contact.name! What is your favorite color? (red/blue) Your number is @(format_urn(contact.urn))`,
				`Red`,
				`Blue`,
				`Quelle est votres couleur preferee? (rouge/blue)`,
				`@input.text`,
				`red`,
				`rouge`,
				`blue`,
				`bleu`,
				`fra`,
				`@(TITLE(results.favorite_color.category_localized)) it is! What is your favorite soda? (pepsi/coke)`,
				`@(TITLE(results.favorite_color.category_localized))! Bien sur! Quelle est votes soda preferee? (pepsi/coke)`,
				`@input.text`,
				`pepsi`,
				`coke coca cola`,
				`http://localhost/?cmd=success`,
				`{ "contact": @(json(contact.uuid)), "soda": @(json(results.soda.value)) }`,
				`Great, you are done and like @results.soda.value! Webhook status was @results.webhook.value`,
				`Parfait, vous avez finis et tu aimes @results.soda.category`,
			},
		},
		{
			"../../test/testdata/runner/all_actions.json",
			"8ca44c09-791d-453a-9799-a70dd3303306",
			[]string{
				`@(format_location(contact.fields.state)) Messages`,
				`@(format_location(contact.fields.state)) Members`,
				`@(replace(lower(contact.name), " ", "_"))`,
				`XXX-YYY-ZZZ`,
				"Here is your activation token",
				"Hi @fields.first_name, Your activation token is @fields.activation_token, your coupon is @(trigger.params.coupons[0].code)",
				"@urns.mailto",
				"test@@example.com",
				`Hi @contact.name, are you ready?`,
				`Hola @contact.name, ¿estás listo?`,
				`Hi @contact.name, are you ready for these attachments?`,
				`image/jpeg:http://s3.amazon.com/bucket/test_en.jpg?a=@(url_encode(format_location(fields.state)))`,
				`Hi @contact.name, are you ready to complete today's survey?`,
				`This is a message to each of @contact.name's urns.`,
				`This is a reply with attachments and quick replies`,
				`image/jpeg:http://s3.amazon.com/bucket/test_en.jpg?a=@(url_encode(format_location(fields.state)))`,
				`Yes`,
				`No`,
				`m`,
				`Jeff Jefferson`,
				`@results.gender.category`,
				`@fields.raw_district`,
				`http://localhost/?cmd=success&name=@(url_encode(contact.name))`,
			},
		},
	}

	for _, tc := range testCases {
		flow, err := test.LoadFlowFromAssets(tc.path, assets.FlowUUID(tc.uuid))
		require.NoError(t, err)

		// try extracting all templates
		templates := flow.ExtractTemplates()
		assert.Equal(t, tc.templates, templates, "extracted templates mismatch for flow %s[uuid=%s]", tc.path, tc.uuid)

		// try rewriting all templates in uppercase
		flow.RewriteTemplates(func(t string) string { return strings.ToUpper(t) })

		// re-extract all templates
		rewritten := flow.ExtractTemplates()

		for i := range templates {
			templates[i] = strings.ToUpper(templates[i])
		}

		assert.Equal(t, templates, rewritten)
	}
}

func TestExtractDependencies(t *testing.T) {
	testCases := []struct {
		path         string
		uuid         string
		dependencies []assets.Reference
	}{
		{
			"../../test/testdata/runner/all_actions.json",
			"8ca44c09-791d-453a-9799-a70dd3303306",
			[]assets.Reference{
				assets.NewLabelReference("3f65d88a-95dc-4140-9451-943e94e06fea", "Spam"),
				assets.NewFieldReference("state", ""),
				assets.NewGroupReference("2aad21f6-30b7-42c5-bd7f-1b720c154817", "Survey Audience"),
				assets.NewFieldReference("activation_token", "Activation Token"),
				assets.NewFieldReference("first_name", ""),
				assets.NewFlowReference("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "Collect Language"),
				flows.NewContactReference("820f5923-3369-41c6-b3cd-af577c0bd4b8", "Bob"),
				assets.NewFieldReference("gender", "Gender"),
				assets.NewFieldReference("raw_district", ""),
				assets.NewFieldReference("district", "District"),
				assets.NewChannelReference("57f1078f-88aa-46f4-a59a-948a5739c03d", "Android Channel"),
			},
		},
		{
			"../../test/testdata/runner/router_tests.json",
			"615b8a0f-588c-4d20-a05f-363b0b4ce6f4",
			[]assets.Reference{
				assets.NewGroupReference("2aad21f6-30b7-42c5-bd7f-1b720c154817", ""),
				assets.NewGroupReference("bf282a79-aa74-4557-9932-22a9b3bce537", ""),
				assets.NewFieldReference("raw_district", ""),
				assets.NewFieldReference("district", "District"),
			},
		},
		{
			"../../test/testdata/runner/dynamic_groups.json",
			"1b462ce8-983a-4393-b133-e15a0efdb70c",
			[]assets.Reference{
				assets.NewFieldReference("gender", "Gender"),
				assets.NewFieldReference("age", "Age"),
			},
		},
	}

	for _, tc := range testCases {
		flow, err := test.LoadFlowFromAssets(tc.path, assets.FlowUUID(tc.uuid))
		require.NoError(t, err)

		// try extracting all dependencies
		assert.Equal(t, tc.dependencies, flow.ExtractDependencies(), "extracted dependencies mismatch for flow %s[uuid=%s]", tc.path, tc.uuid)
	}
}

func TestExtractResults(t *testing.T) {
	testCases := []struct {
		path    string
		uuid    string
		results []*flows.ResultSpec
	}{
		{
			"../../test/testdata/runner/all_actions.json",
			"8ca44c09-791d-453a-9799-a70dd3303306",
			[]*flows.ResultSpec{
				{Key: "gender", Name: "Gender", Categories: []string{"Male"}},
			},
		},
		{
			"../../test/testdata/runner/router_tests.json",
			"615b8a0f-588c-4d20-a05f-363b0b4ce6f4",
			[]*flows.ResultSpec{
				{Key: "urn_check", Name: "URN Check", Categories: []string{"Telegram", "Other"}},
				{Key: "group_check", Name: "Group Check", Categories: []string{"Testers", "Other"}},
				{Key: "district_check", Name: "District Check", Categories: []string{"Valid", "Invalid"}},
			},
		},
		{
			"../../test/testdata/runner/two_questions.json",
			"615b8a0f-588c-4d20-a05f-363b0b4ce6f4",
			[]*flows.ResultSpec{
				{Key: "favorite_color", Name: "Favorite Color", Categories: []string{"Red", "Blue", "Other", "No Response"}},
				{Key: "soda", Name: "Soda", Categories: []string{"Pepsi", "Coke", "Other"}},
			},
		},
	}

	for _, tc := range testCases {
		flow, err := test.LoadFlowFromAssets(tc.path, assets.FlowUUID(tc.uuid))
		require.NoError(t, err)

		// try extracting all dependencies
		assert.Equal(t, tc.results, flow.ExtractResults(), "extracted results mismatch for flow %s[uuid=%s]", tc.path, tc.uuid)
	}
}
