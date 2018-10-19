package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	RegisterType(TypeSetContactChannel, func() flows.Action { return &SetContactChannelAction{} })
}

// TypeSetContactChannel is the type for the set contact channel action
const TypeSetContactChannel string = "set_contact_channel"

// SetContactChannelAction can be used to update or cledar the preferred channel of the current contact.
//
// A [event:contact_channel_changed] event will be created with the set channel.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "set_contact_channel",
//     "channel": {"uuid": "4bb288a0-7fca-4da1-abe8-59a593aff648", "name": "FAcebook Channel"}
//   }
//
// @action set_contact_channel
type SetContactChannelAction struct {
	BaseAction
	onlineAction

	Channel *assets.ChannelReference `json:"channel"`
}

// NewSetContactChannelAction creates a new set channel action
func NewSetContactChannelAction(uuid flows.ActionUUID, channel *assets.ChannelReference) *SetContactChannelAction {
	return &SetContactChannelAction{
		BaseAction: NewBaseAction(TypeSetContactChannel, uuid),
		Channel:    channel,
	}
}

// Validate validates our action is valid and has all the assets it needs
func (a *SetContactChannelAction) Validate(assets flows.SessionAssets, context *flows.ValidationContext) error {
	if a.Channel != nil {
		_, err := assets.Channels().Get(a.Channel.UUID)
		return err
	}
	return nil
}

func (a *SetContactChannelAction) Execute(run flows.FlowRun, step flows.Step) error {
	if run.Contact() == nil {
		a.logError(run, step, fmt.Errorf("can't execute action in session without a contact"))
		return nil
	}

	var channel *flows.Channel
	var err error
	if a.Channel != nil {
		channel, err = run.Session().Assets().Channels().Get(a.Channel.UUID)
		if err != nil {
			return err
		}
	}

	if run.Contact().PreferredChannel() != channel {
		run.Contact().UpdatePreferredChannel(channel)
		a.log(run, step, events.NewContactChannelChangedEvent(channel))
	}
	return nil
}
