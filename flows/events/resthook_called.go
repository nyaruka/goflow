package events

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"

	"github.com/pkg/errors"
)

func init() {
	RegisterType(TypeResthookCalled, func() flows.Event { return &ResthookCalledEvent{} })
}

// TypeResthookCalled is the type for our resthook events
const TypeResthookCalled string = "resthook_called"

// ResthookCalledEvent events are created when a resthook is called. The event contains
// the payload that will be sent to any subscribers of that resthook. Note that this event is
// created regardless of whether there any subscriberes for that resthook.
//
//   {
//     "type": "resthook_called",
//     "created_on": "2006-01-02T15:04:05Z",
//     "resthook": "success",
//     "payload": {
//       "contact:":{
//         "name":"Bob"
//       }
//     }
//   }
//
// @event resthook_called
type ResthookCalledEvent struct {
	BaseEvent

	Resthook string          `json:"resthook"`
	Payload  json.RawMessage `json:"payload"`
}

// NewResthookCalledEvent returns a new webhook called event
func NewResthookCalledEvent(resthook string, payload json.RawMessage) *ResthookCalledEvent {
	return &ResthookCalledEvent{
		BaseEvent: NewBaseEvent(TypeResthookCalled),
		Resthook:  resthook,
		Payload:   payload,
	}
}

// ResthookCalledEnvelope lets us avoid a infinite loop between json.Marshal and ResthookCalledEvent.MarshalJSON
type ResthookCalledEnvelope ResthookCalledEvent

// MarshalJSON marshals this object to JSON
func (e *ResthookCalledEvent) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal((*ResthookCalledEnvelope)(e))
	if err != nil {
		// we generate the payload so there's always a chance it isn't valid JSON
		return nil, errors.Errorf("error marshaling resthook_called event with payload: %s", string(e.Payload))
	}
	return data, nil
}
