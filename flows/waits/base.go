package waits

import (
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

var registeredTypes = map[string](func() flows.Wait){}

// RegisterType registers a new type of wait
func RegisterType(name string, initFunc func() flows.Wait) {
	registeredTypes[name] = initFunc
}

// the base of all wait types
type baseWait struct {
}

// Timeout would return the timeout of this wait for wait types that do that
func (w *baseWait) Timeout() *int { return nil }

// TimeoutOn would return when this wait times out for wait types that do that
func (w *baseWait) TimeoutOn() *time.Time { return nil }

// Begin beings waiting
func (w *baseWait) Begin(run flows.FlowRun) {}

// base of all wait types than can timeout
type baseTimeoutWait struct {
	baseWait

	Timeout_   *int       `json:"timeout,omitempty"`
	TimeoutOn_ *time.Time `json:"timeout_on,omitempty"`
}

// Timeout returns the timeout of this wait in seconds or nil if no timeout is set
func (w *baseTimeoutWait) Timeout() *int { return w.Timeout_ }

// TimeoutOn returns when this wait times out
func (w *baseTimeoutWait) TimeoutOn() *time.Time { return w.TimeoutOn_ }

// Begin beings waiting at this wait
func (w *baseTimeoutWait) Begin(run flows.FlowRun) {
	if w.Timeout_ != nil {
		timeoutOn := utils.Now().Add(time.Second * time.Duration(*w.Timeout_))

		w.TimeoutOn_ = &timeoutOn
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

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadWait reads a wait from the given typed envelope
func ReadWait(envelope *utils.TypedEnvelope) (flows.Wait, error) {
	f := registeredTypes[envelope.Type]
	if f == nil {
		return nil, fmt.Errorf("unknown type: %s", envelope.Type)
	}

	wait := f()
	return wait, utils.UnmarshalAndValidate(envelope.Data, wait)
}
