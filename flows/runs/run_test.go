package runs_test

import (
	"strings"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/runs"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"

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
            "spec_version": "13.0",
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
	uuids.SetGenerator(uuids.NewSeededGenerator(12345))
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	server := test.NewTestHTTPServer(49999)
	defer server.Close()

	session, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyNone)
	require.NoError(t, err)

	flow, err := session.Assets().Flows().Get("50c3706e-fedb-42c0-8eab-dda3335714b7")
	require.NoError(t, err)

	run := session.Runs()[0]

	checkRun := func(r flows.Run) {
		assert.Equal(t, string(flows.RunUUID("e7187099-7d38-4f60-955c-325957214c42")), string(r.UUID()))
		assert.Equal(t, string(flows.RunStatusCompleted), string(r.Status()))
		assert.Equal(t, flow, r.Flow())
		assert.Equal(t, flow.Reference(), r.FlowReference())
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
	uuids.SetGenerator(uuids.NewSeededGenerator(12345))
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	dates.SetNowSource(dates.NewFixedNowSource(time.Date(2018, 9, 13, 13, 36, 30, 123456789, time.UTC)))
	defer dates.SetNowSource(dates.DefaultNowSource)

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
		{`@child.uuid`, `9688d21d-95aa-4bed-afc7-f31b35731a3d`},
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
			`{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00","not_set":null,"state":null}`,
		},
		{
			`@(json(fields))`,
			`{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00","not_set":null,"state":null}`,
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
			`{"categories":["Red"],"categories_localized":["Red"],"category":"Red","category_localized":"Red","created_on":"2018-09-13T13:36:30.123456Z","extra":null,"input":"","name":"Favorite Color","node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03","value":"red","values":["red"]}`,
		},
		{
			`@(json(run.results.favorite_color))`,
			`{"categories":["Red"],"categories_localized":["Red"],"category":"Red","category_localized":"Red","created_on":"2018-09-13T13:36:30.123456Z","extra":null,"input":"","name":"Favorite Color","node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03","value":"red","values":["red"]}`,
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
			`{"activation_token":null,"age":33,"gender":"Female","join_date":null,"not_set":null,"state":null}`,
		},
	}

	for _, tc := range testCases {
		actual, err := run.EvaluateTemplate(tc.template)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual, "template mismatch for %s", tc.template)
	}

	// test with escaping
	evaluated, err := run.EvaluateTemplateText(`gender = @("M\" OR")`, flows.ContactQueryEscaping, true)
	assert.NoError(t, err)
	assert.Equal(t, `gender = "M\" OR"`, evaluated)
}

func TestMissingRelatedRunContext(t *testing.T) {
	// create a run with no parent or child
	sa, err := test.CreateSessionAssets([]byte(sessionAssets), "")
	require.NoError(t, err)

	trigger, err := triggers.ReadTrigger(sa, []byte(sessionTrigger), assets.IgnoreMissing)
	require.NoError(t, err)

	eng := test.NewEngine()
	session, _, err := eng.NewSession(sa, trigger)
	require.NoError(t, err)

	run := session.Runs()[0]

	// since we have no parent, check that it resolves to nil
	val, err := run.EvaluateTemplateValue(`@parent`)
	assert.NoError(t, err)
	assert.Nil(t, val)

	// check that trying to resolve a property of parent is an error
	val, err = run.EvaluateTemplateValue(`@parent.contact`)
	assert.NoError(t, err)
	assert.Equal(t, types.NewXErrorf("null doesn't support lookups"), val)

	// we also have no child, check that it resolves to nil
	val, err = run.EvaluateTemplateValue(`@child`)
	assert.NoError(t, err)
	assert.Nil(t, val)

	// check that trying to resolve a property of child is an error
	val, err = run.EvaluateTemplateValue(`@child.contact`)
	assert.NoError(t, err)
	assert.Equal(t, types.NewXErrorf("null doesn't support lookups"), val)
}

func TestSaveResult(t *testing.T) {
	sa, err := test.CreateSessionAssets([]byte(sessionAssets), "")
	require.NoError(t, err)

	trigger, err := triggers.ReadTrigger(sa, []byte(sessionTrigger), assets.IgnoreMissing)
	require.NoError(t, err)

	eng := test.NewEngine()
	session, _, err := eng.NewSession(sa, trigger)
	require.NoError(t, err)

	run := session.Runs()[0]

	dates.SetNowSource(dates.NewFixedNowSource(time.Date(2020, 4, 20, 12, 39, 30, 123456789, time.UTC)))
	defer dates.SetNowSource(dates.DefaultNowSource)

	// no results means empty object with default of empty string
	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{"__default__": types.XTextEmpty}), flows.Context(session.Environment(), run.Results()))

	run.SaveResult(flows.NewResult("Response 1", "red", "Red", "Rojo", "6d35528e-cae3-4e30-b842-8fe6ed7d5c02", "I like red", nil, dates.Now()))

	// name is snaked
	assert.Equal(t, "red", run.Results().Get("response_1").Value)
	assert.Equal(t, "Red", run.Results().Get("response_1").Category)
	assert.Equal(t, time.Date(2020, 4, 20, 12, 39, 30, 123456789, time.UTC), run.ModifiedOn())

	run.SaveResult(flows.NewResult("Response 1", "blue", "Blue", "Azul", "6d35528e-cae3-4e30-b842-8fe6ed7d5c02", "I like blue", nil, dates.Now()))

	// result is overwritten
	assert.Equal(t, "blue", run.Results().Get("response_1").Value)
	assert.Equal(t, "Blue", run.Results().Get("response_1").Category)
	assert.Equal(t, time.Date(2020, 4, 20, 12, 39, 30, 123456789, time.UTC), run.ModifiedOn())

	// long values should truncated
	run.SaveResult(flows.NewResult("Response 1", strings.Repeat("創", 700), "Blue", "Azul", "6d35528e-cae3-4e30-b842-8fe6ed7d5c02", "I like blue", nil, dates.Now()))

	assert.Equal(t, strings.Repeat("創", 640), run.Results().Get("response_1").Value)
}
