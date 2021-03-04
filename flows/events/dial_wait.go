package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeDialWait, func() flows.Event { return &DialWaitEvent{} })
}

// TypeDialWait is the type of our dial wait event
const TypeDialWait string = "dial_wait"

// DialWaitEvent events are created when a flow pauses waiting for an IVR dial to complete.
//
//   {
//     "type": "dial_wait",
//     "created_on": "2019-01-02T15:04:05Z",
//     "urn": "tel:+593979123456"
//   }
//
// @event dial_wait
type DialWaitEvent struct {
	baseEvent

	URN urns.URN `json:"urn" validate:"required,urn"`
}

// NewDialWait returns a new dial wait with the passed in URN
func NewDialWait(urn urns.URN) *DialWaitEvent {
	return &DialWaitEvent{
		baseEvent: newBaseEvent(TypeDialWait),
		URN:       urn,
	}
}

var _ flows.Event = (*DialWaitEvent)(nil)
