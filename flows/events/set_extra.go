package events

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// TypeSetExtra is the type of our set extra event
const TypeSetExtra string = "set_extra"

// SetExtraEvent events are created to set extra context on a run
//
// ```
//   {
//     "type": "set_extra",
//     "created_on": "2006-01-02T15:04:05Z",
//     "extra": {
//       "name": "Bart",
//       "address": {
//         "street": "742 Evergreen Terrace",
//         "city": "Springfield"
//       }
//     }
//   }
// ```
//
// @event set_extra
type SetExtraEvent struct {
	BaseEvent
	Extra json.RawMessage `json:"extra"`
}

// Type returns the type of this event
func (e *SetExtraEvent) Type() string { return TypeSetExtra }

// Apply applies this event to the given run
func (e *SetExtraEvent) Apply(run flows.FlowRun) error {
	extra := utils.JSONFragment(e.Extra)

	run.SetExtra(extra)
	return nil
}
