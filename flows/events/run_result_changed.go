package events

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeRunResultChanged, func() flows.Event { return &RunResultChanged{} })
}

// TypeRunResultChanged is the type of our run result event
const TypeRunResultChanged string = "run_result_changed"

type PreviousResult struct {
	Value    string `json:"value"`
	Category string `json:"category"`
}

// RunResultChanged events are created when a run result is changed.
//
//	{
//	  "type": "run_result_changed",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "name": "Gender",
//	  "value": "m",
//	  "category": "Male"
//	}
//
// @event run_result_changed
type RunResultChanged struct {
	BaseEvent

	Name     string          `json:"name" validate:"required"`
	Value    string          `json:"value"`
	Category string          `json:"category"`
	Extra    json.RawMessage `json:"extra,omitempty"`
	Previous *PreviousResult `json:"previous,omitempty"`
}

// NewRunResultChanged returns a new save result event for the passed in values
func NewRunResultChanged(result, prev *flows.Result) *RunResultChanged {
	var p *PreviousResult
	if prev != nil {
		p = &PreviousResult{
			Value:    prev.Value,
			Category: prev.Category,
		}
	}

	return &RunResultChanged{
		BaseEvent: NewBaseEvent(TypeRunResultChanged),
		Name:      result.Name,
		Value:     result.Value,
		Category:  result.Category,
		Extra:     result.Extra,
		Previous:  p,
	}
}
