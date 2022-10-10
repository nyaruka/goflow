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

	assertParse := func(a string, expectedType, expectedURL string) {
		actualType, actualURL := utils.Attachment(a).ToParts()
		assert.Equal(t, expectedType, actualType, "content type mismatch for attachment '%s'", a)
		assert.Equal(t, expectedURL, actualURL, "content type mismatch for attachment '%s'", a)
	}

	assertParse("audio:http://test.m4a", "audio", "http://test.m4a")
	assertParse("audio/mp4:http://test.m4a", "audio/mp4", "http://test.m4a")
	assertParse("audio:/file/path", "audio", "/file/path")

	assertParse("application:http://test.pdf", "application", "http://test.pdf")
	assertParse("application/pdf:http://test.pdf", "application/pdf", "http://test.pdf")

	assertParse("geo:-2.90875,-79.0117686", "geo", "-2.90875,-79.0117686")

	// be lenient with invalid attachments
	assertParse("foo", "", "foo")
	assertParse("http://test.jpg", "", "http://test.jpg")
	assertParse("https://test.jpg", "", "https://test.jpg")
	assertParse("HTTPS://test.jpg", "", "HTTPS://test.jpg")
}
