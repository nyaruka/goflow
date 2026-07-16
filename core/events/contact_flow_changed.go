package events

import (
	"github.com/nyaruka/goflow/assets"
)

func init() {
	registerType(TypeContactFlowChanged, func() Event { return &ContactFlowChanged{} })
}

// TypeContactFlowChanged is the type of our contact flow changed event
const TypeContactFlowChanged string = "contact_flow_changed"

// ContactFlowChanged events are created when the current flow of the contact has been changed - i.e. the flow
// in which the contact is currently waiting. The flow will be null if the contact is no longer in a flow.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "contact_flow_changed",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"}
//	}
//
// @event contact_flow_changed
type ContactFlowChanged struct {
	BaseEvent

	Flow *assets.FlowReference `json:"flow"`
}

// NewContactFlowChanged returns a new contact flow changed event
func NewContactFlowChanged(flow *assets.FlowReference) *ContactFlowChanged {
	return &ContactFlowChanged{
		BaseEvent: NewBaseEvent(TypeContactFlowChanged),
		Flow:      flow,
	}
}

var _ Event = (*ContactFlowChanged)(nil)
