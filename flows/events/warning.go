package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeWarning, func() flows.Event { return &Warning{} })
}

// TypeWarning is the type of our warning events
const TypeWarning string = "warning"

// Warning events are created for things like accessing deprecated context values.
//
//	{
//	  "type": "warning",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "text": "deprecated context value accessed: legacy_extra"
//	}
//
// @event warning
type Warning struct {
	BaseEvent

	Text string `json:"text" validate:"required"`
}

// NewWarning returns a new warning event
func NewWarning(text string) *Warning {
	return &Warning{
		BaseEvent: NewBaseEvent(TypeWarning),
		Text:      text,
	}
}
