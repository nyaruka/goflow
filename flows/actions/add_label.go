package actions

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

const ADD_LABEL string = "add_label"

type AddLabelAction struct {
	BaseAction
	Labels []flows.Label `json:"labels"`
}

func (a *AddLabelAction) Type() string { return ADD_LABEL }

func (a *AddLabelAction) Validate() error {
	return utils.ValidateAll(a)
}

func (a *AddLabelAction) Execute(run flows.FlowRun, step flows.Step) error {
	return nil
}
