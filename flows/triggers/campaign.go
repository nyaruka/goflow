package triggers

import (
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
//	      "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
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

func (b *Builder) CampaignFired(event *events.CampaignFired, campaign *flows.Campaign) *CampaignBuilder {
	return &CampaignBuilder{
		t: &Campaign{
			baseTrigger: newBaseTrigger(TypeCampaign, event, b.flow, false, nil),
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

func readCampaign(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &baseEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	t := &Campaign{}
	if err := t.unmarshal(sa, e, missing); err != nil {
		return nil, err
	}

	campEvt := t.event.(*events.CampaignFired)
	t.campaign = sa.Campaigns().Get(campEvt.Campaign.UUID)
	if t.campaign == nil {
		missing(campEvt.Campaign, nil)
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *Campaign) MarshalJSON() ([]byte, error) {
	e := &baseEnvelope{}

	if err := t.marshal(e); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
