package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeContactGroupsChanged, func() flows.Event { return &ContactGroupsChangedEvent{} })
}

// TypeContactGroupsChanged is the type of our groups changed event
const TypeContactGroupsChanged string = "contact_groups_changed"

// ContactGroupsChangedEvent events are created when a contact is added or removed to/from one or more groups.
//
//   {
//     "type": "contact_groups_changed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "groups_added": [{"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Reporters"}],
//     "groups_removed": [{"uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a", "name": "Customers"}]
//   }
//
// @event contact_groups_changed
type ContactGroupsChangedEvent struct {
	BaseEvent

	GroupsAdded   []*assets.GroupReference `json:"groups_added,omitempty" validate:"omitempty,dive"`
	GroupsRemoved []*assets.GroupReference `json:"groups_removed,omitempty" validate:"omitempty,dive"`
}

// NewContactGroupsChangedEvent returns a new contact_groups_changed event
func NewContactGroupsChangedEvent(added []*assets.GroupReference, removed []*assets.GroupReference) *ContactGroupsChangedEvent {
	return &ContactGroupsChangedEvent{
		BaseEvent:     NewBaseEvent(),
		GroupsAdded:   added,
		GroupsRemoved: removed,
	}
}

// Type returns the type of this event
func (e *ContactGroupsChangedEvent) Type() string { return TypeContactGroupsChanged }
