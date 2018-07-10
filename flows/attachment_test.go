package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestAttachment(t *testing.T) {
	env := utils.NewDefaultEnvironment()

	attachment := flows.Attachment("image/jpeg:https://example.com/test.jpg")

	assert.Equal(t, "image/jpeg", attachment.ContentType())
	assert.Equal(t, "https://example.com/test.jpg", attachment.URL())
	assert.Equal(t, "attachment", attachment.Describe())

	assert.Equal(t, types.NewXText("image/jpeg"), attachment.Resolve(env, "content_type"))
	assert.Equal(t, types.NewXText("https://example.com/test.jpg"), attachment.Resolve(env, "url"))
	assert.Equal(t, types.NewXResolveError(attachment, "xxx"), attachment.Resolve(env, "xxx"))
	assert.Equal(t, types.NewXText("https://example.com/test.jpg"), attachment.Reduce(env))
	assert.Equal(t, types.NewXText(`{"content_type":"image/jpeg","url":"https://example.com/test.jpg"}`), attachment.ToXJSON(env))

	// be leniant with invalid attachments
	assert.Equal(t, "", flows.Attachment("foo").ContentType())
	assert.Equal(t, "foo", flows.Attachment("foo").URL())
}

func TestAttachmentList(t *testing.T) {
	env := utils.NewDefaultEnvironment()

	a1 := flows.Attachment("image/jpeg:https://example.com/test.jpg")
	a2 := flows.Attachment("audio/mp3:https://example.com/test.mp3")
	attachments := flows.AttachmentList{a1, a2}

	assert.Equal(t, "attachments", attachments.Describe())
	assert.Equal(t, 2, attachments.Length())
	assert.Equal(t, a2, attachments.Index(1))
	assert.Equal(t, types.NewXArray(a1, a2), attachments.Reduce(env))
	assert.Equal(t, types.NewXText(`[{"content_type":"image/jpeg","url":"https://example.com/test.jpg"},{"content_type":"audio/mp3","url":"https://example.com/test.mp3"}]`), attachments.ToXJSON(env))
}
