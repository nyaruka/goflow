package flows

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
)

// Connection represents a connection to a specific channel using a specific URN
type Connection struct {
	channel *assets.ChannelReference
	urn     urns.URN
}

// NewConnection creates a new connection
func NewConnection(channel *assets.ChannelReference, urn urns.URN) *Connection {
	return &Connection{channel: channel, urn: urn}
}

func (c *Connection) Channel() *assets.ChannelReference { return c.channel }
func (c *Connection) URN() urns.URN                     { return c.urn }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type connectionEnvelope struct {
	Channel *assets.ChannelReference `json:"channel" validate:"required,dive"`
	URN     urns.URN                 `json:"urn" validate:"required,urn"`
}

func (c *Connection) UnmarshalJSON(data []byte) error {
	e := &connectionEnvelope{}
	if err := json.Unmarshal(data, e); err != nil {
		return err
	}

	c.channel = e.Channel
	c.urn = e.URN
	return nil
}

// MarshalJSON marshals this connection into JSON
func (c *Connection) MarshalJSON() ([]byte, error) {
	return json.Marshal(&connectionEnvelope{
		Channel: c.channel,
		URN:     c.urn,
	})
}
