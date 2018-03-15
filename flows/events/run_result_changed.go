package events

import "github.com/nyaruka/goflow/flows"

// TypeRunResultChanged is the type of our run result event
const TypeRunResultChanged string = "run_result_changed"

// RunResultChangedEvent events are created when a result is saved. They contain not only
// the name, value and category of the result, but also the UUID of the node where
// the result was generated.
//
// ```
//   {
//     "type": "run_result_changed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "name": "Gender",
//     "value": "m",
//     "category": "Male",
//     "category_localized": "Homme",
//     "node_uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
//     "input": "M"
//   }
// ```
//
// @event run_result_changed
type RunResultChangedEvent struct {
	baseEvent
	callerOrEngineEvent

	Name              string         `json:"name" validate:"required"`
	Value             string         `json:"value"`
	Category          string         `json:"category"`
	CategoryLocalized string         `json:"category_localized,omitempty"`
	NodeUUID          flows.NodeUUID `json:"node_uuid" validate:"required,uuid4"`
	Operand           string         `json:"operand,omitempty"`
}

// NewRunResultChangedEvent returns a new save result event for the passed in values
func NewRunResultChangedEvent(name string, value string, categoryName string, categoryLocalized string, node flows.NodeUUID, operand string) *RunResultChangedEvent {
	return &RunResultChangedEvent{
		baseEvent:         newBaseEvent(),
		Name:              name,
		Value:             value,
		Category:          categoryName,
		CategoryLocalized: categoryLocalized,
		NodeUUID:          node,
		Operand:           operand,
	}
}

// Type returns the type of this event
func (e *RunResultChangedEvent) Type() string { return TypeRunResultChanged }

// Validate validates our event is valid and has all the assets it needs
func (e *RunResultChangedEvent) Validate(assets flows.SessionAssets) error {
	return nil
}

// Apply applies this event to the given run
func (e *RunResultChangedEvent) Apply(run flows.FlowRun) error {
	run.Results().Save(e.Name, e.Value, e.Category, e.CategoryLocalized, e.NodeUUID, e.Operand, e.baseEvent.CreatedOn())
	return nil
}
