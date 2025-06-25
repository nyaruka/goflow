package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeFailure, func() flows.Event { return &Error{} })
}

// TypeFailure is the type of our error events
const TypeFailure string = "failure"

// Failure events are created when an error occurs during flow execution which prevents continuation of the session.
//
//	{
//	  "type": "failure",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "text": "unable to read flow"
//	}
//
// @event failure
type Failure struct {
	BaseEvent

	Text string `json:"text" validate:"required"`
}

// NewFailure returns a new failure event for the passed in error
func NewFailure(err error) *Failure {
	return &Failure{
		BaseEvent: NewBaseEvent(TypeFailure),
		Text:      err.Error(),
	}
}
