package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeFlowTriggered, func() flows.Event { return &FlowTriggeredEvent{} })
}

// TypeFlowTriggered is the type of our flow triggered event
const TypeFlowTriggered string = "flow_triggered"

// FlowTriggeredEvent events are created when an action has started a sub-flow.
//
//   {
//     "type": "flow_triggered",
//     "created_on": "2006-01-02T15:04:05Z",
//     "flow": {"uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a", "name": "Registration"},
//     "parent_run_uuid": "95eb96df-461b-4668-b168-727f8ceb13dd"
//   }
//
// @event flow_triggered
type FlowTriggeredEvent struct {
	BaseEvent

	Flow          *assets.FlowReference `json:"flow" validate:"required"`
	ParentRunUUID flows.RunUUID         `json:"parent_run_uuid" validate:"omitempty,uuid4"`
}

// NewFlowTriggeredEvent returns a new flow triggered event for the passed in flow and parent run
func NewFlowTriggeredEvent(flow *assets.FlowReference, parentRunUUID flows.RunUUID) *FlowTriggeredEvent {
	return &FlowTriggeredEvent{
		BaseEvent:     NewBaseEvent(),
		Flow:          flow,
		ParentRunUUID: parentRunUUID,
	}
}

// Type returns the type of this event
func (e *FlowTriggeredEvent) Type() string { return TypeFlowTriggered }
