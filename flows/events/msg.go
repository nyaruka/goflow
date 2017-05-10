package events

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
)

// MsgDirection is the direction of a Msg (either in or out)
type MsgDirection string

const (
	// MsgOut represents an outgoing message
	MsgOut MsgDirection = "O"

	// MsgIn represents an incoming message
	MsgIn MsgDirection = "I"
)

// MSG_IN is a constant for incoming messages
const MSG_IN string = "msg_in"

// MSG_OUT is a constant for incoming messages
const MSG_OUT string = "msg_out"

// NewOutgoingMsgEvent creates a new outgoing msg event for the passed in channel, contact and string
func NewOutgoingMsgEvent(channel flows.ChannelUUID, contact flows.ContactUUID, text string) *MsgOutEvent {
	event := MsgOutEvent{Channel: channel, Contact: contact, Text: text}
	return &event
}

// NewIncomingMsgEvent creates a new incoming msg event for the passed in channel, contact and string
func NewIncomingMsgEvent(channel flows.ChannelUUID, contact flows.ContactUUID, text string) *MsgInEvent {
	event := MsgInEvent{Channel: channel, Contact: contact, Text: text}
	return &event
}

// MsgInEvent represents either an incoming message
type MsgInEvent struct {
	ID      int64             `json:"id"` // can be unset/zero for new outgoing msgs
	Channel flows.ChannelUUID `json:"channel"     validate:"nonzero"`
	Contact flows.ContactUUID `json:"contact"     validate:"nonzero"`
	URN     flows.URN         `json:"urn"         validate:"nonzero"`
	Text    string            `json:"text"        validate:"nonzero"`
	BaseEvent
}

// Type returns the type of this event
func (e *MsgInEvent) Type() string { return MSG_IN }

// MsgOutEvent represents either an outgoing message
type MsgOutEvent struct {
	ID      int64             `json:"id"` // can be unset/zero for new outgoing msgs
	Channel flows.ChannelUUID `json:"channel"     validate:"nonzero"`
	Contact flows.ContactUUID `json:"contact"     validate:"nonzero"`
	URN     flows.URN         `json:"urn"         validate:"nonzero"`
	Text    string            `json:"text"        validate:"nonzero"`
	BaseEvent
}

// Type returns the type of this event
func (e *MsgOutEvent) Type() string { return MSG_OUT }

// Resolve resolves the passed in key to a value, returning an error if the key is unknown
func (e *MsgInEvent) Resolve(key string) interface{} {
	switch key {

	case "id":
		return e.ID

	case "direction":
		return MsgIn

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
