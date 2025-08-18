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

type FlowEntered struct {
	BaseEvent

	Flow          *assets.FlowReference `json:"flow" validate:"required"`
	ParentRunUUID flows.RunUUID         `json:"parent_run_uuid" validate:"omitempty,uuid"`
	Terminal      bool                  `json:"terminal"`
}
