package triggers

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// TypeUser is a constant for incoming messages
const TypeUser string = "user"

// UserTrigger is used when a session was triggered manually by a user
//
// ```
//   {
//     "type": "user",
//     "flow": {"uuid": "ea7d8b6b-a4b2-42c1-b9cf-c0370a95a721", "name": "Registration"},
//     "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//   }
// ```
type UserTrigger struct {
	baseTrigger
}

// NewUserTrigger creates a new user trigger
func NewUserTrigger(flow flows.Flow, triggeredOn time.Time) flows.Trigger {
	return &UserTrigger{baseTrigger{flow: flow, triggeredOn: triggeredOn}}
}

// Type returns the type of this trigger
func (t *UserTrigger) Type() string { return TypeUser }

var _ flows.Trigger = (*UserTrigger)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func ReadUserTrigger(session flows.Session, envelope *utils.TypedEnvelope) (flows.Trigger, error) {
	trigger := UserTrigger{}
	e := baseTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(envelope.Data, &e, "trigger[type=user]"); err != nil {
		return nil, err
	}

	if err := readBaseTrigger(session, &trigger.baseTrigger, &e); err != nil {
		return nil, err
	}

	return &trigger, nil
}

func (t *UserTrigger) MarshalJSON() ([]byte, error) {
	var envelope baseTriggerEnvelope

	envelope.TriggeredOn = t.triggeredOn
	envelope.Flow = t.flow.Reference()

	return json.Marshal(envelope)
}
