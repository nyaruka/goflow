package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeIVRPlay, func() flows.Event { return &IVRPlayEvent{} })
}

// TypeIVRPlay is a constant for IVR play events
const TypeIVRPlay string = "ivr_play"

// IVRPlayEvent events are created when an action wants to play an audio recording to the current contact.
// Text is optionally and only used for logging purposes.
//
//   {
//     "type": "ivr_play",
//     "created_on": "2006-01-02T15:04:05Z",
//     "audio_url": "http://uploads.temba.io/2353262.m4a",
//     "text": "Hi John. May we ask you some questions?"
//   }
//
// @event ivr_play
type IVRPlayEvent struct {
	BaseEvent

	AudioURL string `json:"audio_url" validate:"required"`
	Text     string `json:"text,omitempty"`
}

// NewIVRPlayEvent creates a new IVR play event
func NewIVRPlayEvent(audioURL string, text string) *IVRPlayEvent {
	return &IVRPlayEvent{
		BaseEvent: NewBaseEvent(TypeIVRPlay),
		AudioURL:  audioURL,
		Text:      text,
	}
}
