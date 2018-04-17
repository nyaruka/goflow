package engine_test

import (
	"fmt"
	"testing"

	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvaluateTemplateAsString(t *testing.T) {
	tests := []struct {
		template string
		expected string
		hasError bool
	}{
		{"@contact.uuid", "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f", false},
		{"@contact.name", "Ryan Lewis", false},
		{"@contact.first_name", "Ryan", false},
		{"@contact.language", "eng", false},
		{"@contact.timezone", "", false},
		{"@contact.urns", `["tel:+12065551212","twitterid:54784326227#nyaruka","mailto:foo@bar.com"]`, false},
		{"@contact.urns.tel", `["tel:+12065551212"]`, false},
		{"@contact.urns.0", "tel:+12065551212", false},
		{"@(contact.urns[0])", "tel:+12065551212", false},
		{"@contact.urns.0.scheme", "tel", false},
		{"@contact.urns.0.path", "+12065551212", false},
		{"@contact.urns.0.display", "", false},
		{"@contact.urns.0.channel", "My Android Phone", false},
		{"@contact.urns.0.channel.uuid", "57f1078f-88aa-46f4-a59a-948a5739c03d", false},
		{"@contact.urns.0.channel.name", "My Android Phone", false},
		{"@contact.urns.0.channel.address", "+12345671111", false},
		{"@contact.urns.1", "twitterid:54784326227#nyaruka", false},
		{"@contact.urns.1.channel", "", false},
		{"@(format_urn(contact.urns.0))", "(206) 555-1212", false},
		{"@contact.groups", `["Testers","Males"]`, false},
		{"@(length(contact.groups))", "2", false},
		{"@contact.fields", `{"activation_token":"AACC55","age":"23","gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00"}`, false},
		{"@contact.fields.activation_token", "AACC55", false},
		{"@contact.fields.age", "23", false},
		{"@contact.fields.join_date", "2017-12-02T00:00:00.000000-02:00", false},
		{"@contact.fields.favorite_icecream", "", true},
		{"@(is_error(contact.fields.favorite_icecream))", "true", false},
		{"@(length(contact.fields))", "4", false},

		{"@run.input", "Hi there\nhttp://s3.amazon.com/bucket/test.jpg\nhttp://s3.amazon.com/bucket/test.mp3", false},
		{"@run.input.text", "Hi there", false},
		{"@run.input.attachments", `["http://s3.amazon.com/bucket/test.jpg","http://s3.amazon.com/bucket/test.mp3"]`, false},
		{"@run.input.attachments.0", "http://s3.amazon.com/bucket/test.jpg", false},
		{"@run.input.created_on", "2000-01-01T00:00:00.000000Z", false},
		{"@run.input.channel.name", "My Android Phone", false},
		{"@run.status", "completed", false},
		{"@run.results", `{"favorite_color":"red","phone_number":"+12344563452"}`, false},
		{"@run.results.favorite_color", "red", false},
		{"@run.results.favorite_color.category", "Red", false},
		{"@run.results.favorite_icecream", "", true},
		{"@(is_error(run.results.favorite_icecream))", "true", false},
		{"@(length(run.results))", "2", false},

		{"@trigger.params", `{"source": "website","address": {"state": "WA"}}`, false},
		{"@trigger.params.source", "website", false},
		{"@(length(trigger.params.address))", "1", false},

		// non-expressions
		{"bob@nyaruka.com", "bob@nyaruka.com", false},
		{"@twitter_handle", "@twitter_handle", false},
	}

	session, err := test.CreateTestSession(49995, nil)
	require.NoError(t, err)

	run := session.Runs()[0]

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
		{"contact.uuid", `"5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"`},
		{"contact.name", `"Ryan Lewis"`},
		{"contact.urns", `[{"display":"","path":"+12065551212","scheme":"tel"},{"display":"nyaruka","path":"54784326227","scheme":"twitterid"},{"display":"","path":"foo@bar.com","scheme":"mailto"}]`},
		{"contact.urns.0", `{"display":"","path":"+12065551212","scheme":"tel"}`},
		{"contact.fields", `{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00"}`},
		{"contact.fields.age", `23`},
		{"contact", `{"channel":{"address":"+12345671111","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"fields":{"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00"},"groups":[{"name":"Testers","uuid":"b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"},{"name":"Males","uuid":"4f1f98fc-27a7-4a69-bbdb-24744ba739a9"}],"language":"eng","name":"Ryan Lewis","timezone":null,"urns":[{"display":"","path":"+12065551212","scheme":"tel"},{"display":"nyaruka","path":"54784326227","scheme":"twitterid"},{"display":"","path":"foo@bar.com","scheme":"mailto"}],"uuid":"5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"}`},
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
