package events

import "github.com/nyaruka/goflow/flows"

// TypeAddToGroup is the type of our add to group action
const TypeAddToGroup string = "add_to_group"

// AddToGroupEvent events will be created with the groups a contact should be added to.
//
// ```
//   {
//     "type": "add_to_group",
//     "created_on": "2006-01-02T15:04:05Z",
//     "group_uuids": ["b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"]
//   }
// ```
//
// @event add_to_group
type AddToGroupEvent struct {
	BaseEvent
	GroupUUIDs []flows.GroupUUID `json:"group_uuids" validate:"required,min=1,dive,uuid4"`
}

// NewAddToGroupEvent returns a new add to group event
func NewAddToGroupEvent(groups []flows.GroupUUID) *AddToGroupEvent {
	return &AddToGroupEvent{
		BaseEvent:  NewBaseEvent(),
		GroupUUIDs: groups,
	}
}

// Type returns the type of this event
func (e *AddToGroupEvent) Type() string { return TypeAddToGroup }

// Apply applies this event to the given run
func (e *AddToGroupEvent) Apply(run flows.FlowRun) error {
	groupSet, err := run.Session().Assets().GetGroupSet()
	if err != nil {
		return err
	}

	for _, groupUUID := range e.GroupUUIDs {
		group := groupSet.FindByUUID(groupUUID)

		// TODO groups which evaluate to a name match

		if group != nil {
			run.Contact().AddGroup(group)
		}
	}
	return nil
}
