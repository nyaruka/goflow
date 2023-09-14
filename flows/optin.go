package flows

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// OptIn adds some functionality to optin assets.
type OptIn struct {
	assets.OptIn
}

// NewOptIn returns a new optin object from the given optin asset
func NewOptIn(asset assets.OptIn) *OptIn {
	return &OptIn{OptIn: asset}
}

// Asset returns the underlying asset
func (o *OptIn) Asset() assets.OptIn { return o.OptIn }

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
//
// @context optin
func (o *OptIn) Context(env envs.Environment) map[string]types.XValue {
	return map[string]types.XValue{
		"__default__": types.NewXText(string(o.Name())),
		"uuid":        types.NewXText(string(o.UUID())),
		"name":        types.NewXText(string(o.Name())),
	}
}

// OptInAssets provides access to all optin assets
type OptInAssets struct {
	byUUID map[assets.OptInUUID]*OptIn
}

// NewOptInAssets creates a new set of optin assets
func NewOptInAssets(optins []assets.OptIn) *OptInAssets {
	s := &OptInAssets{
		byUUID: make(map[assets.OptInUUID]*OptIn, len(optins)),
	}
	for _, asset := range optins {
		optin := NewOptIn(asset)
		s.byUUID[asset.UUID()] = optin
	}
	return s
}

// Get finds the optin with the given UUID
func (s *OptInAssets) Get(uuid assets.OptInUUID) *OptIn {
	return s.byUUID[uuid]
}
