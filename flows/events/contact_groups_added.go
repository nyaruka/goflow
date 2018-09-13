package events

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeContactGroupsAdded, func() flows.Event { return &ContactGroupsAddedEvent{} })
}

// TypeContactGroupsAdded is the type of our add to group action
const TypeContactGroupsAdded string = "contact_groups_added"

// ContactGroupsAddedEvent events will be created with the groups a contact was added to.
//
//   {
//     "type": "contact_groups_added",
//     "created_on": "2006-01-02T15:04:05Z",
//     "groups": [{"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Reporters"}]
//   }
//
// @event contact_groups_added
type ContactGroupsAddedEvent struct {
	BaseEvent
	callerOrEngineEvent

	Groups []*flows.GroupReference `json:"groups" validate:"required,min=1,dive"`
}

// NewContactGroupsAddedEvent returns a new contact_groups_added event
func NewContactGroupsAddedEvent(groups []*flows.GroupReference) *ContactGroupsAddedEvent {
	return &ContactGroupsAddedEvent{
		BaseEvent: NewBaseEvent(),
		Groups:    groups,
	}
}

// Type returns the type of this event
func (e *ContactGroupsAddedEvent) Type() string { return TypeContactGroupsAdded }

// Validate validates our event is valid and has all the assets it needs
func (e *ContactGroupsAddedEvent) Validate(assets flows.SessionAssets) error {
	for _, group := range e.Groups {
		if _, err := assets.Groups().Get(group.UUID); err != nil {
			return err
		}
	}
	return nil
}

// Apply applies this event to the given run
func (e *ContactGroupsAddedEvent) Apply(run flows.FlowRun) error {
	if run.Contact() == nil {
		return fmt.Errorf("can't apply event in session without a contact")
	}

	assets := run.Session().Assets()

	for _, groupRef := range e.Groups {
		group, err := assets.Groups().Get(groupRef.UUID)
		if err == nil {
			run.Contact().Groups().Add(group)
		}
	}
	return nil
}
