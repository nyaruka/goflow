package triggers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/jsonx"
)

func init() {
	registerType(TypeManual, readManualTrigger)
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

// NewManual creates a new manual trigger
func NewManual(env envs.Environment, flow *assets.FlowReference, contact *flows.Contact, batch bool, params *types.XObject) flows.Trigger {
	return &ManualTrigger{
		baseTrigger: newBaseTrigger(TypeManual, env, flow, contact, nil, batch, params),
	}
}

// NewManualVoice creates a new manual trigger with a channel connection for voice
func NewManualVoice(env envs.Environment, flow *assets.FlowReference, contact *flows.Contact, connection *flows.Connection, batch bool, params *types.XObject) flows.Trigger {
	return &ManualTrigger{
		baseTrigger: newBaseTrigger(TypeManual, env, flow, contact, connection, batch, params),
	}
}

// Context for manual triggers always has non-nil params
func (t *ManualTrigger) Context(env envs.Environment) map[string]types.XValue {
	params := t.params
	if params == nil {
		params = types.XObjectEmpty
	}

	return map[string]types.XValue{
		"type":    types.NewXText(t.type_),
		"params":  params,
		"keyword": nil,
	}
}

var _ flows.Trigger = (*ManualTrigger)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readManualTrigger(sessionAssets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &baseTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	t := &ManualTrigger{}

	if err := t.unmarshal(sessionAssets, e, missing); err != nil {
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

	return jsonx.Marshal(e)
}
