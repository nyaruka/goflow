package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeRedirectEnded, func() flows.Event { return &RedirectEndedEvent{} })
}

// TypeRedirectEnded is the type of our redirect ended event
const TypeRedirectEnded string = "redirect_ended"

// RedirectEndedEvent events are created when a session is resumed after waiting for a redirect which ended successfully.
//
//   {
//     "type": "redirect_ended",
//     "created_on": "2019-01-02T15:04:05Z",
//     "response": "answered"
//   }
//
// @event redirect_ended
type RedirectEndedEvent struct {
	baseEvent

	Response flows.RedirectResponse `json:"response" validate:"required"`
}

// NewRedirectEnded returns a new redirect ended event
func NewRedirectEnded(response flows.RedirectResponse) *RedirectEndedEvent {
	return &RedirectEndedEvent{
		baseEvent: newBaseEvent(TypeRedirectEnded),
		Response:  response,
	}
}

var _ flows.Event = (*RedirectEndedEvent)(nil)
