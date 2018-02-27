package events

import "github.com/nyaruka/goflow/flows"

// TypeGroupAdded is the type of our add to group action
const TypeGroupAdded string = "group_added"

// GroupAddedEvent events will be created with the groups a contact was added to.
//
// ```
//   {
//     "type": "group_added",
//     "created_on": "2006-01-02T15:04:05Z",
//     "groups": [{"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Reporters"}]
//   }
// ```
//
// @event group_added
type GroupAddedEvent struct {
	BaseEvent
	Groups []*flows.GroupReference `json:"groups" validate:"required,min=1,dive"`
}

// NewGroupAddedEvent returns a new add to group event
func NewGroupAddedEvent(groups []*flows.GroupReference) *GroupAddedEvent {
	return &GroupAddedEvent{
		BaseEvent: NewBaseEvent(),
		Groups:    groups,
	}
}

// Type returns the type of this event
func (e *GroupAddedEvent) Type() string { return TypeGroupAdded }

// Apply applies this event to the given run
func (e *GroupAddedEvent) Apply(run flows.FlowRun) error {
	groupSet, err := run.Session().Assets().GetGroupSet()
	if err != nil {
		return err
	}

	for _, groupRef := range e.Groups {
		group := groupSet.FindByUUID(groupRef.UUID)

		if group != nil {
			run.Contact().Groups().Add(group)
		}
	}
	return nil
}
