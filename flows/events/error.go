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

// Error events are created when an error occurs during flow execution.
//
//	{
//	  "type": "error",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "text": "invalid date format: '12th of October'"
//	}
//
// @event error
type Error struct {
	BaseEvent

	Text string `json:"text" validate:"required"`
}

// NewError returns a new error event for the passed in text
func NewError(text string) *Error {
	return &Error{
		BaseEvent: NewBaseEvent(TypeError),
		Text:      text,
	}
}

// NewDependencyError returns an error event for a missing dependency
func NewDependencyError(ref assets.Reference) *Error {
	return NewError(fmt.Sprintf("missing dependency: %s", ref.String()))
}
