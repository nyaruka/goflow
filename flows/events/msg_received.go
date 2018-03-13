package events

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/inputs"
)

// TypeMsgReceived is a constant for incoming messages
const TypeMsgReceived string = "msg_received"

// MsgReceivedEvent events are used for starting flows or resuming flows which are waiting for a message.
// They represent an incoming message for a contact.
//
// ```
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
// ```
//
// @event msg_received
type MsgReceivedEvent struct {
	baseEvent
	callerOnlyEvent

	Msg flows.MsgIn `json:"msg" validate:"required,dive"`
}

// NewMsgReceivedEvent creates a new incoming msg event for the passed in channel, URN and text
func NewMsgReceivedEvent(msg *flows.MsgIn) *MsgReceivedEvent {
	return &MsgReceivedEvent{
		baseEvent: newBaseEvent(),
		Msg:       *msg,
	}
}

// Type returns the type of this event
func (e *MsgReceivedEvent) Type() string { return TypeMsgReceived }

// Validate validates our event is valid and has all the assets it needs
func (e *MsgReceivedEvent) Validate(assets flows.SessionAssets) error {
	if e.Msg.Channel() != nil {
		_, err := assets.GetChannel(e.Msg.Channel().UUID)
		return err
	}
	return nil
}

// Apply applies this event to the given run
func (e *MsgReceivedEvent) Apply(run flows.FlowRun) error {
	var channel flows.Channel
	var err error
	if e.Msg.Channel() != nil {
		channel, err = run.Session().Assets().GetChannel(e.Msg.Channel().UUID)
		if err != nil {
			return err
		}
	}

	// update this run's input
	input := inputs.NewMsgInput(flows.InputUUID(e.Msg.UUID()), channel, e.CreatedOn(), e.Msg.URN(), e.Msg.Text(), e.Msg.Attachments())
	run.SetInput(input)
	run.ResetExpiration(nil)
	return nil
}
