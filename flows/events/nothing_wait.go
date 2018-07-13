package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeNothingWait, func() flows.Event { return &NothingWaitEvent{} })
}

// TypeNothingWait is the type of our nothing wait event
const TypeNothingWait string = "nothing_wait"

// NothingWaitEvent events are created when a flow requests to hand back control to the caller but isn't
// waiting for anything from the caller.
//
//   {
//     "type": "nothing_wait",
//     "created_on": "2006-01-02T15:04:05.234532Z"
//   }
//
// @event nothing_wait
type NothingWaitEvent struct {
	BaseEvent
	engineOnlyEvent
}

// NewNothingWait returns a new nothing wait
func NewNothingWait() *NothingWaitEvent {
	return &NothingWaitEvent{BaseEvent: NewBaseEvent()}
}

// Type returns the type of this event
func (e *NothingWaitEvent) Type() string { return TypeNothingWait }

// Apply applies this event to the given run
func (e *NothingWaitEvent) Apply(run flows.FlowRun) error {
	return nil
}
