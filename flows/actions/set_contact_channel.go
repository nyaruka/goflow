package actions

import (
	"context"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/modifiers"
)

func init() {
	registerType(TypeSetContactChannel, func() flows.Action { return &SetContactChannel{} })
}

// TypeSetContactChannel is the type for the set contact channel action
const TypeSetContactChannel string = "set_contact_channel"

// SetContactChannel can be used to change or clear the preferred channel of the current contact.
//
// Because channel affinity is a property of a contact's URNs, a [event:contact_urns_changed] event will be created if any
// changes are made to the contact's URNs.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "set_contact_channel",
//	  "channel": {"uuid": "4bb288a0-7fca-4da1-abe8-59a593aff648", "name": "Facebook Channel"}
//	}
//
// @action set_contact_channel
type SetContactChannel struct {
	baseAction
	onlineAction

	Channel *assets.ChannelReference `json:"channel" validate:"omitempty"`
}

// NewSetContactChannel creates a new set channel action
func NewSetContactChannel(uuid flows.ActionUUID, channel *assets.ChannelReference) *SetContactChannel {
	return &SetContactChannel{
		baseAction: newBaseAction(TypeSetContactChannel, uuid),
		Channel:    channel,
	}
}

// Execute runs our action
func (a *SetContactChannel) Execute(ctx context.Context, run flows.Run, step flows.Step, log flows.EventLogger) error {
	var channel *flows.Channel
	if a.Channel != nil {
		channel = run.Session().Assets().Channels().Get(a.Channel.UUID)
		if channel == nil {
			log(events.NewDependencyError(a.Channel))
			return nil
		}
	}

	_, err := a.applyModifier(run, modifiers.NewChannel(channel), log)
	return err
}

func (a *SetContactChannel) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	if a.Channel != nil {
		dependency(a.Channel)
	}
}
