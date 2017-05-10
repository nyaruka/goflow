package actions

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

const ADD_TO_GROUP string = "add_to_group"

type AddToGroupAction struct {
	BaseAction
	Group flows.GroupUUID `json:"group"    validate:"nonzero"`
	Name  string          `json:"name"     validate:"nonzero"`
}

func (a *AddToGroupAction) Type() string { return ADD_TO_GROUP }

func (a *AddToGroupAction) Validate() error {
	return utils.ValidateAll(a)
}

func (a *AddToGroupAction) Execute(run flows.FlowRun, step flows.Step) error {
	contact := run.Contact()
	if contact != nil {
		contact.AddGroup(a.Group, a.Name)
	}
	run.AddEvent(step, &events.AddToGroupEvent{Group: a.Group, Name: a.Name})
	return nil
}
