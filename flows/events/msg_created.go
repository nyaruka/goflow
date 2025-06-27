package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeMsgCreated, func() flows.Event { return &MsgCreated{} })
}

// TypeMsgCreated is a constant for incoming messages
const TypeMsgCreated string = "msg_created"

// MsgCreated events are created when an action wants to send a reply to the current contact.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "msg_created",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "msg": {
//	    "uuid": "2d611e17-fb22-457f-b802-b8f7ec5cda5b",
//	    "channel": {"uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf", "name": "Twilio"},
//	    "urn": "tel:+12065551212",
//	    "text": "hi there",
//	    "attachments": ["image/jpeg:https://s3.amazon.com/mybucket/attachment.jpg"]
//	  }
//	}
//
// @event msg_created
type MsgCreated struct {
	BaseEvent

	Msg *flows.MsgOut `json:"msg" validate:"required"`
}

// NewMsgCreated creates a new outgoing msg event to a single contact
func NewMsgCreated(msg *flows.MsgOut) *MsgCreated {
	return &MsgCreated{
		BaseEvent: NewBaseEvent(TypeMsgCreated),
		Msg:       msg,
	}
}
