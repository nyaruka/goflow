package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
)

// TypeSendMsg is a constant for incoming messages
const TypeSendMsg string = "send_msg"

// SendMsgEvent events are created for outgoing messages.
//
// ```
//   {
//     "type": "send_msg",
//     "created_on": "2006-01-02T15:04:05Z",
//     "text": "hi, what's up",
//     "attachments": [],
//     "urns": ["tel:+12065551212"],
//     "contacts": [{"uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a", "name": "Bob"}]
//   }
// ```
//
// @event send_msg
type SendMsgEvent struct {
	BaseEvent
	Text         string                    `json:"text"`
	Attachments  []string                  `json:"attachments,omitempty"`
	QuickReplies []string                  `json:"quick_replies,omitempty"`
	URNs         []urns.URN                `json:"urns,omitempty" validate:"dive,urn"`
	Contacts     []*flows.ContactReference `json:"contacts,omitempty" validate:"dive"`
	Groups       []*flows.GroupReference   `json:"groups,omitempty" validate:"dive"`
}

// NewSendMsgToContactEvent creates a new outgoing msg event to a single contact
func NewSendMsgToContactEvent(text string, attachments []string, contact *flows.ContactReference) *SendMsgEvent {
	event := SendMsgEvent{
		BaseEvent:   NewBaseEvent(),
		Text:        text,
		Attachments: attachments,
		Contacts:    []*flows.ContactReference{contact},
	}
	return &event
}

// NewSendMsgEvent creates a new outgoing msg event for the given recipients
func NewSendMsgEvent(text string, attachments []string, urns []urns.URN, contacts []*flows.ContactReference, groups []*flows.GroupReference) *SendMsgEvent {
	event := SendMsgEvent{
		BaseEvent:   NewBaseEvent(),
		Text:        text,
		Attachments: attachments,
		URNs:        urns,
		Contacts:    contacts,
		Groups:      groups,
	}
	return &event
}

// Type returns the type of this event
func (e *SendMsgEvent) Type() string { return TypeSendMsg }

// Apply applies this event to the given run
func (e *SendMsgEvent) Apply(run flows.FlowRun) error {
	return nil
}
