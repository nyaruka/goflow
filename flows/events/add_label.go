package events

import "github.com/nyaruka/goflow/flows"

// TypeAddLabel is the type of our add label action
const TypeAddLabel string = "add_label"

// AddLabelEvent events will be created with the labels that should be applied to the given input.
//
// ```
//   {
//     "type": "add_label",
//     "created_on": "2006-01-02T15:04:05Z",
//     "input_uuid": "4aef4050-1895-4c80-999a-70368317a4f5",
//     "label_uuids": ["b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"]
//   }
// ```
//
// @event add_label
type AddLabelEvent struct {
	BaseEvent
	InputUUID  flows.InputUUID   `json:"input_uuid" validate:"required,uuid4"`
	LabelUUIDs []flows.LabelUUID `json:"label_uuids" validate:"required,min=1,dive,uuid4"`
}

// NewAddLabelEvent returns a new add to group event
func NewAddLabelEvent(inputUUID flows.InputUUID, labelUUIDs []flows.LabelUUID) *AddLabelEvent {
	return &AddLabelEvent{
		BaseEvent:  NewBaseEvent(),
		InputUUID:  inputUUID,
		LabelUUIDs: labelUUIDs,
	}
}

// Type returns the type of this event
func (e *AddLabelEvent) Type() string { return TypeAddLabel }

// Apply applies this event to the given run
func (e *AddLabelEvent) Apply(run flows.FlowRun) error {
	return nil
}
