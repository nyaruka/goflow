package actions

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

const EMAIL string = "email"

type EmailAction struct {
	BaseAction
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	Emails  []string `json:"emails"`
}

func (a *EmailAction) Type() string { return EMAIL }

func (a *EmailAction) Validate() error {
	return utils.ValidateAll(a)
}

func (a *EmailAction) Execute(run flows.FlowRun, step flows.Step) error {
	return nil
}
