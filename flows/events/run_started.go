package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeRunStarted, func() flows.Event { return &RunStarted{} })
	registerType("flow_entered", func() flows.Event { return &RunStarted{} }) // deprecated
}

// TypeRunStarted is the type of our run started event
const TypeRunStarted string = "run_started"

// RunStarted events are created when an action has entered a sub-flow.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "flow_entered",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "flow": {"uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a", "name": "Registration"},
//	  "parent_uuid": "95eb96df-461b-4668-b168-727f8ceb13dd",
//	  "terminal": false
//	}
//
// @event run_started
type RunStarted struct {
	BaseEvent

	Flow       *assets.FlowReference `json:"flow" validate:"required"`
	ParentUUID flows.RunUUID         `json:"parent_uuid" validate:"omitempty,uuid"`
	Terminal   bool                  `json:"terminal"`
}

// NewRunStarted returns a new run started event for the passed in flow and parent
func NewRunStarted(flow *assets.FlowReference, parentUUID flows.RunUUID, terminal bool) *RunStarted {
	return &RunStarted{
		BaseEvent:  NewBaseEvent(TypeRunStarted),
		Flow:       flow,
		ParentUUID: parentUUID,
		Terminal:   terminal,
	}
}
