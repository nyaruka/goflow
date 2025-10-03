package actions

import (
	"context"
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeEnterFlow, func() flows.Action { return &EnterFlow{} })
}

// TypeEnterFlow is the type for the enter flow action
const TypeEnterFlow string = "enter_flow"

// EnterFlow can be used to start a contact down another flow. The current flow will pause until the subflow exits or expires.
//
// A [event:run_started] event will be created to record that a new run was started.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "enter_flow",
//	  "flow": {"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Collect Language"},
//	  "terminal": false
//	}
//
// @action enter_flow
type EnterFlow struct {
	baseAction
	universalAction

	Flow     *assets.FlowReference `json:"flow" validate:"required"`
	Terminal bool                  `json:"terminal,omitempty"`
}

// NewEnterFlow creates a new start flow action
func NewEnterFlow(uuid flows.ActionUUID, flow *assets.FlowReference, terminal bool) *EnterFlow {
	return &EnterFlow{
		baseAction: newBaseAction(TypeEnterFlow, uuid),
		Flow:       flow,
		Terminal:   terminal,
	}
}

// Execute runs our action
func (a *EnterFlow) Execute(ctx context.Context, run flows.Run, step flows.Step, log flows.EventLogger) error {
	flow, err := run.Session().Assets().Flows().Get(a.Flow.UUID)

	// we ignore other missing asset types but a missing flow means we don't know how to route so we can't continue
	if err != nil {
		a.fail(run, err, log)
		return nil
	}

	if run.Session().Type() != flow.Type() {
		a.fail(run, fmt.Errorf("can't enter %s of type %s from type %s", flow.Reference(false), flow.Type(), run.Session().Type()), log)
		return nil
	}

	run.Session().PushFlow(flow, run, a.Terminal)
	return nil
}

func (a *EnterFlow) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	dependency(a.Flow)
}
