package inputs

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

// TypeMsg is a constant for incoming messages
const TypeMsg string = "msg"

type MsgInput struct {
	BaseInput
	urn  flows.URN
	text string
}

type msgInputEnvelope struct {
	baseInputEnvelope
	URN  flows.URN `json:"urn"  validate:"required"`
	Text string    `json:"text" validate:"required"`
}

func ReadMsgInput(env flows.FlowEnvironment, envelope *utils.TypedEnvelope) (*MsgInput, error) {
	input := MsgInput{}
	i := msgInputEnvelope{}
	err := json.Unmarshal(envelope.Data, &i)
	if err != nil {
		return nil, err
	}

	// TODO channel lookups
	// channel, err := env.GetChannel(i.ChannelUUID)
	//if err != nil {
	//	return nil, err
	//}

	input.BaseInput.SetChannelUUID(i.ChannelUUID)
	input.BaseInput.SetCreatedOn(i.CreatedOn)
	input.urn = i.URN
	input.text = i.Text
	return &input, nil
}

// NewMsgReceivedEvent creates a new incoming msg event for the passed in channel, contact and string
func NewMsgInput(channel flows.ChannelUUID, urn flows.URN, text string) *MsgInput {
	return &MsgInput{BaseInput: BaseInput{channelUUID: channel}, urn: urn, text: text}
}

// Type returns the type of this event
func (i *MsgInput) Type() string { return TypeMsg }

func (i *MsgInput) Event(run flows.FlowRun) flows.Event {
	return events.NewMsgReceivedEvent(i.ChannelUUID(), run.Contact().UUID(), i.urn, i.text)
}

// Resolve resolves the passed in key to a value, returning an error if the key is unknown
func (i *MsgInput) Resolve(key string) interface{} {
	switch key {
	case "urn":
		return i.urn
	case "text":
		return i.text
	}
	return i.BaseInput.Resolve(key)
}

// Default returns our default value if evaluated in a context, our text in our case
func (i *MsgInput) Default() interface{} {
	return i.text
}

// String returns our default value if evaluated in a context, our text in our case
func (i *MsgInput) String() string {
	return i.text
}

var _ flows.Input = (*MsgInput)(nil)
