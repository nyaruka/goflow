package events

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeContactLanguageChanged, func() flows.Event { return &ContactLanguageChangedEvent{} })
}

// TypeContactLanguageChanged is the type of our contact language changed event
const TypeContactLanguageChanged string = "contact_language_changed"

// ContactLanguageChangedEvent events are created when a Language of a contact has been changed
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
	callerOrEngineEvent

	Language string `json:"language"`
}

// NewContactLanguageChangedEvent returns a new contact language changed event
func NewContactLanguageChangedEvent(language string) *ContactLanguageChangedEvent {
	return &ContactLanguageChangedEvent{
		BaseEvent: NewBaseEvent(),
		Language:  language,
	}
}

// Type returns the type of this event
func (e *ContactLanguageChangedEvent) Type() string { return TypeContactLanguageChanged }

// Validate validates our event is valid and has all the assets it needs
func (e *ContactLanguageChangedEvent) Validate(assets flows.SessionAssets) error {
	return nil
}

// Apply applies this event to the given run
func (e *ContactLanguageChangedEvent) Apply(run flows.FlowRun) error {
	if run.Contact() == nil {
		return fmt.Errorf("can't apply event in session without a contact")
	}

	if e.Language != "" {
		lang, err := utils.ParseLanguage(e.Language)
		if err != nil {
			return err
		}
		run.Contact().SetLanguage(lang)
	} else {
		run.Contact().SetLanguage(utils.NilLanguage)
	}

	return nil
}
