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

var templateTests = []struct {
	template string
	expected string
	errorMsg string
}{
	// contact basic properties
	{"@contact.uuid", "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f", ""},
	{"@contact.id", "1234567", ""},
	{"@CONTACT.NAME", "Ryan Lewis", ""},
	{"@contact.name", "Ryan Lewis", ""},
	{"@contact.language", "eng", ""},
	{"@contact.timezone", "America/Guayaquil", ""},

	// contact single URN access
	{"@contact.urn", `tel:+12065551212`, ""},
	{"@(urn_parts(contact.urn).scheme)", `tel`, ""},
	{"@(urn_parts(contact.urn).path)", `+12065551212`, ""},
	{"@(format_urn(contact.urn))", `(206) 555-1212`, ""},

	// contact URN list access
	{"@contact.urns", `[tel:+12065551212, twitterid:54784326227#nyaruka, mailto:foo@bar.com]`, ""},
	{"@(contact.urns[0])", "tel:+12065551212", ""},
	{"@(contact.urns[110])", "", "error evaluating @(contact.urns[110]): index 110 out of range for 3 items"},
	{"@(urn_parts(contact.urns[0]).scheme)", "tel", ""},
	{"@(urn_parts(contact.urns[0]).path)", "+12065551212", ""},
	{"@(urn_parts(contact.urns[0]).display)", "", ""},
	{"@(contact.urns[1])", "twitterid:54784326227#nyaruka", ""},
	{"@(format_urn(contact.urns[0]))", "(206) 555-1212", ""},

	// simplified URN access
	{"@urns", `{ext: , facebook: , fcm: , jiochat: , line: , mailto: mailto:foo@bar.com, tel: tel:+12065551212, telegram: , twitter: , twitterid: twitterid:54784326227#nyaruka, viber: , wechat: , whatsapp: }`, ""},
	{"@urns.tel", `tel:+12065551212`, ""},
	{"@urns.mailto", `mailto:foo@bar.com`, ""},
	{"@urns.viber", ``, ""},
	{"@(format_urn(urns.tel))", "(206) 555-1212", ""},

	// contact groups
	{`@(foreach(contact.groups, extract, "name"))`, `[Testers, Males]`, ""},
	{`@(join(foreach(contact.groups, extract, "name"), "|"))`, `Testers|Males`, ""},
	{`@(count(contact.groups))`, "2", ""},

	// contact fields
	{"@contact.fields", "Activation Token: AACC55\nAge: 23\nGender: Male\nJoin Date: 2017-12-02T00:00:00.000000-02:00", ""},
	{"@contact.fields.activation_token", "AACC55", ""},
	{"@contact.fields.age", "23", ""},
	{"@contact.fields.join_date", "2017-12-02T00:00:00.000000-02:00", ""},
	{"@contact.fields.favorite_icecream", "", "error evaluating @contact.fields.favorite_icecream: object has no property 'favorite_icecream'"},
	{"@(is_error(contact.fields.favorite_icecream))", "true", ""},
	{"@(has_error(contact.fields.favorite_icecream).match)", "object has no property 'favorite_icecream'", ""},
	{"@(count(contact.fields))", "5", ""},

	// simplifed field access
	{"@fields", "Activation Token: AACC55\nAge: 23\nGender: Male\nJoin Date: 2017-12-02T00:00:00.000000-02:00", ""},
	{"@fields.activation_token", "AACC55", ""},
	{"@fields.age", "23", ""},
	{"@fields.join_date", "2017-12-02T00:00:00.000000-02:00", ""},
	{"@fields.favorite_icecream", "", "error evaluating @fields.favorite_icecream: object has no property 'favorite_icecream'"},
	{"@(is_error(fields.favorite_icecream))", "true", ""},
	{"@(has_error(fields.favorite_icecream).match)", "object has no property 'favorite_icecream'", ""},
	{"@(count(fields))", "5", ""},

	{"@input", "Hi there\nhttp://s3.amazon.com/bucket/test.jpg\nhttp://s3.amazon.com/bucket/test.mp3", ""},
	{"@input.text", "Hi there", ""},
	{"@input.attachments", `[image/jpeg:http://s3.amazon.com/bucket/test.jpg, audio/mp3:http://s3.amazon.com/bucket/test.mp3]`, ""},
	{"@(input.attachments[0])", "image/jpeg:http://s3.amazon.com/bucket/test.jpg", ""},
	{"@input.created_on", "2017-12-31T11:35:10.035757-02:00", ""},
	{"@input.channel.name", "My Android Phone", ""},

	{"@results.favorite_color", `red`, ""},
	{"@results.favorite_color.value", "red", ""},
	{"@results.favorite_color.category", "Red", ""},
	{"@results.favorite_color.category_localized", "Red", ""},
	{"@(is_error(results.favorite_icecream))", "true", ""},
	{"@(has_error(results.favorite_icecream).match)", "object has no property 'favorite_icecream'", ""},
	{"@(count(results))", "4", ""},

	{"@run.results.favorite_color", `[red]`, ""},
	{"@run.results.favorite_color.values", "[red]", ""},
	{`@(run.results.favorite_color.values[0])`, `red`, ""},
	{"@run.results.favorite_color.categories", "[Red]", ""},
	{`@(run.results.favorite_color.categories[0])`, `Red`, ""},
	{"@run.results.favorite_icecream", "", "error evaluating @run.results.favorite_icecream: object has no property 'favorite_icecream'"},
	{"@(is_error(run.results.favorite_icecream))", "true", ""},
	{"@(has_error(run.results.favorite_icecream).match)", "object has no property 'favorite_icecream'", ""},
	{"@(count(run.results))", "4", ""},

	{"@run.status", "completed", ""},

	{"@webhook", "{results: [{state: WA}, {state: IN}]}", ""},
	{"@webhook.results", "[{state: WA}, {state: IN}]", ""},
	{"@(webhook.results[1])", "{state: IN}", ""},
	{"@(webhook.results[1].state)", "IN", ""},

	{"@trigger.params", `{address: {state: WA}, source: website}`, ""},
	{"@trigger.params.source", "website", ""},
	{"@(count(trigger.params.address))", "1", ""},

	// migrated split by expressions
	{`@(if(is_error(results.favorite_color.value), "@flow.favorite_color", results.favorite_color.value))`, `red`, ""},
	{`@(if(is_error(legacy_extra["0"].default_city), "@extra.0.default_city", legacy_extra["0"].default_city))`, `@extra.0.default_city`, ""},

	// non-expressions
	{"bob@nyaruka.com", "bob@nyaruka.com", ""},
	{"@twitter_handle", "@twitter_handle", ""},
}

func TestEvaluateTemplate(t *testing.T) {
	utils.SetTimeSource(utils.NewFixedTimeSource(time.Date(2018, 9, 13, 13, 36, 30, 123456789, time.UTC)))
	defer utils.SetTimeSource(utils.DefaultTimeSource)

	server := test.NewTestHTTPServer(0)
	defer server.Close()

	session, _, err := test.CreateTestSession(server.URL, nil)
	require.NoError(t, err)

	run := session.Runs()[0]

	for _, tc := range templateTests {
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
	session, _, err := test.CreateTestSession("http://localhost", nil)
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
		{"contact.urns", `["tel:+12065551212","twitterid:54784326227#nyaruka","mailto:foo@bar.com"]`},
		{"contact.urns[0]", `"tel:+12065551212"`},
		{"contact.fields", `{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00","not_set":null}`},
		{"contact.fields.age", `23`},
		{
			"contact",
			`{
				"channel": {"address":"+12345671111","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},
				"created_on": "2018-06-20T11:40:30.123456Z", 
				"fields": {"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00","not_set":null},
				"first_name": "Ryan",
				"groups": [{"name":"Testers","uuid":"b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"},{"name":"Males","uuid":"4f1f98fc-27a7-4a69-bbdb-24744ba739a9"}],
				"id": "1234567",
				"language": "eng",
				"name": "Ryan Lewis",
				"timezone": "America/Guayaquil",
				"urn": "tel:+12065551212",
				"urns": ["tel:+12065551212","twitterid:54784326227#nyaruka","mailto:foo@bar.com"],
				"uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"
			}`,
		},
		{
			"input",
			`{
				"attachments":["image/jpeg:http://s3.amazon.com/bucket/test.jpg","audio/mp3:http://s3.amazon.com/bucket/test.mp3"],
				"channel":{"address":"+12345671111","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},
				"created_on":"2017-12-31T11:35:10.035757-02:00",
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
					"channel":{"address":"+12345671111","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},
					"created_on":"2018-06-20T11:40:30.123456Z",
					"fields":{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00","not_set":null},
					"first_name":"Ryan",
					"groups":[{"name":"Testers","uuid":"b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"},{"name":"Males","uuid":"4f1f98fc-27a7-4a69-bbdb-24744ba739a9"}],
					"id":"1234567",
					"language":"eng",
					"name":"Ryan Lewis",
					"timezone":"America/Guayaquil",
					"urn":"tel:+12065551212",
					"urns":["tel:+12065551212","twitterid:54784326227#nyaruka","mailto:foo@bar.com"],
					"uuid":"5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"
				},
				"created_on":"2018-04-11T13:24:30.123456Z",
				"exited_on":"2018-04-11T13:24:30.123456Z",
				"flow":{"name":"Registration","revision":123,"uuid":"50c3706e-fedb-42c0-8eab-dda3335714b7"},
				"path":[
					{"arrived_on":"2018-04-11T13:24:30.123456Z","exit_uuid":"d7a36118-0a38-4b35-a7e4-ae89042f0d3c","node_uuid":"72a1f5df-49f9-45df-94c9-d86f7ea064e5","uuid":"692926ea-09d6-4942-bd38-d266ec8d3716"},
					{"arrived_on":"2018-04-11T13:24:30.123456Z","exit_uuid":"100f2d68-2481-4137-a0a3-177620ba3c5f","node_uuid":"3dcccbb4-d29c-41dd-a01f-16d814c9ab82","uuid":"5802813d-6c58-4292-8228-9728778b6c98"},
					{"arrived_on":"2018-04-11T13:24:30.123456Z","exit_uuid":"d898f9a4-f0fc-4ac4-a639-c98c602bb511","node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03","uuid":"970b8069-50f5-4f6f-8f41-6b2d9f33d623"},
					{"arrived_on":"2018-04-11T13:24:30.123456Z","exit_uuid":"9fc5f8b4-2247-43db-b899-ab1ac50ba06c","node_uuid":"c0781400-737f-4940-9a6c-1ec1c3df0325","uuid":"5ecda5fc-951c-437b-a17e-f85e49829fb9"}
				],
				"results":{
					"2factor":{
						"categories":[""],
						"categories_localized":[""],
						"created_on":"2018-04-11T13:24:30.123456Z",
						"extra":null,
						"input":"",
						"name":"2Factor",
						"node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
						"values":["34634624463525"]
					},
					"favorite_color":{
						"categories":["Red"],
						"categories_localized":["Red"],
						"created_on":"2018-04-11T13:24:30.123456Z",
						"extra":null,
						"input":"",
						"name":"Favorite Color",
						"node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
						"values":["red"]
					},
					"phone_number":{
						"categories":[""],
						"categories_localized":[""],
						"created_on":"2018-04-11T13:24:30.123456Z",
						"extra":null,
						"input":"",
						"name":"Phone Number",
						"node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
						"values":["+12344563452"]
					},
					"webhook":{
						"categories":["Success"],
						"categories_localized":["Success"],
						"created_on":"2018-04-11T13:24:30.123456Z",
						"extra":{"results":[{"state":"WA"},{"state":"IN"}]},
						"input":"GET http://127.0.0.1:49992/?content=%7B%22results%22%3A%5B%7B%22state%22%3A%22WA%22%7D%2C%7B%22state%22%3A%22IN%22%7D%5D%7D",
						"name":"webhook",
						"node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
						"values":["200"]
					}
				},
				"status":"completed",
				"uuid":"d2f852ec-7b4e-457f-ae7f-f8b243c49ff5"
			}`,
		},
		{
			"child",
			`{
				"contact": {
					"channel": {
						"address": "+12345671111",
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
					"urn": "tel:+12065551212",
					"urns": [
						"tel:+12065551212",
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
				"results": {
					"age": {
						"category": "Youth",
						"category_localized": "Youth",
						"created_on": "2018-04-11T13:24:30.123456Z",
						"input": "",
						"name": "Age",
						"node_uuid": "d9dba561-b5ee-4f62-ba44-60c4dc242b84",
						"value": "23"
					}
				},
				"run": {
					"contact": {
						"channel": {
							"address": "+12345671111",
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
						"urn": "tel:+12065551212",
						"urns": [
							"tel:+12065551212",
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
							"categories": [
								"Youth"
							],
							"categories_localized": [
								"Youth"
							],
							"created_on": "2018-04-11T13:24:30.123456Z",
							"extra": null,
							"input": "",
							"name": "Age",
							"node_uuid": "d9dba561-b5ee-4f62-ba44-60c4dc242b84",
							"values": [
								"23"
							]
						}
					},
					"status": "completed",
					"uuid": "8720f157-ca1c-432f-9c0b-2014ddc77094"
				},
				"urns": {
					"ext": null,
					"facebook": null,
					"fcm": null,
					"jiochat": null,
					"line": null,
					"mailto": "mailto:foo@bar.com",
					"tel": "tel:+12065551212",
					"telegram": null,
					"twitter": null,
					"twitterid": "twitterid:54784326227#nyaruka",
					"viber": null,
					"wechat": null,
					"whatsapp": null
				}
			}`,
		},
		{
			"parent",
			`{
				"contact": {
					"channel": {
						"address": "+12345671111",
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
					"urn": "tel:+593979111222",
					"urns": [
						"tel:+593979111222"
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
				"results": {
					"role": {
						"category": "Reporter",
						"category_localized": "Reporter",
						"created_on": "2000-01-01T00:00:00.000000Z",
						"input": "a reporter",
						"name": "Role",
						"node_uuid": "385cb848-5043-448e-9123-05cbcf26ad74",
						"value": "reporter"
					}
				},
				"run": {
					"contact": {
						"channel": {
							"address": "+12345671111",
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
						"urn": "tel:+593979111222",
						"urns": [
							"tel:+593979111222"
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
							"categories": [
								"Reporter"
							],
							"categories_localized": [
								"Reporter"
							],
							"created_on": "2000-01-01T00:00:00.000000Z",
							"extra": null,
							"input": "a reporter",
							"name": "Role",
							"node_uuid": "385cb848-5043-448e-9123-05cbcf26ad74",
							"values": [
								"reporter"
							]
						}
					},
					"status": "active",
					"uuid": "4213ac47-93fd-48c4-af12-7da8218ef09d"
				},
				"urns": {
					"ext": null,
					"facebook": null,
					"fcm": null,
					"jiochat": null,
					"line": null,
					"mailto": null,
					"tel": "tel:+593979111222",
					"telegram": null,
					"twitter": null,
					"twitterid": null,
					"viber": null,
					"wechat": null,
					"whatsapp": null
				}
			}`,
		},
		{"trigger", `{"params":{"source":"website","address":{"state":"WA"}},"type":"flow_action"}`},
	}

	server := test.NewTestHTTPServer(49992)
	defer server.Close()
	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)
	defer utils.SetTimeSource(utils.DefaultTimeSource)

	utils.SetUUIDGenerator(utils.NewSeededUUID4Generator(123456))
	utils.SetTimeSource(utils.NewFixedTimeSource(time.Date(2018, 4, 11, 13, 24, 30, 123456000, time.UTC)))

	session, _, err := test.CreateTestSession(server.URL, nil)
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
	session, _, err := test.CreateTestSession("", nil)
	require.NoError(t, err)

	sessionJSON, err := json.Marshal(session)
	require.NoError(t, err)

	// try to read it back but with only the flow assets
	source, err := static.NewSource([]byte(`{
		"flows": [
			{
				"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
				"name": "Registration",
				"spec_version": "13.0",
				"language": "eng",
				"type": "messaging",
				"revision": 123,
				"nodes": []
			},
			{
				"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
				"name": "Collect Age",
				"spec_version": "13.0",
				"language": "eng",
				"type": "messaging",
				"nodes": []
			}
		]
	}`))
	require.NoError(t, err)
	sessionAssets, err := engine.NewSessionAssets(source)

	missingAssets := make([]assets.Reference, 0)
	missing := func(a assets.Reference) { missingAssets = append(missingAssets, a) }

	eng := engine.NewBuilder().WithDefaultUserAgent("test").Build()
	_, err = eng.ReadSession(sessionAssets, sessionJSON, missing)
	require.NoError(t, err)
	assert.Equal(t, 14, len(missingAssets))
	assert.Equal(t, assets.NewChannelReference(assets.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d"), ""), missingAssets[0])
	assert.Equal(t, assets.NewGroupReference(assets.GroupUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Testers"), missingAssets[1])
	assert.Equal(t, assets.NewGroupReference(assets.GroupUUID("4f1f98fc-27a7-4a69-bbdb-24744ba739a9"), "Males"), missingAssets[2])

	// still get error if we're missing flow assets
	emptyAssets, err := engine.NewSessionAssets(static.NewEmptySource())
	require.NoError(t, err)

	_, err = eng.ReadSession(emptyAssets, sessionJSON, missing)
	assert.EqualError(t, err, "unable to read run 0: unable to load flow[uuid=50c3706e-fedb-42c0-8eab-dda3335714b7,name=Registration]: no such flow with UUID '50c3706e-fedb-42c0-8eab-dda3335714b7'")
}

func TestWaitTimeout(t *testing.T) {
	defer utils.SetTimeSource(utils.DefaultTimeSource)

	t1 := time.Date(2018, 4, 11, 13, 24, 30, 123456000, time.UTC)
	utils.SetTimeSource(utils.NewFixedTimeSource(t1))

	sessionAssets, err := ioutil.ReadFile("testdata/timeout_test.json")
	require.NoError(t, err)

	// create our session assets
	sa, err := test.CreateSessionAssets(json.RawMessage(sessionAssets), "")
	require.NoError(t, err)

	flow, err := sa.Flows().Get(assets.FlowUUID("76f0a02f-3b75-4b86-9064-e9195e1b3a02"))
	require.NoError(t, err)

	contact := flows.NewEmptyContact(sa, "Joe", "eng", nil)
	contact.AddURN(flows.NewContactURN(urns.URN("tel:+18005555777"), nil))
	trigger := triggers.NewManualTrigger(nil, flow.Reference(), contact, nil)

	// create session
	eng := engine.NewBuilder().WithDefaultUserAgent("goflow-testing").Build()
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

	_, err = session.Resume(resumes.NewWaitTimeoutResume(nil, nil))
	require.NoError(t, err)

	require.Equal(t, flows.SessionStatusCompleted, session.Status())
	require.Equal(t, 2, len(run.Path()))
	require.Equal(t, 5, len(run.Events()))

	result := run.Results().Get("favorite_color")
	require.Equal(t, "Timeout", result.Category)
	require.Equal(t, "2018-04-11T13:24:30.123456Z", result.Value)
	require.Equal(t, "", result.Input)
}
