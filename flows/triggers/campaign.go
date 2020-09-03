package triggers

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeCampaign, readCampaignTrigger)
}

// TypeCampaign is the type for sessions triggered by campaign events
const TypeCampaign string = "campaign"

// CampaignUUID is the type for campaign UUIDs
type CampaignUUID uuids.UUID

// CampaignEventUUID is the type for campaign event UUIDs
type CampaignEventUUID uuids.UUID

// CampaignReference is a reference to the campaign that triggered the session
type CampaignReference struct {
	UUID CampaignUUID `json:"uuid" validate:"required,uuid4"`
	Name string       `json:"name" validate:"required"`
}

// NewCampaignReference creates a new campaign reference
func NewCampaignReference(uuid CampaignUUID, name string) *CampaignReference {
	return &CampaignReference{UUID: uuid, Name: name}
}

// CampaignEvent describes the specific event in the campaign that triggered the session
type CampaignEvent struct {
	UUID     CampaignEventUUID  `json:"uuid" validate:"required,uuid4"`
	Campaign *CampaignReference `json:"campaign" validate:"required,dive"`
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

var _ flows.Trigger = (*CampaignTrigger)(nil)

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// CampaignBuilder is a builder for campaign type triggers
type CampaignBuilder struct {
	t *CampaignTrigger
}

// Campaign returns a campaign trigger builder
func (b *Builder) Campaign(campaign *CampaignReference, eventUUID CampaignEventUUID) *CampaignBuilder {
	return &CampaignBuilder{
		t: &CampaignTrigger{
			baseTrigger: newBaseTrigger(TypeCampaign, b.environment, b.flow, b.contact, nil, false, nil),
			event:       &CampaignEvent{UUID: eventUUID, Campaign: campaign},
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
