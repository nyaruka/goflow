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
//     "flow": {"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Collect Language"}
//   }
// ```
//
// @action start_flow
type StartFlowAction struct {
	BaseAction
	Flow *flows.FlowReference `json:"flow" validate:"required"`
}

// Type returns the type of this action
func (a *StartFlowAction) Type() string { return TypeStartFlow }

// Validate validates our action is valid
func (a *StartFlowAction) Validate(assets flows.SessionAssets) error {
	_, err := assets.GetFlow(a.Flow.UUID)
	return err
}

// Execute runs our action
func (a *StartFlowAction) Execute(run flows.FlowRun, step flows.Step) ([]flows.Event, error) {

	if run.Session().FlowOnStack(a.Flow.UUID) {
		return []flows.Event{events.NewFatalErrorEvent(fmt.Errorf("flow loop detected, stopping execution before starting flow: %s", a.Flow.UUID))}, nil
	}

	return []flows.Event{events.NewFlowTriggeredEvent(a.Flow.UUID, run.UUID())}, nil
}
