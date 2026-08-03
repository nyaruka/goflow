package events

import (
	"fmt"

	"github.com/nyaruka/goflow/assets"
)

func init() {
	registerType(TypeError, func() Event { return &Error{} })
}

// TypeError is the type of our error events
const TypeError string = "error"

const (
	ErrorCodeActionUnsupported    = "action:unsupported"
	ErrorCodeDependencyMissing    = "dependency:missing"
	ErrorCodeGroupMissing         = "group:missing"
	ErrorCodeLabelMissing         = "label:missing"
	ErrorCodeTimezoneInvalid      = "timezone:invalid"
	ErrorCodeURLInvalid           = "url:invalid"
	ErrorCodeURNInvalid           = "urn:invalid"
	ErrorCodeURNTaken             = "urn:taken"
	ErrorCodeUserMissing          = "user:missing"
	ErrorCodeExpression           = "expression"
	ErrorCodeExpressionTooComplex = "expression:too_complex"
	ErrorCodeWebhookRequestSize   = "webhook:request_size"
)

// Error events are created when an error occurs during flow execution. Some errors have a `code`
// which identifies the type of error, and `extra` values with more details.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "error",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "text": "invalid date format: '12th of October'"
//	}
//
// @event error
type Error struct {
	BaseEvent

	Text  string            `json:"text"           validate:"required"`
	Code  string            `json:"code,omitempty"`
	Extra map[string]string `json:"extra,omitempty"`
}

// NewError returns a new error event for the passed in text
func NewError(text, code string, extra ...string) *Error {
	return &Error{
		BaseEvent: NewBaseEvent(TypeError),
		Text:      text,
		Code:      code,
		Extra:     extraFromPairs(extra),
	}
}

// NewRawError returns a new error event for the passed in error
func NewRawError(err error) *Error {
	return NewError(err.Error(), "")
}

// NewDependencyError returns an error event for a missing dependency
func NewDependencyError(ref assets.Reference) *Error {
	return NewError(fmt.Sprintf("Missing dependency: %s", ref.String()), ErrorCodeDependencyMissing, "type", ref.Type(), "identity", ref.Identity())
}

// NewActionUnsupportedError returns an error event for an action that is no longer supported
func NewActionUnsupportedError(text string) *Error {
	return NewError(text, ErrorCodeActionUnsupported)
}
