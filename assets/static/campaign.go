package static

import (
	"github.com/nyaruka/goflow/assets"
)

// Campaign is a JSON serializable implementation of a campaign asset
type Campaign struct {
	UUID_ assets.CampaignUUID `json:"uuid"  validate:"required,uuid"`
	Name_ string              `json:"name"  validate:"required"`
}

// NewCampaign creates a new campaign
func NewCampaign(uuid assets.CampaignUUID, name string) assets.Campaign {
	return &Campaign{
		UUID_: uuid,
		Name_: name,
	}
}

// UUID returns the UUID of this campaign
func (t *Campaign) UUID() assets.CampaignUUID { return t.UUID_ }

// Name returns the name of this campaign
func (t *Campaign) Name() string { return t.Name_ }
