package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeFlowEntered, func() flows.Event { return &FlowEntered{} })
}

// TypeFlowEntered is the type of our flow entered event
const TypeFlowEntered string = "flow_entered"

// FlowEntered events are created when an action has entered a sub-flow.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "flow_entered",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "flow": {"uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a", "name": "Registration"},
//	  "parent_run_uuid": "95eb96df-461b-4668-b168-727f8ceb13dd",
//	  "terminal": false
//	}
//
// @event flow_entered
type FlowEntered struct {
	BaseEvent

	Flow          *assets.FlowReference `json:"flow" validate:"required"`
	ParentRunUUID flows.RunUUID         `json:"parent_run_uuid" validate:"omitempty,uuid"`
	Terminal      bool                  `json:"terminal"`
}

// NewFlowEntered returns a new flow entered event for the passed in flow and parent run
func NewFlowEntered(flow *assets.FlowReference, parentRunUUID flows.RunUUID, terminal bool) *FlowEntered {
	return &FlowEntered{
		BaseEvent:     NewBaseEvent(TypeFlowEntered),
		Flow:          flow,
		ParentRunUUID: parentRunUUID,
		Terminal:      terminal,
	}
}
