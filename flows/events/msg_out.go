package events

import "github.com/nyaruka/goflow/flows"

// TypeMsgOut is a constant for incoming messages
const TypeMsgOut string = "msg_out"

// MsgOutEvent events are created for each outgoing message. They represent an MT message to a
// contact, urn or group.
//
// ```
//   {
//    "step": "8eebd020-1af5-431c-b943-aa670fc74da9",
//    "created_on": "2006-01-02T15:04:05Z",
//    "type": "msg_out",
//    "urn": "tel:%2B12065551212",
//    "contact": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
//    "text": "hi, what's up"
//   }
// ```
//
// @event msg_out
type MsgOutEvent struct {
	BaseEvent
	URN     flows.URN         `json:"urn,omitempty"`
	Contact flows.ContactUUID `json:"contact,omitempty"`
	Group   flows.GroupUUID   `json:"group,omitempty"`
	Text    string            `json:"text"                  validate:"required"`
}

// NewMsgToContact creates a new outgoing msg event for the passed in channel, contact and string
func NewMsgToContact(contact flows.ContactUUID, text string) *MsgOutEvent {
	event := MsgOutEvent{Contact: contact, Text: text}
	return &event
}

// NewMsgToURN creates a new outgoing msg event for the passed in channel, urn and string
func NewMsgToURN(urn flows.URN, text string) *MsgOutEvent {
	event := MsgOutEvent{URN: urn, Text: text}
	return &event
}

// NewMsgToGroup creates a new outgoing msg event for the passed in channel, group and string
func NewMsgToGroup(group flows.GroupUUID, text string) *MsgOutEvent {
	event := MsgOutEvent{Group: group, Text: text}
	return &event
}

// Type returns the type of this event
func (e *MsgOutEvent) Type() string { return TypeMsgOut }
