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
//     "groups": [{"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Reporters"}]
//   }
// ```
//
// @event remove_from_group
type RemoveFromGroupEvent struct {
	Groups []*flows.GroupReference `json:"groups" validate:"required,min=1,dive"`
	BaseEvent
}

// NewRemoveFromGroupEvent returns a new remove from group event
func NewRemoveFromGroupEvent(groups []*flows.GroupReference) *RemoveFromGroupEvent {
	return &RemoveFromGroupEvent{
		BaseEvent: NewBaseEvent(),
		Groups:    groups,
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

	for _, groupRef := range e.Groups {
		group := groupSet.FindByUUID(groupRef.UUID)

		if group != nil {
			run.Contact().Groups().Remove(group)
		}
	}
	return nil
}
