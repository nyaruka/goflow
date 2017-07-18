package events

import (
	"github.com/nyaruka/goflow/flows"
)

// TypeFlowEntered is the type of our flow enter event
const TypeFlowEntered string = "flow_entered"

// FlowEnteredEvent events are created when a contact first enters a flow.
//
// ```
//   {
//    "type": "flow_entered",
//    "step_uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//    "created_on": "2006-01-02T15:04:05Z",
//    "flow_uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
//    "contact_uuid": "95eb96df-461b-4668-b168-727f8ceb13dd"
//   }
// ```
//
// @event flow_entered
type FlowEnteredEvent struct {
	BaseEvent
	FlowUUID    flows.FlowUUID    `json:"flow_uuid"        validate:"required,uuid4"`
	ContactUUID flows.ContactUUID `json:"contact_uuid"     validate:"required,uuid4"`
}

// NewFlowEnterEvent returns a new flow enter event for the passed in flow and contact
func NewFlowEnterEvent(flow flows.FlowUUID, contact flows.ContactUUID) *FlowEnteredEvent {
	return &FlowEnteredEvent{
		BaseEvent:   NewBaseEvent(),
		FlowUUID:    flow,
		ContactUUID: contact,
	}
}

// Type returns the type of this event
func (e *FlowEnteredEvent) Type() string { return TypeFlowEntered }

// Apply applies this event to the given run
func (e *FlowEnteredEvent) Apply(run flows.FlowRun) {}
