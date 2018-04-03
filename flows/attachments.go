package flows

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/excellent/types"
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

// Resolve resolves the given key when this attachment is referenced in an expression
func (a Attachment) Resolve(key string) interface{} {
	switch key {
	case "content_type":
		return a.ContentType()
	case "url":
		return a.URL()
	}

	return fmt.Errorf("No field '%s' on attachment", key)
}

// Atomize is called when this object needs to be reduced to a primitive
func (a Attachment) Atomize() interface{} { return a.URL() }

var _ types.Atomizable = (Attachment)("")
var _ types.Resolvable = (Attachment)("")

// AttachmentList is a list of attachments
type AttachmentList []Attachment

// Index is called when this object is indexed into in an expression
func (a AttachmentList) Index(index int) interface{} {
	return a[index]
}

// Length is called when the length of this object is requested in an expression
func (a AttachmentList) Length() int {
	return len(a)
}

// Atomize is called when this object needs to be reduced to a primitive
func (a AttachmentList) Atomize() interface{} {
	array := types.NewArray()
	for _, attachment := range a {
		array.Append(attachment)
	}
	return array
}

var _ types.Atomizable = (AttachmentList)(nil)
var _ types.Indexable = (AttachmentList)(nil)
