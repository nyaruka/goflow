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
//    "step_uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//    "created_on": "2006-01-02T15:04:05Z",
//    "type": "save_flow_result",
//    "node_uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
//    "result_name": "Gender",
//    "value": "m",
//    "category": "Make"
//   }
// ```
//
// @event save_flow_result
type SaveFlowResultEvent struct {
	BaseEvent
	NodeUUID          flows.NodeUUID `json:"node_uuid"        validate:"required"`
	ResultName        string         `json:"result_name"      validate:"required"`
	Value             string         `json:"value"`
	Category          string         `json:"category"`
	CategoryLocalized string         `json:"category_localized,omitempty"`
}

// NewSaveFlowResult returns a new save result event for the passed in values
func NewSaveFlowResult(node flows.NodeUUID, name string, value string, categoryName string, categoryLocalized string) *SaveFlowResultEvent {
	return &SaveFlowResultEvent{
		BaseEvent:  NewBaseEvent(),
		NodeUUID:   node,
		ResultName: name,
		Value:      value, Category: categoryName,
		CategoryLocalized: categoryLocalized,
	}
}

// Type returns the type of this event
func (e *SaveFlowResultEvent) Type() string { return TypeSaveFlowResult }

// Apply applies this event to the given run
func (e *SaveFlowResultEvent) Apply(run flows.FlowRun) {
	run.Results().Save(e.NodeUUID, e.ResultName, e.Value, e.Category, e.CategoryLocalized, e.BaseEvent.CreatedOn())
}
