package engine_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sessionAssets = `{
    "channels": [
        {
            "uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d",
            "name": "My Android Phone",
            "address": "+17036975131",
            "schemes": [
                "tel"
            ],
            "roles": [
                "send",
                "receive"
            ],
            "country": "US"
        }
    ],
    "fields": [
        {
			"uuid": "d66a7823-eada-40e5-9a3a-57239d4690bf",
			"key": "gender",
            "name": "Gender",
            "type": "text"
        }
    ],
    "flows": [
        {
            "uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
            "name": "No Related Runs",
            "spec_version": "13.0.0",
            "language": "eng",
            "type": "messaging",
            "revision": 123,
            "nodes": [
                {
                    "uuid": "3dcccbb4-d29c-41dd-a01f-16d814c9ab82",
                    "router": {
                        "type": "switch",
                        "categories": [
                            {
                                "uuid": "d7342563-7c9d-4576-b6d1-0c1f148765d2",
                                "name": "All Responses",
                                "exit_uuid": "37d8813f-1402-4ad2-9cc2-e9054a96525b"
                            }
                        ],
                        "operand": "@input.text",
                        "default_category_uuid": "d7342563-7c9d-4576-b6d1-0c1f148765d2"
                    },
                    "exits": [
                        {
                            "uuid": "37d8813f-1402-4ad2-9cc2-e9054a96525b",
                            "destination_uuid": null
                        }
                    ]
                }
            ]
		}
	]
}`

var sessionContact = `{
	"uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f",
	"id": 1234567,
	"name": "Ryan Lewis",
	"status": "active",
	"language": "eng",
	"timezone": "America/Guayaquil",
	"created_on": "2018-06-20T11:40:30.123456789-00:00",
	"urns": [ "tel:+12065551212"],
	"fields": {
		"gender": {"text": "M"}
	}
}`

var sessionTrigger = `{
    "type": "manual",
    "triggered_on": "2017-12-31T11:31:15.035757258-02:00",
    "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "No Related Runs"}
}`

func TestRuns(t *testing.T) {
	test.MockUniverse()

	server := test.NewTestHTTPServer(49999)
	defer server.Close()

	session, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyNone)
	require.NoError(t, err)

	flow, err := session.Assets().Flows().Get("50c3706e-fedb-42c0-8eab-dda3335714b7")
	require.NoError(t, err)

	checkRuns := func(s flows.Session) {
		r1, r2 := s.Runs()[0], s.Runs()[1]

		assert.Equal(t, flows.RunUUID("01969b47-113b-76f8-9c0b-2014ddc77094"), r1.UUID())
		assert.Equal(t, flows.RunStatusCompleted, r1.Status())
		assert.Equal(t, flow, r1.Flow())
		assert.Equal(t, flow.Reference(true), r1.FlowReference())
		assert.Equal(t, "Parent", r1.Parent().Flow().Name())
		assert.Equal(t, 0, len(r1.Ancestors())) // no parent runs within this session
		assert.True(t, r1.HadInput())

		assert.Equal(t, flows.RunUUID("01969b47-24c3-76f8-8f41-6b2d9f33d623"), r2.UUID())
		assert.Equal(t, flows.RunUUID("01969b47-113b-76f8-9c0b-2014ddc77094"), r2.Parent().UUID())
	}

	checkRuns(session)

	// check we can marshal and marshal the run and get the same values
	sessionJSON, err := jsonx.Marshal(session)
	require.NoError(t, err)

	session2, err := session.Engine().ReadSession(session.Assets(), sessionJSON, session.Environment(), session.Contact(), nil, assets.IgnoreMissing)
	require.NoError(t, err)

	// needed so that prepareForSprint is invoked that parses the parent run summary from the trigger
	session2.Resume(t.Context(), resumes.NewWaitTimeout(events.NewWaitTimedOut()))

	checkRuns(session2)
}

func TestRunContext(t *testing.T) {
	test.MockUniverse()

	// create a run with no parent or child
	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	run := session.Runs()[0]

	testCases := []struct {
		template string
		expected string
	}{
		{`@run`, `Ryan Lewis@Registration`},
		{`@child`, `Ryan Lewis@Collect Age`},
		{`@child.uuid`, `01969b47-24c3-76f8-8f41-6b2d9f33d623`},
		{`@child.run`, `{status: completed}`}, // to be removed in 13.2
		{`@child.contact.name`, `Ryan Lewis`},
		{`@child.flow.name`, "Collect Age"},
		{`@child.status`, "completed"},
		{`@child.fields`, "Activation Token: AACC55\nAge: 23\nGender: Male\nJoin Date: 2017-12-02T00:00:00.000000-02:00"},
		{`@parent`, `Jasmine@Parent`},
		{`@parent.uuid`, `4213ac47-93fd-48c4-af12-7da8218ef09d`},
		{`@parent.run`, `{status: active}`},
		{`@parent.contact.name`, `Jasmine`},
		{`@parent.flow.name`, "Parent"},
		{`@parent.status`, "active"},
		{`@parent.fields`, "Age: 33\nGender: Female"},
		{`@node.uuid`, "c0781400-737f-4940-9a6c-1ec1c3df0325"},
		{`@node.visit_count`, "1"},
		{`@trigger.type`, "flow_action"},
		{`@resume.type`, "msg"},
		{
			`@(json(contact.fields))`,
			`{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00","language":null,"not_set":null,"state":null}`,
		},
		{
			`@(json(fields))`,
			`{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00","language":null,"not_set":null,"state":null}`,
		},
		{
			`@(json(contact.urns))`,
			`["tel:+12024561111","twitterid:54784326227#nyaruka","mailto:foo@bar.com"]`,
		},
		{
			`@(json(urns))`,
			`{"discord":null,"ext":null,"facebook":null,"fcm":null,"freshchat":null,"instagram":null,"jiochat":null,"line":null,"mailto":"mailto:foo@bar.com","rocketchat":null,"slack":null,"tel":"tel:+12024561111","telegram":null,"twitter":null,"twitterid":"twitterid:54784326227#nyaruka","viber":null,"vk":null,"webchat":null,"wechat":null,"whatsapp":null}`,
		},
		{
			`@(json(results.favorite_color))`,
			`{"category":"Red","category_localized":"Red","created_on":"2025-05-04T12:31:17.123456Z","extra":null,"input":"","name":"Favorite Color","node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03","value":"red"}`,
		},
		{
			`@(json(run.results.favorite_color))`,
			`{"category":"Red","category_localized":"Red","created_on":"2025-05-04T12:31:17.123456Z","extra":null,"input":"","name":"Favorite Color","node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03","value":"red"}`,
		},
		{
			`@(json(parent.contact.urns))`,
			`["tel:+12024562222"]`,
		},
		{
			`@(json(parent.urns))`,
			`{"discord":null,"ext":null,"facebook":null,"fcm":null,"freshchat":null,"instagram":null,"jiochat":null,"line":null,"mailto":null,"rocketchat":null,"slack":null,"tel":"tel:+12024562222","telegram":null,"twitter":null,"twitterid":null,"viber":null,"vk":null,"webchat":null,"wechat":null,"whatsapp":null}`,
		},
		{
			`@(json(parent.fields))`,
			`{"activation_token":null,"age":33,"gender":"Female","join_date":null,"language":null,"not_set":null,"state":null}`,
		},
	}

	for _, tc := range testCases {
		log := test.NewEventLog()
		actual, _ := run.EvaluateTemplate(tc.template, log.Log)
		assert.NoError(t, log.Error())
		assert.Equal(t, tc.expected, actual, "template mismatch for %s", tc.template)
	}

	// test with escaping
	log := test.NewEventLog()
	evaluated, _ := run.EvaluateTemplateText(`gender = @("M\" OR")`, flows.ContactQueryEscaping, true, log.Log)
	assert.NoError(t, log.Error())
	assert.Equal(t, `gender = "M\" OR"`, evaluated)
}

func TestMissingRelatedRunContext(t *testing.T) {
	// create a run with no parent or child
	sa, err := test.CreateSessionAssets([]byte(sessionAssets), "")
	require.NoError(t, err)

	contact, err := flows.ReadContact(sa, []byte(sessionContact), assets.IgnoreMissing)
	require.NoError(t, err)

	trigger, err := triggers.Read(sa, []byte(sessionTrigger), assets.IgnoreMissing)
	require.NoError(t, err)

	tz, _ := time.LoadLocation("America/Guayaquil")
	env := envs.NewBuilder().WithAllowedLanguages("eng", "spa").WithTimezone(tz).Build()

	eng := test.NewEngine()
	session, _, err := eng.NewSession(context.Background(), sa, env, contact, trigger, nil)
	require.NoError(t, err)

	run := session.Runs()[0]
	log := test.NewEventLog()

	// since we have no parent, check that it resolves to nil
	val, _ := run.EvaluateTemplateValue(`@parent`, log.Log)
	assert.NoError(t, log.Error())
	assert.Nil(t, val)

	// check that trying to resolve a property of parent is an error
	val, _ = run.EvaluateTemplateValue(`@parent.contact`, log.Log)
	assert.NoError(t, log.Error())
	assert.Equal(t, types.NewXErrorf("null doesn't support lookups"), val)

	// we also have no child, check that it resolves to nil
	val, _ = run.EvaluateTemplateValue(`@child`, log.Log)
	assert.NoError(t, log.Error())
	assert.Nil(t, val)

	// check that trying to resolve a property of child is an error
	val, _ = run.EvaluateTemplateValue(`@child.contact`, log.Log)
	assert.NoError(t, log.Error())
	assert.Equal(t, types.NewXErrorf("null doesn't support lookups"), val)
}

func TestSetResult(t *testing.T) {
	sa, err := test.CreateSessionAssets([]byte(sessionAssets), "")
	require.NoError(t, err)

	contact, err := flows.ReadContact(sa, []byte(sessionContact), assets.IgnoreMissing)
	require.NoError(t, err)

	trigger, err := triggers.Read(sa, []byte(sessionTrigger), assets.IgnoreMissing)
	require.NoError(t, err)

	tz, _ := time.LoadLocation("America/Guayaquil")
	env := envs.NewBuilder().WithAllowedLanguages("eng", "spa").WithTimezone(tz).Build()

	eng := test.NewEngine()
	session, _, err := eng.NewSession(context.Background(), sa, env, contact, trigger, nil)
	require.NoError(t, err)

	run := session.Runs()[0]

	dates.SetNowFunc(dates.NewFixedNow(time.Date(2020, 4, 20, 12, 39, 30, 123456789, time.UTC)))
	defer dates.SetNowFunc(time.Now)

	// no results means empty object with default of empty string
	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{"__default__": types.XTextEmpty}), flows.Context(session.Environment(), run.Results()))

	prev, changed := run.SetResult(flows.NewResult("Response 1", "red", "Red", "Rojo", "6d35528e-cae3-4e30-b842-8fe6ed7d5c02", "I like red", nil, dates.Now()))
	assert.Nil(t, prev)
	assert.True(t, changed)

	// name is snaked
	assert.Equal(t, "red", run.Results().Get("response_1").Value)
	assert.Equal(t, "Red", run.Results().Get("response_1").Category)
	assert.Equal(t, time.Date(2020, 4, 20, 12, 39, 30, 123456789, time.UTC), run.ModifiedOn())

	prev, changed = run.SetResult(flows.NewResult("Response 1", "blue", "Blue", "Azul", "6d35528e-cae3-4e30-b842-8fe6ed7d5c02", "I like blue", nil, dates.Now()))
	if assert.NotNil(t, prev) {
		assert.Equal(t, "Red", prev.Category)
	}
	assert.True(t, changed)

	// result is overwritten
	assert.Equal(t, "blue", run.Results().Get("response_1").Value)
	assert.Equal(t, "Blue", run.Results().Get("response_1").Category)
	assert.Equal(t, time.Date(2020, 4, 20, 12, 39, 30, 123456789, time.UTC), run.ModifiedOn())

	// try saving new result with same value and category again
	prev, changed = run.SetResult(flows.NewResult("Response 1", "blue", "Blue", "Azul", "6f53c6ae-b66e-44dc-af9e-638e26ad05e9", "blue", nil, dates.Now()))
	assert.Nil(t, prev)
	assert.False(t, changed)

	// long values should truncated
	prev, changed = run.SetResult(flows.NewResult("Response 1", strings.Repeat("創", 700), "Blue", "Azul", "6d35528e-cae3-4e30-b842-8fe6ed7d5c02", "I like blue", nil, dates.Now()))
	assert.NotNil(t, prev)
	assert.True(t, changed)

	assert.Equal(t, strings.Repeat("創", 640), run.Results().Get("response_1").Value)
}

func TestTranslation(t *testing.T) {
	msgAction1 := []byte(`{
		"uuid": "0a8467eb-911a-41db-8101-ccf415c48e6a",
		"type": "send_msg",
		"text": "Hello",
		"attachments": [
			"image/jpeg:http://media.com/hello.jpg",
			"audio/mp4:http://media.com/hello.m4a"
		],
		"quick_replies": [
			"yes",
			"no"
		]
	}`)
	msgAction2 := []byte(`{
		"uuid": "0a8467eb-911a-41db-8101-ccf415c48e6a",
		"type": "send_msg",
		"text": "Hello"
	}`)

	tcs := []struct {
		description          string
		envLangs             []i18n.Language
		contactLang          i18n.Language
		msgAction            []byte
		expectedText         string
		expectedAttachments  []utils.Attachment
		expectedQuickReplies []flows.QuickReply
	}{
		{
			description:  "contact language is valid and is flow base language, msg action has all fields",
			envLangs:     []i18n.Language{"eng", "spa"},
			contactLang:  "eng",
			msgAction:    msgAction1,
			expectedText: "Hello",
			expectedAttachments: []utils.Attachment{
				"image/jpeg:http://media.com/hello.jpg",
				"audio/mp4:http://media.com/hello.m4a",
			},
			expectedQuickReplies: []flows.QuickReply{{Text: "yes"}, {Text: "no"}},
		},
		{
			description:  "contact language is valid and translations exist, msg action has all fields",
			envLangs:     []i18n.Language{"eng", "spa"},
			contactLang:  "spa",
			msgAction:    msgAction1,
			expectedText: "Hola",
			expectedAttachments: []utils.Attachment{
				"audio/mp4:http://media.com/hola.m4a",
			},
			expectedQuickReplies: []flows.QuickReply{{Text: "si"}},
		},
		{
			description:  "contact language is allowed but no translations exist, msg action has all fields",
			envLangs:     []i18n.Language{"eng", "spa", "kin"},
			contactLang:  "kin",
			msgAction:    msgAction1,
			expectedText: "Hello",
			expectedAttachments: []utils.Attachment{
				"image/jpeg:http://media.com/hello.jpg",
				"audio/mp4:http://media.com/hello.m4a",
			},
			expectedQuickReplies: []flows.QuickReply{{Text: "yes"}, {Text: "no"}},
		},
		{
			description:  "contact language is not allowed and translations exist, msg action has all fields",
			envLangs:     []i18n.Language{"eng"},
			contactLang:  "spa",
			msgAction:    msgAction1,
			expectedText: "Hello",
			expectedAttachments: []utils.Attachment{
				"image/jpeg:http://media.com/hello.jpg",
				"audio/mp4:http://media.com/hello.m4a",
			},
			expectedQuickReplies: []flows.QuickReply{{Text: "yes"}, {Text: "no"}},
		},
		{
			description:          "contact language is valid and is flow base language, msg action only has text",
			envLangs:             []i18n.Language{"eng", "spa"},
			contactLang:          "eng",
			msgAction:            msgAction2,
			expectedText:         "Hello",
			expectedAttachments:  []utils.Attachment{},
			expectedQuickReplies: []flows.QuickReply{},
		},
		{
			description:  "contact language is valid and translations exist, msg action only has text",
			envLangs:     []i18n.Language{"eng", "spa"},
			contactLang:  "spa",
			msgAction:    msgAction2,
			expectedText: "Hola",
			expectedAttachments: []utils.Attachment{
				"audio/mp4:http://media.com/hola.m4a",
			},
			expectedQuickReplies: []flows.QuickReply{{Text: "si"}},
		},
		{
			description:  "attachments and quick replies translations are single empty strings and should be ignored",
			envLangs:     []i18n.Language{"eng", "fra"},
			contactLang:  "fra",
			msgAction:    msgAction1,
			expectedText: "Bonjour",
			expectedAttachments: []utils.Attachment{
				"image/jpeg:http://media.com/hello.jpg",
				"audio/mp4:http://media.com/hello.m4a",
			},
			expectedQuickReplies: []flows.QuickReply{{Text: "yes"}, {Text: "no"}},
		},
	}

	for _, tc := range tcs {
		assetsJSON, _ := os.ReadFile("testdata/translation_assets.json")
		assetsJSON = test.JSONReplace(assetsJSON, []string{"flows", "[0]", "nodes", "[0]", "actions", "[0]"}, tc.msgAction)

		env := envs.NewBuilder().WithAllowedLanguages(tc.envLangs...).Build()
		_, _, sp := test.NewSessionBuilder().
			WithEnvironment(env).
			WithContact("2efa1803-ae4d-4a58-ba54-b523e53e40f3", 123, "Bob", tc.contactLang, "tel+1234567890").
			WithAssetsJSON(assetsJSON).
			MustBuild()

		require.Len(t, sp.Events(), 1)
		require.Equal(t, "msg_created", sp.Events()[0].Type())
		evt := sp.Events()[0].(*events.MsgCreated)

		assert.Equal(t, tc.expectedText, evt.Msg.Text(), "msg text mismatch in test '%s'", tc.description)
		assert.Equal(t, tc.expectedAttachments, evt.Msg.Attachments(), "attachments mismatch in test case '%s'", tc.description)
		assert.Equal(t, tc.expectedQuickReplies, evt.Msg.QuickReplies(), "quick replies mismatch in test case '%s'", tc.description)
	}
}

const legacySessionAssets = `{
	"flows": [
		{
			"uuid": "92be2b4b-4cdc-413b-b516-7f0fa1dda0df",
			"name": "Cat Fact Loop",
			"spec_version": "14.3.0",
			"type": "messaging",
			"expire_after_minutes": 5,
			"language": "eng",
			"localization": {},
			"revision": 181,
			"nodes": [
				{
					"uuid": "2cc73016-dd92-4b50-a1bc-64ce91577696",
					"actions": [
						{
							"attachments": [],
							"text": "Hi there!",
							"type": "send_msg",
							"quick_replies": [],
							"uuid": "bd298cfa-44cd-4ec3-a339-02fa81d12fff"
						}
					],
					"exits": [
						{
							"uuid": "4ccb9092-727d-4436-942c-e4190c8c89ee",
							"destination_uuid": "ae4ca3e7-c918-44d9-b6e2-1a53f03f89af"
						}
					]
				},
				{
					"uuid": "ae4ca3e7-c918-44d9-b6e2-1a53f03f89af",
					"actions": [
						{
							"uuid": "b4cecf32-cabf-4493-b2c2-b7b735959bd0",
							"headers": {
								"Accept": "application/json"
							},
							"type": "call_webhook",
							"url": "https://catfact.ninja/fact?visit=@node.visit_count",
							"body": "",
							"method": "GET"
						}
					],
					"router": {
						"type": "switch",
						"operand": "@webhook.status",
						"cases": [
							{
								"uuid": "bba135dd-bb3b-4cf4-bf75-93145951884a",
								"type": "has_number_between",
								"arguments": [
									"200",
									"299"
								],
								"category_uuid": "02f8097c-5427-4aaf-bf93-96b6aee8347c"
							}
						],
						"categories": [
							{
								"uuid": "02f8097c-5427-4aaf-bf93-96b6aee8347c",
								"name": "Success",
								"exit_uuid": "d91a77f1-bc03-4b76-b128-d777a13afa0c"
							},
							{
								"uuid": "f32b4e23-57fe-44d8-b087-5d0e98fda22a",
								"name": "Failure",
								"exit_uuid": "95eeebf1-779c-49ca-bfbd-5fecce14ae6b"
							}
						],
						"default_category_uuid": "f32b4e23-57fe-44d8-b087-5d0e98fda22a",
						"result_name": ""
					},
					"exits": [
						{
							"uuid": "d91a77f1-bc03-4b76-b128-d777a13afa0c",
							"destination_uuid": "cf1ca532-fa73-4ebb-bd2e-055eefd19b84"
						},
						{
							"uuid": "95eeebf1-779c-49ca-bfbd-5fecce14ae6b",
							"destination_uuid": null
						}
					]
				},
				{
					"uuid": "cf1ca532-fa73-4ebb-bd2e-055eefd19b84",
					"actions": [
						{
							"attachments": [],
							"text": "@webhook.json.fact. Want another?",
							"type": "send_msg",
							"quick_replies": [],
							"uuid": "454a3ee2-e559-4306-b32c-b56fd1eae636"
						}
					],
					"exits": [
						{
							"uuid": "579d9c19-e35f-4dc4-b5f7-16d74993067c",
							"destination_uuid": "07282700-46ab-4f2e-be88-56f12145017f"
						}
					]
				},
				{
					"uuid": "07282700-46ab-4f2e-be88-56f12145017f",
					"actions": [],
					"router": {
						"type": "switch",
						"default_category_uuid": "9ed246aa-3687-4d66-87f9-57dbe1f272f7",
						"cases": [
							{
								"arguments": [
									"yes"
								],
								"type": "has_any_word",
								"uuid": "0df7f8f6-a2c2-4106-93b1-401ab97149dd",
								"category_uuid": "b696aca3-bd90-4fdb-adf4-f8e2e99dfe93"
							}
						],
						"categories": [
							{
								"uuid": "b696aca3-bd90-4fdb-adf4-f8e2e99dfe93",
								"name": "Yes",
								"exit_uuid": "c1604427-ef00-4ccc-b67d-a964b61d3904"
							},
							{
								"uuid": "9ed246aa-3687-4d66-87f9-57dbe1f272f7",
								"name": "Other",
								"exit_uuid": "864b7e84-7d4a-4019-a6e4-4552021ca6b6"
							}
						],
						"operand": "@input.text",
						"wait": {
							"type": "msg"
						},
						"result_name": "Result 1"
					},
					"exits": [
						{
							"uuid": "c1604427-ef00-4ccc-b67d-a964b61d3904",
							"destination_uuid": "2cc73016-dd92-4b50-a1bc-64ce91577696"
						},
						{
							"uuid": "864b7e84-7d4a-4019-a6e4-4552021ca6b6",
							"destination_uuid": null
						}
					]
				}
			]
		}
	],
	"channels": [
		{
			"uuid": "bbda279d-2f7e-414e-8b04-78e5785acd9b",
			"name": "Test",
			"schemes": ["tel"],
			"roles": ["send", "receive"]
		}
	]
}`

const legacySessionWithRunEvents = `{
  "uuid": "0198be5c-fb22-7301-a936-d71dfa8fd6fd",
  "type": "messaging",
  "created_on": "2025-08-18T18:07:01.41019759Z",
  "trigger": {
    "type": "manual",
    "flow": {
      "uuid": "92be2b4b-4cdc-413b-b516-7f0fa1dda0df",
      "name": "Scratch"
    },
    "params": {},
    "triggered_on": "2025-08-18T18:07:01.409792877Z",
    "origin": "ui"
  },
  "contact_uuid": "a98b5f89-d5e8-4bda-9092-d31031379b38",
  "runs": [
    {
      "uuid": "0198be5c-fb22-738f-8cd2-d36bc351c876",
      "flow": {
        "uuid": "92be2b4b-4cdc-413b-b516-7f0fa1dda0df",
        "name": "Cat Fact Loop",
        "revision": 181
      },
      "path": [
        {
          "uuid": "b34e390b-e1ed-4a11-bc05-5f19d677d0f0",
          "node_uuid": "2cc73016-dd92-4b50-a1bc-64ce91577696",
          "exit_uuid": "4ccb9092-727d-4436-942c-e4190c8c89ee",
          "arrived_on": "2025-08-18T18:07:01.410240384Z"
        },
        {
          "uuid": "62b0c768-4eba-4786-ac62-e7f36585aecd",
          "node_uuid": "ae4ca3e7-c918-44d9-b6e2-1a53f03f89af",
          "exit_uuid": "d91a77f1-bc03-4b76-b128-d777a13afa0c",
          "arrived_on": "2025-08-18T18:07:01.410332112Z"
        },
        {
          "uuid": "592d45c0-bfa5-4504-a361-db8c191af30d",
          "node_uuid": "cf1ca532-fa73-4ebb-bd2e-055eefd19b84",
          "exit_uuid": "579d9c19-e35f-4dc4-b5f7-16d74993067c",
          "arrived_on": "2025-08-18T18:07:01.792095432Z"
        },
        {
          "uuid": "1bccf6bb-ba20-4de6-8983-83d9fa8e153b",
          "node_uuid": "07282700-46ab-4f2e-be88-56f12145017f",
          "exit_uuid": "c1604427-ef00-4ccc-b67d-a964b61d3904",
          "arrived_on": "2025-08-18T18:07:01.792216054Z"
        },
        {
          "uuid": "a1715024-66fc-460e-8a09-e7203956abd4",
          "node_uuid": "2cc73016-dd92-4b50-a1bc-64ce91577696",
          "exit_uuid": "4ccb9092-727d-4436-942c-e4190c8c89ee",
          "arrived_on": "2025-08-18T18:07:28.937648758Z"
        },
        {
          "uuid": "aa540932-0c61-49d5-8554-52ffb5d788b0",
          "node_uuid": "ae4ca3e7-c918-44d9-b6e2-1a53f03f89af",
          "exit_uuid": "d91a77f1-bc03-4b76-b128-d777a13afa0c",
          "arrived_on": "2025-08-18T18:07:28.937770744Z"
        },
        {
          "uuid": "7ae0e662-91be-4ab4-a2ef-8e9fd6329bec",
          "node_uuid": "cf1ca532-fa73-4ebb-bd2e-055eefd19b84",
          "exit_uuid": "579d9c19-e35f-4dc4-b5f7-16d74993067c",
          "arrived_on": "2025-08-18T18:07:29.128785141Z"
        },
        {
          "uuid": "330ccadd-1923-4bf9-97fb-482b09432a52",
          "node_uuid": "07282700-46ab-4f2e-be88-56f12145017f",
          "exit_uuid": "c1604427-ef00-4ccc-b67d-a964b61d3904",
          "arrived_on": "2025-08-18T18:07:29.129655511Z"
        },
        {
          "uuid": "c2b28222-cce3-4a2f-a696-5f8ddf6131cf",
          "node_uuid": "2cc73016-dd92-4b50-a1bc-64ce91577696",
          "exit_uuid": "4ccb9092-727d-4436-942c-e4190c8c89ee",
          "arrived_on": "2025-08-18T18:07:43.218035659Z"
        },
        {
          "uuid": "12332d6d-5b22-417e-bfab-a184f52cc7f5",
          "node_uuid": "ae4ca3e7-c918-44d9-b6e2-1a53f03f89af",
          "exit_uuid": "d91a77f1-bc03-4b76-b128-d777a13afa0c",
          "arrived_on": "2025-08-18T18:07:43.218127Z"
        },
        {
          "uuid": "a63bebbf-4327-4c1d-9372-2cd1de085c0d",
          "node_uuid": "cf1ca532-fa73-4ebb-bd2e-055eefd19b84",
          "exit_uuid": "579d9c19-e35f-4dc4-b5f7-16d74993067c",
          "arrived_on": "2025-08-18T18:07:43.421723692Z"
        },
        {
          "uuid": "11651d3c-5789-4cc1-a64b-72c0325dcc8b",
          "node_uuid": "07282700-46ab-4f2e-be88-56f12145017f",
          "arrived_on": "2025-08-18T18:07:43.421860998Z"
        }
      ],
      "events": [
        {
          "uuid": "0198be5c-fb22-7507-9f14-40a6d0e9f795",
          "type": "msg_created",
          "created_on": "2025-08-18T18:07:01.410329837Z",
          "step_uuid": "b34e390b-e1ed-4a11-bc05-5f19d677d0f0",
          "msg": {
            "urn": "telegram:742307595?channel=bbda279d-2f7e-414e-8b04-78e5785acd9b&id=600#rowanseymour",
            "channel": {
              "uuid": "bbda279d-2f7e-414e-8b04-78e5785acd9b",
              "name": "Staging Test"
            },
            "text": "Hi there!",
            "locale": "eng-AF"
          }
        },
        {
          "uuid": "0198be5c-fc9f-7df0-ba2c-6afe85e5cf4d",
          "type": "webhook_called",
          "created_on": "2025-08-18T18:07:01.791914195Z",
          "step_uuid": "62b0c768-4eba-4786-ac62-e7f36585aecd",
          "url": "https://catfact.ninja/fact?visit=1",
          "status_code": 200,
          "request": "GET /fact?visit=1 HTTP/1.1\r\nHost: catfact.ninja\r\nUser-Agent: RapidProMailroom/10.3.42\r\nAccept: application/json\r\nX-Mailroom-Mode: normal\r\nAccept-Encoding: gzip\r\n\r\n",
          "response": "HTTP/2.0 200 OK\r\nAccess-Control-Allow-Origin: *\r\nAlt-Svc: h3=\":443\"; ma=86400\r\nCache-Control: no-cache, private\r\nCf-Cache-Status: DYNAMIC\r\nCf-Ray: 971359333fd5803d-FRA\r\nContent-Type: application/json\r\nDate: Mon, 18 Aug 2025 18:07:01 GMT\r\nNel: {\"report_to\":\"cf-nel\",\"success_fraction\":0.0,\"max_age\":604800}\r\nReport-To: {\"group\":\"cf-nel\",\"max_age\":604800,\"endpoints\":[{\"url\":\"https://a.nel.cloudflare.com/report/v4?s=vUMS646LIul3nX6tFeIT4ov%2FihgA51lEK0apRwZh90DMZCp%2FEAt0MBIj97gXQ9TLIb6ElD332Z%2BIcKUIircDjo2%2BmSLGVVZpNmoj7Kk%3D\"}]}\r\nServer: cloudflare\r\nSet-Cookie: XSRF-TOKEN=%3D; SameSite=Lax; Secure; Path=/; Max-Age=7200; Expires=Mon, 18 Aug 2025 20:07:01 GMT\r\nSet-Cookie: catfacts_session=%3D; HttpOnly; SameSite=Lax; Secure; Path=/; Max-Age=7200; Expires=Mon, 18 Aug 2025 20:07:01 GMT\r\nSet-Cookie: Smo1Ok9CvZe9Vc3RDF9VbunMKL5pCQzKstyAcszn=; HttpOnly; SameSite=Lax; Secure; Path=/; Max-Age=7200; Expires=Mon, 18 Aug 2025 20:07:01 GMT\r\nX-Content-Type-Options: nosniff\r\nX-Frame-Options: SAMEORIGIN\r\nX-Ratelimit-Limit: 100\r\nX-Ratelimit-Remaining: 99\r\nX-Xss-Protection: 1; mode=block\r\n\r\n{\"fact\":\"All cats need taurine in their diet to avoid blindness. Cats must also have fat in their diet as they are unable to produce it on their own.\",\"length\":140}",
          "elapsed_ms": 381,
          "retries": 0,
          "status": "success"
        },
		{
          "uuid": "0198be5c-fc9f-7df0-cccc-6afe85e5cf4d",
          "type": "run_result_changed",
          "created_on": "2025-08-18T18:07:43.421452406Z",
          "step_uuid": "62b0c768-4eba-4786-ac62-e7f36585aecd",
		  "name": "Fact",
		  "value": "200",
          "extra": {"fact": "All cats need taurine in their diet to avoid blindness. Cats must also have fat in their diet as they are unable to produce it on their own.", "length": 140}
        },
        {
          "uuid": "0198be5c-fca0-7344-bbef-ef98a05f5084",
          "type": "msg_created",
          "created_on": "2025-08-18T18:07:01.792214431Z",
          "step_uuid": "592d45c0-bfa5-4504-a361-db8c191af30d",
          "msg": {
            "urn": "telegram:742307595?channel=bbda279d-2f7e-414e-8b04-78e5785acd9b&id=600#rowanseymour",
            "channel": {
              "uuid": "bbda279d-2f7e-414e-8b04-78e5785acd9b",
              "name": "Staging Test"
            },
            "text": "All cats need taurine in their diet to avoid blindness. Cats must also have fat in their diet as they are unable to produce it on their own.. Want another?",
            "locale": "eng-AF"
          }
        },
        {
          "uuid": "0198be5c-fca0-7357-be09-21a4dc1470e8",
          "type": "msg_wait",
          "created_on": "2025-08-18T18:07:01.792219294Z",
          "step_uuid": "1bccf6bb-ba20-4de6-8983-83d9fa8e153b",
          "expires_on": "2025-08-18T18:12:01.792217982Z"
        },
        {
          "uuid": "0198be5d-6586-7ea0-afab-a46b3b2afe55",
          "type": "msg_received",
          "created_on": "2025-08-18T18:07:28.801262658Z",
          "msg": {
            "urn": "telegram:742307595#rowanseymour",
            "channel": {
              "uuid": "bbda279d-2f7e-414e-8b04-78e5785acd9b",
              "name": "Staging Test"
            },
            "text": "yes",
            "external_id": "9125"
          }
        },
        {
          "uuid": "0198be5d-66a9-79dc-921b-7af273ebe797",
          "type": "run_result_changed",
          "created_on": "2025-08-18T18:07:28.937646736Z",
          "step_uuid": "1bccf6bb-ba20-4de6-8983-83d9fa8e153b",
          "name": "Result 1",
          "value": "yes",
          "category": "Yes"
        },
        {
          "uuid": "0198be5d-66a9-7ba6-907d-00e1159d0c2a",
          "type": "msg_created",
          "created_on": "2025-08-18T18:07:28.937763759Z",
          "step_uuid": "a1715024-66fc-460e-8a09-e7203956abd4",
          "msg": {
            "urn": "telegram:742307595?channel=bbda279d-2f7e-414e-8b04-78e5785acd9b&id=600#rowanseymour",
            "channel": {
              "uuid": "bbda279d-2f7e-414e-8b04-78e5785acd9b",
              "name": "Staging Test"
            },
            "text": "Hi there!",
            "locale": "eng-AF"
          }
        },
        {
          "uuid": "0198be5d-6767-7bdb-b564-cab6ba09bc73",
          "type": "webhook_called",
          "created_on": "2025-08-18T18:07:29.127777809Z",
          "step_uuid": "aa540932-0c61-49d5-8554-52ffb5d788b0",
          "url": "https://catfact.ninja/fact?visit=2",
          "status_code": 200,
          "request": "GET /fact?visit=2 HTTP/1.1\r\nHost: catfact.ninja\r\nUser-Agent: RapidProMailroom/10.3.42\r\nAccept: application/json\r\nX-Mailroom-Mode: normal\r\nAccept-Encoding: gzip\r\n\r\n",
          "response": "HTTP/2.0 200 OK\r\nAccess-Control-Allow-Origin: *\r\nAlt-Svc: h3=\":443\"; ma=86400\r\nCache-Control: no-cache, private\r\nCf-Cache-Status: DYNAMIC\r\nCf-Ray: 971359de2cb9803d-FRA\r\nContent-Type: application/json\r\nDate: Mon, 18 Aug 2025 18:07:29 GMT\r\nNel: {\"report_to\":\"cf-nel\",\"success_fraction\":0.0,\"max_age\":604800}\r\nReport-To: {\"group\":\"cf-nel\",\"max_age\":604800,\"endpoints\":[{\"url\":\"https://a.nel.cloudflare.com/report/v4?s=H2u8%2Bc9S05pRSeF7LV4MIBO3tOtpQu1HwsfjdTauwCiV%2FxGrXWRdwWmr4ICkHjED4EG8d30TVyd4q82tER51tEcZ3l7J%2BiXyBAjSlls%3D\"}]}\r\nServer: cloudflare\r\nSet-Cookie: XSRF-TOKEN=%3D; SameSite=Lax; Secure; Path=/; Max-Age=7200; Expires=Mon, 18 Aug 2025 20:07:29 GMT\r\nSet-Cookie: catfacts_session=%3D; HttpOnly; SameSite=Lax; Secure; Path=/; Max-Age=7200; Expires=Mon, 18 Aug 2025 20:07:29 GMT\r\nSet-Cookie: XUZCCLqZh2i3HptP7hxbHTdfqGpajEZWaFkGRiGl=; HttpOnly; SameSite=Lax; Secure; Path=/; Max-Age=7200; Expires=Mon, 18 Aug 2025 20:07:29 GMT\r\nX-Content-Type-Options: nosniff\r\nX-Frame-Options: SAMEORIGIN\r\nX-Ratelimit-Limit: 100\r\nX-Ratelimit-Remaining: 98\r\nX-Xss-Protection: 1; mode=block\r\n\r\n{\"fact\":\"In an average year, cat owners in the United States spend over $2 billion on cat food.\",\"length\":86}",
          "elapsed_ms": 189,
          "retries": 0,
          "status": "success"
        },
        {
          "uuid": "0198be5d-6769-79f3-8ceb-fd7a908ea0c9",
          "type": "msg_created",
          "created_on": "2025-08-18T18:07:29.129652561Z",
          "step_uuid": "7ae0e662-91be-4ab4-a2ef-8e9fd6329bec",
          "msg": {
            "urn": "telegram:742307595?channel=bbda279d-2f7e-414e-8b04-78e5785acd9b&id=600#rowanseymour",
            "channel": {
              "uuid": "bbda279d-2f7e-414e-8b04-78e5785acd9b",
              "name": "Staging Test"
            },
            "text": "In an average year, cat owners in the United States spend over $2 billion on cat food.. Want another?",
            "locale": "eng-AF"
          }
        },
        {
          "uuid": "0198be5d-6769-7a11-bf2a-58d09f81ab66",
          "type": "msg_wait",
          "created_on": "2025-08-18T18:07:29.129660074Z",
          "step_uuid": "330ccadd-1923-4bf9-97fb-482b09432a52",
          "expires_on": "2025-08-18T18:12:29.12965883Z"
        },
        {
          "uuid": "0198be5d-9d25-733f-b086-cc90e147e4d7",
          "type": "msg_received",
          "created_on": "2025-08-18T18:07:43.098348379Z",
          "msg": {
            "urn": "telegram:742307595#rowanseymour",
            "channel": {
              "uuid": "bbda279d-2f7e-414e-8b04-78e5785acd9b",
              "name": "Staging Test"
            },
            "text": "yes",
            "external_id": "9128"
          }
        },
        {
          "uuid": "0198be5d-9e72-71e9-b10a-896228c3016f",
          "type": "msg_created",
          "created_on": "2025-08-18T18:07:43.218125626Z",
          "step_uuid": "c2b28222-cce3-4a2f-a696-5f8ddf6131cf",
          "msg": {
            "urn": "telegram:742307595?channel=bbda279d-2f7e-414e-8b04-78e5785acd9b&id=600#rowanseymour",
            "channel": {
              "uuid": "bbda279d-2f7e-414e-8b04-78e5785acd9b",
              "name": "Staging Test"
            },
            "text": "Hi there!",
            "locale": "eng-AF"
          }
        },
        {
          "uuid": "0198be5d-9f3d-76e5-96ff-f40cae763043",
          "type": "webhook_called",
          "created_on": "2025-08-18T18:07:43.421452406Z",
          "step_uuid": "12332d6d-5b22-417e-bfab-a184f52cc7f5",
          "url": "https://catfact.ninja/fact?visit=3",
          "status_code": 200,
          "request": "GET /fact?visit=3 HTTP/1.1\r\nHost: catfact.ninja\r\nUser-Agent: RapidProMailroom/10.3.42\r\nAccept: application/json\r\nX-Mailroom-Mode: normal\r\nAccept-Encoding: gzip\r\n\r\n",
          "response": "HTTP/2.0 200 OK\r\nAccess-Control-Allow-Origin: *\r\nAlt-Svc: h3=\":443\"; ma=86400\r\nCache-Control: no-cache, private\r\nCf-Cache-Status: DYNAMIC\r\nCf-Ray: 97135a376bea803d-FRA\r\nContent-Type: application/json\r\nDate: Mon, 18 Aug 2025 18:07:43 GMT\r\nNel: {\"report_to\":\"cf-nel\",\"success_fraction\":0.0,\"max_age\":604800}\r\nReport-To: {\"group\":\"cf-nel\",\"max_age\":604800,\"endpoints\":[{\"url\":\"https://a.nel.cloudflare.com/report/v4?s=4mRqij4DyZW2OihHSb6aB5jt6ICRqlO9DLdg7HktG4uyJ2ZE5YFkY620HVz9oWfGVip0m90ab1Ed1bLktskSYITkkyd3APOzs%2F0n9iI%3D\"}]}\r\nServer: cloudflare\r\nSet-Cookie: XSRF-TOKEN=%3D; SameSite=Lax; Secure; Path=/; Max-Age=7200; Expires=Mon, 18 Aug 2025 20:07:43 GMT\r\nSet-Cookie: catfacts_session=%3D; HttpOnly; SameSite=Lax; Secure; Path=/; Max-Age=7200; Expires=Mon, 18 Aug 2025 20:07:43 GMT\r\nSet-Cookie: 7qxUSwBR5f7LcXMUg9mcrrXaQdtdZa5zqobEKlu3=; HttpOnly; SameSite=Lax; Secure; Path=/; Max-Age=7200; Expires=Mon, 18 Aug 2025 20:07:43 GMT\r\nVary: Accept-Encoding\r\nX-Content-Type-Options: nosniff\r\nX-Frame-Options: SAMEORIGIN\r\nX-Ratelimit-Limit: 100\r\nX-Ratelimit-Remaining: 98\r\nX-Xss-Protection: 1; mode=block\r\n\r\n{\"fact\":\"Kittens who are taken along on short, trouble-free car trips to town tend to make good passengers when they get older. They get used to the sounds and motions of traveling and make less connection between the car and the visits to the vet.\",\"length\":239}",
          "elapsed_ms": 203,
          "retries": 0,
          "status": "success"
        },
		{
          "uuid": "0198be5d-9f3d-76e5-9777-f40cae763043",
          "type": "run_result_changed",
          "created_on": "2025-08-18T18:07:43.421452406Z",
          "step_uuid": "12332d6d-5b22-417e-bfab-a184f52cc7f5",
		  "name": "Fact",
		  "value": "200",
          "extra": {"fact": "Kittens who are taken along on short, trouble-free car trips to town tend to make good passengers when they get older. They get used to the sounds and motions of traveling and make less connection between the car and the visits to the vet.", "length": 239}
        },
        {
          "uuid": "0198be5d-9f3d-7d1a-8af2-8d65d3ea65d3",
          "type": "msg_created",
          "created_on": "2025-08-18T18:07:43.421859061Z",
          "step_uuid": "a63bebbf-4327-4c1d-9372-2cd1de085c0d",
          "msg": {
            "urn": "telegram:742307595?channel=bbda279d-2f7e-414e-8b04-78e5785acd9b&id=600#rowanseymour",
            "channel": {
              "uuid": "bbda279d-2f7e-414e-8b04-78e5785acd9b",
              "name": "Staging Test"
            },
            "text": "Kittens who are taken along on short, trouble-free car trips to town tend to make good passengers when they get older. They get used to the sounds and motions of traveling and make less connection between the car and the visits to the vet.. Want another?",
            "locale": "eng-AF"
          }
        },
        {
          "uuid": "0198be5d-9f3d-7d2e-84ab-d35b90424008",
          "type": "msg_wait",
          "created_on": "2025-08-18T18:07:43.421864146Z",
          "step_uuid": "11651d3c-5789-4cc1-a64b-72c0325dcc8b",
          "expires_on": "2025-08-18T18:12:43.421862781Z"
        }
      ],
      "results": {
        "result_1": {
          "name": "Result 1",
          "value": "yes",
          "category": "Yes",
          "node_uuid": "07282700-46ab-4f2e-be88-56f12145017f",
          "input": "yes",
          "created_on": "2025-08-18T18:07:43.21803134Z"
        }
      },
      "status": "waiting",
      "created_on": "2025-08-18T18:07:01.410231973Z",
      "modified_on": "2025-08-18T18:07:43.421864759Z",
      "exited_on": null
    }
  ],
  "status": "waiting",
  "input": {
    "type": "msg",
    "uuid": "0198be5d-9d25-733f-b086-cc90e147e4d7",
    "channel": {
      "uuid": "bbda279d-2f7e-414e-8b04-78e5785acd9b",
      "name": "Staging Test"
    },
    "created_on": "2025-08-18T18:07:43.098348379Z",
    "urn": "telegram:742307595#rowanseymour",
    "text": "yes",
    "external_id": "9128"
  }
}`

func TestLegacySessionsWithRunEvents(t *testing.T) {
	env := envs.NewBuilder().Build()
	source, err := static.NewSource([]byte(legacySessionAssets))
	require.NoError(t, err)

	// create our engine session
	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	contact, err := flows.ReadContact(sa, []byte(`{
		"uuid": "a98b5f89-d5e8-4bda-9092-d31031379b38",
		"name": "Ryan Lewis",
        "status": "active",
		"language": "eng",
		"timezone": "America/Guayaquil",
		"urns": [
			"facebook:1234567890"
		],
		"groups": [],
		"fields": {},
		"created_on": "2018-06-20T11:40:30.123456789-00:00"
    }`), assets.PanicOnMissing)
	require.NoError(t, err)

	eng := engine.NewBuilder().Build()
	session, err := eng.ReadSession(sa, []byte(legacySessionWithRunEvents), env, contact, nil, assets.PanicOnMissing)
	require.NoError(t, err)

	assert.Equal(t, 3, session.Sprints())
	assert.True(t, session.Runs()[0].HadInput())
	assert.Equal(t, "{\"fact\": \"Kittens who are taken along on short, trouble-free car trips to town tend to make good passengers when they get older. They get used to the sounds and motions of traveling and make less connection between the car and the visits to the vet.\", \"length\": 239}", string(session.Runs()[0].Webhook().ResponseJSON))
}
