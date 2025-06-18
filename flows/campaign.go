package flows

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// Campaign adds some functionality to campaign assets.
type Campaign struct {
	assets.Campaign
}

// NewCampaign returns a new campaign object from the given campaign asset
func NewCampaign(asset assets.Campaign) *Campaign {
	return &Campaign{Campaign: asset}
}

// Asset returns the underlying asset
func (c *Campaign) Asset() assets.Campaign { return c.Campaign }

// Reference returns a reference to this campaign
func (c *Campaign) Reference() *assets.CampaignReference {
	if c == nil {
		return nil
	}
	return assets.NewCampaignReference(c.UUID(), c.Name())
}

// Context returns the properties available in expressions
//
//	uuid:text -> the UUID of the campaign
//	name:text -> the name of the campaign
//
// @context campaign
func (c *Campaign) Context(env envs.Environment) map[string]types.XValue {
	return map[string]types.XValue{
		"__default__": types.NewXText(string(c.Name())),
		"uuid":        types.NewXText(string(c.UUID())),
		"name":        types.NewXText(string(c.Name())),
	}
}

// CampaignAssets provides access to all campaign assets
type CampaignAssets struct {
	byUUID map[assets.CampaignUUID]*Campaign
}

// NewCampaignAssets creates a new set of campaign assets
func NewCampaignAssets(campaigns []assets.Campaign) *CampaignAssets {
	s := &CampaignAssets{
		byUUID: make(map[assets.CampaignUUID]*Campaign, len(campaigns)),
	}
	for _, asset := range campaigns {
		campaign := NewCampaign(asset)
		s.byUUID[asset.UUID()] = campaign
	}
	return s
}

// Get finds the campaign with the given UUID
func (s *CampaignAssets) Get(uuid assets.CampaignUUID) *Campaign {
	return s.byUUID[uuid]
}
