package utils

import "strings"

// Attachment is a media attachment on a message in the format <content-type>:<url>. Content type may be a full
// media type or may omit the subtype when it is unknown.
//
// Examples:
//  - image/jpeg:http://s3.amazon.com/bucket/test.jpg
//  - image:http://s3.amazon.com/bucket/test.jpg
//
type Attachment string

// ToParts splits an attachment string into content-type and URL
func (a Attachment) ToParts() (string, string) {
	offset := strings.Index(string(a), ":")
	if offset >= 0 {
		return string(a[:offset]), string(a[offset+1:])
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
