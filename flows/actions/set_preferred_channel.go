package actions

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

const SET_PREFERRED_CHANNEL string = "set_preferred_channel"

type PreferredChannelAction struct {
	BaseAction
	Name    string            `json:"name"`
	Channel flows.ChannelUUID `json:"channel"`
}

func (a *PreferredChannelAction) Type() string { return SET_PREFERRED_CHANNEL }

func (a *PreferredChannelAction) Validate() error {
	return utils.ValidateAll(a)
}

func (a *PreferredChannelAction) Execute(run flows.FlowRun, step flows.Step) error {
	return nil
}
