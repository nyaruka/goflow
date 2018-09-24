package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeError, func() flows.Event { return &ErrorEvent{} })
}

// TypeError is the type of our error events
const TypeError string = "error"

// ErrorEvent events will be created whenever an error is encountered during flow execution. This
// can vary from template evaluation errors to invalid actions.
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

	Text  string `json:"text" validate:"required"`
	Fatal bool   `json:"fatal"`
}

// NewErrorEvent returns a new error event for the passed in error
func NewErrorEvent(err error) *ErrorEvent {
	return &ErrorEvent{
		BaseEvent: NewBaseEvent(),
		Text:      err.Error(),
	}
}

// NewFatalErrorEvent returns a new fatal error event for the passed in error
func NewFatalErrorEvent(err error) *ErrorEvent {
	return &ErrorEvent{
		BaseEvent: NewBaseEvent(),
		Text:      err.Error(),
		Fatal:     true,
	}
}

// Type returns the type of this event
func (e *ErrorEvent) Type() string { return TypeError }
