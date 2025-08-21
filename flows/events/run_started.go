package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeRunStarted, func() flows.Event { return &RunStarted{} })
}

// TypeRunStarted is the type of our run started event
const TypeRunStarted string = "run_started"

// RunStarted events are created when an action has entered a sub-flow.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "run_started",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "run_uuid": "0198bf2d-a96d-7d36-8a4f-154eecc0d798",
//	  "parent_uuid": "95eb96df-461b-4668-b168-727f8ceb13dd",
//	  "flow": {"uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a", "name": "Registration"},
//	  "terminal": false
//	}
//
// @event run_started
type RunStarted struct {
	BaseEvent

	Flow       *assets.FlowReference `json:"flow"                  validate:"required"`
	RunUUID    flows.RunUUID         `json:"run_uuid"              validate:"required,uuid"`
	ParentUUID flows.RunUUID         `json:"parent_uuid,omitempty" validate:"omitempty,uuid"`
	Terminal   bool                  `json:"terminal,omitempty"`
}

// NewRunStarted returns a new run started event
func NewRunStarted(run flows.Run, terminal bool) *RunStarted {
	var parentUUID flows.RunUUID
	if run.Parent() != nil {
		parentUUID = run.Parent().UUID()
	}

	return &RunStarted{
		BaseEvent:  NewBaseEvent(TypeRunStarted),
		Flow:       run.FlowReference(),
		RunUUID:    run.UUID(),
		ParentUUID: parentUUID,
		Terminal:   terminal,
	}
}
