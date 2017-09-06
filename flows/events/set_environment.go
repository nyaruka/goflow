package events

import (
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// TypeSetEnvironment is the type of our set environment event
const TypeSetEnvironment string = "set_environment"

// SetEnvironmentEvent events are created to set the environment on a session
//
// ```
//   {
//     "type": "set_environment",
//     "created_on": "2006-01-02T15:04:05Z",
//     "date_format": "yyyy-MM-dd",
//     "time_format": "hh:mm",
//     "timezone": "Africa/Kigali",
//     "languages": ["eng", "fra"]
//   }
// ```
//
// @event set_environment
type SetEnvironmentEvent struct {
	BaseEvent
	DateFormat utils.DateFormat   `json:"date_format"`
	TimeFormat utils.TimeFormat   `json:"time_format"`
	Timezone   string             `json:"timezone"`
	Languages  utils.LanguageList `json:"languages"`
}

// Type returns the type of this event
func (e *SetEnvironmentEvent) Type() string { return TypeSetEnvironment }

// Apply applies this event to the given run
func (e *SetEnvironmentEvent) Apply(run flows.FlowRun, step flows.Step) error {
	tz, err := time.LoadLocation(e.Timezone)
	if err != nil {
		return err
	}

	env := utils.NewEnvironment(e.DateFormat, e.TimeFormat, tz, e.Languages)

	run.Session().SetEnvironment(env)
	return nil
}
