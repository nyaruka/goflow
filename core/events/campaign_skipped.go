package events

import (
	"github.com/nyaruka/goflow/assets"
)

func init() {
	registerType(TypeCampaignSkipped, func() Event { return &CampaignSkipped{} })
}

// TypeCampaignSkipped is our type for the campaign skipped event
const TypeCampaignSkipped string = "campaign_skipped"

// CampaignSkipped events are created when a campaign point that should have fired for the contact was skipped because
// they were already in a flow.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "campaign_skipped",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "campaign": {
//	    "uuid": "58e9b092-fe42-4173-876c-ff45a14a24fe",
//	    "name": "Reminders"
//	  },
//	  "point_uuid": "77410ef5-1b3d-4571-b8c3-a692b07d2d09"
//	}
//
// @event campaign_skipped
type CampaignSkipped struct {
	BaseEvent

	Campaign  *assets.CampaignReference `json:"campaign"   validate:"required"`
	PointUUID assets.CampaignPointUUID  `json:"point_uuid" validate:"required"`
}

// NewCampaignSkipped returns a new campaign skipped event
func NewCampaignSkipped(campaign *assets.CampaignReference, pointUUID assets.CampaignPointUUID) *CampaignSkipped {
	return &CampaignSkipped{
		BaseEvent: NewBaseEvent(TypeCampaignSkipped),
		Campaign:  campaign,
		PointUUID: pointUUID,
	}
}
