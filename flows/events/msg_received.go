package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeMsgReceived, func() flows.Event { return &MsgReceivedEvent{} })
}

// TypeMsgReceived is a constant for incoming messages
const TypeMsgReceived string = "msg_received"

// MsgReceivedEvent events are sent by the caller to tell the engine that a message was received from
// the contact and that it should try to resume the session.
//
//   {
//     "type": "msg_received",
//     "created_on": "2006-01-02T15:04:05Z",
//     "msg": {
//       "uuid": "2d611e17-fb22-457f-b802-b8f7ec5cda5b",
//       "channel": {"uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf", "name": "Twilio"},
//       "urn": "tel:+12065551212",
//       "text": "hi there",
//       "attachments": ["https://s3.amazon.com/mybucket/attachment.jpg"]
//     }
//   }
//
// @event msg_received
type MsgReceivedEvent struct {
	BaseEvent

	Msg flows.MsgIn `json:"msg" validate:"required,dive"`
}

// NewMsgReceived creates a new incoming msg event for the passed in channel, URN and text
func NewMsgReceived(msg *flows.MsgIn) *MsgReceivedEvent {
	return &MsgReceivedEvent{
		BaseEvent: NewBaseEvent(TypeMsgReceived),
		Msg:       *msg,
	}
}

var _ flows.Event = (*MsgReceivedEvent)(nil)
