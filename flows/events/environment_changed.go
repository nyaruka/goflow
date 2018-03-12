package events

import (
	"encoding/json"

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
//     "environment": {
//       "date_format": "yyyy-MM-dd",
//       "time_format": "hh:mm",
//       "timezone": "Africa/Kigali",
//       "languages": ["eng", "fra"]
//     }
//   }
// ```
//
// @event environment_changed
type EnvironmentChangedEvent struct {
	BaseEvent
	callerOnlyEvent

	Environment json.RawMessage `json:"environment"`
}

// Type returns the type of this event
func (e *EnvironmentChangedEvent) Type() string { return TypeEnvironmentChanged }

// Validate validates our event is valid and has all the assets it needs
func (e *EnvironmentChangedEvent) Validate(assets flows.SessionAssets) error {
	return nil
}

// Apply applies this event to the given run
func (e *EnvironmentChangedEvent) Apply(run flows.FlowRun) error {
	env, err := utils.ReadEnvironment(e.Environment)
	if err != nil {
		return err
	}

	run.Session().SetEnvironment(env)
	return nil
}
