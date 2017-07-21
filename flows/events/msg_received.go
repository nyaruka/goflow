package events

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/inputs"
)

// TypeMsgReceived is a constant for incoming messages
const TypeMsgReceived string = "msg_received"

// MsgReceivedEvent events are used for resuming flows or starting flows. They represent an MO
// message for a contact.
//
// ```
//   {
//    "type": "msg_received",
//    "step_uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//    "created_on": "2006-01-02T15:04:05Z",
//    "urn": "tel:+12065551212",
//    "channel_uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf",
//    "contact_uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
//    "text": "hi there"
//   }
// ```
//
// @event msg_received
type MsgReceivedEvent struct {
	BaseEvent
	ChannelUUID flows.ChannelUUID `json:"channel_uuid" validate:"uuid4"`
	URN         flows.URN         `json:"urn"          validate:"required"`
	ContactUUID flows.ContactUUID `json:"contact_uuid" validate:"required,uuid4"`
	Text        string            `json:"text"`
}

// NewMsgReceivedEvent creates a new incoming msg event for the passed in channel, contact and string
func NewMsgReceivedEvent(channel flows.ChannelUUID, contact flows.ContactUUID, urn flows.URN, text string) *MsgReceivedEvent {
	return &MsgReceivedEvent{
		BaseEvent:   NewBaseEvent(),
		ChannelUUID: channel,
		ContactUUID: contact,
		URN:         urn,
		Text:        text,
	}
}

// Type returns the type of this event
func (e *MsgReceivedEvent) Type() string { return TypeMsgReceived }

// Apply applies this event to the given run
func (e *MsgReceivedEvent) Apply(run flows.FlowRun) {
	channel, _ := run.Environment().GetChannel(e.ChannelUUID)

	// update this run's input
	run.SetInput(inputs.NewMsgInput(channel, e.CreatedOn(), e.URN, e.Text))
}
