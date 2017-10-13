package actions

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

const TypeSetPreferredChannel string = "set_preferred_channel"

type PreferredChannelAction struct {
	BaseAction
	Channel *flows.ChannelReference `json:"channel"`
}

func (a *PreferredChannelAction) Type() string { return TypeSetPreferredChannel }

func (a *PreferredChannelAction) Validate(assets flows.SessionAssets) error {
	_, err := assets.GetChannel(a.Channel.UUID)
	return err
}

func (a *PreferredChannelAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	// this is a no-op if we have no contact
	if run.Contact() == nil {
		return nil
	}

	log.Add(events.NewPreferredChannel(a.Channel))
	return nil
}
