package events

import "github.com/nyaruka/goflow/flows"

// TypeFlowWait is the type of our flow wait event
const TypeFlowWait string = "flow_wait"

// FlowWaitEvent events are created when a flow pauses waiting for a subflow to exit.
//
// ```
//   {
//    "step": "8eebd020-1af5-431c-b943-aa670fc74da9",
//    "created_on": "2006-01-02T15:04:05Z",
//    "type": "flow_wait",
//    "flow": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a"
//   }
// ```
//
// @event flow_wait
type FlowWaitEvent struct {
	BaseEvent
	Flow flows.FlowUUID `json:"flow"`
}

// NewFlowWait returns a new flow wait event
func NewFlowWait(flow flows.FlowUUID) *FlowWaitEvent {
	return &FlowWaitEvent{Flow: flow}
}

// Type returns the type of this event
func (e *FlowWaitEvent) Type() string { return TypeFlowWait }
