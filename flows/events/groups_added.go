package events

import "github.com/nyaruka/goflow/flows"

// TypeGroupsAdded is the type of our add to group action
const TypeGroupsAdded string = "groups_added"

// GroupsAddedEvent events will be created with the groups a contact was added to.
//
// ```
//   {
//     "type": "groups_added",
//     "created_on": "2006-01-02T15:04:05Z",
//     "groups": [{"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Reporters"}]
//   }
// ```
//
// @event groups_added
type GroupsAddedEvent struct {
	BaseEvent
	Groups []*flows.GroupReference `json:"groups" validate:"required,min=1,dive"`
}

// NewGroupsAddedEvent returns a new groups_added event
func NewGroupsAddedEvent(groups []*flows.GroupReference) *GroupsAddedEvent {
	return &GroupsAddedEvent{
		BaseEvent: NewBaseEvent(),
		Groups:    groups,
	}
}

// Type returns the type of this event
func (e *GroupsAddedEvent) Type() string { return TypeGroupsAdded }

// Apply applies this event to the given run
func (e *GroupsAddedEvent) Apply(run flows.FlowRun) error {
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
