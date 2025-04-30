package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeWarning, func() flows.Event { return &WarningEvent{} })
}

// TypeWarning is the type of our warning events
const TypeWarning string = "warning"

// WarningEvent events are created for things like accessing deprecated context values.
//
//	{
//	  "uuid": "019688A6-41d2-7366-958a-630e35c62431",
//	  "type": "warning",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "text": "deprecated context value accessed: legacy_extra"
//	}
//
// @event warning
type WarningEvent struct {
	BaseEvent

	Text string `json:"text" validate:"required"`
}

// NewWarning returns a new warning event
func NewWarning(text string) *WarningEvent {
	return &WarningEvent{
		BaseEvent: NewBaseEvent(TypeWarning),
		Text:      text,
	}
}
