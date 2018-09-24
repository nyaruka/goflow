package events

import (
	"github.com/nyaruka/goflow/assets"
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

	Groups []*assets.GroupReference `json:"groups" validate:"required,min=1,dive"`
}

// NewContactGroupsAddedEvent returns a new contact_groups_added event
func NewContactGroupsAddedEvent(groups []*assets.GroupReference) *ContactGroupsAddedEvent {
	return &ContactGroupsAddedEvent{
		BaseEvent: NewBaseEvent(),
		Groups:    groups,
	}
}

// Type returns the type of this event
func (e *ContactGroupsAddedEvent) Type() string { return TypeContactGroupsAdded }
