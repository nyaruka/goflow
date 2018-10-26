package flows

import (
	"strings"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// Attachment is a media attachment on a message, and it has the following properties which can be accessed:
//
//  * `content_type` the MIME type of the attachment
//  * `url` the URL of the attachment
//
// Examples:
//
//   @input.attachments.0.content_type -> image/jpeg
//   @input.attachments.0.url -> http://s3.amazon.com/bucket/test.jpg
//   @(json(input.attachments.0)) -> {"content_type":"image/jpeg","url":"http://s3.amazon.com/bucket/test.jpg"}
//
// @context attachment
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
func (a Attachment) Resolve(env utils.Environment, key string) types.XValue {
	switch strings.ToLower(key) {
	case "content_type":
		return types.NewXText(a.ContentType())
	case "url":
		return types.NewXText(a.URL())
	}

	return types.NewXResolveError(a, key)
}

// Describe returns a representation of this type for error messages
func (a Attachment) Describe() string { return "attachment" }

// Reduce is called when this object needs to be reduced to a primitive
func (a Attachment) Reduce(env utils.Environment) types.XPrimitive { return types.NewXText(a.URL()) }

// ToXJSON is called when this type is passed to @(json(...))
func (a Attachment) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, a, "content_type", "url").ToXJSON(env)
}

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

// Describe returns a representation of this type for error messages
func (a AttachmentList) Describe() string { return "attachments" }

// Reduce is called when this object needs to be reduced to a primitive
func (a AttachmentList) Reduce(env utils.Environment) types.XPrimitive {
	array := types.NewXArray()
	for _, attachment := range a {
		array.Append(attachment)
	}
	return array
}

// ToXJSON is called when this type is passed to @(json(...))
func (a AttachmentList) ToXJSON(env utils.Environment) types.XText { return a.Reduce(env).ToXJSON(env) }

var _ types.XValue = (AttachmentList)(nil)
var _ types.XIndexable = (AttachmentList)(nil)
