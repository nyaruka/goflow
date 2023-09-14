package flows

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/pkg/errors"
)

// OptIn adds some functionality to optin assets.
type OptIn struct {
	assets.OptIn

	Channel *Channel
}

// NewOptIn returns a new optin object from the given optin asset
func NewOptIn(channels *ChannelAssets, asset assets.OptIn) (*OptIn, error) {
	ch := channels.Get(asset.Channel().UUID)
	if ch == nil {
		return nil, errors.Errorf("no such channel with UUID %s", asset.Channel().UUID)
	}

	return &OptIn{Channel: ch}, nil
}

// Asset returns the underlying asset
func (o *OptIn) Asset() assets.OptIn { return o.OptIn }

// OptInAssets provides access to all optin assets
type OptInAssets struct {
	byUUID map[assets.OptInUUID]*OptIn
}

// NewOptInAssets creates a new set of optin assets
func NewOptInAssets(channels *ChannelAssets, optins []assets.OptIn) (*OptInAssets, []assets.OptIn) {
	broken := make([]assets.OptIn, 0)
	s := &OptInAssets{
		byUUID: make(map[assets.OptInUUID]*OptIn, len(optins)),
	}
	for _, asset := range optins {
		optin, err := NewOptIn(channels, asset)
		if err != nil {
			broken = append(broken, asset)
		} else {
			s.byUUID[asset.UUID()] = optin
		}
	}
	return s, broken
}

// Get finds the optin with the given UUID
func (s *OptInAssets) Get(uuid assets.OptInUUID) *OptIn {
	return s.byUUID[uuid]
}
