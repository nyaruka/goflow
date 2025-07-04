package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeDialEnded, func() flows.Event { return &DialEnded{} })
}

// TypeDialEnded is the type of our dial ended event
const TypeDialEnded string = "dial_ended"

// DialEnded events are created when a session is resumed after waiting for a dial.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "dial_ended",
//	  "created_on": "2019-01-02T15:04:05Z",
//	  "dial": {
//	    "status": "answered",
//	    "duration": 10
//	  }
//	}
//
// @event dial_ended
type DialEnded struct {
	BaseEvent

	Dial *flows.Dial `json:"dial" validate:"required"`
}

// NewDialEnded returns a new dial ended event
func NewDialEnded(dial *flows.Dial) *DialEnded {
	return &DialEnded{
		BaseEvent: NewBaseEvent(TypeDialEnded),
		Dial:      dial,
	}
}

var _ flows.Event = (*DialEnded)(nil)
