package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
)

// TypeMsgSent is a constant for outgoing message events
const TypeMsgSent string = "msg_sent"

// MsgSentEvent events are created for outgoing messages.
//
// ```
//   {
//     "type": "msg_sent",
//     "created_on": "2006-01-02T15:04:05Z",
//     "text": "hi, what's up",
//     "attachments": [],
//     "quick_replies": ["Doing Fine", "Got 99 problems"],
//     "urns": ["tel:+12065551212"],
//     "contacts": [{"uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a", "name": "Bob"}]
//   }
// ```
//
// @event msg_sent
type MsgSentEvent struct {
	BaseEvent
	Text         string                    `json:"text"`
	Attachments  []string                  `json:"attachments,omitempty"`
	QuickReplies []string                  `json:"quick_replies,omitempty"`
	URNs         []urns.URN                `json:"urns,omitempty" validate:"dive,urn"`
	Contacts     []*flows.ContactReference `json:"contacts,omitempty" validate:"dive"`
	Groups       []*flows.GroupReference   `json:"groups,omitempty" validate:"dive"`
}

// NewMsgSentToContactEvent creates a new outgoing msg event to a single contact
func NewMsgSentToContactEvent(text string, attachments []string, quickReples []string, contact *flows.ContactReference) *MsgSentEvent {
	event := MsgSentEvent{
		BaseEvent:    NewBaseEvent(),
		Text:         text,
		Attachments:  attachments,
		QuickReplies: quickReples,
		Contacts:     []*flows.ContactReference{contact},
	}
	return &event
}

// NewMsgSentEvent creates a new outgoing msg event for the given recipients
func NewMsgSentEvent(text string, attachments []string, quickReples []string, urns []urns.URN, contacts []*flows.ContactReference, groups []*flows.GroupReference) *MsgSentEvent {
	event := MsgSentEvent{
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
func (e *MsgSentEvent) Type() string { return TypeMsgSent }

// Apply applies this event to the given run
func (e *MsgSentEvent) Apply(run flows.FlowRun) error {
	return nil
}
