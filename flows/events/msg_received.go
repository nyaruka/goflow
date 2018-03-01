package events

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/inputs"
)

// TypeMsgReceived is a constant for incoming messages
const TypeMsgReceived string = "msg_received"

// MsgReceivedEvent events are used for starting flows or resuming flows which are waiting for a message.
// They represent an MO message for a contact.
//
// ```
//   {
//     "type": "msg_received",
//     "created_on": "2006-01-02T15:04:05Z",
//     "msg": {
//       "uuid": "2d611e17-fb22-457f-b802-b8f7ec5cda5b",
//       "channel": {"uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf", "name": "Twilio"},
//       "created_on": "2006-01-02T15:04:04.454789Z",
//       "urn": "tel:+12065551212",
//       "text": "hi there",
//       "attachments": ["https://s3.amazon.com/mybucket/attachment.jpg"]
//     }
//   }
// ```
//
// @event msg_received
type MsgReceivedEvent struct {
	BaseEvent
	Msg json.RawMessage `json:"msg" validate:"required"`
}

// NewMsgReceivedEvent creates a new incoming msg event for the passed in channel, URN and text
func NewMsgReceivedEvent(msg *inputs.MsgInput) *MsgReceivedEvent {
	// these events are only generated in tests
	msgBytes, _ := json.Marshal(msg)

	return &MsgReceivedEvent{
		BaseEvent: NewBaseEvent(),
		Msg:       msgBytes,
	}
}

// Type returns the type of this event
func (e *MsgReceivedEvent) Type() string { return TypeMsgReceived }

// Apply applies this event to the given run
func (e *MsgReceivedEvent) Apply(run flows.FlowRun) error {
	msgInput, err := inputs.ReadMsgInput(run.Session(), e.Msg)
	if err != nil {
		return err
	}

	// update this run's input
	run.SetInput(msgInput)
	run.ResetExpiration(nil)
	return nil
}
