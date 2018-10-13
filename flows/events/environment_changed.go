package events

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeEnvironmentChanged, func() flows.Event { return &EnvironmentChangedEvent{} })
}

// TypeEnvironmentChanged is the type of our environment changed event
const TypeEnvironmentChanged string = "environment_changed"

// EnvironmentChangedEvent events are sent by the caller to tell the engine to update the session environment.
//
//   {
//     "type": "environment_changed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "environment": {
//       "date_format": "YYYY-MM-DD",
//       "time_format": "hh:mm",
//       "timezone": "Africa/Kigali",
//       "default_language": "eng",
//       "allowed_languages": ["eng", "fra"]
//     }
//   }
//
// @event environment_changed
type EnvironmentChangedEvent struct {
	BaseEvent

	Environment json.RawMessage `json:"environment"`
}

// NewEnvironmentChangedEvent creates a new environment changed event
func NewEnvironmentChangedEvent(env utils.Environment) *EnvironmentChangedEvent {
	marshalled, _ := json.Marshal(env)
	return &EnvironmentChangedEvent{
		BaseEvent:   NewBaseEvent(TypeEnvironmentChanged),
		Environment: marshalled,
	}
}

var _ flows.Event = (*EnvironmentChangedEvent)(nil)
