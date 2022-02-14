package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeDialEnded, func() flows.Event { return &DialEndedEvent{} })
}

// TypeDialEnded is the type of our dial ended event
const TypeDialEnded string = "dial_ended"

// DialEndedEvent events are created when a session is resumed after waiting for a dial.
//
//   {
//     "type": "dial_ended",
//     "created_on": "2019-01-02T15:04:05Z",
//     "dial": {
//       "status": "answered",
//       "duration": 10
//     }
//   }
//
// @event dial_ended
type DialEndedEvent struct {
	BaseEvent

	Dial *flows.Dial `json:"dial" validate:"required,dive"`
}

// NewDialEnded returns a new dial ended event
func NewDialEnded(dial *flows.Dial) *DialEndedEvent {
	return &DialEndedEvent{
		BaseEvent: NewBaseEvent(TypeDialEnded),
		Dial:      dial,
	}
}

var _ flows.Event = (*DialEndedEvent)(nil)
