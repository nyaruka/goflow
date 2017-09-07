package events

import "github.com/nyaruka/goflow/flows"

// TypeRemoveFromGroup is the type fo our remove from group action
const TypeRemoveFromGroup string = "remove_from_group"

// RemoveFromGroupEvent events are created when a contact is removed from one or more
// groups.
//
// ```
//   {
//     "type": "remove_from_group",
//     "created_on": "2006-01-02T15:04:05Z",
//     "group_uuids": ["b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"]
//   }
// ```
//
// @event remove_from_group
type RemoveFromGroupEvent struct {
	GroupUUIDs []flows.GroupUUID `json:"group_uuids" validate:"required,min=1,dive,uuid4"`
	BaseEvent
}

// NewRemoveFromGroupEvent returns a new remove from group event
func NewRemoveFromGroupEvent(groups []flows.GroupUUID) *RemoveFromGroupEvent {
	return &RemoveFromGroupEvent{
		BaseEvent:  NewBaseEvent(),
		GroupUUIDs: groups,
	}
}

// Type returns the type of this event
func (e *RemoveFromGroupEvent) Type() string { return TypeRemoveFromGroup }

// Apply applies this event to the given run
func (e *RemoveFromGroupEvent) Apply(run flows.FlowRun) error {
	groupSet, err := run.Session().Assets().GetGroupSet()
	if err != nil {
		return err
	}

	for _, groupUUID := range e.GroupUUIDs {
		group := groupSet.FindByUUID(groupUUID)

		if group != nil {
			run.Contact().RemoveGroup(group)
		}
	}
	return nil
}
