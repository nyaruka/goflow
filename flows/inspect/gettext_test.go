package inspect_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/inspect"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/stretchr/testify/assert"
)

func TestLocalizedText(t *testing.T) {
	sendMsg := actions.NewSendMsg(
		flows.ActionUUID("7a463f01-2bf4-4ea6-8d7b-3f743d19f27a"),
		"Hi there",
		[]string{"image:https://example.com/test.jpg", "audio:https://example.com/test.mp3"},
		[]string{"Yes", "No"},
		false,
	)

	extracted := make(map[string][]string)

	inspect.LocalizedText(sendMsg, func(uuid uuids.UUID, property string, translated []string) {
		extracted[property] = translated
	})

	assert.Equal(t, map[string][]string{
		"attachments":   []string{"image:https://example.com/test.jpg", "audio:https://example.com/test.mp3"},
		"quick_replies": []string{"Yes", "No"},
		"text":          []string{"Hi there"},
	}, extracted)
}
