package events

import "github.com/nyaruka/goflow/flows"

// TypeLabelAdded is the type of our add label action
const TypeLabelAdded string = "label_added"

// LabelAddedEvent events will be created with the labels that were applied to the given input.
//
// ```
//   {
//     "type": "label_added",
//     "created_on": "2006-01-02T15:04:05Z",
//     "input_uuid": "4aef4050-1895-4c80-999a-70368317a4f5",
//     "labels": [{"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Spam"}]
//   }
// ```
//
// @event label_added
type LabelAddedEvent struct {
	BaseEvent
	InputUUID flows.InputUUID         `json:"input_uuid" validate:"required,uuid4"`
	Labels    []*flows.LabelReference `json:"labels" validate:"required,min=1,dive"`
}

// NewLabelAddedEvent returns a new add to group event
func NewLabelAddedEvent(inputUUID flows.InputUUID, labels []*flows.LabelReference) *LabelAddedEvent {
	return &LabelAddedEvent{
		BaseEvent: NewBaseEvent(),
		InputUUID: inputUUID,
		Labels:    labels,
	}
}

// Type returns the type of this event
func (e *LabelAddedEvent) Type() string { return TypeLabelAdded }

// Apply applies this event to the given run
func (e *LabelAddedEvent) Apply(run flows.FlowRun) error {
	return nil
}
