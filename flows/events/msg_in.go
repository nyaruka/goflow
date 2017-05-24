package events

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
)

// TypeMsgIn is a constant for incoming messages
const TypeMsgIn string = "msg_in"

// MsgInEvent events are used for resuming flows or starting flows. They represent an MO
// message for a contact.
//
// ```
//   {
//    "step": "8eebd020-1af5-431c-b943-aa670fc74da9",
//    "created_on": "2006-01-02T15:04:05Z",
//    "type": "msg_in",
//    "urn": "tel:+12065551212",
//    "channel": "61602f3e-f603-4c70-8a8f-c477505bf4bf",
//    "contact": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
//    "text": "hi there"
//   }
// ```
//
// @event msg_in
type MsgInEvent struct {
	BaseEvent
	ID      int64             `json:"id"`
	Channel flows.ChannelUUID `json:"channel,omitempty"     validate:"required"`
	URN     flows.URN         `json:"urn,omitempty"         validate:"required"`
	Contact flows.ContactUUID `json:"contact,omitempty"     validate:"required"`
	Text    string            `json:"text"                  validate:"required"`
}

// NewIncomingMsgEvent creates a new incoming msg event for the passed in channel, contact and string
func NewIncomingMsgEvent(channel flows.ChannelUUID, contact flows.ContactUUID, text string) *MsgInEvent {
	event := MsgInEvent{Channel: channel, Contact: contact, Text: text}
	return &event
}

// Type returns the type of this event
func (e *MsgInEvent) Type() string { return TypeMsgIn }

// Resolve resolves the passed in key to a value, returning an error if the key is unknown
func (e *MsgInEvent) Resolve(key string) interface{} {
	switch key {

	case "id":
		return e.ID

	case "direction":
		return flows.MsgIn

	case "channel":
		return e.Channel

	case "contact":
		return e.Contact

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
func (e *MsgInEvent) Default() interface{} {
	return e.Text
}

var _ flows.Input = (*MsgInEvent)(nil)
