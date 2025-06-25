package modifiers

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeChannel, readChannel)
}

// TypeChannel is the type of our channel modifier
const TypeChannel string = "channel"

// Channel modifies the preferred channel of a contact
type Channel struct {
	baseModifier

	channel *flows.Channel
}

// NewChannel creates a new channel modifier
func NewChannel(channel *flows.Channel) *Channel {
	return &Channel{
		baseModifier: newBaseModifier(TypeChannel),
		channel:      channel,
	}
}

// Apply applies this modification to the given contact
func (m *Channel) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) bool {
	if m.channel != nil && !m.channel.HasRole(assets.ChannelRoleSend) {
		log(events.NewError("can't set channel that can't send as the preferred channel"))

	} else if contact.UpdatePreferredChannel(m.channel) {
		// if URNs change in anyway, generate a URNs changed event
		log(events.NewContactURNsChanged(contact.URNs().RawURNs()))
		return true
	}
	return false
}

var _ flows.Modifier = (*Channel)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type channelEnvelope struct {
	utils.TypedEnvelope

	Channel *assets.ChannelReference `json:"channel" validate:"omitempty"`
}

func readChannel(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &channelEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	var channel *flows.Channel
	if e.Channel != nil {
		channel = sa.Channels().Get(e.Channel.UUID)
		if channel == nil {
			missing(e.Channel, nil)
			return nil, ErrNoModifier // nothing left to modify without the channel
		}
	}
	return NewChannel(channel), nil
}

func (m *Channel) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&channelEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Channel:       m.channel.Reference(),
	})
}
