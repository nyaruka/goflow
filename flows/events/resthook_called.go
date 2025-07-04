package events

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeResthookCalled, func() flows.Event { return &ResthookCalled{} })
}

// TypeResthookCalled is the type for our resthook events
const TypeResthookCalled string = "resthook_called"

// ResthookCalled events are created when a resthook is called. The event contains
// the payload that will be sent to any subscribers of that resthook. Note that this event is
// created regardless of whether there any subscriberes for that resthook.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "resthook_called",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "resthook": "success",
//	  "payload": {
//	    "contact:":{
//	      "name":"Bob"
//	    }
//	  }
//	}
//
// @event resthook_called
type ResthookCalled struct {
	BaseEvent

	Resthook string          `json:"resthook"`
	Payload  json.RawMessage `json:"payload"`
}

// NewResthookCalled returns a new webhook called event
func NewResthookCalled(resthook string, payload json.RawMessage) *ResthookCalled {
	return &ResthookCalled{
		BaseEvent: NewBaseEvent(TypeResthookCalled),
		Resthook:  resthook,
		Payload:   payload,
	}
}
