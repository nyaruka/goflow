package runs_test

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
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/runs"
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

var sessionTrigger = `{
    "type": "manual",
    "triggered_on": "2017-12-31T11:31:15.035757258-02:00",
    "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "No Related Runs"},
    "contact": {
        "uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f",
        "id": 1234567,
        "name": "Ryan Lewis",
        "language": "eng",
        "timezone": "America/Guayaquil",
        "created_on": "2018-06-20T11:40:30.123456789-00:00",
        "urns": [ "tel:+12065551212"],
        "fields": {
            "gender": {"text": "M"}
        }
    },
    "environment": {
        "date_format": "YYYY-MM-DD",
        "allowed_languages": [
            "eng", 
            "spa"
        ],
        "redaction_policy": "none",
        "time_format": "hh:mm",
        "timezone": "America/Guayaquil"
    }
}`

func TestRun(t *testing.T) {
	test.MockUniverse()

	server := test.NewTestHTTPServer(49999)
	defer server.Close()

	session, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyNone)
	require.NoError(t, err)

	flow, err := session.Assets().Flows().Get("50c3706e-fedb-42c0-8eab-dda3335714b7")
	require.NoError(t, err)

	run := session.Runs()[0]

	checkRun := func(r flows.Run) {
		assert.Equal(t, flows.RunUUID("01969b47-113b-76f8-9c0b-2014ddc77094"), r.UUID())
		assert.Equal(t, flows.RunStatusCompleted, r.Status())
		assert.Equal(t, flow, r.Flow())
		assert.Equal(t, flow.Reference(true), r.FlowReference())
		assert.Equal(t, 10, len(r.Events()))
		assert.Equal(t, "Parent", r.Parent().Flow().Name())
		assert.Equal(t, 0, len(r.Ancestors())) // no parent runs within this session
		assert.True(t, r.ReceivedInput())
	}

	checkRun(run)

	// check we can marshal and marshal the run and get the same values
	runJSON, err := jsonx.Marshal(run)
	require.NoError(t, err)

	run2, err := runs.ReadRun(session, runJSON, assets.IgnoreMissing)
	require.NoError(t, err)

	checkRun(run2)
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
		{`@child.uuid`, `01969b47-24c3-76f8-8228-9728778b6c98`},
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

	trigger, err := triggers.ReadTrigger(sa, []byte(sessionTrigger), assets.IgnoreMissing)
	require.NoError(t, err)

	eng := test.NewEngine()
	session, _, err := eng.NewSession(context.Background(), sa, trigger)
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

	trigger, err := triggers.ReadTrigger(sa, []byte(sessionTrigger), assets.IgnoreMissing)
	require.NoError(t, err)

	eng := test.NewEngine()
	session, _, err := eng.NewSession(context.Background(), sa, trigger)
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
		evt := sp.Events()[0].(*events.MsgCreatedEvent)

		assert.Equal(t, tc.expectedText, evt.Msg.Text(), "msg text mismatch in test '%s'", tc.description)
		assert.Equal(t, tc.expectedAttachments, evt.Msg.Attachments(), "attachments mismatch in test case '%s'", tc.description)
		assert.Equal(t, tc.expectedQuickReplies, evt.Msg.QuickReplies(), "quick replies mismatch in test case '%s'", tc.description)
	}
}
