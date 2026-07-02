package flows

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/events"
)

// Call represents a call over a specific channel and URN
type Call struct {
	uuid    CallUUID
	channel *Channel
	urn     urns.URN
}

// NewCall creates a new call
func NewCall(uuid CallUUID, channel *Channel, urn urns.URN) *Call {
	return &Call{uuid: uuid, channel: channel, urn: urn}
}

func (c *Call) UUID() CallUUID    { return c.uuid }
func (c *Call) Channel() *Channel { return c.channel }
func (c *Call) URN() urns.URN     { return c.urn }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadCall reads a call from the passed in envelope.
func ReadCall(sa SessionAssets, e *events.CallEnvelope, missing assets.MissingCallback) *Call {
	var channel *Channel
	if e.Channel != nil {
		channel = sa.Channels().Get(e.Channel.UUID)
		if channel == nil {
			missing(e.Channel, nil)
		}
	}

	return &Call{uuid: e.UUID, channel: channel, urn: e.URN}
}

// Marshal marshals a call into an envelope.
func (c *Call) Marshal() *CallEnvelope {
	return &CallEnvelope{
		UUID:    c.uuid,
		Channel: c.channel.Reference(),
		URN:     c.urn,
	}
}
