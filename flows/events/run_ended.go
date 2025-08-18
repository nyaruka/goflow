package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeRunEnded, func() flows.Event { return &RunEnded{} })
}

// TypeRunEnded is the type of our run ended event
const TypeRunEnded string = "run_ended"

// RunEnded events are created when a run ends.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "run_ended",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "run_uuid": "0198bf2d-a96d-7d36-8a4f-154eecc0d798",
//	  "flow": {"uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a", "name": "Registration"},
//	  "status": "completed"
//	}
//
// @event run_ended
type RunEnded struct {
	BaseEvent

	RunUUID flows.RunUUID         `json:"run_uuid"    validate:"required,uuid"`
	Flow    *assets.FlowReference `json:"flow"        validate:"required"`
	Status  flows.RunStatus       `json:"status"      validate:"required"`
}

// NewRunEnded returns a new run ended event
func NewRunEnded(runUUID flows.RunUUID, flow *assets.FlowReference, status flows.RunStatus) *RunEnded {
	return &RunEnded{
		BaseEvent: NewBaseEvent(TypeRunEnded),
		RunUUID:   runUUID,
		Flow:      flow,
		Status:    status,
	}
}
