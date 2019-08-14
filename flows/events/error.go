package events

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeError, func() flows.Event { return &ErrorEvent{} })
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

// NewErrorEvent returns a new error event for the passed in error
func NewErrorEvent(err error) *ErrorEvent {
	return &ErrorEvent{
		BaseEvent: NewBaseEvent(TypeError),
		Text:      err.Error(),
	}
}

// NewErrorEventf returns a new error event for the passed in format string and args
func NewErrorEventf(format string, a ...interface{}) *ErrorEvent {
	return &ErrorEvent{
		BaseEvent: NewBaseEvent(TypeError),
		Text:      fmt.Sprintf(format, a...),
	}
}
