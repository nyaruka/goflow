package waits

import (
	"time"

	"github.com/nyaruka/goflow/flows"
)

// Base of all wait types
type BaseWait struct {
}

func (w *BaseWait) HasTimedOut() bool {
	return false
}

func (w *BaseWait) begin(run flows.FlowRun) {
	run.SetStatus(flows.RunStatusWaiting)
}

// Base of all wait types than can timeout
type TimeoutWait struct {
	BaseWait

	Timeout   *int       `json:"timeout,omitempty"`
	TimeoutOn *time.Time `json:"timeout_on,omitempty"`
}

func (w *TimeoutWait) begin(run flows.FlowRun) {
	if w.Timeout != nil {
		timeoutOn := time.Now().UTC().Add(time.Second * time.Duration(*w.Timeout))

		w.TimeoutOn = &timeoutOn
		w.Timeout = nil
	}

	w.BaseWait.begin(run)
}

func (w *TimeoutWait) HasTimedOut() bool {
	return w.TimeoutOn != nil && time.Now().After(*w.TimeoutOn)
}
