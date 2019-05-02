package runs_test

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows/engine"
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
            "address": "+12345671111",
            "schemes": ["tel"],
            "roles": ["send", "receive"]
        }
    ],
    "fields": [
        {
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
                        "wait": {
                            "type": "msg",
                            "timeout": {
                                "seconds": 600,
                                "category_uuid": "0680b01f-ba0b-48f4-a688-d2f963130126"
                            }
                        },
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
        "default_language": "eng",
        "allowed_languages": [
            "eng", 
            "spa"
        ],
        "redaction_policy": "none",
        "time_format": "hh:mm",
        "timezone": "America/Guayaquil"
    }
}`

func TestRunContext(t *testing.T) {
	utils.SetTimeSource(utils.NewFixedTimeSource(time.Date(2018, 9, 13, 13, 36, 30, 123456789, time.UTC)))
	defer utils.SetTimeSource(utils.DefaultTimeSource)

	// create a run with no parent or child
	session, _, err := test.CreateTestSession("", nil)
	require.NoError(t, err)

	run := session.Runs()[0]

	testCases := []struct {
		template string
		expected string
	}{
		{`@run`, `Ryan Lewis@Registration`},
		{`@child`, `Ryan Lewis@Collect Age`},
		{`@child.run`, `Ryan Lewis@Collect Age`},
		{`@child.contact.name`, `Ryan Lewis`},
		{`@child.run.contact.name`, `Ryan Lewis`},
		{`@child.fields`, "Activation Token: AACC55\nAge: 23\nGender: Male\nJoin Date: 2017-12-02T00:00:00.000000-02:00"},
		{`@parent`, `Jasmine@Parent`},
		{`@parent.run`, `Jasmine@Parent`},
		{`@parent.contact.name`, `Jasmine`},
		{`@parent.run.contact.name`, `Jasmine`},
		{`@parent.fields`, "Age: 33\nGender: Female"},
		{
			`@(json(contact.fields))`,
			`{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00","not_set":null}`,
		},
		{
			`@(json(fields))`,
			`{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00","not_set":null}`,
		},
		{
			`@(json(contact.urns))`,
			`["tel:+12065551212","twitterid:54784326227#nyaruka","mailto:foo@bar.com"]`,
		},
		{
			`@(json(urns))`,
			`{"ext":null,"facebook":null,"fcm":null,"jiochat":null,"line":null,"mailto":"mailto:foo@bar.com","tel":"tel:+12065551212","telegram":null,"twitter":null,"twitterid":"twitterid:54784326227#nyaruka","viber":null,"wechat":null,"whatsapp":null}`,
		},
		{
			`@(json(results.favorite_color))`,
			`{"category":"Red","category_localized":"Red","created_on":"2018-09-13T13:36:30.123456Z","input":"","name":"Favorite Color","node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03","value":"red"}`,
		},
		{
			`@(json(run.results.favorite_color))`,
			`{"categories":["Red"],"categories_localized":["Red"],"created_on":"2018-09-13T13:36:30.123456Z","extra":null,"input":"","name":"Favorite Color","node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03","values":["red"]}`,
		},
		{
			`@(json(parent.contact.urns))`,
			`["tel:+593979111222"]`,
		},
		{
			`@(json(parent.urns))`,
			`{"ext":null,"facebook":null,"fcm":null,"jiochat":null,"line":null,"mailto":null,"tel":"tel:+593979111222","telegram":null,"twitter":null,"twitterid":null,"viber":null,"wechat":null,"whatsapp":null}`,
		},
		{
			`@(json(parent.fields))`,
			`{"activation_token":null,"age":33,"gender":"Female","join_date":null,"not_set":null}`,
		},
	}

	for _, tc := range testCases {
		actual, err := run.EvaluateTemplate(tc.template)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestMissingRelatedRunContext(t *testing.T) {
	// create a run with no parent or child
	sa, err := test.CreateSessionAssets([]byte(sessionAssets), "")
	require.NoError(t, err)

	trigger, err := triggers.ReadTrigger(sa, []byte(sessionTrigger), assets.IgnoreMissing)
	require.NoError(t, err)

	eng := engine.NewBuilder().WithDefaultUserAgent("goflow-testing").Build()
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
	assert.Equal(t, types.NewXErrorf("null has no property 'contact'"), val)

	// we also have no child, check that it resolves to nil
	val, err = run.EvaluateTemplateValue(`@child`)
	assert.NoError(t, err)
	assert.Nil(t, val)

	// check that trying to resolve a property of child is an error
	val, err = run.EvaluateTemplateValue(`@child.contact`)
	assert.NoError(t, err)
	assert.Equal(t, types.NewXErrorf("null has no property 'contact'"), val)
}
