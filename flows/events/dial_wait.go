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
//	{
//	  "uuid": "019688A6-41d2-7366-958a-630e35c62431",
//	  "type": "dial_wait",
//	  "created_on": "2019-01-02T15:04:05Z",
//	  "urn": "tel:+593979123456",
//	  "dial_limit_seconds": 60,
//	  "call_limit_seconds": 120,
//	  "expires_on": "2022-02-02T13:27:30Z"
//	}
//
// @event dial_wait
type DialWaitEvent struct {
	BaseEvent

	URN              urns.URN `json:"urn" validate:"required,urn"`
	DialLimitSeconds int      `json:"dial_limit_seconds"`
	CallLimitSeconds int      `json:"call_limit_seconds"`

	// when this wait expires and the whole run can be expired
	ExpiresOn time.Time `json:"expires_on,omitempty"`
}

// NewDialWait returns a new dial wait with the passed in URN
func NewDialWait(urn urns.URN, dialLimitSeconds, callLimitSeconds int, expiresOn time.Time) *DialWaitEvent {
	return &DialWaitEvent{
		BaseEvent:        NewBaseEvent(TypeDialWait),
		URN:              urn,
		DialLimitSeconds: dialLimitSeconds,
		CallLimitSeconds: callLimitSeconds,
		ExpiresOn:        expiresOn,
	}
}

var _ flows.Event = (*DialWaitEvent)(nil)
