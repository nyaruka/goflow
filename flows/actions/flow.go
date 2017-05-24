package actions

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

// TypeFlow is the type of our flow action
const TypeFlow string = "flow"

// FlowAction can be used to start a contact down another flow. The current flow will pause until the subflow exits or expires.
//
// A `flow_enter` event will be created when the flow is started, a `flow_exit` event will be created upon the subflows exit.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "flow",
//     "name": "Collect Language",
//     "flow": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"
//   }
//
// @action flow
type FlowAction struct {
	BaseAction
	Name string         `json:"name"         validate:"required"`
	Flow flows.FlowUUID `json:"flow"         validate:"required,uuid4"`
}

// Type returns the type of this action
func (a *FlowAction) Type() string { return TypeFlow }

// Validate validates our action is valid
func (a *FlowAction) Validate() error {
	return utils.ValidateAll(a)
}

// Execute runs our action
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
