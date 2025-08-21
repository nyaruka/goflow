package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeCallReceived, func() flows.Event { return &CallReceived{} })
}

// TypeCallReceived is the type of our call received event
const TypeCallReceived string = "call_received"

// CallReceived events are created when an incoming call is received.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "call_received",
//	  "created_on": "2019-01-02T15:04:05Z",
//	  "call": {
//	    "uuid": "0198ce92-ff2f-7b07-b158-b21ab168ebba",
//	    "urn": "tel:+1234567890",
//	    "channel": {"uuid": "2e2b43c7-88e9-43d6-b291-ebf9a24c2f86", "name": "Twilio"}
//	  }
//	}
//
// @event call_received
type CallReceived struct {
	BaseEvent

	Call *flows.CallEnvelope `json:"call" validate:"required"`
}

// NewCallReceived returns a new call received event
func NewCallReceived(call *flows.Call) *CallReceived {
	return &CallReceived{
		BaseEvent: NewBaseEvent(TypeCallReceived),
		Call:      call.Marshal(),
	}
}

var _ flows.Event = (*CallReceived)(nil)
