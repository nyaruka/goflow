package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	RegisterType(TypeStartFlow, func() flows.Action { return &StartFlowAction{} })
}

// TypeStartFlow is the type for the start flow action
const TypeStartFlow string = "start_flow"

// StartFlowAction can be used to start a contact down another flow. The current flow will pause until the subflow exits or expires.
//
// A [event:flow_triggered] event will be created to record that the flow was started.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "start_flow",
//     "flow": {"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Collect Language"}
//   }
//
// @action start_flow
type StartFlowAction struct {
	BaseAction
	universalAction

	Flow *flows.FlowReference `json:"flow" validate:"required"`
}

// Type returns the type of this action
func (a *StartFlowAction) Type() string { return TypeStartFlow }

// Validate validates our action is valid and has all the assets it needs
func (a *StartFlowAction) Validate(assets flows.SessionAssets) error {
	// check we have the flow
	_, err := assets.GetFlow(a.Flow.UUID)
	return err
}

// Execute runs our action
func (a *StartFlowAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	if run.Session().FlowOnStack(a.Flow.UUID) {
		log.Add(events.NewFatalErrorEvent(fmt.Errorf("flow loop detected, stopping execution before starting flow: %s", a.Flow.UUID)))
		return nil
	}

	log.Add(events.NewFlowTriggeredEvent(a.Flow, run.UUID()))
	return nil
}
