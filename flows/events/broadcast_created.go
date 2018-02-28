package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
)

// TypeBroadcastCreated is a constant for outgoing message events
const TypeBroadcastCreated string = "broadcast_created"

// BroadcastCreatedEvent events are created for outgoing messages.
//
// ```
//   {
//     "type": "broadcast_created",
//     "created_on": "2006-01-02T15:04:05Z",
//     "text": "hi, what's up",
//     "attachments": [],
//     "quick_replies": ["Doing Fine", "Got 99 problems"],
//     "urns": ["tel:+12065551212"],
//     "contacts": [{"uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a", "name": "Bob"}]
//   }
// ```
//
// @event broadcast_created
type BroadcastCreatedEvent struct {
	BaseEvent
	Text         string                    `json:"text"`
	Attachments  []string                  `json:"attachments,omitempty"`
	QuickReplies []string                  `json:"quick_replies,omitempty"`
	URNs         []urns.URN                `json:"urns,omitempty" validate:"dive,urn"`
	Contacts     []*flows.ContactReference `json:"contacts,omitempty" validate:"dive"`
	Groups       []*flows.GroupReference   `json:"groups,omitempty" validate:"dive"`
}

// NewMsgCreatedEvent creates a new outgoing msg event to a single contact
func NewMsgCreatedEvent(text string, attachments []string, quickReples []string, contact *flows.ContactReference) *BroadcastCreatedEvent {
	event := BroadcastCreatedEvent{
		BaseEvent:    NewBaseEvent(),
		Text:         text,
		Attachments:  attachments,
		QuickReplies: quickReples,
		Contacts:     []*flows.ContactReference{contact},
	}
	return &event
}

// NewBroadcastCreatedEvent creates a new outgoing msg event for the given recipients
func NewBroadcastCreatedEvent(text string, attachments []string, quickReples []string, urns []urns.URN, contacts []*flows.ContactReference, groups []*flows.GroupReference) *BroadcastCreatedEvent {
	event := BroadcastCreatedEvent{
		BaseEvent:    NewBaseEvent(),
		Text:         text,
		Attachments:  attachments,
		QuickReplies: quickReples,
		URNs:         urns,
		Contacts:     contacts,
		Groups:       groups,
	}
	return &event
}

// Type returns the type of this event
func (e *BroadcastCreatedEvent) Type() string { return TypeBroadcastCreated }

// Apply applies this event to the given run
func (e *BroadcastCreatedEvent) Apply(run flows.FlowRun) error {
	return nil
}
