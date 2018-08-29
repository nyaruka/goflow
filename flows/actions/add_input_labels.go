package actions

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	RegisterType(TypeAddInputLabels, func() flows.Action { return &AddInputLabelsAction{} })
}

// TypeAddInputLabels is the type for the add label action
const TypeAddInputLabels string = "add_input_labels"

// AddInputLabelsAction can be used to add labels to the last user input on a flow. An [event:input_labels_added] event
// will be created with the labels added when this action is encountered. If there is
// no user input at that point then this action will be ignored.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "add_input_labels",
//     "labels": [{
//       "uuid": "3f65d88a-95dc-4140-9451-943e94e06fea",
//       "name": "Spam"
//     }]
//   }
//
// @action add_input_labels
type AddInputLabelsAction struct {
	BaseAction
	universalAction

	Labels []*flows.LabelReference `json:"labels" validate:"required,dive"`
}

// Type returns the type of this action
func (a *AddInputLabelsAction) Type() string { return TypeAddInputLabels }

// Validate validates our action is valid and has all the assets it needs
func (a *AddInputLabelsAction) Validate(assets flows.SessionAssets) error {
	// check we have all labels
	return a.validateLabels(assets, a.Labels)
}

// Execute runs the labeling action
func (a *AddInputLabelsAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	// only generate event if run has input
	input := run.Input()
	if input == nil {
		return nil
	}

	labels, err := a.resolveLabels(run, step, a.Labels, log)
	if err != nil {
		return err
	}

	labelRefs := make([]*flows.LabelReference, 0, len(labels))
	for _, label := range labels {
		labelRefs = append(labelRefs, label.Reference())
	}

	if len(labelRefs) > 0 {
		log.Add(events.NewInputLabelsAddedEvent(input.UUID(), labelRefs))
	}

	return nil
}
