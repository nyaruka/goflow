package waits

import (
	"time"

	"github.com/nyaruka/goflow/flows"
)

type BaseWait struct {
	Timeout   *int       `json:"timeout,omitempty"`
	TimeoutOn *time.Time `json:"timeout_on,omitempty"`
}

func (w *BaseWait) begin(run flows.FlowRun) {
	if w.Timeout != nil {
		timeoutOn := time.Now().UTC().Add(time.Second * time.Duration(*w.Timeout))

		w.TimeoutOn = &timeoutOn
		w.Timeout = nil
	}

	run.SetStatus(flows.RunStatusWaiting)
}

func (w *BaseWait) HasTimedOut() bool {
	return w.TimeoutOn != nil && time.Now().After(*w.TimeoutOn)
}
