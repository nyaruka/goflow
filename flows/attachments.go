package flows

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/utils"
)

// Attachment is a media attachment on a message
type Attachment string

// ContentType returns the MIME type of this attachment
func (a Attachment) ContentType() string {
	offset := strings.Index(string(a), ":")
	if offset >= 0 {
		return string(a[:offset])
	}
	return ""
}

// URL returns the full URL of this attachment
func (a Attachment) URL() string {
	offset := strings.Index(string(a), ":")
	if offset >= 0 {
		return string(a[offset+1:])
	}
	return string(a)
}

func (a Attachment) Resolve(key string) interface{} {
	switch key {

	case "content_type":
		return a.ContentType()

	case "url":
		return a.URL()
	}

	return fmt.Errorf("No field '%s' on attachment", key)
}

func (a Attachment) Default() interface{} { return a }
func (a Attachment) String() string       { return a.URL() }

var _ utils.VariableResolver = (Attachment)("")
