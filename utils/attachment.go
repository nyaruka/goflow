package utils

import (
	"regexp"
	"strings"
)

// Attachment is a media attachment on a message in the format <content-type>:<url>. Content type may be a full
// media type or may omit the subtype when it is unknown.
//
// Examples:
//   - image/jpeg:http://s3.amazon.com/bucket/test.jpg
//   - image:http://s3.amazon.com/bucket/test.jpg
type Attachment string

// UnavailableType is the pseudo content type we use for attachments that couldn't be fetched
const UnavailableType = "unavailable"

// we allow outgoing attachments to have types like "image"
var contentTypeRegex = regexp.MustCompile(`^(image|audio|video|application|geo|unavailable|(\w+/[-+.\w]+))$`)

// ToParts splits an attachment string into content-type and URL
func (a Attachment) ToParts() (string, string) {
	offset := strings.Index(string(a), ":")
	if offset >= 0 {
		t, u := strings.ToLower(string(a[:offset])), string(a[offset+1:])
		if contentTypeRegex.MatchString(t) {
			return t, u
		}
	}
	return "", string(a)
}

// ContentType returns the MIME type of this attachment
func (a Attachment) ContentType() string {
	contentType, _ := a.ToParts()
	return contentType
}

// URL returns the full URL of this attachment
func (a Attachment) URL() string {
	_, url := a.ToParts()
	return url
}

// IsValidAttachment returns whether the given string is a valid attachment
func IsValidAttachment(s string) bool {
	typ, url := Attachment(s).ToParts()
	return typ != "" && url != ""
}
