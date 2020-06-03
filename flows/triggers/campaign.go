package triggers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/jsonx"
)

func init() {
	registerType(TypeCampaign, readCampaignTrigger)
}

// TypeCampaign is the type for sessions triggered by campaign events
const TypeCampaign string = "campaign"

// CampaignReference is a reference to the campaign that triggered the session
type CampaignReference struct {
	UUID string `json:"uuid" validate:"required,uuid4"`
	Name string `json:"name" validate:"required"`
}

// NewCampaignReference creates a new campaign reference
func NewCampaignReference(uuid, name string) *CampaignReference {
	return &CampaignReference{UUID: uuid, Name: name}
}

// CampaignEvent describes the specific event in the campaign that triggered the session
type CampaignEvent struct {
	UUID     string             `json:"uuid" validate:"required,uuid4"`
	Campaign *CampaignReference `json:"campaign" validate:"required,dive"`
}

// NewCampaignEvent creates a new campaign event
func NewCampaignEvent(uuid string, campaign *CampaignReference) *CampaignEvent {
	return &CampaignEvent{UUID: uuid, Campaign: campaign}
}

// CampaignTrigger is used when a session was triggered by a campaign event
//
//   {
//     "type": "campaign",
//     "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
//     "contact": {
//       "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
//       "name": "Bob",
//       "created_on": "2018-01-01T12:00:00.000000Z"
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

// NewCampaign creates a new campaign trigger with the passed in values
func NewCampaign(env envs.Environment, flow *assets.FlowReference, contact *flows.Contact, event *CampaignEvent) *CampaignTrigger {
	return &CampaignTrigger{
		baseTrigger: newBaseTrigger(TypeCampaign, env, flow, contact, nil, false, nil),
		event:       event,
	}
}

var _ flows.Trigger = (*CampaignTrigger)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type campaignTriggerEnvelope struct {
	baseTriggerEnvelope
	Event *CampaignEvent `json:"event" validate:"required,dive"`
}

func readCampaignTrigger(sessionAssets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &campaignTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	t := &CampaignTrigger{
		event: e.Event,
	}
	if err := t.unmarshal(sessionAssets, &e.baseTriggerEnvelope, missing); err != nil {
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

	return jsonx.Marshal(e)
}
