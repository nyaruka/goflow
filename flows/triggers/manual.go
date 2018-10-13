package triggers

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/goflow/assets"
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
//   {
//     "type": "manual",
//     "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
//     "contact": {
//       "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
//       "name": "Bob",
//       "created_on": "2018-01-01T12:00:00.000000Z"
//     },
//     "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//   }
//
// @trigger manual
type ManualTrigger struct {
	baseTrigger
}

// NewManualTrigger creates a new manual trigger
func NewManualTrigger(env utils.Environment, contact *flows.Contact, flow *assets.FlowReference, params types.XValue, triggeredOn time.Time) flows.Trigger {
	return &ManualTrigger{
		baseTrigger: newBaseTrigger(TypeManual, env, flow, contact, params, triggeredOn),
	}
}

var _ flows.Trigger = (*ManualTrigger)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadManualTrigger reads a manual trigger
func ReadManualTrigger(session flows.Session, data json.RawMessage) (flows.Trigger, error) {
	e := &baseTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	t := &ManualTrigger{}

	if err := t.unmarshal(session, e); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *ManualTrigger) MarshalJSON() ([]byte, error) {
	e := &baseTriggerEnvelope{}

	if err := t.marshal(e); err != nil {
		return nil, err
	}

	return json.Marshal(e)
}
