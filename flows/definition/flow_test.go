package definition_test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/definition/migrations"
	"github.com/nyaruka/goflow/flows/routers"
	"github.com/nyaruka/goflow/flows/routers/waits"
	"github.com/nyaruka/goflow/flows/routers/waits/hints"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/jsonx"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsVersionSupported(t *testing.T) {
	assert.False(t, definition.IsVersionSupported(semver.MustParse("10.5")))
	assert.True(t, definition.IsVersionSupported(semver.MustParse("11.0")))
	assert.True(t, definition.IsVersionSupported(semver.MustParse("11.9")))
	assert.True(t, definition.IsVersionSupported(semver.MustParse("13.0.0")))
	assert.True(t, definition.IsVersionSupported(semver.MustParse("13.3.0")))
	assert.False(t, definition.IsVersionSupported(semver.MustParse("14.0.0")))
}

func TestBrokenFlows(t *testing.T) {
	testCases := []struct {
		path string
		err  string
	}{
		{
			"exitless_node.json",
			"unable to read node: field 'exits' must have a minimum of 1 items",
		},
		{
			"exitless_category.json",
			"unable to read router: unable to read category: field 'exit_uuid' is required",
		},
		{
			"duplicate_node_uuid.json",
			"node UUID a58be63b-907d-4a1a-856b-0bb5579d7507 isn't unique",
		},
		{
			"invalid_action_by_tag.json",
			"unable to read action: field 'text' is required",
		},
		{
			"invalid_action_by_method.json",
			"invalid node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: invalid action[uuid=e5a03dde-3b2f-4603-b5d0-d927f6bcc361, type=call_webhook]: header '\"$?' is not a valid HTTP header",
		},
		{
			"invalid_timeout_category.json",
			"invalid node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: invalid router: timeout category 13fea3d4-b925-495b-b593-1c9e905e700d is not a valid category",
		},
		{
			"invalid_default_exit.json",
			"invalid node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: invalid router: default category 37d8813f-1402-4ad2-9cc2-e9054a96525b is not a valid category",
		},
		{
			"invalid_case_category.json",
			"invalid node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: invalid router: case category 37d8813f-1402-4ad2-9cc2-e9054a96525b is not a valid category",
		},
		{
			"invalid_exit_dest.json",
			"invalid node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: destination 714f1409-486e-4e8e-bb08-23e2943ef9f6 of exit[uuid=37d8813f-1402-4ad2-9cc2-e9054a96525b] isn't a known node",
		},
	}

	for _, tc := range testCases {
		assetsJSON, err := ioutil.ReadFile("testdata/broken_flows/" + tc.path)
		require.NoError(t, err)

		sa, err := test.CreateSessionAssets(assetsJSON, "")
		require.NoError(t, err, "unable to load assets: %s", tc.path)

		_, err = sa.Flows().Get("76f0a02f-3b75-4b86-9064-e9195e1b3a02")

		if tc.err != "" {
			assert.EqualError(t, err, tc.err, "read error mismatch for %s", tc.path)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestNewFlow(t *testing.T) {
	var flowDef = fmt.Sprintf(`
{
    "uuid": "8ca44c09-791d-453a-9799-a70dd3303306", 
    "name": "Test Flow",
    "spec_version": "%s",
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
}`, definition.CurrentSpecVersion)

	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	flow, err := definition.NewFlow(
		assets.FlowUUID("8ca44c09-791d-453a-9799-a70dd3303306"),
		"Test Flow",          // name
		envs.Language("eng"), // base language
		flows.FlowTypeMessaging,
		123, // revision
		30,  // expires after minutes
		definition.NewLocalization(),
		[]flows.Node{
			definition.NewNode(
				flows.NodeUUID("a58be63b-907d-4a1a-856b-0bb5579d7507"),
				[]flows.Action{
					actions.NewSendMsg(
						flows.ActionUUID("76112ef2-790e-4b5b-84cb-e910f191a335"),
						"Do you like beer?",
						nil,
						nil,
						false,
					),
				},
				routers.NewSwitch(
					waits.NewMsgWait(nil, hints.NewImageHint()),
					"Response 1",
					[]flows.Category{
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
						routers.NewCase(uuids.UUID("9f593e22-7886-4c08-a52f-0e8780504d75"), "has_any_word", []string{"yes", "yeah"}, flows.CategoryUUID("97b9451c-2856-475b-af38-32af68100897")),
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
					actions.NewAddInputLabels(
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
	require.NoError(t, err)

	marshaled, err := jsonx.Marshal(flow)
	assert.NoError(t, err)

	test.AssertEqualJSON(t, []byte(flowDef), marshaled, "flow definition mismatch")

	// check in expressions
	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("Test Flow"),
		"name":        types.NewXText("Test Flow"),
		"revision":    types.NewXNumberFromInt(123),
		"uuid":        types.NewXText("8ca44c09-791d-453a-9799-a70dd3303306"),
	}), flows.Context(session.Environment(), flow))

	// check inspection
	info := flow.Inspect(session.Assets())
	infoJSON, _ := jsonx.Marshal(info)

	test.AssertEqualJSON(t, []byte(`{
		"dependencies": [
			{
				"key": "gender",
				"name": "",
				"type": "field"
			},
			{
				"uuid": "3f65d88a-95dc-4140-9451-943e94e06fea",
				"name": "Spam",
				"type": "label"
			}
		],
		"issues": [],
		"parent_refs": [],
		"results": [
			{
				"categories": [
					"Yes",
					"No"
				],
				"key": "response_1",
				"name": "Response 1",
				"node_uuids": [
					"a58be63b-907d-4a1a-856b-0bb5579d7507"
				]
			}
		],
		"waiting_exits": [
			"023a5c10-d74a-4fad-9560-990caead8170",
			"8943c032-2a91-456c-8080-2a249f1b420c"
		]
	}`), infoJSON, "inspection mismatch")
}

func TestEmptyFlow(t *testing.T) {
	env := envs.NewBuilder().Build()
	flow, err := test.LoadFlowFromAssets(env, "../../test/testdata/runner/empty.json", "76f0a02f-3b75-4b86-9064-e9195e1b3a02")
	require.NoError(t, err)

	marshaled, err := jsonx.Marshal(flow)
	require.NoError(t, err)

	expected := fmt.Sprintf(`{
		"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
		"name": "Empty Flow",
		"revision": 0,
		"spec_version": "%s",
		"type": "messaging",
		"expire_after_minutes": 0,
		"language": "eng",
		"localization": {},
		"nodes": []
  	}`, definition.CurrentSpecVersion)
	test.AssertEqualJSON(t, []byte(expected), marshaled, "flow definition mismatch")

	info := flow.Inspect(nil)
	infoJSON, _ := jsonx.Marshal(info)

	test.AssertEqualJSON(t, []byte(`{
		"dependencies": [],
		"issues": [],
		"parent_refs": [],
		"results": [],
		"waiting_exits": []
	}`), infoJSON, "inspection mismatch")
}

func TestReadFlow(t *testing.T) {
	// try reading something without a flow header
	_, err := definition.ReadFlow([]byte(`{"nodes":[]}`), nil)
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
	}`), nil)
	assert.EqualError(t, err, fmt.Sprintf("spec version 2000.0.0 is newer than this library (%s)", definition.CurrentSpecVersion))

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
	}`), nil)
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
  }`), nil)
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
	  }`), nil)
	assert.NoError(t, err)
	test.AssertEqualJSON(t, []byte(`{
		"nodes": {
			"b1c5f247-565d-4a7a-8763-c59abbed0a57": {
				"type": "execute_actions"
			}
		},
		"stickies": {}
	}`), flow.UI(), "ui mismatch for read flow")

	// try reading a legacy definition
	flow, err = definition.ReadFlow([]byte(`{
		"base_language": "eng",
		"entry": "10e483a8-5ffb-4c4f-917b-d43ce86c1d65", 
		"flow_type": "M",
		"action_sets": [{
			"uuid": "10e483a8-5ffb-4c4f-917b-d43ce86c1d65",
			"y": 100, 
			"x": 100, 
			"destination": null, 
			"exit_uuid": "cfcf5cef-49f9-41a6-886b-f466575a3045",
			"actions": []
		}],
		"metadata": {
			"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
			"name": "TestFlow"
		}
	}`), &migrations.Config{})
	assert.NoError(t, err)
	assert.Equal(t, assets.FlowUUID("50c3706e-fedb-42c0-8eab-dda3335714b7"), flow.UUID())
	assert.Equal(t, "TestFlow", flow.Name())
	assert.Equal(t, flows.FlowTypeMessaging, flow.Type())
	assert.Equal(t, 1, len(flow.Nodes()))
}

func TestExtractTemplatesAndLocalizables(t *testing.T) {
	env := envs.NewBuilder().Build()

	testCases := []struct {
		path         string
		uuid         string
		templates    []string
		localizables []string
	}{
		{
			"../../test/testdata/runner/two_questions.json",
			"615b8a0f-588c-4d20-a05f-363b0b4ce6f4",
			[]string{
				`Hi @contact.name! What is your favorite color? (red/blue) Your number is @(format_urn(contact.urn))`,
				`Quelle est votres couleur preferee? (rouge/blue)`,
				`Red`,
				`Blue`,
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
			[]string{
				"Hi @contact.name! What is your favorite color? (red/blue) Your number is @(format_urn(contact.urn))",
				"Red",
				"Blue",
				"red",
				"blue",
				"Red",
				"Blue",
				"Other",
				"No Response",
				"@(TITLE(results.favorite_color.category_localized)) it is! What is your favorite soda? (pepsi/coke)",
				"pepsi",
				"coke coca cola",
				"Pepsi",
				"Coke",
				"Other",
				"Great, you are done and like @results.soda.value! Webhook status was @results.webhook.value",
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
				"@urns.mailto",
				"test@@example.com",
				"Here is your activation token",
				"Hi @fields.first_name, Your activation token is @fields.activation_token, your coupon is @(trigger.params.coupons[0].code)",
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
			[]string{
				"Here is your activation token",
				"Hi @fields.first_name, Your activation token is @fields.activation_token, your coupon is @(trigger.params.coupons[0].code)",
				"Hi @contact.name, are you ready?",
				"Hi @contact.name, are you ready for these attachments?",
				"image/jpeg:http://s3.amazon.com/bucket/test_en.jpg?a=@(url_encode(format_location(fields.state)))",
				"Hi @contact.name, are you ready to complete today's survey?",
				"This is a message to each of @contact.name's urns.",
				"This is a reply with attachments and quick replies",
				"image/jpeg:http://s3.amazon.com/bucket/test_en.jpg?a=@(url_encode(format_location(fields.state)))",
				"Yes",
				"No",
				"Male",
			},
		},
	}

	for _, tc := range testCases {
		flow, err := test.LoadFlowFromAssets(env, tc.path, assets.FlowUUID(tc.uuid))
		require.NoError(t, err)

		// try extracting all templates
		templates := flow.ExtractTemplates()
		assert.Equal(t, tc.templates, templates, "extracted templates mismatch for flow %s[uuid=%s]", tc.path, tc.uuid)

		// try extracting all localizable text
		localizables := flow.ExtractLocalizables()
		assert.Equal(t, tc.localizables, localizables, "extracted localizables mismatch for flow %s[uuid=%s]", tc.path, tc.uuid)
	}
}

func TestInspection(t *testing.T) {
	env := envs.NewBuilder().Build()

	testCases := []struct {
		path string
		uuid string
	}{
		{
			"../../test/testdata/runner/all_actions.json",
			"8ca44c09-791d-453a-9799-a70dd3303306",
		},
		{
			"../../test/testdata/runner/router_tests.json",
			"615b8a0f-588c-4d20-a05f-363b0b4ce6f4",
		},
		{
			"../../test/testdata/runner/dynamic_groups.json",
			"1b462ce8-983a-4393-b133-e15a0efdb70c",
		},
		{
			"../../test/testdata/runner/two_questions.json",
			"615b8a0f-588c-4d20-a05f-363b0b4ce6f4",
		},
		{
			"../../test/testdata/runner/triggered.json",
			"ce902e6f-bc0a-40cf-a58c-1e300d15ec85",
		},
		{
			"../../test/testdata/runner/missing_dependencies.json",
			"447efb41-c1e2-44f9-b906-4ed6b5031e59",
		},
	}

	for _, tc := range testCases {
		sa, err := test.LoadSessionAssets(env, tc.path)
		require.NoError(t, err)

		flow, err := sa.Flows().Get(assets.FlowUUID(tc.uuid))
		require.NoError(t, err)

		actualInfo := flow.Inspect(sa)
		actualJSON, _ := jsonx.MarshalPretty(actualInfo)

		testDataPath := "testdata/inspection/" + tc.path[strings.LastIndex(tc.path, "/"):]

		if !test.UpdateSnapshots {
			expectedJSON, err := ioutil.ReadFile(testDataPath)
			require.NoError(t, err)
			test.AssertEqualJSON(t, expectedJSON, actualJSON, "inspection mismatch for flow %s[uuid=%s]", tc.path, tc.uuid)
		} else {
			err := ioutil.WriteFile(testDataPath, actualJSON, 0666)
			require.NoError(t, err)
		}
	}
}

func TestChangeLanguage(t *testing.T) {
	env := envs.NewBuilder().Build()

	flow, err := test.LoadFlowFromAssets(env, "testdata/change_language.json", "19cad1f2-9110-4271-98d4-1b968bf19410")
	require.NoError(t, err)

	assertLanguageChange := func(lang envs.Language) {
		copy, err := flow.ChangeLanguage(lang)
		assert.NoError(t, err)

		marshaled, err := jsonx.MarshalPretty(copy)
		require.NoError(t, err)
		test.AssertSnapshot(t, "change_language_to_"+string(lang), string(marshaled))

		// check flow is valid by reading it back
		_, err = definition.ReadFlow(marshaled, nil)
		assert.NoError(t, err)
	}

	assertLanguageChange("spa") // has a complete translation
	assertLanguageChange("ara") // missing translations will be left in eng
	assertLanguageChange("kin") // everything is missing and will be left in eng
}
