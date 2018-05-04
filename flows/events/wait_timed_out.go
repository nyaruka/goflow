package events

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
)

// TypeWaitTimedOut is the type of our wait timed out events
const TypeWaitTimedOut string = "wait_timed_out"

// WaitTimedOutEvent events are sent by the caller when a wait has timed out - i.e. they are sent instead of
// the item that the wait was waiting for
//
//   {
//     "type": "wait_timed_out",
//     "created_on": "2006-01-02T15:04:05Z"
//   }
//
// @event wait_timed_out
type WaitTimedOutEvent struct {
	baseEvent
	callerOnlyEvent
}

// Type returns the type of this event
func (e *WaitTimedOutEvent) Type() string { return TypeWaitTimedOut }

// Validate validates our event is valid and has all the assets it needs
func (e *WaitTimedOutEvent) Validate(assets flows.SessionAssets) error {
	return nil
}

// Apply applies this event to the given run
func (e *WaitTimedOutEvent) Apply(run flows.FlowRun) error {
	if run.Status() != flows.RunStatusWaiting {
		return fmt.Errorf("wait_timed_out events can only be applied to waiting runs")
	}

	run.SetInput(nil)
	return nil
}
