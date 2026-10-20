package events

import (
	"github.com/nyaruka/goflow/assets"
)

func init() {
	registerType(TypeCampaignSkipped, func() Event { return &CampaignSkipped{} })
}

// TypeCampaignSkipped is our type for the campaign skipped event
const TypeCampaignSkipped string = "campaign_skipped"

// CampaignSkippedReason is the reason a campaign point was skipped for a contact
type CampaignSkippedReason string

const (
	CampaignSkippedReasonExistingSession CampaignSkippedReason = "existing_session"
	CampaignSkippedReasonContactStatus   CampaignSkippedReason = "contact_status"
)

// CampaignSkipped events are created when a campaign point that should have fired for the contact was skipped.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "campaign_skipped",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "campaign": {
//	    "uuid": "58e9b092-fe42-4173-876c-ff45a14a24fe",
//	    "name": "Reminders"
//	  },
//	  "point_uuid": "77410ef5-1b3d-4571-b8c3-a692b07d2d09",
//	  "reason": "existing_session"
//	}
//
// @event campaign_skipped
type CampaignSkipped struct {
	BaseEvent

	Campaign  *assets.CampaignReference `json:"campaign"   validate:"required"`
	PointUUID assets.CampaignPointUUID  `json:"point_uuid" validate:"required"`
	Reason    CampaignSkippedReason     `json:"reason"     validate:"required,eq=existing_session|eq=contact_status"`
}

// NewCampaignSkipped returns a new campaign skipped event
func NewCampaignSkipped(campaign *assets.CampaignReference, pointUUID assets.CampaignPointUUID, reason CampaignSkippedReason) *CampaignSkipped {
	return &CampaignSkipped{
		BaseEvent: NewBaseEvent(TypeCampaignSkipped),
		Campaign:  campaign,
		PointUUID: pointUUID,
		Reason:    reason,
	}
}
