package flows

import (
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
func (a Attachment) Resolve(key string) types.XValue {
	switch key {
	case "content_type":
		return types.NewXString(a.ContentType())
	case "url":
		return types.NewXString(a.URL())
	}

	return types.NewXResolveError(a, key)
}

// Reduce is called when this object needs to be reduced to a primitive
func (a Attachment) Reduce() types.XPrimitive { return types.NewXString(a.URL()) }

func (a Attachment) ToXJSON() types.XString { return types.NewXString("TODO") }

var _ types.XValue = (Attachment)("")
var _ types.XResolvable = (Attachment)("")

// AttachmentList is a list of attachments
type AttachmentList []Attachment

// Index is called when this object is indexed into in an expression
func (a AttachmentList) Index(index int) types.XValue {
	return a[index]
}

// Length is called when the length of this object is requested in an expression
func (a AttachmentList) Length() int {
	return len(a)
}

// Reduce is called when this object needs to be reduced to a primitive
func (a AttachmentList) Reduce() types.XPrimitive {
	array := types.NewXArray()
	for _, attachment := range a {
		array.Append(attachment)
	}
	return array
}

func (a AttachmentList) ToXJSON() types.XString { return types.NewXString("TODO") }

var _ types.XValue = (AttachmentList)(nil)
var _ types.XIndexable = (AttachmentList)(nil)
