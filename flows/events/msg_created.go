package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeMsgCreated, func() flows.Event { return &MsgCreatedEvent{} })
}

// TypeMsgCreated is a constant for incoming messages
const TypeMsgCreated string = "msg_created"

// MsgCreatedEvent events are created when an action wants to send a reply to the current contact.
//
//   {
//     "type": "msg_created",
//     "created_on": "2006-01-02T15:04:05Z",
//     "msg": {
//       "uuid": "2d611e17-fb22-457f-b802-b8f7ec5cda5b",
//       "channel": {"uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf", "name": "Twilio"},
//       "urn": "tel:+12065551212",
//       "text": "hi there",
//       "attachments": ["image/jpeg:https://s3.amazon.com/mybucket/attachment.jpg"]
//     }
//   }
//
// @event msg_created
type MsgCreatedEvent struct {
	BaseEvent

	Msg *flows.MsgOut `json:"msg" validate:"required,dive"`
}

// NewMsgCreated creates a new outgoing msg event to a single contact
func NewMsgCreated(msg *flows.MsgOut) *MsgCreatedEvent {
	return &MsgCreatedEvent{
		BaseEvent: NewBaseEvent(TypeMsgCreated),
		Msg:       msg,
	}
}
