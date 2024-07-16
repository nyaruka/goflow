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

	assertParse := func(a string, expectedType, expectedURL string, isValid bool) {
		actualType, actualURL := utils.Attachment(a).ToParts()
		assert.Equal(t, expectedType, actualType, "content type mismatch for attachment '%s'", a)
		assert.Equal(t, expectedURL, actualURL, "content type mismatch for attachment '%s'", a)
		assert.Equal(t, isValid, utils.IsValidAttachment(a), "is valid mismatch for attachment '%s'", a)
	}

	assertParse("audio:http://test.m4a", "audio", "http://test.m4a", true)
	assertParse("audio/mp4:http://test.m4a", "audio/mp4", "http://test.m4a", true)
	assertParse("audio:/file/path", "audio", "/file/path", true)

	assertParse("application:http://test.pdf", "application", "http://test.pdf", true)
	assertParse("application/pdf:http://test.pdf", "application/pdf", "http://test.pdf", true)

	assertParse("geo:-2.90875,-79.0117686", "geo", "-2.90875,-79.0117686", true)

	assertParse("unavailable:http://bad.link", "unavailable", "http://bad.link", true)

	// be lenient with invalid attachments
	assertParse("", "", "", false)
	assertParse("foo", "", "foo", false)
	assertParse("http://test.jpg", "", "http://test.jpg", false)
	assertParse("https://test.jpg", "", "https://test.jpg", false)
	assertParse("HTTPS://test.jpg", "", "HTTPS://test.jpg", false)
	assertParse(":http://test.jpg", "", ":http://test.jpg", false)
}
