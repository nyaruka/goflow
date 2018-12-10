package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeIVRPlay, func() flows.Event { return &IVRPlayEvent{} })
}

// TypeIVRSay is a constant for IVR say events
const TypeIVRSay string = "ivr_say"

// IVRSayEvent events are created when an action wants to say a message to the current contact using TTS.
//
//   {
//     "type": "ivr_say",
//     "created_on": "2006-01-02T15:04:05Z",
//     "text": "Hi John. May we ask you some questions?"
//   }
//
// @event ivr_say
type IVRSayEvent struct {
	BaseEvent

	Text string `json:"text" validate:"required"`
}

// NewIVRSayEvent creates a new IVR say event
func NewIVRSayEvent(text string) *IVRSayEvent {
	return &IVRSayEvent{
		BaseEvent: NewBaseEvent(TypeIVRSay),
		Text:      text,
	}
}
