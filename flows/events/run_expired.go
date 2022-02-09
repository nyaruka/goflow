package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeRunExpired, func() flows.Event { return &RunExpiredEvent{} })
}

// TypeRunExpired is the type of our flow expired event
const TypeRunExpired string = "run_expired"

// RunExpiredEvent events are sent by the caller to tell the engine that a run has expired.
//
//   {
//     "type": "run_expired",
//     "created_on": "2006-01-02T15:04:05Z",
//     "run_uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a"
//   }
//
// @event run_expired
type RunExpiredEvent struct {
	BaseEvent

	RunUUID flows.RunUUID `json:"run_uuid"    validate:"required,uuid4"`
}

// NewRunExpired creates a new run expired event
func NewRunExpired(run flows.Run) *RunExpiredEvent {
	return &RunExpiredEvent{
		BaseEvent: NewBaseEvent(TypeRunExpired),
		RunUUID:   run.UUID(),
	}
}

var _ flows.Event = (*RunExpiredEvent)(nil)
