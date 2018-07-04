package actions

import (
	"github.com/nyaruka/goflow/flows"
)

// TypeCallResthook is the type for the call resthook action
const TypeCallResthook string = "call_resthook"

// CallResthookAction can be used to call a resthook.
//
// A `resthook_called` event will be created based on the results of the HTTP call.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "call_resthook",
//     "resthook": "new-registration"
//   }
//
// @action call_resthook
type CallResthookAction struct {
	BaseAction
	Resthook string `json:"resthook" validate:"required"`
}

// Type returns the type of this action
func (a *CallResthookAction) Type() string { return TypeCallResthook }

// Validate validates our action is valid and has all the assets it needs
func (a *CallResthookAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute runs this action
func (a *CallResthookAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	return nil
}
