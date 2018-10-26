package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/assets"
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
//     "flow": {"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Collect Language"},
//     "terminal": false
//   }
//
// @action start_flow
type StartFlowAction struct {
	BaseAction
	universalAction

	Flow     *assets.FlowReference `json:"flow" validate:"required"`
	Terminal bool                  `json:"terminal"`
}

// NewStartFlowAction creates a new start flow action
func NewStartFlowAction(uuid flows.ActionUUID, flow *assets.FlowReference, terminal bool) *StartFlowAction {
	return &StartFlowAction{
		BaseAction: NewBaseAction(TypeStartFlow, uuid),
		Flow:       flow,
		Terminal:   terminal,
	}
}

// Validate validates our action is valid and has all the assets it needs
func (a *StartFlowAction) Validate(assets flows.SessionAssets, context *flows.ValidationContext) error {

	// check the flow exists and that it's valid
	return a.validateFlow(assets, a.Flow, context)
}

// Execute runs our action
func (a *StartFlowAction) Execute(run flows.FlowRun, step flows.Step) error {
	flow, err := run.Session().Assets().Flows().Get(a.Flow.UUID)
	if err != nil {
		return err
	}

	if !run.Session().CanEnterFlow(flow) {
		a.fatalError(run, step, fmt.Errorf("flow loop detected, stopping execution before starting flow: %s", a.Flow.UUID))
		return nil
	}

	run.Session().PushFlow(flow, run, a.Terminal)
	a.log(run, step, events.NewFlowTriggeredEvent(a.Flow, run.UUID(), a.Terminal))
	return nil
}
