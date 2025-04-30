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
//	{
//	  "uuid": "019688A6-41d2-7366-958a-630e35c62431",
//	  "type": "failure",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "text": "unable to read flow"
//	}
//
// @event failure
type FailureEvent struct {
	BaseEvent

	Text string `json:"text" validate:"required"`
}

// NewFailure returns a new failure event for the passed in error
func NewFailure(err error) *FailureEvent {
	return &FailureEvent{
		BaseEvent: NewBaseEvent(TypeFailure),
		Text:      err.Error(),
	}
}
