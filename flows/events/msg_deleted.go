package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeMsgDeleted, func() flows.Event { return &MsgDeleted{} })
}

const TypeMsgDeleted string = "msg_deleted"

// MsgDeleted events describe the deletion of an incoming message.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "msg_deleted",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "msg_uuid": "2d611e17-fb22-457f-b802-b8f7ec5cda5b"
//	}
//
// @event msg_deleted
type MsgDeleted struct {
	BaseEvent

	MsgUUID   flows.EventUUID `json:"msg_uuid" validate:"required"`
	ByContact bool            `json:"by_contact,omitempty"`
}

// NewMsgDeleted creates a new msg deleted event
func NewMsgDeleted(msgUUID flows.EventUUID, byContact bool) *MsgDeleted {
	return &MsgDeleted{
		BaseEvent: NewBaseEvent(TypeMsgDeleted),
		MsgUUID:   msgUUID,
		ByContact: byContact,
	}
}

var _ flows.Event = (*MsgDeleted)(nil)
