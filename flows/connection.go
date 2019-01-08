package flows

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
)

// Connection represents a connection to a specific channel using a specific URN
type Connection struct {
	channel *assets.ChannelReference `json:"channel" validate:"required,dive"`
	urn     urns.URN                 `json:"urn" validate:"required,urn"`
}

// NewConnection creates a new connection
func NewConnection(channel *assets.ChannelReference, urn urns.URN) *Connection {
	return &Connection{channel: channel, urn: urn}
}

func (c *Connection) Channel() *assets.ChannelReference { return c.channel }
func (c *Connection) URN() urns.URN                     { return c.urn }
