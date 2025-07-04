package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeIVRCreated, func() flows.Event { return &IVRCreated{} })
}

// TypeIVRCreated is a constant for IVR created events
const TypeIVRCreated string = "ivr_created"

// IVRCreated events are created when an action wants to send an IVR response to the current contact.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "ivr_created",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "msg": {
//	    "uuid": "2d611e17-fb22-457f-b802-b8f7ec5cda5b",
//	    "channel": {"uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf", "name": "Twilio"},
//	    "urn": "tel:+12065551212",
//	    "text": "hi there",
//	    "attachments": ["audio:https://s3.amazon.com/mybucket/attachment.m4a"],
//	    "locale": "eng-US"
//	  }
//	}
//
// @event ivr_created
type IVRCreated struct {
	BaseEvent

	Msg *flows.MsgOut `json:"msg" validate:"required"`
}

// NewIVRCreated creates a new IVR created event
func NewIVRCreated(msg *flows.MsgOut) *IVRCreated {
	return &IVRCreated{
		BaseEvent: NewBaseEvent(TypeIVRCreated),
		Msg:       msg,
	}
}
