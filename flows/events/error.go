package events

import (
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeError, func() flows.Event { return &Error{} })
}

// TypeError is the type of our error events
const TypeError string = "error"

const (
	ErrorCodeDependencyMissing = "dependency_missing"
	ErrorCodeURNTaken          = "urn_taken"
)

// Error events are created when an error occurs during flow execution.
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

	Text string `json:"text"           validate:"required"`
	Code string `json:"code,omitempty"`
}

// NewError returns a new error event for the passed in text
func NewError(text, code string) *Error {
	return &Error{
		BaseEvent: NewBaseEvent(TypeError),
		Text:      text,
		Code:      code,
	}
}

// NewRawError returns a new error event for the passed in error
func NewRawError(err error) *Error {
	return NewError(err.Error(), "")
}

// NewDependencyError returns an error event for a missing dependency
func NewDependencyError(ref assets.Reference) *Error {
	return NewError(fmt.Sprintf("Missing dependency: %s", ref.String()), ErrorCodeDependencyMissing)
}
