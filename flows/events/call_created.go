package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeCallCreated, func() flows.Event { return &CallCreated{} })
}

// TypeCallCreated is the type of our call created event
const TypeCallCreated string = "call_created"

// CallCreated events are created when an outgoing call is created.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "call_created",
//	  "created_on": "2019-01-02T15:04:05Z",
//	  "call": {
//	    "uuid": "0198ce92-ff2f-7b07-b158-b21ab168ebba",
//	    "urn": "tel:+1234567890",
//	    "channel": {"uuid": "2e2b43c7-88e9-43d6-b291-ebf9a24c2f86", "name": "Twilio"}
//	  }
//	}
//
// @event call_created
type CallCreated struct {
	BaseEvent

	Call *flows.CallEnvelope `json:"call" validate:"required"`
}

// NewCallCreated returns a new call created event
func NewCallCreated(call *flows.Call) *CallCreated {
	return &CallCreated{
		BaseEvent: NewBaseEvent(TypeCallCreated),
		Call:      call.Marshal(),
	}
}

var _ flows.Event = (*CallCreated)(nil)
