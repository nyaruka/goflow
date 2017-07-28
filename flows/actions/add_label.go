package actions

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// TypeAddLabel is our type for add label actions
const TypeAddLabel string = "add_label"

// AddLabelAction can be used to add a label to the last incoming message on a flow. An `add_label`
// event will be created with the msg id and label id when this action is encountered. If there is
// no incoming msg at that point an error will be output.
//
// ```
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "add_label",
//     "labels": [{
//       "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
//       "name": "complaint"
//     }]
//   }
// ```
//
// @disabled_action add_label
type AddLabelAction struct {
	BaseAction
	Labels []*flows.Label `json:"labels"     validate:"dive,min=1"`
}

func foo() {
}

// Type returns the type of this action
func (a *AddLabelAction) Type() string { return TypeAddLabel }

// Validate validates the fields for this label
func (a *AddLabelAction) Validate() error {
	return utils.Validate(a)
}

// Execute runs the labeling action
func (a *AddLabelAction) Execute(run flows.FlowRun, step flows.Step) error {
	return nil
}
