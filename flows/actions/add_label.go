package actions

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// TypeAddLabel is our type for add label actions
const TypeAddLabel string = "add_label"

// AddLabelAction can be used to add a label to the last user input on a flow. An `add_label` event
// will be created with the input UUID and label UUIDs when this action is encountered. If there is
// no user input at that point then this action will be ignored.
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
	Labels []*flows.LabelReference `json:"labels" validate:"required,min=1,dive"`
}

func foo() {
}

// Type returns the type of this action
func (a *AddLabelAction) Type() string { return TypeAddLabel }

// Validate validates the fields for this label
func (a *AddLabelAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute runs the labeling action
func (a *AddLabelAction) Execute(run flows.FlowRun, step flows.Step) ([]flows.Event, error) {
	// only generate event if run has input
	input := run.Input()
	if input == nil {
		return nil, nil
	}

	log := make([]flows.Event, 0)

	labels, err := resolveLabels(run, step, a, a.Labels, log)
	if err != nil {
		return log, err
	}

	labelUUIDs := make([]flows.LabelUUID, 0, len(labels))
	for _, label := range labels {
		labelUUIDs = append(labelUUIDs, label.UUID())
	}

	if len(labelUUIDs) > 0 {
		log = append(log, events.NewAddLabelEvent(input.UUID(), labelUUIDs))
	}

	return log, nil
}
