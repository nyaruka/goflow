package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRequest struct {
	Trigger *utils.TypedEnvelope   `json:"trigger"`
	Events  []*utils.TypedEnvelope `json:"events"`
}

func TestEvaluateTemplateAsString(t *testing.T) {
	tests := []struct {
		template string
		expected string
		hasError bool
	}{
		{"@contact.uuid", "ba96bf7f-bc2a-4873-a7c7-254d1927c4e3", false},
		{"@contact.name", "Ben Haggerty", false},
		{"@contact.first_name", "Ben", false},
		{"@contact.language", "eng", false},
		{"@contact.timezone", "America/Guayaquil", false},
		{"@contact.urns", `["tel:+12065551212","facebook:1122334455667788","mailto:ben@macklemore"]`, false},
		{"@contact.urns.tel", `["tel:+12065551212"]`, false},
		{"@contact.urns.0", "tel:+12065551212", false},
		{"@(contact.urns[0])", "tel:+12065551212", false},
		{"@contact.urns.0.scheme", "tel", false},
		{"@contact.urns.0.path", "+12065551212", false},
		{"@contact.urns.0.display", "", false},
		{"@contact.urns.0.channel", "Nexmo", false},
		{"@contact.urns.0.channel.uuid", "57f1078f-88aa-46f4-a59a-948a5739c03d", false},
		{"@contact.urns.0.channel.name", "Nexmo", false},
		{"@contact.urns.0.channel.address", "+12345671111", false},
		{"@contact.urns.1", "facebook:1122334455667788", false},
		{"@contact.urns.1.channel", "", false},
		{"@(format_urn(contact.urns.0))", "(206) 555-1212", false},
		{"@contact.groups", `["Azuay State","Survey Audience"]`, false},
		{"@(length(contact.groups))", "2", false},
		{"@contact.fields", `{"activation_token":"","age":"23","first_name":"Bob","gender":"","joined":"2018-03-27T10:30:00.123456+02:00","state":"Azuay"}`, false},
		{"@contact.fields.first_name", "Bob", false},
		{"@contact.fields.age", "23", false},
		{"@contact.fields.joined", "2018-03-27T10:30:00.123456+02:00", false},
		{"@contact.fields.state", "Azuay", false},
		{"@contact.fields.favorite_icecream", "", true},
		{"@(is_error(contact.fields.favorite_icecream))", "true", false},
		{"@(length(contact.fields))", "6", false},

		{"@run.input", "Hi there\nhttp://s3.amazon.com/bucket/test_en.jpg?a=Azuay", false},
		{"@run.input.text", "Hi there", false},
		{"@run.input.attachments", `["http://s3.amazon.com/bucket/test_en.jpg?a=Azuay"]`, false},
		{"@run.input.attachments.0", "http://s3.amazon.com/bucket/test_en.jpg?a=Azuay", false},
		{"@run.input.created_on", "2000-01-01T00:00:00.000000Z", false},
		{"@run.input.channel.name", "Nexmo", false},
		{"@run.status", "waiting", false},
		{"@run.results", `{"favorite_color":"red"}`, false},
		{"@run.results.favorite_color", "red", false},
		{"@run.results.favorite_color.category", "Red", false},
		{"@run.results.favorite_icecream", "", true},
		{"@(is_error(run.results.favorite_icecream))", "true", false},
		{"@(length(run.results))", "1", false},
		{"@run.exited_on", "", false},

		{"@trigger.params", "{\n            \"coupons\": [\n                {\n                    \"code\": \"AAA-BBB-CCC\",\n                    \"expiration\": \"2000-01-01T00:00:00.000000000-00:00\"\n                }\n            ]\n        }", false},
		{"@trigger.params.coupons.0.code", "AAA-BBB-CCC", false},
		{"@(length(trigger.params.coupons))", "1", false},

		// non-expressions
		{"bob@nyaruka.com", "bob@nyaruka.com", false},
		{"@twitter_handle", "@twitter_handle", false},
	}

	run := createTestSession(t)

	for _, test := range tests {
		eval, err := run.EvaluateTemplateAsString(test.template, false)
		if test.hasError {
			assert.Error(t, err, "expected error evaluating template '%s'", test.template)
		} else {
			assert.NoError(t, err, "unexpected error evaluating template '%s'", test.template)
			assert.Equal(t, test.expected, eval, "Actual '%s' does not match expected '%s' evaluating template: '%s'", eval, test.expected, test.template)
		}
	}
}

func TestContextToJSON(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"contact.uuid", `"ba96bf7f-bc2a-4873-a7c7-254d1927c4e3"`},
		{"contact.name", `"Ben Haggerty"`},
		{"contact.urns", `[{"display":"","path":"+12065551212","scheme":"tel"},{"display":"","path":"1122334455667788","scheme":"facebook"},{"display":"","path":"ben@macklemore","scheme":"mailto"}]`},
		{"contact.urns.0", `{"display":"","path":"+12065551212","scheme":"tel"}`},
		{"contact.fields", `{"activation_token":"","age":23,"first_name":"Bob","gender":"","joined":"2018-03-27T10:30:00.123456+02:00","state":"Azuay"}`},
		{"contact.fields.age", `23`},
		{"contact", `{"channel":{"address":"+12345671111","name":"Nexmo","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"fields":{"activation_token":"","age":23,"first_name":"Bob","gender":"","joined":"2018-03-27T10:30:00.123456+02:00","state":"Azuay"},"groups":[{"name":"Azuay State","uuid":"d7ff4872-9238-452f-9d38-2f558fea89e0"},{"name":"Survey Audience","uuid":"2aad21f6-30b7-42c5-bd7f-1b720c154817"}],"language":"eng","name":"Ben Haggerty","timezone":"America/Guayaquil","urns":[{"display":"","path":"+12065551212","scheme":"tel"},{"display":"","path":"1122334455667788","scheme":"facebook"},{"display":"","path":"ben@macklemore","scheme":"mailto"}],"uuid":"ba96bf7f-bc2a-4873-a7c7-254d1927c4e3"}`},
		{"run.input", `{"attachments":[{"content_type":"image/jpeg","url":"http://s3.amazon.com/bucket/test_en.jpg?a=Azuay"}],"channel":{"address":"+12345671111","name":"Nexmo","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"created_on":"2000-01-01T00:00:00.000000Z","text":"Hi there","type":"msg","urn":{"display":"","path":"+12065551212","scheme":"tel"},"uuid":"84f8a3cf-0f2c-4881-9502-2d7b114bf01f"}`},

		// TODO add way to mock call calls to Now() so we can have deterministic tests without doing text substitution of dates?
		//{"run", `{"contact":{"channel":{"address":"+12345671111","name":"Nexmo","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"fields":{"activation_token":"","age":23,"first_name":"Bob","gender":"","joined":"2018-03-27T10:30:00.123456+02:00","state":"Azuay"},"groups":[{"name":"Azuay State","uuid":"d7ff4872-9238-452f-9d38-2f558fea89e0"},{"name":"Survey Audience","uuid":"2aad21f6-30b7-42c5-bd7f-1b720c154817"}],"language":"eng","name":"Ben Haggerty","timezone":"America/Guayaquil","urns":[{"display":"","path":"+12065551212","scheme":"tel"},{"display":"","path":"1122334455667788","scheme":"facebook"},{"display":"","path":"ben@macklemore","scheme":"mailto"}],"uuid":"ba96bf7f-bc2a-4873-a7c7-254d1927c4e3"},"created_on":"2018-04-12T16:46:45.641842Z","exited_on":null,"flow":{"name":"Test Flow","uuid":"76f0a02f-3b75-4b86-9064-e9195e1b3a02"},"input":{"attachments":[{"content_type":"image/jpeg","url":"http://s3.amazon.com/bucket/test_en.jpg?a=Azuay"}],"channel":{"address":"+12345671111","name":"Nexmo","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"created_on":"2000-01-01T00:00:00.000000Z","text":"Hi there","type":"msg","urn":{"display":"","path":"+12065551212","scheme":"tel"},"uuid":"84f8a3cf-0f2c-4881-9502-2d7b114bf01f"},"results":{"favorite_color":{"category":"Red","category_localized":"Red","created_on":"2018-04-12T16:46:45.641921Z","name":"Favorite Color","value":"red"}},"status":"waiting","uuid":"a5c743ed-1373-44c2-905a-3df24d418889","webhook":null}`},
		//{"child", `{"contact":{"channel":{"address":"+12345671111","name":"Nexmo","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"fields":{"activation_token":"","age":23,"first_name":"Bob","gender":"","joined":"2018-03-27T10:30:00.123456+02:00","state":"Azuay"},"groups":[{"name":"Azuay State","uuid":"d7ff4872-9238-452f-9d38-2f558fea89e0"},{"name":"Survey Audience","uuid":"2aad21f6-30b7-42c5-bd7f-1b720c154817"}],"language":"eng","name":"Ben Haggerty","timezone":"America/Guayaquil","urns":[{"display":"","path":"+12065551212","scheme":"tel"},{"display":"","path":"1122334455667788","scheme":"facebook"},{"display":"","path":"ben@macklemore","scheme":"mailto"}],"uuid":"ba96bf7f-bc2a-4873-a7c7-254d1927c4e3"},"flow":{"name":"Collect Language","uuid":"b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"},"results":{},"status":"completed","uuid":"9ec74122-b46f-47de-b232-dc82c51a6808"}`},

		{"parent", `null`},
		{"trigger", `{"params":{"coupons":[{"code":"AAA-BBB-CCC","expiration":"2000-01-01T00:00:00.000000000-00:00"}]},"type":"manual"}`},
	}

	run := createTestSession(t)

	for _, test := range tests {
		template := fmt.Sprintf("@(json(%s))", test.path)
		eval, err := run.EvaluateTemplateAsString(template, false)

		assert.NoError(t, err, "unexpected error evaluating template '%s'", template)
		assert.Equal(t, test.expected, eval, "to_json returned unexpected value for template '%s'", template)
	}
}

func createTestSession(t *testing.T) flows.FlowRun {
	assetsJSON, err := ioutil.ReadFile("testdata/assets.json")
	require.NoError(t, err)

	// build our session
	assetCache := NewAssetCache(100, 5, "testing/1.0")
	err = assetCache.Include(assetsJSON)
	require.NoError(t, err)

	session := NewSession(assetCache, NewMockAssetServer())
	require.NoError(t, err)

	// read trigger from file
	requestJSON, err := ioutil.ReadFile("testdata/trigger.json")
	require.NoError(t, err)

	testRequest := testRequest{}
	err = json.Unmarshal(requestJSON, &testRequest)
	require.NoError(t, err)

	trigger, err := triggers.ReadTrigger(session, testRequest.Trigger)
	require.NoError(t, err)

	initialEvents, err := events.ReadEvents(testRequest.Events)
	require.NoError(t, err)

	session.Start(trigger, initialEvents)

	// check for unexpected errors in the session
	for _, event := range session.Events() {
		fmt.Printf("event %s\n", event.Type())
		require.NotEqual(t, event.Type(), events.TypeError)
	}

	// sanity check session looks as it should
	require.Equal(t, 2, len(session.Runs()))
	run, child := session.Runs()[0], session.Runs()[1]

	require.Equal(t, "Test Flow", run.Flow().Name())
	require.Equal(t, 1, run.Results().Length())
	require.Equal(t, flows.RunStatusWaiting, run.Status())

	require.Equal(t, "Collect Language", child.Flow().Name())
	require.Equal(t, flows.RunStatusCompleted, child.Status())
	return run
}
