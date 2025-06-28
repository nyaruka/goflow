package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeInputLabelsAdded, func() flows.Event { return &InputLabelsAdded{} })
}

// TypeInputLabelsAdded is the type of our add label action
const TypeInputLabelsAdded string = "input_labels_added"

// InputLabelsAdded events are created when an action wants to add labels to the current input.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "input_labels_added",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "input_uuid": "4aef4050-1895-4c80-999a-70368317a4f5",
//	  "labels": [{"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Spam"}]
//	}
//
// @event input_labels_added
type InputLabelsAdded struct {
	BaseEvent

	InputUUID flows.InputUUID          `json:"input_uuid" validate:"required,uuid"`
	Labels    []*assets.LabelReference `json:"labels" validate:"required,min=1,dive"`
}

// NewInputLabelsAdded returns a new labels added event
func NewInputLabelsAdded(inputUUID flows.InputUUID, labels []*flows.Label) *InputLabelsAdded {
	return &InputLabelsAdded{
		BaseEvent: NewBaseEvent(TypeInputLabelsAdded),
		InputUUID: inputUUID,
		Labels:    labelsToReferences(labels),
	}
}

// converts a slice of labels to a slice of references
func labelsToReferences(labels []*flows.Label) []*assets.LabelReference {
	refs := make([]*assets.LabelReference, len(labels))
	for i := range labels {
		refs[i] = labels[i].Reference()
	}
	return refs
}
