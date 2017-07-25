package events

import "github.com/nyaruka/goflow/flows"

// TypeSendMsg is a constant for incoming messages
const TypeSendMsg string = "send_msg"

// SendMsgEvent events are created for each outgoing message. They represent an MT message to a
// contact, urn or group.
//
// ```
//   {
//     "type": "send_msg",
//     "created_on": "2006-01-02T15:04:05Z",
//     "urn": "tel:%2B12065551212",
//     "contact_uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
//     "text": "hi, what's up",
//     "attachments": []
//   }
// ```
//
// @event send_msg
type SendMsgEvent struct {
	BaseEvent
	URN         flows.URN         `json:"urn,omitempty"`
	ContactUUID flows.ContactUUID `json:"contact_uuid,omitempty"`
	GroupUUID   flows.GroupUUID   `json:"group_uuid,omitempty"`
	Text        string            `json:"text"                      validate:"required"`
	Attachments []string          `json:"attachments,omitempty"`
}

// NewSendMsgToContact creates a new outgoing msg event for the passed in channel, contact and string
func NewSendMsgToContact(contact flows.ContactUUID, text string, attachments []string) *SendMsgEvent {
	event := SendMsgEvent{BaseEvent: NewBaseEvent(), ContactUUID: contact, Text: text, Attachments: attachments}
	return &event
}

// NewSendMsgToURN creates a new outgoing msg event for the passed in channel, urn and string
func NewSendMsgToURN(urn flows.URN, text string, attachments []string) *SendMsgEvent {
	event := SendMsgEvent{BaseEvent: NewBaseEvent(), URN: urn, Text: text, Attachments: attachments}
	return &event
}

// NewSendMsgToGroup creates a new outgoing msg event for the passed in channel, group and string
func NewSendMsgToGroup(group flows.GroupUUID, text string, attachments []string) *SendMsgEvent {
	event := SendMsgEvent{BaseEvent: NewBaseEvent(), GroupUUID: group, Text: text, Attachments: attachments}
	return &event
}

// Type returns the type of this event
func (e *SendMsgEvent) Type() string { return TypeSendMsg }

// Apply applies this event to the given run
func (e *SendMsgEvent) Apply(run flows.FlowRun) {}
