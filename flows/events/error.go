package events

import (
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeError, func() flows.Event { return &ErrorEvent{} })
}

// TypeError is the type of our error events
const TypeError string = "error"

// ErrorEvent events are created when an error occurs during flow execution.
//
//   {
//     "type": "error",
//     "created_on": "2006-01-02T15:04:05Z",
//     "text": "invalid date format: '12th of October'"
//   }
//
// @event error
type ErrorEvent struct {
	BaseEvent

	Text string `json:"text" validate:"required"`
}

// NewError returns a new error event for the passed in error
func NewError(err error) *ErrorEvent {
	return NewErrorf(err.Error())
}

// NewErrorf returns a new error event for the passed in format string and args
func NewErrorf(format string, a ...interface{}) *ErrorEvent {
	return &ErrorEvent{
		BaseEvent: NewBaseEvent(TypeError),
		Text:      fmt.Sprintf(format, a...),
	}
}

// NewDependencyError returns an error event for a missing dependency
func NewDependencyError(ref assets.Reference) *ErrorEvent {
	return NewErrorf("missing dependency: %s", ref.String())
}
