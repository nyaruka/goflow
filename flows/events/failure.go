package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeFailure, func() flows.Event { return &ErrorEvent{} })
}

// TypeFailure is the type of our error events
const TypeFailure string = "failure"

// FailureEvent events are created when an error occurs during flow execution which prevents continuation of the session.
//
//   {
//     "type": "failure",
//     "created_on": "2006-01-02T15:04:05Z",
//     "text": "unable to read flow"
//   }
//
// @event failure
type FailureEvent struct {
	baseEvent

	Text string `json:"text" validate:"required"`
}

// NewFailure returns a new failure event for the passed in error
func NewFailure(err error) *FailureEvent {
	return &FailureEvent{
		baseEvent: newBaseEvent(TypeFailure),
		Text:      err.Error(),
	}
}
