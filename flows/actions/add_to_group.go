package actions

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

const ADD_TO_GROUP string = "add_to_group"

type AddToGroupAction struct {
	BaseAction
	Groups []*flows.Group `json:"groups"    validate:"nonzero"`
}

func (a *AddToGroupAction) Type() string { return ADD_TO_GROUP }

func (a *AddToGroupAction) Validate() error {
	return utils.ValidateAll(a)
}

func (a *AddToGroupAction) Execute(run flows.FlowRun, step flows.Step) error {
	contact := run.Contact()
	if contact != nil {
		for _, group := range a.Groups {
			contact.AddGroup(group.UUID(), group.Name())
			run.AddEvent(step, &events.AddToGroupEvent{Group: group.UUID(), Name: group.Name()})
		}
	}

	return nil
}
