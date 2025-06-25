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
	registerType(TypeCampaign, readCampaign)
}

// TypeCampaign is the type for sessions triggered by campaigns.
const TypeCampaign string = "campaign"

// Campaign is used when a session was triggered by a campaign.
//
//	{
//	  "type": "campaign",
//	  "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
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
type Campaign struct {
	baseTrigger
	event    *events.CampaignFired
	campaign *flows.Campaign
}

// Context for manual triggers always has non-nil params
func (t *Campaign) Context(env envs.Environment) map[string]types.XValue {
	c := t.context()
	c.campaign = flows.Context(env, t.campaign)
	return c.asMap()
}

var _ flows.Trigger = (*Campaign)(nil)

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// CampaignBuilder is a builder for campaign type triggers
type CampaignBuilder struct {
	t *Campaign
}

// Campaign returns a campaign trigger builder
func (b *Builder) Campaign(campaign *flows.Campaign, event *events.CampaignFired) *CampaignBuilder {
	return &CampaignBuilder{
		t: &Campaign{
			baseTrigger: newBaseTrigger(TypeCampaign, b.flow, false, nil),
			event:       event,
			campaign:    campaign,
		},
	}
}

// Build builds the trigger
func (b *CampaignBuilder) Build() *Campaign {
	return b.t
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type campaignEnvelope struct {
	baseEnvelope

	Event json.RawMessage `json:"event" validate:"required"`
}

func readCampaign(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &campaignEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	// older sessions will have events that aren't really events so fix 'em
	e.Event, _ = jsonparser.Set(e.Event, []byte(`"campaign_fired"`), "type")
	e.Event, _ = jsonparser.Set(e.Event, jsonx.MustMarshal(e.TriggeredOn), "created_on")

	event, err := events.Read(e.Event)
	if err != nil {
		return nil, fmt.Errorf("error reading campaign trigger event: %w", err)
	}

	campEvt := event.(*events.CampaignFired)
	campaign := sa.Campaigns().Get(campEvt.Campaign.UUID)
	if campaign == nil {
		missing(campEvt.Campaign, nil)
	}

	t := &Campaign{
		event:    campEvt,
		campaign: campaign,
	}
	if err := t.unmarshal(sa, &e.baseEnvelope, missing); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *Campaign) MarshalJSON() ([]byte, error) {
	me, err := json.Marshal(t.event)
	if err != nil {
		return nil, fmt.Errorf("error marshaling campaign trigger event: %w", err)
	}

	e := &campaignEnvelope{
		Event: me,
	}

	if err := t.marshal(&e.baseEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
