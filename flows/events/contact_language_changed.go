package events

import (
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeContactLanguageChanged, func() flows.Event { return &ContactLanguageChangedEvent{} })
}

// TypeContactLanguageChanged is the type of our contact language changed event
const TypeContactLanguageChanged string = "contact_language_changed"

// ContactLanguageChangedEvent events are created when the language of the contact has been changed.
//
//   {
//     "type": "contact_language_changed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "language": "eng"
//   }
//
// @event contact_language_changed
type ContactLanguageChangedEvent struct {
	BaseEvent

	Language string `json:"language"`
}

// NewContactLanguageChanged returns a new contact language changed event
func NewContactLanguageChanged(language envs.Language) *ContactLanguageChangedEvent {
	return &ContactLanguageChangedEvent{
		BaseEvent: NewBaseEvent(TypeContactLanguageChanged),
		Language:  string(language),
	}
}
