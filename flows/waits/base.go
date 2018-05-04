package waits

import (
	"time"

	"github.com/nyaruka/goflow/flows"
)

// the base of all wait types
type baseWait struct {
}

// Timeout would return the timeout of this wait for wait types that do that
func (w *baseWait) Timeout() *int {
	return nil
}

// HasTimedOut returns whether this wait has timed out
func (w *baseWait) HasTimedOut() bool {
	return false
}

// Begin beings waiting
func (w *baseWait) Begin(run flows.FlowRun) {
	run.SetStatus(flows.RunStatusWaiting)
}

func (w *baseWait) Resume(run flows.FlowRun) {
	run.SetStatus(flows.RunStatusActive)
}

func (w *baseWait) ResumeByTimeOut(run flows.FlowRun) {
	w.Resume(run)
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

// HasTimedOut returns whether this wait has timed out
func (w *baseTimeoutWait) HasTimedOut() bool {
	return w.TimeoutOn != nil && time.Now().After(*w.TimeoutOn)
}
