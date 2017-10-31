package events

import "github.com/nyaruka/goflow/flows"

// TypeSaveFlowResult is the type of our save result event
const TypeSaveFlowResult string = "save_flow_result"

// SaveFlowResultEvent events are created when a result is saved. They contain not only
// the name, value and category of the result, but also the UUID of the node where
// the result was saved.
//
// ```
//   {
//     "type": "save_flow_result",
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
// @event save_flow_result
type SaveFlowResultEvent struct {
	BaseEvent
	Name              string         `json:"name" validate:"required"`
	Value             string         `json:"value"`
	Category          string         `json:"category"`
	CategoryLocalized string         `json:"category_localized,omitempty"`
	NodeUUID          flows.NodeUUID `json:"node_uuid" validate:"required,uuid4"`
	Input             string         `json:"input,omitempty"`
}

// NewSaveFlowResult returns a new save result event for the passed in values
func NewSaveFlowResult(name string, value string, categoryName string, categoryLocalized string, node flows.NodeUUID, input string) *SaveFlowResultEvent {
	return &SaveFlowResultEvent{
		BaseEvent:         NewBaseEvent(),
		Name:              name,
		Value:             value,
		Category:          categoryName,
		CategoryLocalized: categoryLocalized,
		NodeUUID:          node,
		Input:             input,
	}
}

// Type returns the type of this event
func (e *SaveFlowResultEvent) Type() string { return TypeSaveFlowResult }

// Apply applies this event to the given run
func (e *SaveFlowResultEvent) Apply(run flows.FlowRun) error {
	run.Results().Save(e.Name, e.Value, e.Category, e.CategoryLocalized, e.NodeUUID, e.Input, e.BaseEvent.CreatedOn())
	return nil
}
