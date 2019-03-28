package utils_test

import (
	"testing"

	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestAttachment(t *testing.T) {
	attachment := utils.Attachment("image/jpeg:https://example.com/test.jpg")

	assert.Equal(t, "image/jpeg", attachment.ContentType())
	assert.Equal(t, "https://example.com/test.jpg", attachment.URL())

	// be lenient with invalid attachments
	assert.Equal(t, "", utils.Attachment("foo").ContentType())
	assert.Equal(t, "foo", utils.Attachment("foo").URL())
}
