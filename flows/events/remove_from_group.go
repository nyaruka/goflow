package events

import "github.com/nyaruka/goflow/flows"

// TypeRemoveFromGroup is the type fo our remove from group action
const TypeRemoveFromGroup string = "remove_from_group"

// RemoveFromGroupEvent events are created when a contact is removed from one or more
// groups.
//
// ```
//   {
//    "step_uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//    "created_on": "2006-01-02T15:04:05Z",
//    "type": "remove_from_group",
//    "groups": [{
//	     "name": "Survey Audience",
//	     "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"
//	  }]
//   }
// ```
//
// @event remove_from_group
type RemoveFromGroupEvent struct {
	Groups []*flows.Group `json:"groups"  validate:"required,min=1"`
	BaseEvent
}

// NewRemoveFromGroup returns a new remove from group event
func NewRemoveFromGroup(groups []*flows.Group) *RemoveFromGroupEvent {
	return &RemoveFromGroupEvent{
		BaseEvent: NewBaseEvent(),
		Groups:    groups,
	}
}

// Type returns the type of this event
func (e *RemoveFromGroupEvent) Type() string { return TypeRemoveFromGroup }

// Apply applies this event to the given run
func (e *RemoveFromGroupEvent) Apply(run flows.FlowRun) {
	for _, group := range e.Groups {
		run.Contact().RemoveGroup(group.UUID())
	}
}
