package triggers

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// TypeManual is the type for manually triggered sessions
const TypeManual string = "manual"

// ManualTrigger is used when a session was triggered manually by a user
//
// ```
//   {
//     "type": "manual",
//     "flow": {"uuid": "ea7d8b6b-a4b2-42c1-b9cf-c0370a95a721", "name": "Registration"},
//     "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//   }
// ```
type ManualTrigger struct {
	baseTrigger
}

// NewManualTrigger creates a new manual trigger
func NewManualTrigger(flow flows.Flow, triggeredOn time.Time) flows.Trigger {
	return &ManualTrigger{baseTrigger{flow: flow, triggeredOn: triggeredOn}}
}

// Type returns the type of this trigger
func (t *ManualTrigger) Type() string { return TypeManual }

var _ flows.Trigger = (*ManualTrigger)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func ReadManualTrigger(session flows.Session, envelope *utils.TypedEnvelope) (flows.Trigger, error) {
	trigger := ManualTrigger{}
	e := baseTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(envelope.Data, &e, "trigger[type=manual]"); err != nil {
		return nil, err
	}

	if err := readBaseTrigger(session, &trigger.baseTrigger, &e); err != nil {
		return nil, err
	}

	return &trigger, nil
}

func (t *ManualTrigger) MarshalJSON() ([]byte, error) {
	var envelope baseTriggerEnvelope

	envelope.TriggeredOn = t.triggeredOn
	envelope.Flow = t.flow.Reference()

	return json.Marshal(envelope)
}
