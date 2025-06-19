package triggers

import (
	"encoding/json"
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeCampaign, readCampaignTrigger)
}

// TypeCampaign is the type for sessions triggered by campaigns.
const TypeCampaign string = "campaign"

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
//	      "type": "campaign_fired",
//	      "created_on": "2006-01-02T15:04:05Z",
//	      "campaign": {"uuid": "58e9b092-fe42-4173-876c-ff45a14a24fe", "name": "Reminders"},
//	      "point_uuid": "34d16dbd-476d-4b77-bac3-9f3d597848cc"
//	  },
//	  "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//	}
//
// @trigger campaign
type CampaignTrigger struct {
	baseTrigger
	event    *events.CampaignFiredEvent
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
func (b *Builder) Campaign(campaign *flows.Campaign, event *events.CampaignFiredEvent) *CampaignBuilder {
	return &CampaignBuilder{
		t: &CampaignTrigger{
			baseTrigger: newBaseTrigger(TypeCampaign, b.environment, b.flow, b.contact, nil, false, nil),
			event:       event,
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
	Event json.RawMessage `json:"event" validate:"required"`
}

func readCampaignTrigger(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &campaignTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	// older sessions will have events that aren't really events so fix 'em
	e.Event, _ = jsonparser.Set(e.Event, []byte(`"campaign_fired"`), "type")
	e.Event, _ = jsonparser.Set(e.Event, jsonx.MustMarshal(e.TriggeredOn), "created_on")

	event, err := events.ReadEvent(e.Event)
	if err != nil {
		return nil, fmt.Errorf("error reading campaign trigger event: %w", err)
	}

	campEvt := event.(*events.CampaignFiredEvent)
	campaign := sa.Campaigns().Get(campEvt.Campaign.UUID)
	if campaign == nil {
		missing(campEvt.Campaign, nil)
	}

	t := &CampaignTrigger{
		event:    campEvt,
		campaign: campaign,
	}
	if err := t.unmarshal(sa, &e.baseTriggerEnvelope, missing); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *CampaignTrigger) MarshalJSON() ([]byte, error) {
	me, err := json.Marshal(t.event)
	if err != nil {
		return nil, fmt.Errorf("error marshaling campaign trigger event: %w", err)
	}

	e := &campaignTriggerEnvelope{
		Event: me,
	}

	if err := t.marshal(&e.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
