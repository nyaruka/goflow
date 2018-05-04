package waits

import (
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// the base of all wait types
type baseWait struct {
}

// Timeout would return the timeout of this wait for wait types that do that
func (w *baseWait) Timeout() *int {
	return nil
}

// Begin beings waiting
func (w *baseWait) Begin(run flows.FlowRun) {
	run.SetStatus(flows.RunStatusWaiting)
}

// base of all wait types than can timeout
type baseTimeoutWait struct {
	baseWait

	Timeout_  *int       `json:"timeout,omitempty"`
	TimeoutOn *time.Time `json:"timeout_on,omitempty"`
}

// Timeout returns the timeout of this wait in seconds or nil if no timeout is set
func (w *baseTimeoutWait) Timeout() *int {
	return w.Timeout_
}

// Begin beings waiting at this wait
func (w *baseTimeoutWait) Begin(run flows.FlowRun) {
	if w.Timeout_ != nil {
		timeoutOn := time.Now().UTC().Add(time.Second * time.Duration(*w.Timeout_))

		w.TimeoutOn = &timeoutOn
	}

	w.baseWait.Begin(run)
}

// CanResume returns true if a wait timed out event has been received
func (w *baseTimeoutWait) CanResume(callerEvents []flows.Event) bool {
	return containsEventOfType(callerEvents, events.TypeWaitTimedOut)
}

// utility function to look for an event of a given type
func containsEventOfType(events []flows.Event, eventType string) bool {
	for _, event := range events {
		if event.Type() == eventType {
			return true
		}
	}
	return false
}
