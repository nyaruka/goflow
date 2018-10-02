package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeWaitTimedOut, func() flows.Event { return &WaitTimedOutEvent{} })
}

// TypeWaitTimedOut is the type of our wait timed out events
const TypeWaitTimedOut string = "wait_timed_out"

// WaitTimedOutEvent events are sent by the caller when a wait has timed out - i.e. they are sent instead of
// the item that the wait was waiting for.
//
//   {
//     "type": "wait_timed_out",
//     "created_on": "2006-01-02T15:04:05Z"
//   }
//
// @event wait_timed_out
type WaitTimedOutEvent struct {
	BaseEvent
}

// NewWaitTimedOutEvent creates a new wait timed out event
func NewWaitTimedOutEvent() *WaitTimedOutEvent {
	return &WaitTimedOutEvent{BaseEvent: NewBaseEvent()}
}

// Type returns the type of this event
func (e *WaitTimedOutEvent) Type() string { return TypeWaitTimedOut }

var _ flows.Event = (*WaitTimedOutEvent)(nil)
