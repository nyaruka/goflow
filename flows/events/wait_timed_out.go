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
	wait := run.Session().Wait()

	if run.Status() != flows.RunStatusWaiting || wait == nil {
		return fmt.Errorf("can only be applied to waiting runs")
	}

	if wait.Timeout() == nil || wait.TimeoutOn() == nil {
		return fmt.Errorf("can only be applied when session wait has timeout")
	}

	if e.CreatedOn().Before(*wait.TimeoutOn()) {
		return fmt.Errorf("can't apply before wait has timed out")
	}

	run.SetInput(nil)
	return nil
}
