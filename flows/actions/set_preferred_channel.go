package actions

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

const TypeSetPreferredChannel string = "set_preferred_channel"

type PreferredChannelAction struct {
	BaseAction
	Name        string            `json:"name"`
	ChannelUUID flows.ChannelUUID `json:"channel_uuid"`
}

func (a *PreferredChannelAction) Type() string { return TypeSetPreferredChannel }

func (a *PreferredChannelAction) Validate() error {
	return utils.ValidateAll(a)
}

func (a *PreferredChannelAction) Execute(run flows.FlowRun, step flows.Step) error {
	return nil
}
