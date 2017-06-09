package events

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
)

// TypeMsgReceived is a constant for incoming messages
const TypeMsgReceived string = "msg_received"

// MsgReceivedEvent events are used for resuming flows or starting flows. They represent an MO
// message for a contact.
//
// ```
//   {
//    "type": "msg_received",
//    "step_uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//    "created_on": "2006-01-02T15:04:05Z",
//    "urn": "tel:+12065551212",
//    "channel_uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf",
//    "contact_uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
//    "text": "hi there"
//   }
// ```
//
// @event msg_received
type MsgReceivedEvent struct {
	BaseEvent
	ID          int64             `json:"id"`
	ChannelUUID flows.ChannelUUID `json:"channel_uuid"     validate:"required,uuid4"`
	URN         flows.URN         `json:"urn"              validate:"required"`
	ContactUUID flows.ContactUUID `json:"contact_uuid"     validate:"required,uuid4"`
	Text        string            `json:"text"             validate:"required"`
}

// NewMsgReceivedEvent creates a new incoming msg event for the passed in channel, contact and string
func NewMsgReceivedEvent(channel flows.ChannelUUID, contact flows.ContactUUID, urn flows.URN, text string) *MsgReceivedEvent {
	event := MsgReceivedEvent{ChannelUUID: channel, ContactUUID: contact, URN: urn, Text: text}
	return &event
}

// Type returns the type of this event
func (e *MsgReceivedEvent) Type() string { return TypeMsgReceived }

// Resolve resolves the passed in key to a value, returning an error if the key is unknown
func (e *MsgReceivedEvent) Resolve(key string) interface{} {
	switch key {

	case "id":
		return e.ID

	case "direction":
		return flows.MsgIn

	case "channel_uuid":
		return e.ChannelUUID

	case "contact_uuid":
		return e.ContactUUID

	case "urn":
		return e.URN

	case "text":
		return e.Text

	case "created_on":
		return e.CreatedOn

	}
	return fmt.Errorf("No such field '%s' on Msg event", key)
}

// Default returns our default value if evaluated in a context, our text in our case
func (e *MsgReceivedEvent) Default() interface{} {
	return e.Text
}

var _ flows.Input = (*MsgReceivedEvent)(nil)
