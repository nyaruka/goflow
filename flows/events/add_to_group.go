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
//     "groups": [{"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Reporters"}]
//   }
// ```
//
// @event add_to_group
type AddToGroupEvent struct {
	BaseEvent
	Groups []*flows.GroupReference `json:"groups" validate:"required,min=1,dive"`
}

// NewAddToGroupEvent returns a new add to group event
func NewAddToGroupEvent(groups []*flows.GroupReference) *AddToGroupEvent {
	return &AddToGroupEvent{
		BaseEvent: NewBaseEvent(),
		Groups:    groups,
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

	for _, groupRef := range e.Groups {
		group := groupSet.FindByUUID(groupRef.UUID)

		if group != nil {
			run.Contact().Groups().Add(group)
		}
	}
	return nil
}
