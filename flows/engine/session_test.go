package engine_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

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
		{"@contact.urns", `["tel:+12065551212?channel=57f1078f-88aa-46f4-a59a-948a5739c03d","twitterid:54784326227#nyaruka","mailto:foo@bar.com"]`, ""},
		{"@contact.urns.tel", `["tel:+12065551212?channel=57f1078f-88aa-46f4-a59a-948a5739c03d"]`, ""},
		{"@contact.urns.xxx", "", "error evaluating @contact.urns.xxx: no such URN scheme 'xxx'"},
		{"@contact.urns.0", "tel:+12065551212?channel=57f1078f-88aa-46f4-a59a-948a5739c03d", ""},
		{"@(contact.urns[0])", "tel:+12065551212?channel=57f1078f-88aa-46f4-a59a-948a5739c03d", ""},
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
		{"@contact.fields", `{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00-02:00"}`, ""},
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

	session, err := test.CreateTestSession("http://localhost", nil)
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
		{"run", `{"contact":{"channel":{"address":"+12345671111","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"created_on":"2018-06-20T11:40:30.123456Z","fields":{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00"},"groups":[{"name":"Testers","uuid":"b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"},{"name":"Males","uuid":"4f1f98fc-27a7-4a69-bbdb-24744ba739a9"}],"language":"eng","name":"Ryan Lewis","timezone":"America/Guayaquil","urns":[{"display":"","path":"+12065551212","scheme":"tel"},{"display":"nyaruka","path":"54784326227","scheme":"twitterid"},{"display":"","path":"foo@bar.com","scheme":"mailto"}],"uuid":"5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"},"created_on":"2018-04-11T13:24:30.123456Z","exited_on":"2018-04-11T13:24:30.123456Z","flow":{"name":"Registration","revision":123,"uuid":"50c3706e-fedb-42c0-8eab-dda3335714b7"},"input":{"attachments":[{"content_type":"image/jpeg","url":"http://s3.amazon.com/bucket/test.jpg"},{"content_type":"audio/mp3","url":"http://s3.amazon.com/bucket/test.mp3"}],"channel":{"address":"+12345671111","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"created_on":"2000-01-01T00:00:00.000000Z","text":"Hi there","type":"msg","urn":{"display":"","path":"+12065551212","scheme":"tel"},"uuid":"9bf91c2b-ce58-4cef-aacc-281e03f69ab5"},"results":{"favorite_color":{"category":"Red","category_localized":"Red","created_on":"2018-04-11T13:24:30.123456Z","input":null,"name":"Favorite Color","node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03","value":"red"},"phone_number":{"category":"","category_localized":"","created_on":"2018-04-11T13:24:30.123456Z","input":null,"name":"Phone Number","node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03","value":"+12344563452"}},"status":"completed","uuid":"d2f852ec-7b4e-457f-ae7f-f8b243c49ff5","webhook":{"json":{"results":[{"state":"WA"},{"state":"IN"}]},"request":"GET /?cmd=echo&content=%7B%22results%22%3A%5B%7B%22state%22%3A%22WA%22%7D%2C%7B%22state%22%3A%22IN%22%7D%5D%7D HTTP/1.1\r\nHost: 127.0.0.1:49992\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n","response":"HTTP/1.1 200 OK\r\nContent-Length: 43\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{\"results\":[{\"state\":\"WA\"},{\"state\":\"IN\"}]}","status":"success","status_code":200,"url":"http://127.0.0.1:49992/?cmd=echo&content=%7B%22results%22%3A%5B%7B%22state%22%3A%22WA%22%7D%2C%7B%22state%22%3A%22IN%22%7D%5D%7D"}}`},
		{"child", `{"contact":{"channel":{"address":"+12345671111","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"created_on":"2018-06-20T11:40:30.123456Z","fields":{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00"},"groups":[{"name":"Testers","uuid":"b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"},{"name":"Males","uuid":"4f1f98fc-27a7-4a69-bbdb-24744ba739a9"}],"language":"eng","name":"Ryan Lewis","timezone":"America/Guayaquil","urns":[{"display":"","path":"+12065551212","scheme":"tel"},{"display":"nyaruka","path":"54784326227","scheme":"twitterid"},{"display":"","path":"foo@bar.com","scheme":"mailto"}],"uuid":"5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"},"flow":{"name":"Collect Age","revision":0,"uuid":"b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"},"results":{"age":{"category":"Youth","category_localized":"Youth","created_on":"2018-04-11T13:24:30.123456Z","input":null,"name":"Age","node_uuid":"d9dba561-b5ee-4f62-ba44-60c4dc242b84","value":"23"}},"status":"completed","uuid":"8720f157-ca1c-432f-9c0b-2014ddc77094"}`},
		{"parent", `{"contact":{"channel":null,"created_on":"0001-01-01T00:00:00.000000Z","fields":{"activation_token":"","age":33,"gender":"Female","join_date":null},"groups":[],"language":"spa","name":"Jasmine","timezone":null,"urns":[],"uuid":"c59b0033-e748-4240-9d4c-e85eb6800151"},"flow":{"name":"Parent","revision":0,"uuid":"fece6eac-9127-4343-9269-56e88f391562"},"results":{"role":{"category":"Reporter","category_localized":"Reporter","created_on":"2000-01-01T00:00:00.000000Z","input":"a reporter","name":"Role","node_uuid":"385cb848-5043-448e-9123-05cbcf26ad74","value":"reporter"}},"status":"active","uuid":"4213ac47-93fd-48c4-af12-7da8218ef09d"}`},
		{"trigger", `{"params":{"source":"website","address":{"state":"WA"}},"type":"flow_action"}`},
	}

	server, err := test.NewTestHTTPServer(49992)
	require.NoError(t, err)

	defer server.Close()
	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)
	defer utils.SetTimeSource(utils.DefaultTimeSource)

	utils.SetUUIDGenerator(utils.NewSeededUUID4Generator(123456))
	utils.SetTimeSource(utils.NewFixedTimeSource(time.Date(2018, 4, 11, 13, 24, 30, 123456000, time.UTC)))

	session, err := test.CreateTestSession(server.URL, nil)
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
	defer utils.SetTimeSource(utils.DefaultTimeSource)

	t1 := time.Date(2018, 4, 11, 13, 24, 30, 123456000, time.UTC)
	t2 := t1.Add(time.Minute * 10)
	utils.SetTimeSource(utils.NewFixedTimeSource(t1))

	sessionAssets, err := ioutil.ReadFile("testdata/timeout_test.json")
	require.NoError(t, err)

	// create our engine session
	session, err := test.CreateSession(json.RawMessage(sessionAssets))
	require.NoError(t, err)

	flow, err := session.Assets().Flows().Get(assets.FlowUUID("76f0a02f-3b75-4b86-9064-e9195e1b3a02"))
	require.NoError(t, err)

	contact := flows.NewEmptyContact("Joe", "eng", nil)
	contact.AddURN(urns.URN("tel:+18005555777"))
	trigger := triggers.NewManualTrigger(nil, contact, flow.Reference(), nil, time.Now())

	err = session.Start(trigger, nil)
	require.NoError(t, err)

	require.Equal(t, 1, len(session.Runs()[0].Path()))
	run := session.Runs()[0]

	require.Equal(t, 2, len(run.Events()))
	require.Equal(t, "msg_created", run.Events()[0].Type())
	require.Equal(t, "msg_wait", run.Events()[1].Type())

	waitEvent := run.Events()[1].(*events.MsgWaitEvent)
	require.Equal(t, &t2, waitEvent.TimeoutOn)
	timeoutOn := *waitEvent.TimeoutOn

	// try to resume without any event - we should remain in waiting state
	session.Resume([]flows.Event{})
	require.NoError(t, err)

	require.Equal(t, flows.SessionStatusWaiting, session.Status())
	require.Equal(t, 1, len(run.Path()))
	require.Equal(t, 2, len(run.Events()))

	// mock our current time to be 10 seconds after the wait times out
	utils.SetTimeSource(utils.NewFixedTimeSource(t2.Add(time.Second * 10)))

	// should be able to resume with a timed out event in the future
	timeoutEvent := events.NewWaitTimedOutEvent()
	timeoutEvent.CreatedOn_ = timeoutOn.Add(time.Second * 60)

	session.Resume([]flows.Event{timeoutEvent})
	require.NoError(t, err)

	require.Equal(t, flows.SessionStatusCompleted, session.Status())
	require.Equal(t, 2, len(run.Path()))
	require.Equal(t, 5, len(run.Events()))

	result := run.Results().Get("favorite_color")
	require.Equal(t, "Timeout", result.Category)
	require.Equal(t, "2018-04-11T13:35:30.123456Z", result.Value)
	require.Nil(t, result.Input)
}
