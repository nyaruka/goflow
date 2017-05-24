package events

import "github.com/nyaruka/goflow/flows"

// TypeSaveResult is the type of our save result event
const TypeSaveResult string = "save_result"

// SaveResultEvent events are created when a result is saved. They contain not only
// the name, value and category of the result, but also the UUID of the node where
// the result was saved.
//
// ```
//   {
//    "step": "8eebd020-1af5-431c-b943-aa670fc74da9",
//    "created_on": "2006-01-02T15:04:05Z",
//    "type": "save_result",
//    "node": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
//    "name": "Gender",
//    "value": "m",
//    "category": "Make"
//   }
// ```
//
// @event save_result
type SaveResultEvent struct {
	BaseEvent
	Node     flows.NodeUUID `json:"node"        validate:"required"`
	Name     string         `json:"name"        validate:"required"`
	Value    string         `json:"value"       validate:"required"`
	Category string         `json:"category"    validate:"required"`
}

// NewSaveResult returns a new save result event for the passed in values
func NewSaveResult(node flows.NodeUUID, name string, value string, category string) *SaveResultEvent {
	return &SaveResultEvent{Node: node, Name: name, Value: value, Category: category}
}

// Type returns the type of this event
func (e *SaveResultEvent) Type() string { return TypeSaveResult }
