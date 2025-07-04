package events

import (
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeDialWait, func() flows.Event { return &DialWait{} })
}

// TypeDialWait is the type of our dial wait event
const TypeDialWait string = "dial_wait"

// DialWait events are created when a flow pauses waiting for an IVR dial to complete.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "dial_wait",
//	  "created_on": "2019-01-02T15:04:05Z",
//	  "urn": "tel:+593979123456",
//	  "dial_limit_seconds": 60,
//	  "call_limit_seconds": 120,
//	  "expires_on": "2022-02-02T13:27:30Z"
//	}
//
// @event dial_wait
type DialWait struct {
	BaseEvent

	URN              urns.URN `json:"urn" validate:"required,urn"`
	DialLimitSeconds int      `json:"dial_limit_seconds"`
	CallLimitSeconds int      `json:"call_limit_seconds"`

	// when this wait expires and the whole run can be expired
	ExpiresOn time.Time `json:"expires_on,omitempty"`
}

// NewDialWait returns a new dial wait with the passed in URN
func NewDialWait(urn urns.URN, dialLimitSeconds, callLimitSeconds int, expiresOn time.Time) *DialWait {
	return &DialWait{
		BaseEvent:        NewBaseEvent(TypeDialWait),
		URN:              urn,
		DialLimitSeconds: dialLimitSeconds,
		CallLimitSeconds: callLimitSeconds,
		ExpiresOn:        expiresOn,
	}
}

var _ flows.Event = (*DialWait)(nil)
