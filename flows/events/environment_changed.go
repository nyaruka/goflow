package events

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeEnvironmentChanged, func() flows.Event { return &EnvironmentChangedEvent{} })
}

// TypeEnvironmentChanged is the type of our environment changed event
const TypeEnvironmentChanged string = "environment_changed"

// EnvironmentChangedEvent events are created to set the environment on a session
//
//   {
//     "type": "environment_changed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "environment": {
//       "date_format": "YYYY-MM-DD",
//       "time_format": "hh:mm",
//       "timezone": "Africa/Kigali",
//       "languages": ["eng", "fra"]
//     }
//   }
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
		return fmt.Errorf("unable to read environment: %s", err)
	}

	run.Session().SetEnvironment(env)
	return nil
}
