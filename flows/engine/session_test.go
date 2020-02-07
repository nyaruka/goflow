package engine_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/urns"
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
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var templateTests = []struct {
	template   string
	expected   string
	errorMsg   string
	redactURNs bool
}{
	// contact basic properties
	{"@contact.uuid", "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f", "", false},
	{"@contact.id", "1234567", "", false},
	{"@CONTACT.NAME", "Ryan Lewis", "", false},
	{"@contact.name", "Ryan Lewis", "", false},
	{"@contact.language", "eng", "", false},
	{"@contact.timezone", "America/Guayaquil", "", false},

	// contact single URN access
	{"@contact.urn", `tel:+12024561111`, "", false},
	{"@(urn_parts(contact.urn).scheme)", `tel`, "", false},
	{"@(urn_parts(contact.urn).path)", `+12024561111`, "", false},
	{"@(format_urn(contact.urn))", `(202) 456-1111`, "", false},

	// with URN redaction
	{"@contact.urn", `tel:********`, "", true},
	{"@(urn_parts(contact.urn).scheme)", `tel`, "", true},
	{"@(urn_parts(contact.urn).path)", `********`, "", true},
	{"@(format_urn(contact.urn))", `********`, "", true},

	// contact URN list access
	{"@contact.urns", `[tel:+12024561111, twitterid:54784326227#nyaruka, mailto:foo@bar.com]`, "", false},
	{"@(contact.urns[0])", "tel:+12024561111", "", false},
	{"@(contact.urns[110])", "", "error evaluating @(contact.urns[110]): index 110 out of range for 3 items", false},
	{"@(urn_parts(contact.urns[0]).scheme)", "tel", "", false},
	{"@(urn_parts(contact.urns[0]).path)", "+12024561111", "", false},
	{"@(urn_parts(contact.urns[0]).display)", "", "", false},
	{"@(contact.urns[1])", "twitterid:54784326227#nyaruka", "", false},
	{"@(format_urn(contact.urns[0]))", "(202) 456-1111", "", false},

	// with URN redaction
	{"@contact.urns", `[tel:********, twitterid:********, mailto:********]`, "", true},
	{"@(contact.urns[0])", `tel:********`, "", true},

	// simplified URN access
	{"@urns", `{ext: , facebook: , fcm: , freshchat: , jiochat: , line: , mailto: mailto:foo@bar.com, tel: tel:+12024561111, telegram: , twitter: , twitterid: twitterid:54784326227#nyaruka, viber: , vk: , wechat: , whatsapp: }`, "", false},
	{"@urns.tel", `tel:+12024561111`, "", false},
	{"@urns.mailto", `mailto:foo@bar.com`, "", false},
	{"@urns.viber", ``, "", false},
	{"@(format_urn(urns.tel))", "(202) 456-1111", "", false},

	// with URN redaction
	{"@urns.tel", `tel:********`, "", true},
	{"@urns.viber", ``, "", true},

	// contact groups
	{`@(foreach(contact.groups, extract, "name"))`, `[Testers, Males]`, "", false},
	{`@(join(foreach(contact.groups, extract, "name"), "|"))`, `Testers|Males`, "", false},
	{`@(count(contact.groups))`, "2", "", false},

	// contact fields
	{"@contact.fields", "Activation Token: AACC55\nAge: 23\nGender: Male\nJoin Date: 2017-12-02T00:00:00.000000-02:00", "", false},
	{"@contact.fields.activation_token", "AACC55", "", false},
	{"@contact.fields.age", "23", "", false},
	{"@contact.fields.join_date", "2017-12-02T00:00:00.000000-02:00", "", false},
	{"@contact.fields.favorite_icecream", "", "error evaluating @contact.fields.favorite_icecream: object has no property 'favorite_icecream'", false},
	{"@(is_error(contact.fields.favorite_icecream))", "true", "", false},
	{"@(has_error(contact.fields.favorite_icecream).match)", "object has no property 'favorite_icecream'", "", false},
	{"@(count(contact.fields))", "5", "", false},

	// simplifed field access
	{"@fields", "Activation Token: AACC55\nAge: 23\nGender: Male\nJoin Date: 2017-12-02T00:00:00.000000-02:00", "", false},
	{"@fields.activation_token", "AACC55", "", false},
	{"@fields.age", "23", "", false},
	{"@fields.join_date", "2017-12-02T00:00:00.000000-02:00", "", false},
	{"@fields.favorite_icecream", "", "error evaluating @fields.favorite_icecream: object has no property 'favorite_icecream'", false},
	{"@(is_error(fields.favorite_icecream))", "true", "", false},
	{"@(has_error(fields.favorite_icecream).match)", "object has no property 'favorite_icecream'", "", false},
	{"@(count(fields))", "5", "", false},

	{"@input", "Hi there\nhttp://s3.amazon.com/bucket/test.jpg\nhttp://s3.amazon.com/bucket/test.mp3", "", false},
	{"@input.text", "Hi there", "", false},
	{"@input.attachments", `[image/jpeg:http://s3.amazon.com/bucket/test.jpg, audio/mp3:http://s3.amazon.com/bucket/test.mp3]`, "", false},
	{"@(input.attachments[0])", "image/jpeg:http://s3.amazon.com/bucket/test.jpg", "", false},
	{"@input.created_on", "2017-12-31T11:35:10.035757-02:00", "", false},
	{"@input.channel.name", "My Android Phone", "", false},

	{"@results.favorite_color", `red`, "", false},
	{"@results.favorite_color.value", "red", "", false},
	{"@results.favorite_color.category", "Red", "", false},
	{"@results.favorite_color.category_localized", "Red", "", false},
	{"@(is_error(results.favorite_icecream))", "true", "", false},
	{"@(has_error(results.favorite_icecream).match)", "object has no property 'favorite_icecream'", "", false},
	{"@(count(results))", "5", "", false},

	{"@run.results.favorite_color", `red`, "", false},
	{"@run.results.favorite_color.value", "red", "", false},
	{"@run.results.favorite_color.values", "[red]", "", false},
	{`@(run.results.favorite_color.values[0])`, `red`, "", false},
	{"@run.results.favorite_color.category", "Red", "", false},
	{"@run.results.favorite_color.categories", "[Red]", "", false},
	{`@(run.results.favorite_color.categories[0])`, `Red`, "", false},
	{"@run.results.favorite_color.category_localized", "Red", "", false},
	{"@run.results.favorite_color.categories_localized", "[Red]", "", false},
	{"@run.results.favorite_icecream", "", "error evaluating @run.results.favorite_icecream: object has no property 'favorite_icecream'", false},
	{"@(is_error(run.results.favorite_icecream))", "true", "", false},
	{"@(has_error(run.results.favorite_icecream).match)", "object has no property 'favorite_icecream'", "", false},
	{"@(count(run.results))", "5", "", false},

	{"@run.status", "completed", "", false},

	{"@webhook", "{results: [{state: WA}, {state: IN}]}", "", false},
	{"@webhook.results", "[{state: WA}, {state: IN}]", "", false},
	{"@(webhook.results[1])", "{state: IN}", "", false},
	{"@(webhook.results[1].state)", "IN", "", false},

	{"@trigger.params", `{address: {state: WA}, source: website}`, "", false},
	{"@trigger.params.source", "website", "", false},
	{"@(count(trigger.params.address))", "1", "", false},

	// migrated split by expressions
	{`@(if(is_error(results.favorite_color.value), "@flow.favorite_color", results.favorite_color.value))`, `red`, "", false},
	{`@(if(is_error(legacy_extra["0"].default_city), "@extra.0.default_city", legacy_extra["0"].default_city))`, `@extra.0.default_city`, "", false},

	// non-expressions
	{"bob@nyaruka.com", "bob@nyaruka.com", "", false},
	{"@twitter_handle", "@twitter_handle", "", false},
}

func TestEvaluateTemplate(t *testing.T) {
	dates.SetNowSource(dates.NewFixedNowSource(time.Date(2018, 9, 13, 13, 36, 30, 123456789, time.UTC)))
	defer dates.SetNowSource(dates.DefaultNowSource)

	server := test.NewTestHTTPServer(0)
	defer server.Close()

	sessionWithURNs, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyNone)
	require.NoError(t, err)
	sessionWithoutURNs, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyURNs)
	require.NoError(t, err)

	for _, tc := range templateTests {
		var run flows.FlowRun
		if tc.redactURNs {
			run = sessionWithoutURNs.Runs()[0]
		} else {
			run = sessionWithURNs.Runs()[0]
		}

		eval, err := run.EvaluateTemplate(tc.template)

		var actualErrorMsg string
		if err != nil {
			actualErrorMsg = err.Error()
		}

		assert.Equal(t, tc.expected, eval, "output mismatch evaluating template: '%s'", tc.template)
		assert.Equal(t, tc.errorMsg, actualErrorMsg, "error mismatch evaluating template: '%s'", tc.template)
	}
}

func BenchmarkEvaluateTemplate(b *testing.B) {
	session, _, err := test.CreateTestSession("http://localhost", envs.RedactionPolicyNone)
	require.NoError(b, err)

	run := session.Runs()[0]

	for n := 0; n < b.N; n++ {
		for _, tc := range templateTests {
			run.EvaluateTemplate(tc.template)
		}
	}
}

func TestContextToJSON(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"contact.uuid", `"5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"`},
		{"contact.name", `"Ryan Lewis"`},
		{"contact.urns", `["tel:+12024561111","twitterid:54784326227#nyaruka","mailto:foo@bar.com"]`},
		{"contact.urns[0]", `"tel:+12024561111"`},
		{"contact.fields", `{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00","not_set":null}`},
		{"contact.fields.age", `23`},
		{
			"contact",
			`{
				"channel": {"address":"+17036975131","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},
				"created_on": "2018-06-20T11:40:30.123456Z", 
				"fields": {"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00","not_set":null},
				"first_name": "Ryan",
				"groups": [{"name":"Testers","uuid":"b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"},{"name":"Males","uuid":"4f1f98fc-27a7-4a69-bbdb-24744ba739a9"}],
				"id": "1234567",
				"language": "eng",
				"name": "Ryan Lewis",
				"timezone": "America/Guayaquil",
				"urn": "tel:+12024561111",
				"urns": ["tel:+12024561111","twitterid:54784326227#nyaruka","mailto:foo@bar.com"],
				"uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"
			}`,
		},
		{
			"input",
			`{
				"attachments":["image/jpeg:http://s3.amazon.com/bucket/test.jpg","audio/mp3:http://s3.amazon.com/bucket/test.mp3"],
				"channel":{"address":"+17036975131","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},
				"created_on":"2017-12-31T11:35:10.035757-02:00",
				"external_id":"",
				"text":"Hi there",
				"type":"msg",
				"urn":"tel:+12065551212",
				"uuid":"9bf91c2b-ce58-4cef-aacc-281e03f69ab5"
			}`,
		},
		{
			"run",
			`{
				"contact": {
					"channel":{"address":"+17036975131","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},
					"created_on":"2018-06-20T11:40:30.123456Z",
					"fields":{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00","not_set":null},
					"first_name":"Ryan",
					"groups":[{"name":"Testers","uuid":"b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"},{"name":"Males","uuid":"4f1f98fc-27a7-4a69-bbdb-24744ba739a9"}],
					"id":"1234567",
					"language":"eng",
					"name":"Ryan Lewis",
					"timezone":"America/Guayaquil",
					"urn":"tel:+12024561111",
					"urns":["tel:+12024561111","twitterid:54784326227#nyaruka","mailto:foo@bar.com"],
					"uuid":"5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"
				},
				"created_on":"2018-04-11T13:24:30.123456Z",
				"exited_on":"2018-04-11T13:24:30.123456Z",
				"flow":{"name":"Registration","revision":123,"uuid":"50c3706e-fedb-42c0-8eab-dda3335714b7"},
				"path":[
					{"arrived_on":"2018-04-11T13:24:30.123456Z","exit_uuid":"d7a36118-0a38-4b35-a7e4-ae89042f0d3c","node_uuid":"72a1f5df-49f9-45df-94c9-d86f7ea064e5","uuid":"8720f157-ca1c-432f-9c0b-2014ddc77094"},
					{"arrived_on":"2018-04-11T13:24:30.123456Z","exit_uuid":"100f2d68-2481-4137-a0a3-177620ba3c5f","node_uuid":"3dcccbb4-d29c-41dd-a01f-16d814c9ab82","uuid":"970b8069-50f5-4f6f-8f41-6b2d9f33d623"},
					{"arrived_on":"2018-04-11T13:24:30.123456Z","exit_uuid":"d898f9a4-f0fc-4ac4-a639-c98c602bb511","node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03","uuid":"5ecda5fc-951c-437b-a17e-f85e49829fb9"},
					{"arrived_on":"2018-04-11T13:24:30.123456Z","exit_uuid":"9fc5f8b4-2247-43db-b899-ab1ac50ba06c","node_uuid":"c0781400-737f-4940-9a6c-1ec1c3df0325","uuid":"312d3af0-a565-4c96-ba00-bd7f0d08e671"}
				],
				"results":{
					"2factor":{
						"category":"",
						"categories":[""],
						"category_localized":"",
						"categories_localized":[""],
						"created_on":"2018-04-11T13:24:30.123456Z",
						"extra":null,
						"input":"",
						"name":"2Factor",
						"node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
						"value":"34634624463525",
						"values":["34634624463525"]
					},
					"favorite_color":{
						"category":"Red",
						"categories":["Red"],
						"category_localized":"Red",
						"categories_localized":["Red"],
						"created_on":"2018-04-11T13:24:30.123456Z",
						"extra":null,
						"input":"",
						"name":"Favorite Color",
						"node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
						"value":"red",
						"values":["red"]
					},
					"intent": {
						"categories": [
							"Success"
						],
						"categories_localized": [
							"Success"
						],
						"category": "Success",
						"category_localized": "Success",
						"created_on": "2018-04-11T13:24:30.123456Z",
						"extra": {
							"entities": {
								"location": [
									{
										"confidence": 1,
										"value": "Quito"
									}
								]
							},
							"intents": [
								{
									"confidence": 0.5,
									"name": "book_flight"
								},
								{
									"confidence": 0.25,
									"name": "book_hotel"
								}
							]
						},
						"input": "Hi there",
						"name": "intent",
						"node_uuid": "f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
						"value": "book_flight",
						"values": [
							"book_flight"
						]
					},
					"phone_number":{
						"category":"",
						"categories":[""],
						"category_localized":"",
						"categories_localized":[""],
						"created_on":"2018-04-11T13:24:30.123456Z",
						"extra":null,
						"input":"",
						"name":"Phone Number",
						"node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
						"value":"+12344563452",
						"values":["+12344563452"]
					},
					"webhook":{
						"category":"Success",
						"categories":["Success"],
						"category_localized":"Success",
						"categories_localized":["Success"],
						"created_on":"2018-04-11T13:24:30.123456Z",
						"extra":{"results":[{"state":"WA"},{"state":"IN"}]},
						"input":"GET http://127.0.0.1:49992/?content=%7B%22results%22%3A%5B%7B%22state%22%3A%22WA%22%7D%2C%7B%22state%22%3A%22IN%22%7D%5D%7D",
						"name":"webhook",
						"node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
						"value":"200",
						"values":["200"]
					}
				},
				"status":"completed",
				"uuid":"692926ea-09d6-4942-bd38-d266ec8d3716"
			}`,
		},
		{
			"child",
			`{
				"contact": {
					"channel": {
						"address": "+17036975131",
						"name": "My Android Phone",
						"uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d"
					},
					"created_on": "2018-06-20T11:40:30.123456Z",
					"fields": {
						"activation_token": "AACC55",
						"age": 23,
						"gender": "Male",
						"join_date": "2017-12-02T00:00:00.000000-02:00",
						"not_set": null
					},
					"first_name": "Ryan",
					"groups": [
						{
							"name": "Testers",
							"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"
						},
						{
							"name": "Males",
							"uuid": "4f1f98fc-27a7-4a69-bbdb-24744ba739a9"
						}
					],
					"id": "1234567",
					"language": "eng",
					"name": "Ryan Lewis",
					"timezone": "America/Guayaquil",
					"urn": "tel:+12024561111",
					"urns": [
						"tel:+12024561111",
						"twitterid:54784326227#nyaruka",
						"mailto:foo@bar.com"
					],
					"uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"
				},
				"fields": {
					"activation_token": "AACC55",
					"age": 23,
					"gender": "Male",
					"join_date": "2017-12-02T00:00:00.000000-02:00",
					"not_set": null
				},
				"flow": {
					"name": "Collect Age",
					"revision": 0,
					"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"
				},
				"results": {
					"age": {
						"category": "Youth",
						"categories": ["Youth"],
						"category_localized": "Youth",
						"categories_localized": ["Youth"],
						"created_on": "2018-04-11T13:24:30.123456Z",
						"extra": null,
						"input": "",
						"name": "Age",
						"node_uuid": "d9dba561-b5ee-4f62-ba44-60c4dc242b84",
						"value": "23",
						"values": ["23"]
					}
				},
				"run": {
					"contact": {
						"channel": {
							"address": "+17036975131",
							"name": "My Android Phone",
							"uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d"
						},
						"created_on": "2018-06-20T11:40:30.123456Z",
						"fields": {
							"activation_token": "AACC55",
							"age": 23,
							"gender": "Male",
							"join_date": "2017-12-02T00:00:00.000000-02:00",
							"not_set": null
						},
						"first_name": "Ryan",
						"groups": [
							{
								"name": "Testers",
								"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"
							},
							{
								"name": "Males",
								"uuid": "4f1f98fc-27a7-4a69-bbdb-24744ba739a9"
							}
						],
						"id": "1234567",
						"language": "eng",
						"name": "Ryan Lewis",
						"timezone": "America/Guayaquil",
						"urn": "tel:+12024561111",
						"urns": [
							"tel:+12024561111",
							"twitterid:54784326227#nyaruka",
							"mailto:foo@bar.com"
						],
						"uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"
					},
					"flow": {
						"name": "Collect Age",
						"revision": 0,
						"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"
					},
					"results": {
						"age": {						
							"category": "Youth",
							"categories": ["Youth"],
							"category_localized": "Youth",
							"categories_localized": ["Youth"],
							"created_on": "2018-04-11T13:24:30.123456Z",
							"extra": null,
							"input": "",
							"name": "Age",
							"node_uuid": "d9dba561-b5ee-4f62-ba44-60c4dc242b84",
							"value": "23",
							"values": ["23"]
						}
					},
					"status": "completed",
					"uuid": "c34b6c7d-fa06-4563-92a3-d648ab64bccb"
				},
				"status": "completed",
				"urns": {
					"ext": null,
					"facebook": null,
					"fcm": null,
					"freshchat": null,
					"jiochat": null,
					"line": null,
					"mailto": "mailto:foo@bar.com",
					"tel": "tel:+12024561111",
					"telegram": null,
					"twitter": null,
					"twitterid": "twitterid:54784326227#nyaruka",
					"viber": null,
					"vk": null,
					"wechat": null,
					"whatsapp": null
				},
				"uuid": "c34b6c7d-fa06-4563-92a3-d648ab64bccb"
			}`,
		},
		{
			"parent",
			`{
				"contact": {
					"channel": {
						"address": "+17036975131",
						"name": "My Android Phone",
						"uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d"
					},
					"created_on": "2018-01-01T12:00:00.000000Z",
					"fields": {
						"activation_token": null,
						"age": 33,
						"gender": "Female",
						"join_date": null,
						"not_set": null
					},
					"first_name": "Jasmine",
					"groups": [],
					"id": "0",
					"language": "spa",
					"name": "Jasmine",
					"timezone": null,
					"urn": "tel:+12024562222",
					"urns": [
						"tel:+12024562222"
					],
					"uuid": "c59b0033-e748-4240-9d4c-e85eb6800151"
				},
				"fields": {
					"activation_token": null,
					"age": 33,
					"gender": "Female",
					"join_date": null,
					"not_set": null
				},
				"flow": {
					"name": "Parent",
					"revision": 0,
					"uuid": "fece6eac-9127-4343-9269-56e88f391562"
				},
				"results": {
					"role": {
						"category": "Reporter",
						"categories": ["Reporter"],
						"category_localized": "Reporter",
						"categories_localized": ["Reporter"],
						"created_on": "2000-01-01T00:00:00.000000Z",
						"extra": null,
						"input": "a reporter",
						"name": "Role",
						"node_uuid": "385cb848-5043-448e-9123-05cbcf26ad74",
						"value": "reporter",
						"values": ["reporter"]
					}
				},
				"run": {
					"contact": {
						"channel": {
							"address": "+17036975131",
							"name": "My Android Phone",
							"uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d"
						},
						"created_on": "2018-01-01T12:00:00.000000Z",
						"fields": {
							"activation_token": null,
							"age": 33,
							"gender": "Female",
							"join_date": null,
							"not_set": null
						},
						"first_name": "Jasmine",
						"groups": [],
						"id": "0",
						"language": "spa",
						"name": "Jasmine",
						"timezone": null,
						"urn": "tel:+12024562222",
						"urns": [
							"tel:+12024562222"
						],
						"uuid": "c59b0033-e748-4240-9d4c-e85eb6800151"
					},
					"flow": {
						"name": "Parent",
						"revision": 0,
						"uuid": "fece6eac-9127-4343-9269-56e88f391562"
					},
					"results": {
						"role": {
							"category": "Reporter",
							"categories": ["Reporter"],
							"category_localized": "Reporter",
							"categories_localized": ["Reporter"],
							"created_on": "2000-01-01T00:00:00.000000Z",
							"extra": null,
							"input": "a reporter",
							"name": "Role",
							"node_uuid": "385cb848-5043-448e-9123-05cbcf26ad74",
							"value": "reporter",
							"values": ["reporter"]
						}
					},
					"status": "active",
					"uuid": "4213ac47-93fd-48c4-af12-7da8218ef09d"
				},
				"status": "active",
				"urns": {
					"ext": null,
					"facebook": null,
					"fcm": null,
					"freshchat": null,
					"jiochat": null,
					"line": null,
					"mailto": null,
					"tel": "tel:+12024562222",
					"telegram": null,
					"twitter": null,
					"twitterid": null,
					"viber": null,
					"vk": null,
					"wechat": null,
					"whatsapp": null
				},
				"uuid": "4213ac47-93fd-48c4-af12-7da8218ef09d"
			}`,
		},
		{"trigger", `{"type":"flow_action","params":{"source":"website","address":{"state":"WA"}},"keyword": null}`},
	}

	server := test.NewTestHTTPServer(49992)
	defer server.Close()
	defer uuids.SetGenerator(uuids.DefaultGenerator)
	defer dates.SetNowSource(dates.DefaultNowSource)

	uuids.SetGenerator(uuids.NewSeededGenerator(123456))
	dates.SetNowSource(dates.NewFixedNowSource(time.Date(2018, 4, 11, 13, 24, 30, 123456000, time.UTC)))

	session, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyNone)
	require.NoError(t, err)

	run := session.Runs()[0]

	for _, tc := range tests {
		template := fmt.Sprintf("@(json(%s))", tc.path)
		actualJSON, err := run.EvaluateTemplate(template)

		assert.NoError(t, err, "unexpected error evaluating template '%s'", template)

		test.AssertEqualJSON(t, []byte(tc.expected), []byte(actualJSON), "json(...) mismatch for test %s", template)
	}
}

func TestReadWithMissingAssets(t *testing.T) {
	// create standard test session and marshal to JSON
	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	sessionJSON, err := json.Marshal(session)
	require.NoError(t, err)

	// try to read it back but with no assets
	sessionAssets, err := engine.NewSessionAssets(static.NewEmptySource(), nil)

	missingAssets := make([]assets.Reference, 0)
	missing := func(a assets.Reference, err error) { missingAssets = append(missingAssets, a) }

	eng := engine.NewBuilder().Build()
	_, err = eng.ReadSession(sessionAssets, sessionJSON, missing)
	require.NoError(t, err)
	assert.Equal(t, 16, len(missingAssets))
	assert.Equal(t, assets.NewChannelReference(assets.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d"), ""), missingAssets[0])
	assert.Equal(t, assets.NewGroupReference(assets.GroupUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Testers"), missingAssets[1])
	assert.Equal(t, assets.NewGroupReference(assets.GroupUUID("4f1f98fc-27a7-4a69-bbdb-24744ba739a9"), "Males"), missingAssets[2])
	assert.Equal(t, assets.NewFlowReference(assets.FlowUUID("50c3706e-fedb-42c0-8eab-dda3335714b7"), "Registration"), missingAssets[13])
	assert.Equal(t, assets.NewFlowReference(assets.FlowUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Collect Age"), missingAssets[14])
}

func TestResumeWithMissingFlowAssets(t *testing.T) {
	assetsJSON, err := ioutil.ReadFile("../../test/testdata/runner/subflow.json")
	require.NoError(t, err)

	sa, err := test.CreateSessionAssets(assetsJSON, "")
	require.NoError(t, err)

	env := envs.NewBuilder().Build()
	contact := flows.NewEmptyContact(sa, "Bob", envs.NilLanguage, nil)
	trigger := triggers.NewManual(env, assets.NewFlowReference(assets.FlowUUID("76f0a02f-3b75-4b86-9064-e9195e1b3a02"), "Parent Flow"), contact, nil)

	// run session to wait in child flow
	eng := engine.NewBuilder().Build()
	session, _, err := eng.NewSession(sa, trigger)
	require.NoError(t, err)
	assert.Equal(t, flows.SessionStatusWaiting, session.Status())

	// can't directly modify a session's assets but can reload it with different assets
	sessionJSON, err := json.Marshal(session)
	require.NoError(t, err)

	// change the UUID of the child flow so it will effectively be missing
	assetsWithoutChildFlow := test.JSONReplace(assetsJSON, []string{"flows", "[1]", "uuid"}, []byte(`"653a3fa3-ff59-4a89-93c3-a8b9486ec479"`))
	sa, err = test.CreateSessionAssets(assetsWithoutChildFlow, "")
	require.NoError(t, err)

	session, err = eng.ReadSession(sa, sessionJSON, assets.IgnoreMissing)
	require.NoError(t, err)

	_, err = session.Resume(resumes.NewMsg(env, contact, flows.NewMsgIn(flows.MsgUUID(uuids.New()), urns.NilURN, nil, "Hello", nil)))

	// should have an errored session
	assert.NoError(t, err)
	assert.Equal(t, flows.SessionStatusFailed, session.Status())

	// change the UUID of the parent flow so it will effectively be missing
	assetsWithoutParentFlow := test.JSONReplace(assetsJSON, []string{"flows", "[0]", "uuid"}, []byte(`"653a3fa3-ff59-4a89-93c3-a8b9486ec479"`))
	sa, err = test.CreateSessionAssets(assetsWithoutParentFlow, "")
	require.NoError(t, err)

	session, err = eng.ReadSession(sa, sessionJSON, assets.IgnoreMissing)
	require.NoError(t, err)

	_, err = session.Resume(resumes.NewMsg(env, contact, flows.NewMsgIn(flows.MsgUUID(uuids.New()), urns.NilURN, nil, "Hello", nil)))

	// should have an errored session
	assert.NoError(t, err)
	assert.Equal(t, flows.SessionStatusFailed, session.Status())
}

func TestWaitTimeout(t *testing.T) {
	defer dates.SetNowSource(dates.DefaultNowSource)

	t1 := time.Date(2018, 4, 11, 13, 24, 30, 123456000, time.UTC)
	dates.SetNowSource(dates.NewFixedNowSource(t1))

	sessionAssets, err := ioutil.ReadFile("testdata/timeout_test.json")
	require.NoError(t, err)

	// create our session assets
	sa, err := test.CreateSessionAssets(json.RawMessage(sessionAssets), "")
	require.NoError(t, err)

	flow, err := sa.Flows().Get(assets.FlowUUID("76f0a02f-3b75-4b86-9064-e9195e1b3a02"))
	require.NoError(t, err)

	contact := flows.NewEmptyContact(sa, "Joe", "eng", nil)
	contact.AddURN(flows.NewContactURN(urns.URN("tel:+18005555777"), nil))
	trigger := triggers.NewManual(nil, flow.Reference(), contact, nil)

	// create session
	eng := test.NewEngine()
	session, sprint, err := eng.NewSession(sa, trigger)
	require.NoError(t, err)

	require.Equal(t, 1, len(session.Runs()[0].Path()))
	run := session.Runs()[0]

	require.Equal(t, 2, len(sprint.Events()))
	require.Equal(t, "msg_created", sprint.Events()[0].Type())
	require.Equal(t, "msg_wait", sprint.Events()[1].Type())

	// check our wait has a timeout
	waitEvent := run.Events()[1].(*events.MsgWaitEvent)
	require.Equal(t, 600, *waitEvent.TimeoutSeconds)

	_, err = session.Resume(resumes.NewWaitTimeout(nil, nil))
	require.NoError(t, err)

	require.Equal(t, flows.SessionStatusCompleted, session.Status())
	require.Equal(t, 2, len(run.Path()))
	require.Equal(t, 5, len(run.Events()))

	result := run.Results().Get("favorite_color")
	require.Equal(t, "Timeout", result.Category)
	require.Equal(t, "2018-04-11T13:24:30.123456Z", result.Value)
	require.Equal(t, "", result.Input)
}

func TestCurrentContext(t *testing.T) {
	sessionAssets, err := ioutil.ReadFile("../../test/testdata/runner/subflow_loop_with_wait.json")
	require.NoError(t, err)

	// create our session assets
	sa, err := test.CreateSessionAssets(json.RawMessage(sessionAssets), "")
	require.NoError(t, err)

	flow, err := sa.Flows().Get(assets.FlowUUID("76f0a02f-3b75-4b86-9064-e9195e1b3a02"))
	require.NoError(t, err)

	contact := flows.NewEmptyContact(sa, "Joe", "eng", nil)
	trigger := triggers.NewManual(nil, flow.Reference(), contact, nil)

	// create a waiting session
	eng := test.NewEngine()
	session, _, err := eng.NewSession(sa, trigger)
	assert.Equal(t, string(flows.SessionStatusWaiting), string(session.Status()))

	context := session.CurrentContext()
	assert.NotNil(t, context)

	runContext, _ := context.Get("run")
	flowContext, _ := runContext.(*types.XObject).Get("flow")
	flowName, _ := flowContext.(*types.XObject).Get("name")
	assert.Equal(t, types.NewXText("Child flow"), flowName)

	// check we can marshal it
	_, err = json.Marshal(context)
	assert.NoError(t, err)

	// end it
	session.Resume(resumes.NewRunExpiration(nil, nil))
	assert.Equal(t, flows.SessionStatusCompleted, session.Status())

	// can still get context of completed session
	context = session.CurrentContext()
	assert.NotNil(t, context)

	runContext, _ = context.Get("run")
	flowContext, _ = runContext.(*types.XObject).Get("flow")
	flowName, _ = flowContext.(*types.XObject).Get("name")
	assert.Equal(t, types.NewXText("Parent Flow"), flowName)
}
