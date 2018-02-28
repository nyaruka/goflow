package events

import "github.com/nyaruka/goflow/flows"

// TypeContactGroupsRemoved is the type fo our remove from group action
const TypeContactGroupsRemoved string = "contact_groups_removed"

// ContactGroupsRemovedEvent events are created when a contact has been removed from one or more
// groups.
//
// ```
//   {
//     "type": "contact_groups_removed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "groups": [{"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Reporters"}]
//   }
// ```
//
// @event contact_groups_removed
type ContactGroupsRemovedEvent struct {
	Groups []*flows.GroupReference `json:"groups" validate:"required,min=1,dive"`
	BaseEvent
}

// NewContactGroupsRemovedEvent returns a new remove from group event
func NewContactGroupsRemovedEvent(groups []*flows.GroupReference) *ContactGroupsRemovedEvent {
	return &ContactGroupsRemovedEvent{
		BaseEvent: NewBaseEvent(),
		Groups:    groups,
	}
}

// Type returns the type of this event
func (e *ContactGroupsRemovedEvent) Type() string { return TypeContactGroupsRemoved }

// Apply applies this event to the given run
func (e *ContactGroupsRemovedEvent) Apply(run flows.FlowRun) error {
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
