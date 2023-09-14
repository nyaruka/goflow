package static

import (
	"github.com/nyaruka/goflow/assets"
)

// OptIn is a JSON serializable implementation of an optin asset
type OptIn struct {
	UUID_    assets.OptInUUID         `json:"uuid" validate:"required,uuid"`
	Name_    string                   `json:"name" validate:"required"`
	Channel_ *assets.ChannelReference `json:"channel" validate:"required,dive"`
}

// NewOptIn creates a new topic
func NewOptIn(uuid assets.OptInUUID, name string, channel *assets.ChannelReference) assets.OptIn {
	return &OptIn{
		UUID_:    uuid,
		Name_:    name,
		Channel_: channel,
	}
}

// UUID returns the UUID of this ticketer
func (t *OptIn) UUID() assets.OptInUUID { return t.UUID_ }

// Name returns the name of this ticketer
func (t *OptIn) Name() string { return t.Name_ }

// Channel returns a reference to this optin's channel
func (c *OptIn) Channel() *assets.ChannelReference { return c.Channel_ }
