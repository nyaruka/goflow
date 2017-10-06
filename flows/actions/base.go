package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/flows/events"

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
func resolveGroups(run flows.FlowRun, step flows.Step, action flows.Action, references []*flows.GroupReference, log []flows.Event) ([]*flows.Group, error) {
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
			evaluatedGroupName, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), ref.NameMatch)
			if err != nil {
				log = append(log, events.NewErrorEvent(err))
			} else {
				// look up the set of all groups to see if such a group exists
				group = groupSet.FindByName(evaluatedGroupName)
				if group == nil {
					log = append(log, events.NewErrorEvent(fmt.Errorf("no such group with name '%s'", evaluatedGroupName)))
				}
			}
		}

		if group != nil {
			groups = append(groups, group)
		}
	}

	return groups, nil
}

// helper function for actions that have a set of label references that must be resolved to actual labels
func resolveLabels(run flows.FlowRun, step flows.Step, action flows.Action, references []*flows.LabelReference, log []flows.Event) ([]*flows.Label, error) {
	labelSet, err := run.Session().Assets().GetLabelSet()
	if err != nil {
		return nil, err
	}

	labels := make([]*flows.Label, 0, len(references))

	for _, ref := range references {
		var label *flows.Label

		if ref.UUID != "" {
			// label is a fixed label with a UUID
			label = labelSet.FindByUUID(ref.UUID)
			if label == nil {
				return nil, fmt.Errorf("no such label with UUID '%s'", ref.UUID)
			}
		} else {
			// label is an expression that evaluates to an existing label's name
			evaluatedLabelName, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), ref.NameMatch)
			if err != nil {
				log = append(log, events.NewErrorEvent(err))
			} else {
				// look up the set of all labels to see if such a label exists
				label = labelSet.FindByName(evaluatedLabelName)
				if label == nil {
					log = append(log, events.NewErrorEvent(fmt.Errorf("no such label with name '%s'", evaluatedLabelName)))
				}
			}
		}

		if label != nil {
			labels = append(labels, label)
		}
	}

	return labels, nil
}
