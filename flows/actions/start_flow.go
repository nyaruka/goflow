package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
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
//     "flow_uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
//     "flow_name": "Collect Language"
//   }
// ```
//
// @action start_flow
type StartFlowAction struct {
	BaseAction
	FlowUUID flows.FlowUUID `json:"flow_uuid"    validate:"required,uuid4"`
	FlowName string         `json:"flow_name"    validate:"required"`
}

// Type returns the type of this action
func (a *StartFlowAction) Type() string { return TypeStartFlow }

// Validate validates our action is valid
func (a *StartFlowAction) Validate(assets flows.AssetStore) error {
	_, err := assets.GetFlow(a.FlowUUID)
	return err
}

// Execute runs our action
func (a *StartFlowAction) Execute(run flows.FlowRun, step flows.Step) error {

	if run.Session().FlowOnStack(a.FlowUUID) {
		run.AddFatalError(step, a, fmt.Errorf("flow loop detected, stopping execution before starting flow: %s", a.FlowUUID))
		return nil
	}

	run.ApplyEvent(step, a, events.NewFlowTriggeredEvent(a.FlowUUID, run.UUID()))
	return nil
}
