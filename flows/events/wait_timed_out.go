package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeWaitTimedOut, func() flows.Event { return &WaitTimedOutEvent{} })
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

// NewWaitTimedOut creates a new wait timed out event
func NewWaitTimedOut() *WaitTimedOutEvent {
	return &WaitTimedOutEvent{BaseEvent: NewBaseEvent(TypeWaitTimedOut)}
}

var _ flows.Event = (*WaitTimedOutEvent)(nil)
