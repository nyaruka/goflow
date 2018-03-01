package events

import (
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// TypeEnvironmentChanged is the type of our environment changed event
const TypeEnvironmentChanged string = "environment_changed"

// EnvironmentChangedEvent events are created to set the environment on a session
//
// ```
//   {
//     "type": "environment_changed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "date_format": "yyyy-MM-dd",
//     "time_format": "hh:mm",
//     "timezone": "Africa/Kigali",
//     "languages": ["eng", "fra"]
//   }
// ```
//
// @event environment_changed
type EnvironmentChangedEvent struct {
	BaseEvent
	DateFormat utils.DateFormat   `json:"date_format"`
	TimeFormat utils.TimeFormat   `json:"time_format"`
	Timezone   string             `json:"timezone"`
	Languages  utils.LanguageList `json:"languages"`
}

// Type returns the type of this event
func (e *EnvironmentChangedEvent) Type() string { return TypeEnvironmentChanged }

// AllowedOrigin determines where this event type can originate
func (e *EnvironmentChangedEvent) AllowedOrigin() flows.EventOrigin { return flows.EventOriginCaller }

// Apply applies this event to the given run
func (e *EnvironmentChangedEvent) Apply(run flows.FlowRun) error {
	tz, err := time.LoadLocation(e.Timezone)
	if err != nil {
		return err
	}

	env := utils.NewEnvironment(e.DateFormat, e.TimeFormat, tz, e.Languages)

	run.Session().SetEnvironment(env)
	return nil
}
