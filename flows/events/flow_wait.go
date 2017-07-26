package events

import "github.com/nyaruka/goflow/flows"

// TypeFlowWait is the type of our flow wait event
const TypeFlowWait string = "flow_wait"

// FlowWaitEvent events are created when a flow pauses waiting for a subflow to exit.
//
// ```
//   {
//     "type": "flow_wait",
//     "created_on": "2006-01-02T15:04:05Z",
//     "flow_uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a"
//   }
// ```
//
// @event flow_wait
type FlowWaitEvent struct {
	BaseEvent
	FlowUUID flows.FlowUUID `json:"flow_uuid"`
}

// NewFlowWait returns a new flow wait event
func NewFlowWait(flow flows.FlowUUID) *FlowWaitEvent {
	return &FlowWaitEvent{
		BaseEvent: NewBaseEvent(),
		FlowUUID:  flow,
	}
}

// Type returns the type of this event
func (e *FlowWaitEvent) Type() string { return TypeFlowWait }

// Apply applies this event to the given run
func (e *FlowWaitEvent) Apply(run flows.FlowRun) {}
