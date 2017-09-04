package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
)

// BaseAction is our base action type
type BaseAction struct {
	UUID_ flows.ActionUUID `json:"uuid"    validate:"required,uuid4"`
}

func NewBaseAction(uuid flows.ActionUUID) BaseAction {
	return BaseAction{UUID_: uuid}
}

func (a *BaseAction) UUID() flows.ActionUUID { return a.UUID_ }

// helper function for actions that have a set of group references that must be resolved to actual groups
func resolveGroups(run flows.FlowRun, step flows.Step, action flows.Action, references []*flows.GroupReference) ([]*flows.Group, error) {
	groupSet, err := run.Session().Assets().GetGroupSet()
	if err != nil {
		return nil, err
	}

	groups := make([]*flows.Group, 0, len(references))

	for _, ref := range references {
		var group *flows.Group

		if ref.UUID != "" {
			// group is a fixed group with a UUID
			group = groupSet.FindByUUID(ref.UUID)
			if group == nil {
				return nil, fmt.Errorf("no such group with UUID '%s'", ref.UUID)
			}
		} else {
			// group is an expression that evaluates to an existing group's name
			// evaluate the expression to get the group name
			evaluatedGroupName, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), ref.Name)
			if err != nil {
				run.AddError(step, action, err)
			} else {
				// look up the set of all groups to see if such a group exists
				group = groupSet.FindByName(evaluatedGroupName)
				if group == nil {
					run.AddError(step, action, fmt.Errorf("no such group with name '%s'", evaluatedGroupName))
				}
			}
		}

		if group != nil {
			groups = append(groups, group)
		}
	}

	return groups, nil
}
