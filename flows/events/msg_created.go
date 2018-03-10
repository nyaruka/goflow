package events

import (
	"github.com/nyaruka/goflow/flows"
)

// TypeMsgCreated is a constant for incoming messages
const TypeMsgCreated string = "msg_created"

// MsgCreatedEvent events are used for replies to the session contact.
//
// ```
//   {
//     "type": "msg_created",
//     "created_on": "2006-01-02T15:04:05Z",
//     "msg": {
//       "uuid": "2d611e17-fb22-457f-b802-b8f7ec5cda5b",
//       "channel": {"uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf", "name": "Twilio"},
//       "urn": "tel:+12065551212",
//       "text": "hi there",
//       "attachments": ["https://s3.amazon.com/mybucket/attachment.jpg"]
//     }
//   }
// ```
//
// @event msg_created
type MsgCreatedEvent struct {
	BaseEvent
	Msg flows.MsgOut `json:"msg" validate:"required,dive"`
}

// NewMsgCreatedEvent creates a new outgoing msg event to a single contact
func NewMsgCreatedEvent(msg *flows.MsgOut) *MsgCreatedEvent {
	return &MsgCreatedEvent{
		BaseEvent: NewBaseEvent(),
		Msg:       *msg,
	}
}

// Type returns the type of this event
func (e *MsgCreatedEvent) Type() string { return TypeMsgCreated }

// AllowedOrigin determines where this event type can originate
func (e *MsgCreatedEvent) AllowedOrigin() flows.EventOrigin { return flows.EventOriginEngine }

// Validate validates our event is valid and has all the assets it needs
func (e *MsgCreatedEvent) Validate(assets flows.SessionAssets) error {
	return nil
}

// Apply applies this event to the given run
func (e *MsgCreatedEvent) Apply(run flows.FlowRun) error {
	return nil
}
