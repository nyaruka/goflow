package flows

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/pkg/errors"
)

// OptIn adds some functionality to optin assets.
type OptIn struct {
	assets.OptIn

	channel *Channel
}

// NewOptIn returns a new optin object from the given optin asset
func NewOptIn(channels *ChannelAssets, asset assets.OptIn) (*OptIn, error) {
	ch := channels.Get(asset.Channel().UUID)
	if ch == nil {
		return nil, errors.Errorf("no such channel with UUID %s", asset.Channel().UUID)
	}

	return &OptIn{OptIn: asset, channel: ch}, nil
}

// Asset returns the underlying asset
func (o *OptIn) Asset() assets.OptIn { return o.OptIn }

// Channel returns the associated channel
func (o *OptIn) Channel() *Channel { return o.channel }

// Reference returns a reference to this optin
func (o *OptIn) Reference() *assets.OptInReference {
	if o == nil {
		return nil
	}
	return assets.NewOptInReference(o.UUID(), o.Name())
}

// Context returns the properties available in expressions
//
//	uuid:text -> the UUID of the optin
//	name:text -> the name of the optin
//	channel:channel -> the channel of the optin
//
// @context ticket
func (o *OptIn) Context(env envs.Environment) map[string]types.XValue {
	return map[string]types.XValue{
		"__default__": types.NewXText(string(o.Name())),
		"uuid":        types.NewXText(string(o.UUID())),
		"name":        types.NewXText(string(o.Name())),
		"channel":     Context(env, o.channel),
	}
}

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
