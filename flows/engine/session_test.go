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

func TestEvaluateTemplate(t *testing.T) {
	tests := []struct {
		template string
		expected string
		errorMsg string
	}{
		// contact basic properties
		{"@contact.uuid", "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f", ""},
		{"@contact.id", "1234567", ""},
		{"@CONTACT.NAME", "Ryan Lewis", ""},
		{"@contact.name", "Ryan Lewis", ""},
		{"@contact.first_name", "Ryan", ""},
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
		{`@(extract(contact.groups, "name"))`, `[Testers, Males]`, ""},
		{`@(join(extract(contact.groups, "name"), "|"))`, `Testers|Males`, ""},
		{`@(length(contact.groups))`, "2", ""},

		// contact fields
		{"@contact.fields", "{activation_token: AACC55, age: 23, gender: Male, join_date: 2017-12-02T00:00:00.000000-02:00, not_set: }", ""},
		{"@contact.fields.activation_token", "AACC55", ""},
		{"@contact.fields.age", "23", ""},
		{"@contact.fields.join_date", "2017-12-02T00:00:00.000000-02:00", ""},
		{"@contact.fields.favorite_icecream", "", "error evaluating @contact.fields.favorite_icecream: map has no property 'favorite_icecream'"},
		{"@(is_error(contact.fields.favorite_icecream))", "true", ""},
		{"@(length(contact.fields))", "5", ""},

		// simplifed field access
		{"@fields", "{activation_token: AACC55, age: 23, gender: Male, join_date: 2017-12-02T00:00:00.000000-02:00, not_set: }", ""},
		{"@fields.activation_token", "AACC55", ""},
		{"@fields.age", "23", ""},
		{"@fields.join_date", "2017-12-02T00:00:00.000000-02:00", ""},
		{"@fields.favorite_icecream", "", "error evaluating @fields.favorite_icecream: map has no property 'favorite_icecream'"},
		{"@(is_error(fields.favorite_icecream))", "true", ""},
		{"@(length(fields))", "5", ""},

		{"@input", "Hi there\nhttp://s3.amazon.com/bucket/test.jpg\nhttp://s3.amazon.com/bucket/test.mp3", ""},
		{"@input.text", "Hi there", ""},
		{"@input.attachments", `[http://s3.amazon.com/bucket/test.jpg, http://s3.amazon.com/bucket/test.mp3]`, ""},
		{"@(input.attachments[0])", "http://s3.amazon.com/bucket/test.jpg", ""},
		{"@input.created_on", "2017-12-31T11:35:10.035757-02:00", ""},
		{"@input.channel.name", "My Android Phone", ""},

		{"@results", "{2Factor: 34634624463525, Favorite Color: red, Phone Number: +12344563452}", ""},
		{"@results.favorite_color", "red", ""},
		{"@results.favorite_color.category", "Red", ""},
		{"@results.favorite_icecream", "", "error evaluating @results.favorite_icecream: no such run result 'favorite_icecream'"},
		{"@(is_error(results.favorite_icecream))", "true", ""},
		{"@(length(results))", "3", ""},

		{"@run.status", "completed", ""},
		{"@run.results.favorite_color", "red", ""},

		{"@trigger.params", `{"source": "website","address": {"state": "WA"}}`, ""},
		{"@trigger.params.source", "website", ""},
		{"@(length(trigger.params.address))", "1", ""},

		// migrated split by expressions
		{`@(if(is_error(results.favorite_color), "@flow.favorite_color", results.favorite_color))`, `red`, ""},
		{`@(if(is_error(legacy_extra.0.default_city), "@extra.0.default_city", legacy_extra.0.default_city))`, `@extra.0.default_city`, ""},

		// non-expressions
		{"bob@nyaruka.com", "bob@nyaruka.com", ""},
		{"@twitter_handle", "@twitter_handle", ""},
	}

	session, _, err := test.CreateTestSession("http://localhost", nil)
	require.NoError(t, err)

	run := session.Runs()[0]

	for _, test := range tests {
		eval, err := run.EvaluateTemplate(test.template)

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
		{"contact.urns", `["tel:+12065551212","twitterid:54784326227#nyaruka","mailto:foo@bar.com"]`},
		{"contact.urns[0]", `"tel:+12065551212"`},
		{"contact.fields", `{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00","not_set":null}`},
		{"contact.fields.age", `23`},
		{"contact", `{"channel":{"address":"+12345671111","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"created_on":"2018-06-20T11:40:30.123456Z","fields":{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00","not_set":null},"groups":[{"name":"Testers","uuid":"b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"},{"name":"Males","uuid":"4f1f98fc-27a7-4a69-bbdb-24744ba739a9"}],"language":"eng","name":"Ryan Lewis","timezone":"America/Guayaquil","urns":["tel:+12065551212","twitterid:54784326227#nyaruka","mailto:foo@bar.com"],"uuid":"5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"}`},
		{"input", `{"attachments":[{"content_type":"image/jpeg","url":"http://s3.amazon.com/bucket/test.jpg"},{"content_type":"audio/mp3","url":"http://s3.amazon.com/bucket/test.mp3"}],"channel":{"address":"+12345671111","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"created_on":"2017-12-31T11:35:10.035757-02:00","text":"Hi there","type":"msg","urn":"tel:+12065551212","uuid":"9bf91c2b-ce58-4cef-aacc-281e03f69ab5"}`},
		{"run", `{"contact":{"channel":{"address":"+12345671111","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"created_on":"2018-06-20T11:40:30.123456Z","fields":{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00","not_set":null},"groups":[{"name":"Testers","uuid":"b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"},{"name":"Males","uuid":"4f1f98fc-27a7-4a69-bbdb-24744ba739a9"}],"language":"eng","name":"Ryan Lewis","timezone":"America/Guayaquil","urns":["tel:+12065551212","twitterid:54784326227#nyaruka","mailto:foo@bar.com"],"uuid":"5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"},"created_on":"2018-04-11T13:24:30.123456Z","exited_on":"2018-04-11T13:24:30.123456Z","flow":{"name":"Registration","revision":123,"uuid":"50c3706e-fedb-42c0-8eab-dda3335714b7"},"results":{"2factor":{"category":"","category_localized":"","created_on":"2018-04-11T13:24:30.123456Z","input":null,"name":"2Factor","node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03","value":"34634624463525"},"favorite_color":{"category":"Red","category_localized":"Red","created_on":"2018-04-11T13:24:30.123456Z","input":null,"name":"Favorite Color","node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03","value":"red"},"phone_number":{"category":"","category_localized":"","created_on":"2018-04-11T13:24:30.123456Z","input":null,"name":"Phone Number","node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03","value":"+12344563452"},"webhook":{"category":"Success","category_localized":"Success","created_on":"2018-04-11T13:24:30.123456Z","input":"GET http://127.0.0.1:49992/?content=%7B%22results%22%3A%5B%7B%22state%22%3A%22WA%22%7D%2C%7B%22state%22%3A%22IN%22%7D%5D%7D","name":"webhook","node_uuid":"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03","value":"200"}},"status":"completed","uuid":"d2f852ec-7b4e-457f-ae7f-f8b243c49ff5"}`},
		{"child", `{"contact":{"channel":{"address":"+12345671111","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"created_on":"2018-06-20T11:40:30.123456Z","fields":{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00","not_set":null},"groups":[{"name":"Testers","uuid":"b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"},{"name":"Males","uuid":"4f1f98fc-27a7-4a69-bbdb-24744ba739a9"}],"language":"eng","name":"Ryan Lewis","timezone":"America/Guayaquil","urns":["tel:+12065551212","twitterid:54784326227#nyaruka","mailto:foo@bar.com"],"uuid":"5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"},"flow":{"name":"Collect Age","revision":0,"uuid":"b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"},"results":{"age":{"category":"Youth","category_localized":"Youth","created_on":"2018-04-11T13:24:30.123456Z","input":null,"name":"Age","node_uuid":"d9dba561-b5ee-4f62-ba44-60c4dc242b84","value":"23"}},"status":"completed","uuid":"8720f157-ca1c-432f-9c0b-2014ddc77094"}`},
		{"parent", `{"contact":{"channel":{"address":"+12345671111","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"created_on":"2018-01-01T12:00:00.000000Z","fields":{"activation_token":null,"age":33,"gender":"Female","join_date":null,"not_set":null},"groups":[],"language":"spa","name":"Jasmine","timezone":null,"urns":["tel:+593979111222"],"uuid":"c59b0033-e748-4240-9d4c-e85eb6800151"},"flow":{"name":"Parent","revision":0,"uuid":"fece6eac-9127-4343-9269-56e88f391562"},"results":{"role":{"category":"Reporter","category_localized":"Reporter","created_on":"2000-01-01T00:00:00.000000Z","input":"a reporter","name":"Role","node_uuid":"385cb848-5043-448e-9123-05cbcf26ad74","value":"reporter"}},"status":"active","uuid":"4213ac47-93fd-48c4-af12-7da8218ef09d"}`},
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

	for _, test := range tests {
		template := fmt.Sprintf("@(json(%s))", test.path)
		eval, err := run.EvaluateTemplate(template)

		assert.NoError(t, err, "unexpected error evaluating template '%s'", template)
		assert.Equal(t, test.expected, eval, "json() returned unexpected value for template '%s'", template)
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
				"spec_version": "12.0",
				"language": "eng",
				"type": "messaging",
				"revision": 123,
				"nodes": []
			},
			{
				"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
				"name": "Collect Age",
				"spec_version": "12.0",
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
	t2 := t1.Add(time.Minute * 10)
	utils.SetTimeSource(utils.NewFixedTimeSource(t1))

	sessionAssets, err := ioutil.ReadFile("testdata/timeout_test.json")
	require.NoError(t, err)

	// create our engine session
	session, err := test.CreateSession(json.RawMessage(sessionAssets), "")
	require.NoError(t, err)

	flow, err := session.Assets().Flows().Get(assets.FlowUUID("76f0a02f-3b75-4b86-9064-e9195e1b3a02"))
	require.NoError(t, err)

	contact := flows.NewEmptyContact(session.Assets(), "Joe", "eng", nil)
	contact.AddURN(flows.NewContactURN(urns.URN("tel:+18005555777"), nil))
	trigger := triggers.NewManualTrigger(nil, flow.Reference(), contact, nil)

	sprint, err := session.Start(trigger)
	require.NoError(t, err)

	require.Equal(t, 1, len(session.Runs()[0].Path()))
	run := session.Runs()[0]

	require.Equal(t, 2, len(sprint.Events()))
	require.Equal(t, "msg_created", sprint.Events()[0].Type())
	require.Equal(t, "msg_wait", sprint.Events()[1].Type())

	// check that our timeout is 10 minutes in the future
	waitEvent := run.Events()[1].(*events.MsgWaitEvent)
	require.Equal(t, &t2, waitEvent.TimeoutOn)

	// should fail with error event if we try to timeout immediately
	sprint, err = session.Resume(resumes.NewWaitTimeoutResume(nil, nil))
	require.NoError(t, err)
	require.Equal(t, 1, len(sprint.Events()))
	require.Equal(t, "error", sprint.Events()[0].Type())

	// mock our current time to be 10 seconds after the wait times out
	utils.SetTimeSource(utils.NewFixedTimeSource(t2.Add(time.Second * 10)))

	_, err = session.Resume(resumes.NewWaitTimeoutResume(nil, nil))
	require.NoError(t, err)

	require.Equal(t, flows.SessionStatusCompleted, session.Status())
	require.Equal(t, 2, len(run.Path()))
	require.Equal(t, 5, len(run.Events()))

	result := run.Results().Get("favorite_color")
	require.Equal(t, "Timeout", result.Category)
	require.Equal(t, "2018-04-11T13:34:40.123456Z", result.Value)
	require.Nil(t, result.Input)
}
