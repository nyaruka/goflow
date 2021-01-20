package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeRedirectWait, func() flows.Event { return &RedirectWaitEvent{} })
}

// TypeRedirectWait is the type of our redirect wait event
const TypeRedirectWait string = "redirect_wait"

// RedirectWaitEvent events are created when a flow pauses waiting for a redirect such as
// as forwarding an IVR call to another number.
//
//   {
//     "type": "redirect_wait",
//     "created_on": "2019-01-02T15:04:05Z"
//   }
//
// @event redirect_wait
type RedirectWaitEvent struct {
	baseEvent
}

// NewRedirectWait returns a new msg wait with the passed in timeout
func NewRedirectWait() *RedirectWaitEvent {
	return &RedirectWaitEvent{
		baseEvent: newBaseEvent(TypeRedirectWait),
	}
}

var _ flows.Event = (*RedirectWaitEvent)(nil)
