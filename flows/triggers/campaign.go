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
	RegisterType(TypeCampaign, ReadCampaignTrigger)
}

// TypeCampaign is the type for sessions triggered by campaign events
const TypeCampaign string = "campaign"

// Campaign describes the campaign that triggered the session
type Campaign struct {
	UUID string `json:"uuid" validate:"required,uuid4"`
	Name string `json:"name" validate:"required"`
}

// CampaignEvent describes the specific event in the campaign that triggered the session
type CampaignEvent struct {
	UUID     string   `json:"uuid" validate:"required,uuid4"`
	Campaign Campaign `json:"campaign" validate:"required,dive"`
}

// CampaignTrigger is used when a session was triggered by a campaign event
//
//   {
//     "type": "campaign",
//     "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
//     "contact": {
//       "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
//       "name": "Bob"
//     },
//     "event": {
//         "uuid": "34d16dbd-476d-4b77-bac3-9f3d597848cc",
//         "campaign": {"uuid": "58e9b092-fe42-4173-876c-ff45a14a24fe", "name": "New Mothers"}
//     },
//     "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//   }
//
// @trigger campaign
type CampaignTrigger struct {
	baseTrigger
	event *CampaignEvent
}

// NewCampaignTrigger creates a new campaign trigger with the passed in values
func NewCampaignTrigger(env utils.Environment, flow *assets.FlowReference, contact *flows.Contact, event *CampaignEvent, triggeredOn time.Time) *CampaignTrigger {
	return &CampaignTrigger{
		baseTrigger: baseTrigger{
			environment: env,
			flow:        flow,
			contact:     contact,
			triggeredOn: triggeredOn,
		},
		event: event,
	}
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
	Event *CampaignEvent `json:"event" validate:"required,dive"`
}

// ReadCampaignTrigger reads a campaign trigger
func ReadCampaignTrigger(session flows.Session, data json.RawMessage) (flows.Trigger, error) {
	e := &campaignTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	t := &CampaignTrigger{
		event: e.Event,
	}
	if err := t.unmarshal(session, &e.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *CampaignTrigger) MarshalJSON() ([]byte, error) {
	e := &campaignTriggerEnvelope{
		Event: t.event,
	}

	if err := t.marshal(&e.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	return json.Marshal(e)
}
