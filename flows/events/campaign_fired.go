package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeCampaignFired, func() flows.Event { return &CampaignFired{} })
}

// TypeCampaignFired is our type for the campaign fired event
const TypeCampaignFired string = "campaign_fired"

// CampaignFired events are created when a campaign has been fired.
//
//	{
//	  "type": "campaign_fired",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "campaign": {
//	    "uuid": "58e9b092-fe42-4173-876c-ff45a14a24fe",
//	    "name": "Reminders"
//	  },
//	  "point_uuid": "77410ef5-1b3d-4571-b8c3-a692b07d2d09"
//	}
//
// @event campaign_fired
type CampaignFired struct {
	BaseEvent

	Campaign  *assets.CampaignReference `json:"campaign" validate:"required"`
	PointUUID assets.CampaignPointUUID  `json:"point_uuid"` // TODO make required
}

// NewCampaignFired returns a new campaign fired event
func NewCampaignFired(campaign *flows.Campaign, pointUUID assets.CampaignPointUUID) *CampaignFired {
	return &CampaignFired{
		BaseEvent: NewBaseEvent(TypeCampaignFired),
		Campaign:  campaign.Reference(),
		PointUUID: pointUUID,
	}
}
