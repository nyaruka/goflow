package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

// TypeStartFlow is the type of our flow action
const TypeStartFlow string = "start_flow"

// StartFlowAction can be used to start a contact down another flow. The current flow will pause until the subflow exits or expires.
//
// A `flow_entered` event will be created when the flow is started, a `flow_exited` event will be created upon the subflows exit.
//
// ```
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "start_flow",
//     "flow_name": "Collect Language",
//     "flow_uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"
//   }
// ```
//
// @action start_flow
type StartFlowAction struct {
	BaseAction
	FlowName string         `json:"flow_name"    validate:"required"`
	FlowUUID flows.FlowUUID `json:"flow_uuid"    validate:"required,uuid4"`
}

// Type returns the type of this action
func (a *StartFlowAction) Type() string { return TypeStartFlow }

// Validate validates our action is valid
func (a *StartFlowAction) Validate() error {
	return utils.ValidateAll(a)
}

// Execute runs our action
func (a *StartFlowAction) Execute(run flows.FlowRun, step flows.Step) error {
	// lookup our flow
	flow, err := run.Environment().GetFlow(a.FlowUUID)
	if err != nil {
		run.AddError(step, err)
		return err
	}

	// how many times have we started this flow in this session without exiting?
	startCount := 0
	for _, evt := range run.Session().Events() {
		enter, isEnter := evt.(*events.FlowEnteredEvent)
		if isEnter && enter.FlowUUID == a.FlowUUID {
			startCount++
			continue
		}

		exit, isExit := evt.(*events.FlowExitedEvent)
		if isExit && exit.FlowUUID == a.FlowUUID {
			startCount--
			continue
		}
	}

	// we don't allow recursion, you can't call back into yourself
	if startCount > 0 {
		return fmt.Errorf("flow loop detected, stopping execution before starting flow: %s", a.FlowUUID)
	}

	// log our event
	run.AddEvent(step, events.NewFlowEnterEvent(a.FlowUUID, run.Contact().UUID()))

	// start it for our current contact
	_, err = engine.StartFlow(run.Environment(), flow, run.Contact(), run, nil, nil)

	// if we received an error, shortcut out, this session is horked
	if err != nil {
		return err
	}

	// same thing if our child ended as an error, session is horked
	if run.Child().Status() == flows.StatusErrored {
		run.AddEvent(step, events.NewFlowExitedEvent(run.Child()))
		return fmt.Errorf("child run for flow '%s' ended in error, ending execution", a.FlowUUID)
	}

	// did we complete?
	if run.Child().Status() != flows.StatusActive {
		// add our exit event
		run.AddEvent(step, events.NewFlowExitedEvent(run.Child()))
	}

	return nil
}
