package events

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeRunResultChanged, func() flows.Event { return &RunResultChangedEvent{} })
}

// TypeRunResultChanged is the type of our run result event
const TypeRunResultChanged string = "run_result_changed"

// RunResultChangedEvent events are created when a run result is saved. They contain not only
// the name, value and category of the result, but also the UUID of the node where
// the result was generated.
//
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
//
// @event run_result_changed
type RunResultChangedEvent struct {
	BaseEvent

	Name              string          `json:"name" validate:"required"`
	Value             string          `json:"value"`
	Category          string          `json:"category"`
	CategoryLocalized string          `json:"category_localized,omitempty"`
	Input             string          `json:"input,omitempty"`
	Extra             json.RawMessage `json:"extra,omitempty"`
}

// NewRunResultChanged returns a new save result event for the passed in values
func NewRunResultChanged(result *flows.Result) *RunResultChangedEvent {
	return &RunResultChangedEvent{
		BaseEvent:         NewBaseEvent(TypeRunResultChanged),
		Name:              result.Name,
		Value:             result.Value,
		Category:          result.Category,
		CategoryLocalized: result.CategoryLocalized,
		Input:             result.Input,
		Extra:             result.Extra,
	}
}
