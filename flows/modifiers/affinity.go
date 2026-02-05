package modifiers

import (
	"context"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeAffinity, readAffinity)
}

// TypeAffinity is the type of our affinity modifier
const TypeAffinity string = "affinity"

// Affinity modifies the preferred URN and channel of a contact
type Affinity struct {
	baseModifier

	urn     urns.URN
	channel *flows.Channel
}

// NewAffinity creates a new affinity modifier
func NewAffinity(urn urns.URN, channel *flows.Channel) *Affinity {
	return &Affinity{
		baseModifier: newBaseModifier(TypeAffinity),
		urn:          urn,
		channel:      channel,
	}
}

// Apply applies this modification to the given contact
func (m *Affinity) Apply(ctx context.Context, eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventLogger) (bool, error) {
	if contact.SetAffinity(m.urn, m.channel) {
		// if URNs change in anyway, generate a URNs changed event
		log(events.NewContactURNsChanged(contact.URNs().Encode()))
		return true, nil
	}
	return false, nil
}

var _ flows.Modifier = (*Affinity)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type affinityEnvelope struct {
	utils.TypedEnvelope

	URN     urns.URN                 `json:"urn"     validate:"required,urn"`
	Channel *assets.ChannelReference `json:"channel" validate:"required"`
}

func readAffinity(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &affinityEnvelope{}
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

	return NewAffinity(e.URN, channel), nil
}

func (m *Affinity) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&affinityEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		URN:           m.urn,
		Channel:       m.channel.Reference(),
	})
}
