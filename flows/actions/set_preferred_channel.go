package actions

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

const TypeSetPreferredChannel string = "set_preferred_channel"

type PreferredChannelAction struct {
	BaseAction
	ChannelUUID flows.ChannelUUID `json:"channel_uuid"`
	ChannelName string            `json:"channel_name"`
}

func (a *PreferredChannelAction) Type() string { return TypeSetPreferredChannel }

func (a *PreferredChannelAction) Validate(assets flows.Assets) error {
	_, err := assets.GetChannel(a.ChannelUUID)
	return err
}

func (a *PreferredChannelAction) Execute(run flows.FlowRun, step flows.Step) error {
	if run.Contact() == nil {
		return nil
	}

	run.ApplyEvent(step, a, events.NewPreferredChannel(a.ChannelUUID, a.ChannelName))
	return nil
}
