package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeWaitExpired, func() flows.Event { return &WaitExpiredEvent{} })
}

// TypeWaitExpired is the type of our wait expired event
const TypeWaitExpired string = "wait_expired"

// WaitExpiredEvent events are sent by the caller to tell the engine that a wait has expired.
//
//	{
//	  "type": "wait_expired",
//	  "created_on": "2006-01-02T15:04:05Z"
//	}
//
// @event wait_expired
type WaitExpiredEvent struct {
	BaseEvent
}

// NewWaitExpired creates a new wait expired event
func NewWaitExpired() *WaitExpiredEvent {
	return &WaitExpiredEvent{
		BaseEvent: NewBaseEvent(TypeWaitExpired),
	}
}

var _ flows.Event = (*WaitExpiredEvent)(nil)
