package actions

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

const REMOVE_FROM_GROUP string = "remove_from_group"

type RemoveFromGroupAction struct {
	BaseAction
	Groups []*flows.Group `json:"groups"`
}

func (a *RemoveFromGroupAction) Type() string { return REMOVE_FROM_GROUP }

func (a *RemoveFromGroupAction) Validate() error {
	return utils.ValidateAll(a)
}

func (a *RemoveFromGroupAction) Execute(run flows.FlowRun, step flows.Step) error {
	contact := run.Contact()
	if contact != nil {
		// no groups in our action means remove all
		if len(a.Groups) == 0 {
			for _, group := range contact.Groups() {
				contact.RemoveGroup(group.UUID())
				run.AddEvent(step, &events.RemoveFromGroupEvent{Group: group.UUID(), Name: group.Name()})
			}
		} else {
			for _, group := range a.Groups {
				contact.RemoveGroup(group.UUID())
				run.AddEvent(step, &events.RemoveFromGroupEvent{Group: group.UUID(), Name: group.Name()})
			}
		}
	}

	return nil
}
