package triggers

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeManual, ReadManualTrigger)
}

// TypeManual is the type for manually triggered sessions
const TypeManual string = "manual"

// ManualTrigger is used when a session was triggered manually by a user
//
// ```
//   {
//     "type": "manual",
//     "flow": {"uuid": "ea7d8b6b-a4b2-42c1-b9cf-c0370a95a721", "name": "Registration"},
//     "contact": {
//       "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
//       "name": "Bob"
//     },
//     "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//   }
// ```
type ManualTrigger struct {
	baseTrigger
}

// NewManualTrigger creates a new manual trigger
func NewManualTrigger(env utils.Environment, contact *flows.Contact, flow flows.Flow, params types.XValue, triggeredOn time.Time) flows.Trigger {
	return &ManualTrigger{baseTrigger{environment: env, contact: contact, flow: flow, triggeredOn: triggeredOn}}
}

// Type returns the type of this trigger
func (t *ManualTrigger) Type() string { return TypeManual }

// Resolve resolves the given key when this trigger is referenced in an expression
func (t *ManualTrigger) Resolve(env utils.Environment, key string) types.XValue {
	switch key {
	case "type":
		return types.NewXText(TypeManual)
	}

	return t.baseTrigger.Resolve(env, key)
}

// ToXJSON is called when this type is passed to @(json(...))
func (t *ManualTrigger) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, t, "type", "params").ToXJSON(env)
}

var _ flows.Trigger = (*ManualTrigger)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadManualTrigger reads a manual trigger
func ReadManualTrigger(session flows.Session, data json.RawMessage) (flows.Trigger, error) {
	trigger := ManualTrigger{}
	e := baseTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, &e); err != nil {
		return nil, err
	}

	if err := unmarshalBaseTrigger(session, &trigger.baseTrigger, &e); err != nil {
		return nil, err
	}

	return &trigger, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *ManualTrigger) MarshalJSON() ([]byte, error) {
	var envelope baseTriggerEnvelope

	if err := marshalBaseTrigger(&t.baseTrigger, &envelope); err != nil {
		return nil, err
	}

	return json.Marshal(envelope)
}
