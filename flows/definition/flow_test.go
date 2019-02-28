package definition_test

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/routers"
	"github.com/nyaruka/goflow/flows/waits"
	"github.com/nyaruka/goflow/flows/waits/hints"
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
}

func TestFlowValidation(t *testing.T) {
	session, _, err := test.CreateTestSession("", nil)
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

func TestNewFlow(t *testing.T) {
	var flowDef = `{
		"uuid": "8ca44c09-791d-453a-9799-a70dd3303306", 
		"name": "Test Flow",
		"spec_version": "12.0.0",
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
					"type": "msg",
					"hint": {
						"type": "image"
					}
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

	session, _, err := test.CreateTestSession("", nil)
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
				waits.NewMsgWait(nil, hints.NewImageHint()),
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

func TestReadFlow(t *testing.T) {
	// try reading something without a flow header
	_, err := definition.ReadFlow([]byte(`{"nodes":[]}`))
	assert.EqualError(t, err, "unable to read flow header: field 'uuid' is required, field 'name' is required, field 'spec_version' is required")

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
	assert.EqualError(t, err, "spec version 2000.0.0 is newer than this library (12.0.0)")

	// try reading a definition with a newer minor version
	_, err = definition.ReadFlow([]byte(`{
		"uuid": "8ca44c09-791d-453a-9799-a70dd3303306", 
		"name": "Test Flow",
		"spec_version": "12.9999",
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
		"spec_version": "12.0",
		"language": "eng",
		"revision": 123,
		"expire_after_minutes": 30,
		"nodes": []
	}`))
	assert.EqualError(t, err, "unable to read flow: field 'type' is required")
}

func TestExtractAndRewriteTemplates(t *testing.T) {
	testCases := []struct {
		path      string
		uuid      string
		templates []string
	}{
		{
			"../../test/testdata/flows/two_questions.json",
			"615b8a0f-588c-4d20-a05f-363b0b4ce6f4",
			[]string{
				`Hi @contact.name! What is your favorite color? (red/blue) Your number is @(format_urn(contact.urn))`,
				`Red`,
				`Blue`,
				`Quelle est votres couleur preferee? (rouge/blue)`,
				`@input`,
				`red`,
				`rouge`,
				`blue`,
				`bleu`,
				`fra`,
				`@(TITLE(results.favorite_color.category_localized)) it is! What is your favorite soda? (pepsi/coke)`,
				`@(TITLE(results.favorite_color.category_localized))! Bien sur! Quelle est votes soda preferee? (pepsi/coke)`,
				`@input`,
				`pepsi`,
				`coke coca cola`,
				`http://localhost/?cmd=success`,
				`{ "contact": @(json(contact.uuid)), "soda": @(json(results.soda.value)) }`,
				`Great, you are done and like @results.soda! Webhook status was @results.webhook.value`,
				`Parfait, vous avez finis et tu aimes @results.soda.category`,
			},
		},
		{
			"../../test/testdata/flows/all_actions.json",
			"8ca44c09-791d-453a-9799-a70dd3303306",
			[]string{
				`@(format_location(contact.fields.state)) Messages`,
				`@(format_location(contact.fields.state)) Members`,
				`@(replace(lower(contact.name), " ", "_"))`,
				`XXX-YYY-ZZZ`,
				"Here is your activation token",
				"Hi @contact.fields.first_name, Your activation token is @contact.fields.activation_token, your coupon is @(trigger.params.coupons[0].code)",
				"@(contact.urns.mailto[0])",
				"test@@example.com",
				`Hi @contact.name, are you ready?`,
				`Hola @contact.name, ¿estás listo?`,
				`Hi @contact.name, are you ready for these attachments?`,
				`image/jpeg:http://s3.amazon.com/bucket/test_en.jpg?a=@(url_encode(format_location(contact.fields.state)))`,
				`Hi @contact.name, are you ready to complete today's survey?`,
				`This is a message to each of @contact.name's urns.`,
				`This is a reply with attachments and quick replies`,
				`image/jpeg:http://s3.amazon.com/bucket/test_en.jpg?a=@(url_encode(format_location(contact.fields.state)))`,
				`Yes`,
				`No`,
				`m`,
				`Jeff Jefferson`,
				`@results.gender.category`,
				`@contact.fields.raw_district`,
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

		for t := range templates {
			templates[t] = strings.ToUpper(templates[t])
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
			"../../test/testdata/flows/all_actions.json",
			"8ca44c09-791d-453a-9799-a70dd3303306",
			[]assets.Reference{
				assets.NewFieldReference("state", ""),
				assets.NewLabelReference("3f65d88a-95dc-4140-9451-943e94e06fea", "Spam"),
				assets.NewGroupReference("2aad21f6-30b7-42c5-bd7f-1b720c154817", "Survey Audience"),
				assets.NewFieldReference("activation_token", "Activation Token"),
				assets.NewFieldReference("first_name", ""),
				assets.NewFlowReference("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "Collect Language"),
				assets.NewFieldReference("gender", "Gender"),
				assets.NewFieldReference("raw_district", ""),
				assets.NewFieldReference("district", "District"),
			},
		},
		{
			"../../test/testdata/flows/router_tests.json",
			"615b8a0f-588c-4d20-a05f-363b0b4ce6f4",
			[]assets.Reference{
				assets.NewGroupReference("ade39253-0371-4dde-9df4-95245fe6a3a8", ""),
				assets.NewFieldReference("raw_district", ""),
				assets.NewFieldReference("district", "District"),
			},
		},
		{
			"../../test/testdata/flows/dynamic_groups.json",
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
