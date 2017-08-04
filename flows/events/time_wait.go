package events

import "github.com/nyaruka/goflow/flows"

// TypeTimeWait is the type of our time wait event
const TypeTimeWait string = "time_wait"

// TimeWaitEvent events are created when a flow requests to be paused for a certain amount of time.
// ```
//   {
//     "type": "time_wait",
//     "created_on": "2006-01-02T15:04:05Z",
//     "timeout": 300
//   }
// ```
//
// @event time_wait
type TimeWaitEvent struct {
	Timeout int `json:"timeout"`
	BaseEvent
}

// NewTimeWait returns a new time wait with the passed in timeout
func NewTimeWait(timeout int) *TimeWaitEvent {
	return &TimeWaitEvent{
		BaseEvent: NewBaseEvent(),
		Timeout:   timeout,
	}
}

// Type returns the type of this event
func (e *TimeWaitEvent) Type() string { return TypeMsgWait }

// Apply applies this event to the given run
func (e *TimeWaitEvent) Apply(run flows.FlowRun) {}
