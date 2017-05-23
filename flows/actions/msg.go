package actions

import (
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

const MSG string = "msg"

type MsgAction struct {
	BaseAction
	Text string `json:"text"         validate:"required"`
}

func (a *MsgAction) Type() string { return MSG }

func (a *MsgAction) Validate() error {
	return utils.ValidateAll(a)
}

func (a *MsgAction) Execute(run flows.FlowRun, step flows.Step) error {
	text, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), run.GetText(flows.UUID(a.Uuid), "text", a.Text))
	if err != nil {
		run.AddError(step, err)
	}
	run.AddEvent(step, events.NewOutgoingMsgEvent(run.ChannelUUID(), run.Contact().UUID(), text))
	return nil
}
