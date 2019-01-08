package flows

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// Connection represents a connection to a specific channel using a specific URN
type Connection struct {
	channel *Channel
	urn     urns.URN
}

// NewConnection creates a new connection
func NewConnection(channel *Channel, urn urns.URN) *Connection {
	return &Connection{channel: channel, urn: urn}
}

func (c *Connection) Channel() *Channel { return c.channel }
func (c *Connection) URN() urns.URN     { return c.urn }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type connectionEnvelope struct {
	Channel *assets.ChannelReference `json:"channel" validate:"required,dive"`
	URN     urns.URN                 `json:"urn" validate:"urn"`
}

// ReadConnection decodes a connection from the passed in JSON
func ReadConnection(assets SessionAssets, data json.RawMessage) (*Connection, error) {
	e := &connectionEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	channel, err := assets.Channels().Get(e.Channel.UUID)
	if err != nil {
		return nil, err
	}

	return &Connection{channel: channel, urn: e.URN}, nil
}

// MarshalJSON marshals this connection into JSON
func (c *Connection) MarshalJSON() ([]byte, error) {
	return json.Marshal(&connectionEnvelope{Channel: c.channel.Reference(), URN: c.urn})
}
