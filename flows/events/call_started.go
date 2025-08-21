package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeCallStarted, func() flows.Event { return &CallStarted{} })
}

// TypeCallStarted is the type of our call started event
const TypeCallStarted string = "call_started"

// CallStarted events are created when a session is resumed after waiting for a dial.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "call_started",
//	  "created_on": "2019-01-02T15:04:05Z",
//	  "call": {
//	    "uuid": "0198ce92-ff2f-7b07-b158-b21ab168ebba",
//	    "urn": "tel:+1234567890",
//	    "channel": {"uuid": "2e2b43c7-88e9-43d6-b291-ebf9a24c2f86", "name": "Twilio"}
//	  }
//	}
//
// @event call_started
type CallStarted struct {
	BaseEvent

	Call *flows.CallEnvelope `json:"call" validate:"required"`
}

// NewCallStarted returns a new call started event
func NewCallStarted(call *flows.Call) *CallStarted {
	return &CallStarted{
		BaseEvent: NewBaseEvent(TypeCallStarted),
		Call:      call.Marshal(),
	}
}

var _ flows.Event = (*CallStarted)(nil)
