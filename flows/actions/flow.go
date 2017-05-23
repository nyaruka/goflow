package actions

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

const FLOW string = "flow"

type FlowAction struct {
	BaseAction
	Name string         `json: "name"`
	Flow flows.FlowUUID `json:"flow"         validate:"required"`
}

func (a *FlowAction) Type() string { return FLOW }

func (a *FlowAction) Validate() error {
	return utils.ValidateAll(a)
}

func (a *FlowAction) Execute(run flows.FlowRun, step flows.Step) error {
	// lookup our flow
	flow, err := run.Environment().GetFlow(a.Flow)
	if err != nil {
		run.AddError(step, err)
		return err
	}

	// log our event
	run.AddEvent(step, events.NewFlowEnterEvent(a.Flow, run.Contact().UUID()))

	// start it for our current contact
	_, err = engine.StartFlow(run.Environment(), flow, run.Contact(), run)

	// log any error we receive
	if err != nil {
		run.AddError(step, err)
		return err
	}
	return nil
}
