package events

import (
	"time"

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
//     "urn": "tel:+593979123456",
//     "expires_on": "2022-02-02T13:27:30Z"
//   }
//
// @event dial_wait
type DialWaitEvent struct {
	BaseEvent

	URN urns.URN `json:"urn" validate:"required,urn"`

	// when this wait expires and the whole run can be expired
	ExpiresOn *time.Time `json:"expires_on,omitempty"`
}

// NewDialWait returns a new dial wait with the passed in URN
func NewDialWait(urn urns.URN, expiresOn *time.Time) *DialWaitEvent {
	return &DialWaitEvent{
		BaseEvent: NewBaseEvent(TypeDialWait),
		URN:       urn,
		ExpiresOn: expiresOn,
	}
}

var _ flows.Event = (*DialWaitEvent)(nil)
