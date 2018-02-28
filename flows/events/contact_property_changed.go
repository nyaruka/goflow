package events

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// TypeContactPropertyChanged is the type of our update contact event
const TypeContactPropertyChanged string = "contact_property_changed"

// ContactPropertyChangedEvent events are created when a property of a contact has been changed
//
// ```
//   {
//     "type": "contact_property_changed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "property": "language",
//     "value": "eng"
//   }
// ```
//
// @event contact_property_changed
type ContactPropertyChangedEvent struct {
	BaseEvent
	Property string `json:"property" validate:"required,eq=name|eq=language"`
	Value    string `json:"value"`
}

// NewContactPropertyChangedEvent returns a new contact property changed event
func NewContactPropertyChangedEvent(property string, value string) *ContactPropertyChangedEvent {
	return &ContactPropertyChangedEvent{
		BaseEvent: NewBaseEvent(),
		Property:  property,
		Value:     value,
	}
}

// Type returns the type of this event
func (e *ContactPropertyChangedEvent) Type() string { return TypeContactPropertyChanged }

// Apply applies this event to the given run
func (e *ContactPropertyChangedEvent) Apply(run flows.FlowRun) error {
	// if this is either name or language, we save directly to the contact
	if e.Property == "name" {
		run.Contact().SetName(e.Value)
	} else {
		if e.Value != "" {
			lang, err := utils.ParseLanguage(e.Value)
			if err != nil {
				return err
			}
			run.Contact().SetLanguage(lang)
		} else {
			run.Contact().SetLanguage(utils.NilLanguage)
		}
	}

	return nil
}
