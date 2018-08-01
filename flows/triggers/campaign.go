package triggers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeCampaign, ReadCampaignTrigger)
}

// TypeCampaign is the type for sessions triggered by campaign events
const TypeCampaign string = "campaign"

// Campaign describes the campaign that triggered the session
type Campaign struct {
	UUID string `json:"uuid" validate:"required,uuid4"`
	Name string `json:"name" validate:"required"`
}

// CampaignTrigger is used when a session was triggered by a campaign event
//
// ```
//   {
//     "type": "campaign",
//     "flow": {"uuid": "ea7d8b6b-a4b2-42c1-b9cf-c0370a95a721", "name": "Registration"},
//     "contact": {
//       "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
//       "name": "Bob"
//     },
//     "campaign": {"uuid": "58e9b092-fe42-4173-876c-ff45a14a24fe", "name": "New Mothers"},
//     "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//   }
// ```
type CampaignTrigger struct {
	baseTrigger
	Campaign *Campaign
}

// Type returns the type of this trigger
func (t *CampaignTrigger) Type() string { return TypeCampaign }

// Resolve resolves the given key when this trigger is referenced in an expression
func (t *CampaignTrigger) Resolve(env utils.Environment, key string) types.XValue {
	switch key {
	case "type":
		return types.NewXText(TypeCampaign)
	}

	return t.baseTrigger.Resolve(env, key)
}

// ToXJSON is called when this type is passed to @(json(...))
func (t *CampaignTrigger) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, t, "type", "params").ToXJSON(env)
}

var _ flows.Trigger = (*CampaignTrigger)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type campaignTriggerEnvelope struct {
	baseTriggerEnvelope
	Campaign *Campaign `json:"campaign" validate:"required,dive"`
}

// ReadCampaignTrigger reads a campaign trigger
func ReadCampaignTrigger(session flows.Session, data json.RawMessage) (flows.Trigger, error) {
	trigger := &CampaignTrigger{}
	e := campaignTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, &e); err != nil {
		return nil, err
	}

	if err := unmarshalBaseTrigger(session, &trigger.baseTrigger, &e.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	trigger.Campaign = e.Campaign

	return trigger, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *CampaignTrigger) MarshalJSON() ([]byte, error) {
	var envelope campaignTriggerEnvelope

	if err := marshalBaseTrigger(&t.baseTrigger, &envelope.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	envelope.Campaign = t.Campaign

	return json.Marshal(envelope)
}
