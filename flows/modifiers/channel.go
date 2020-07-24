package modifiers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/jsonx"
)

func init() {
	registerType(TypeChannel, readChannelModifier)
}

// TypeChannel is the type of our channel modifier
const TypeChannel string = "channel"

// ChannelModifier modifies the preferred channel of a contact
type ChannelModifier struct {
	baseModifier

	channel *flows.Channel
}

// NewChannel creates a new channel modifier
func NewChannel(channel *flows.Channel) *ChannelModifier {
	return &ChannelModifier{
		baseModifier: newBaseModifier(TypeChannel),
		channel:      channel,
	}
}

// Apply applies this modification to the given contact
func (m *ChannelModifier) Apply(env envs.Environment, assets flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) {
	// if URNs change in anyway, generate a URNs changed event
	if contact.UpdatePreferredChannel(m.channel) {
		log(events.NewContactURNsChanged(contact.URNs().RawURNs()))
	}
}

var _ flows.Modifier = (*ChannelModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type channelModifierEnvelope struct {
	utils.TypedEnvelope
	Channel *assets.ChannelReference `json:"channel" validate:"omitempty,dive"`
}

func readChannelModifier(assets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &channelModifierEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	var channel *flows.Channel
	if e.Channel != nil {
		channel = assets.Channels().Get(e.Channel.UUID)
		if channel == nil {
			missing(e.Channel, nil)
			return nil, ErrNoModifier // nothing left to modify without the channel
		}
	}
	return NewChannel(channel), nil
}

func (m *ChannelModifier) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&channelModifierEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Channel:       m.channel.Reference(),
	})
}
