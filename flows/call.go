package flows

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
)

// Call represents a call over a specific channel and URN
type Call struct {
	channel *assets.ChannelReference
	urn     urns.URN
}

// NewCall creates a new call
func NewCall(channel *assets.ChannelReference, urn urns.URN) *Call {
	return &Call{channel: channel, urn: urn}
}

// Channel returns a reference to the channel
func (c *Call) Channel() *assets.ChannelReference { return c.channel }

// URN returns the URN
func (c *Call) URN() urns.URN { return c.urn }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type callEnvelope struct {
	Channel *assets.ChannelReference `json:"channel" validate:"required"`
	URN     urns.URN                 `json:"urn" validate:"required,urn"`
}

// UnmarshalJSON unmarshals a call from JSON
func (c *Call) UnmarshalJSON(data []byte) error {
	e := &callEnvelope{}
	if err := jsonx.Unmarshal(data, e); err != nil {
		return err
	}

	c.channel = e.Channel
	c.urn = e.URN
	return nil
}

// MarshalJSON marshals this call into JSON
func (c *Call) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&callEnvelope{
		Channel: c.channel,
		URN:     c.urn,
	})
}
