package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeContactLanguageChanged, func() flows.Event { return &ContactLanguageChangedEvent{} })
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

// NewContactLanguageChangedEvent returns a new contact language changed event
func NewContactLanguageChangedEvent(language string) *ContactLanguageChangedEvent {
	return &ContactLanguageChangedEvent{
		BaseEvent: NewBaseEvent(TypeContactLanguageChanged),
		Language:  language,
	}
}
