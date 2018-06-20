package engine_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvaluateTemplateAsString(t *testing.T) {
	tests := []struct {
		template string
		expected string
		errorMsg string
	}{
		{"@contact.uuid", "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f", ""},
		{"@contact.id", "1234567", ""},
		{"@contact.name", "Ryan Lewis", ""},
		{"@contact.first_name", "Ryan", ""},
		{"@contact.language", "eng", ""},
		{"@contact.timezone", "America/Guayaquil", ""},
		{"@contact.urns", `["tel:+12065551212","twitterid:54784326227#nyaruka","mailto:foo@bar.com"]`, ""},
		{"@contact.urns.tel", `["tel:+12065551212"]`, ""},
		{"@contact.urns.xxx", "", "error evaluating @contact.urns.xxx: no such URN scheme 'xxx'"},
		{"@contact.urns.0", "tel:+12065551212", ""},
		{"@(contact.urns[0])", "tel:+12065551212", ""},
		{"@(contact.urns[110])", "", "error evaluating @(contact.urns[110]): index 110 out of range for 3 items"},
		{"@contact.urns.0.scheme", "tel", ""},
		{"@contact.urns.0.path", "+12065551212", ""},
		{"@contact.urns.0.display", "", ""},
		{"@contact.urns.0.channel", "My Android Phone", ""},
		{"@contact.urns.0.channel.uuid", "57f1078f-88aa-46f4-a59a-948a5739c03d", ""},
		{"@contact.urns.0.channel.name", "My Android Phone", ""},
		{"@contact.urns.0.channel.address", "+12345671111", ""},
		{"@contact.urns.1", "twitterid:54784326227#nyaruka", ""},
		{"@contact.urns.1.channel", "", ""},
		{"@(format_urn(contact.urns.0))", "(206) 555-1212", ""},
		{"@contact.groups", `["Testers","Males"]`, ""},
		{"@(join(contact.groups, \",\"))", `Testers,Males`, ""},
		{"@(length(contact.groups))", "2", ""},
		{"@contact.fields", `{"activation_token":"AACC55","age":"23","gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00"}`, ""},
		{"@contact.fields.activation_token", "AACC55", ""},
		{"@contact.fields.age", "23", ""},
		{"@contact.fields.join_date", "2017-12-02T00:00:00.000000-02:00", ""},
		{"@contact.fields.favorite_icecream", "", "error evaluating @contact.fields.favorite_icecream: no such contact field 'favorite_icecream'"},
		{"@(is_error(contact.fields.favorite_icecream))", "true", ""},
		{"@(length(contact.fields))", "4", ""},

		{"@run.input", "Hi there\nhttp://s3.amazon.com/bucket/test.jpg\nhttp://s3.amazon.com/bucket/test.mp3", ""},
		{"@run.input.text", "Hi there", ""},
		{"@run.input.attachments", `["http://s3.amazon.com/bucket/test.jpg","http://s3.amazon.com/bucket/test.mp3"]`, ""},
		{"@run.input.attachments.0", "http://s3.amazon.com/bucket/test.jpg", ""},
		{"@run.input.created_on", "2000-01-01T00:00:00.000000Z", ""},
		{"@run.input.channel.name", "My Android Phone", ""},
		{"@run.status", "completed", ""},
		{"@run.results", `{"favorite_color":"red","phone_number":"+12344563452"}`, ""},
		{"@run.results.favorite_color", "red", ""},
		{"@run.results.favorite_color.category", "Red", ""},
		{"@run.results.favorite_icecream", "", "error evaluating @run.results.favorite_icecream: no such run result 'favorite_icecream'"},
		{"@(is_error(run.results.favorite_icecream))", "true", ""},
		{"@(length(run.results))", "2", ""},

		{"@trigger.params", `{"source": "website","address": {"state": "WA"}}`, ""},
		{"@trigger.params.source", "website", ""},
		{"@(length(trigger.params.address))", "1", ""},

		// non-expressions
		{"bob@nyaruka.com", "bob@nyaruka.com", ""},
		{"@twitter_handle", "@twitter_handle", ""},
	}

	session, err := test.CreateTestSession(49995, nil)
	require.NoError(t, err)

	run := session.Runs()[0]

	for _, test := range tests {
		eval, err := run.EvaluateTemplateAsString(test.template, false)

		var actualErrorMsg string
		if err != nil {
			actualErrorMsg = err.Error()
		}

		assert.Equal(t, test.expected, eval, "output mismatch evaluating template: '%s'", test.template)
		assert.Equal(t, test.errorMsg, actualErrorMsg, "error mismatch evaluating template: '%s'", test.template)
	}
}

func TestContextToJSON(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"contact.uuid", `"5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"`},
		{"contact.name", `"Ryan Lewis"`},
		{"contact.urns", `[{"display":"","path":"+12065551212","scheme":"tel"},{"display":"nyaruka","path":"54784326227","scheme":"twitterid"},{"display":"","path":"foo@bar.com","scheme":"mailto"}]`},
		{"contact.urns.0", `{"display":"","path":"+12065551212","scheme":"tel"}`},
		{"contact.fields", `{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00"}`},
		{"contact.fields.age", `23`},
		{"contact", `{"channel":{"address":"+12345671111","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"created_on":"2018-06-20T11:40:30.123456Z","fields":{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00"},"groups":[{"name":"Testers","uuid":"b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"},{"name":"Males","uuid":"4f1f98fc-27a7-4a69-bbdb-24744ba739a9"}],"language":"eng","name":"Ryan Lewis","timezone":"America/Guayaquil","urns":[{"display":"","path":"+12065551212","scheme":"tel"},{"display":"nyaruka","path":"54784326227","scheme":"twitterid"},{"display":"","path":"foo@bar.com","scheme":"mailto"}],"uuid":"5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"}`},
		{"run.input", `{"attachments":[{"content_type":"image/jpeg","url":"http://s3.amazon.com/bucket/test.jpg"},{"content_type":"audio/mp3","url":"http://s3.amazon.com/bucket/test.mp3"}],"channel":{"address":"+12345671111","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"created_on":"2000-01-01T00:00:00.000000Z","text":"Hi there","type":"msg","urn":{"display":"","path":"+12065551212","scheme":"tel"},"uuid":"9bf91c2b-ce58-4cef-aacc-281e03f69ab5"}`},

		// TODO add way to mock call calls to Now() so we can have deterministic tests without doing text substitution of dates?
		//{"run", `{"contact":{"channel":{"address":"+12345671111","name":"Nexmo","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"fields":{"activation_token":"","age":23,"first_name":"Bob","gender":"","joined":"2018-03-27T10:30:00.123456+02:00","state":"Azuay"},"groups":[{"name":"Azuay State","uuid":"d7ff4872-9238-452f-9d38-2f558fea89e0"},{"name":"Survey Audience","uuid":"2aad21f6-30b7-42c5-bd7f-1b720c154817"}],"language":"eng","name":"Ben Haggerty","timezone":"America/Guayaquil","urns":[{"display":"","path":"+12065551212","scheme":"tel"},{"display":"","path":"1122334455667788","scheme":"facebook"},{"display":"","path":"ben@macklemore","scheme":"mailto"}],"uuid":"ba96bf7f-bc2a-4873-a7c7-254d1927c4e3"},"created_on":"2018-04-12T16:46:45.641842Z","exited_on":null,"flow":{"name":"Test Flow","uuid":"76f0a02f-3b75-4b86-9064-e9195e1b3a02"},"input":{"attachments":[{"content_type":"image/jpeg","url":"http://s3.amazon.com/bucket/test_en.jpg?a=Azuay"}],"channel":{"address":"+12345671111","name":"Nexmo","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"created_on":"2000-01-01T00:00:00.000000Z","text":"Hi there","type":"msg","urn":{"display":"","path":"+12065551212","scheme":"tel"},"uuid":"84f8a3cf-0f2c-4881-9502-2d7b114bf01f"},"results":{"favorite_color":{"category":"Red","category_localized":"Red","created_on":"2018-04-12T16:46:45.641921Z","name":"Favorite Color","value":"red"}},"status":"waiting","uuid":"a5c743ed-1373-44c2-905a-3df24d418889","webhook":null}`},
		//{"child", `{"contact":{"channel":{"address":"+12345671111","name":"Nexmo","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"fields":{"activation_token":"","age":23,"first_name":"Bob","gender":"","joined":"2018-03-27T10:30:00.123456+02:00","state":"Azuay"},"groups":[{"name":"Azuay State","uuid":"d7ff4872-9238-452f-9d38-2f558fea89e0"},{"name":"Survey Audience","uuid":"2aad21f6-30b7-42c5-bd7f-1b720c154817"}],"language":"eng","name":"Ben Haggerty","timezone":"America/Guayaquil","urns":[{"display":"","path":"+12065551212","scheme":"tel"},{"display":"","path":"1122334455667788","scheme":"facebook"},{"display":"","path":"ben@macklemore","scheme":"mailto"}],"uuid":"ba96bf7f-bc2a-4873-a7c7-254d1927c4e3"},"flow":{"name":"Collect Language","uuid":"b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"},"results":{},"status":"completed","uuid":"9ec74122-b46f-47de-b232-dc82c51a6808"}`},
		//{"parent", `null`},

		{"trigger", `{"params":{"source":"website","address":{"state":"WA"}},"type":"flow_action"}`},
	}

	session, err := test.CreateTestSession(49996, nil)
	require.NoError(t, err)

	run := session.Runs()[0]

	for _, test := range tests {
		template := fmt.Sprintf("@(json(%s))", test.path)
		eval, err := run.EvaluateTemplateAsString(template, false)

		assert.NoError(t, err, "unexpected error evaluating template '%s'", template)
		assert.Equal(t, test.expected, eval, "json() returned unexpected value for template '%s'", template)
	}
}

func TestWaitTimeout(t *testing.T) {
	sessionAssets, err := ioutil.ReadFile("testdata/timeout_test.json")
	require.NoError(t, err)

	// create our engine session
	session, err := test.CreateSession(json.RawMessage(sessionAssets))
	require.NoError(t, err)

	flow, err := session.Assets().GetFlow(flows.FlowUUID("76f0a02f-3b75-4b86-9064-e9195e1b3a02"))
	require.NoError(t, err)

	contact := flows.NewContact("Joe", "eng", nil)
	contact.AddURN(urns.URN("tel:+18005555777"))
	trigger := triggers.NewManualTrigger(nil, contact, flow, nil, time.Now())

	err = session.Start(trigger, nil)
	require.NoError(t, err)

	require.Equal(t, 1, len(session.Runs()[0].Path()))
	run := session.Runs()[0]

	require.Equal(t, 2, len(run.Events()))
	require.Equal(t, "msg_created", run.Events()[0].Type())
	require.Equal(t, "msg_wait", run.Events()[1].Type())

	waitEvent := run.Events()[1].(*events.MsgWaitEvent)
	require.NotNil(t, waitEvent.TimeoutOn)
	timeoutOn := *waitEvent.TimeoutOn

	// try to resume without any event - we should remain in waiting state
	session.Resume([]flows.Event{})
	require.NoError(t, err)

	require.Equal(t, flows.SessionStatusWaiting, session.Status())
	require.Equal(t, 1, len(run.Path()))
	require.Equal(t, 2, len(run.Events()))

	// mock our current time to be 10 seconds after the wait times out
	testEnv := session.Environment().(*test.TestEnvironment)
	testEnv.SetNow(timeoutOn.Add(time.Second * 10))

	// now we should be able to resume
	timeoutEvent := events.NewWaitTimedOutEvent()
	timeoutEvent.CreatedOn_ = time.Date(2018, 5, 4, 15, 2, 30, 0, time.UTC)

	session.Resume([]flows.Event{timeoutEvent})
	require.NoError(t, err)

	require.Equal(t, flows.SessionStatusCompleted, session.Status())
	require.Equal(t, 2, len(run.Path()))
	require.Equal(t, 5, len(run.Events()))

	result := run.Results().Get("favorite_color")
	require.Equal(t, "Timeout", result.Category)
	require.Equal(t, "2018-05-04T15:02:30.000000Z", result.Value)
	require.Nil(t, result.Input)
}
