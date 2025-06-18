package triggers

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeCampaign, readCampaignTrigger)
}

// TypeCampaign is the type for sessions triggered by campaigns.
const TypeCampaign string = "campaign"

// CampaignEvent describes the specific point in the campaign that triggered the session.
type CampaignEvent struct {
	UUID     assets.CampaignPointUUID  `json:"uuid" validate:"required,uuid"`
	Campaign *assets.CampaignReference `json:"campaign" validate:"required"`
}

// CampaignTrigger is used when a session was triggered by a campaign.
//
//	{
//	  "type": "campaign",
//	  "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
//	  "contact": {
//	    "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
//	    "name": "Bob",
//	    "created_on": "2018-01-01T12:00:00.000000Z"
//	  },
//	  "event": {
//	      "uuid": "34d16dbd-476d-4b77-bac3-9f3d597848cc",
//	      "campaign": {"uuid": "58e9b092-fe42-4173-876c-ff45a14a24fe", "name": "Reminders"}
//	  },
//	  "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//	}
//
// @trigger campaign
type CampaignTrigger struct {
	baseTrigger
	event    *CampaignEvent
	campaign *flows.Campaign
}

// Context for manual triggers always has non-nil params
func (t *CampaignTrigger) Context(env envs.Environment) map[string]types.XValue {
	c := t.context()
	c.campaign = flows.Context(env, t.campaign)
	return c.asMap()
}

var _ flows.Trigger = (*CampaignTrigger)(nil)

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// CampaignBuilder is a builder for campaign type triggers
type CampaignBuilder struct {
	t *CampaignTrigger
}

// Campaign returns a campaign trigger builder
func (b *Builder) Campaign(campaign *flows.Campaign, pointUUID assets.CampaignPointUUID) *CampaignBuilder {
	return &CampaignBuilder{
		t: &CampaignTrigger{
			baseTrigger: newBaseTrigger(TypeCampaign, b.environment, b.flow, b.contact, nil, false, nil),
			event:       &CampaignEvent{UUID: pointUUID, Campaign: campaign.Reference()},
			campaign:    campaign,
		},
	}
}

// Build builds the trigger
func (b *CampaignBuilder) Build() *CampaignTrigger {
	return b.t
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type campaignTriggerEnvelope struct {
	baseTriggerEnvelope
	Event *CampaignEvent `json:"event" validate:"required"`
}

func readCampaignTrigger(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &campaignTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	campaign := sa.Campaigns().Get(e.Event.Campaign.UUID)
	if campaign == nil {
		missing(e.Event.Campaign, nil)
	}

	t := &CampaignTrigger{
		event:    e.Event,
		campaign: campaign,
	}
	if err := t.unmarshal(sa, &e.baseTriggerEnvelope, missing); err != nil {
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
