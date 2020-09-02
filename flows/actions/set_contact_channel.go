package actions

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/modifiers"
)

func init() {
	registerType(TypeSetContactChannel, func() flows.Action { return &SetContactChannelAction{} })
}

// TypeSetContactChannel is the type for the set contact channel action
const TypeSetContactChannel string = "set_contact_channel"

// SetContactChannelAction can be used to change or clear the preferred channel of the current contact.
//
// Because channel affinity is a property of a contact's URNs, a [event:contact_urns_changed] event will be created if any
// changes are made to the contact's URNs.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "set_contact_channel",
//     "channel": {"uuid": "4bb288a0-7fca-4da1-abe8-59a593aff648", "name": "Facebook Channel"}
//   }
//
// @action set_contact_channel
type SetContactChannelAction struct {
	baseAction
	onlineAction

	Channel *assets.ChannelReference `json:"channel" validate:"omitempty,dive"`
}

// NewSetContactChannel creates a new set channel action
func NewSetContactChannel(uuid flows.ActionUUID, channel *assets.ChannelReference) *SetContactChannelAction {
	return &SetContactChannelAction{
		baseAction: newBaseAction(TypeSetContactChannel, uuid),
		Channel:    channel,
	}
}

// Execute runs our action
func (a *SetContactChannelAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	contact := run.Contact()
	if contact == nil {
		logEvent(events.NewErrorf("can't execute action in session without a contact"))
		return nil
	}

	var channel *flows.Channel
	if a.Channel != nil {
		channel = run.Session().Assets().Channels().Get(a.Channel.UUID)
		if channel == nil {
			logEvent(events.NewDependencyError(a.Channel))
			return nil
		}
	}

	a.applyModifier(run, modifiers.NewChannel(channel), logModifier, logEvent)
	return nil
}
