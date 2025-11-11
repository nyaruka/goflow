package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeMsgStatusChanged, func() flows.Event { return &MsgStatusChanged{} })
}

const TypeMsgStatusChanged string = "msg_status_changed"

// MsgStatusChanged events describe an outgoing message status change.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "msg_status_changed",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "msg_uuid": "2d611e17-fb22-457f-b802-b8f7ec5cda5b",
//	  "status": "delivered"
//	}
//
// @event msg_status_changed
type MsgStatusChanged struct {
	BaseEvent

	MsgUUID flows.EventUUID `json:"msg_uuid" validate:"required"`
	Status  string          `json:"status"   validate:"required"`
	Reason  string          `json:"reason,omitempty"`
}

// NewMsgStatusChanged creates a new message status changed event
func NewMsgStatusChanged(msgUUID flows.EventUUID, status string, reason string) *MsgStatusChanged {
	return &MsgStatusChanged{
		BaseEvent: NewBaseEvent(TypeMsgStatusChanged),
		MsgUUID:   msgUUID,
		Status:    status,
		Reason:    reason,
	}
}

var _ flows.Event = (*MsgStatusChanged)(nil)
