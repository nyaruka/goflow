package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
)

// CampaignUUID is the UUID of a campaign
type CampaignUUID uuids.UUID

// CampaignPointUUID is the type for campaign point UUIDs
type CampaignPointUUID uuids.UUID

// Campaign is a campaign of events.
//
//	{
//	  "uuid": "58e9b092-fe42-4173-876c-ff45a14a24fe",
//	  "name": "Reminders"
//	}
//
// @asset campaign
type Campaign interface {
	UUID() CampaignUUID
	Name() string
}

// CampaignReference is used to reference a campaign
type CampaignReference struct {
	UUID CampaignUUID `json:"uuid" validate:"required,uuid"`
	Name string       `json:"name"`
}

// NewCampaignReference creates a new campaign reference with the given UUID and name
func NewCampaignReference(uuid CampaignUUID, name string) *CampaignReference {
	return &CampaignReference{UUID: uuid, Name: name}
}

// Type returns the name of the asset type
func (r *CampaignReference) Type() string {
	return "campaign"
}

// GenericUUID returns the untyped UUID
func (r *CampaignReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

// Identity returns the unique identity of the asset
func (r *CampaignReference) Identity() string {
	return string(r.UUID)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *CampaignReference) Variable() bool {
	return false
}

func (r *CampaignReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ UUIDReference = (*CampaignReference)(nil)
