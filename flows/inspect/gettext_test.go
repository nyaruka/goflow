package inspect_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/inspect"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
)

func TestLocalizableText(t *testing.T) {
	sendMsg := actions.NewSendMsg(
		flows.ActionUUID("7a463f01-2bf4-4ea6-8d7b-3f743d19f27a"),
		"Hi there",
		[]string{"image:https://example.com/test.jpg", "audio:https://example.com/test.mp3"},
		[]string{"Yes", "No"},
		false,
	)

	extracted := make(map[string][]string)

	inspect.LocalizableText(sendMsg, func(uuid uuids.UUID, property string, vals []string, write func([]string)) {
		extracted[property] = vals

		write([]string{"foo", "bar"})
	})

	assert.Equal(t, map[string][]string{
		"attachments":   {"image:https://example.com/test.jpg", "audio:https://example.com/test.mp3"},
		"quick_replies": {"Yes", "No"},
		"text":          {"Hi there"},
	}, extracted)

	data := jsonx.MustMarshal(sendMsg)
	test.AssertEqualJSON(t, []byte(`{
		"uuid": "7a463f01-2bf4-4ea6-8d7b-3f743d19f27a",
		"type": "send_msg",
		"text": "foo",
		"attachments": [
			"foo",
			"bar"
		],
		"quick_replies": [
			"foo",
			"bar"
		]
	}`), data, "JSON mismatch")
}
