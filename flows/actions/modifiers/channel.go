package modifiers

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	RegisterType(TypeChannel, func() Modifier { return &ChannelModifier{} })
}

// TypeChannel is the type of our channel modifier
const TypeChannel string = "channel"

// ChannelModifier modifies the channel of a contact
type ChannelModifier struct {
	baseModifier

	Channel *flows.Channel
}

// NewChannelModifier creates a new channel modifier
func NewChannelModifier(channel *flows.Channel) *ChannelModifier {
	return &ChannelModifier{
		baseModifier: newBaseModifier(TypeChannel),
		Channel:      channel,
	}
}

// Apply applies this modification to the given contact
func (m *ChannelModifier) Apply(assets flows.SessionAssets, contact *flows.Contact, log func(flows.Event)) bool {
	// if URNs change in anyway, generate a URNs changed event
	if contact.UpdatePreferredChannel(m.Channel) {
		log(events.NewContactURNsChangedEvent(contact.URNs().RawURNs()))
		return true
	}
	return false
}

var _ Modifier = (*ChannelModifier)(nil)
