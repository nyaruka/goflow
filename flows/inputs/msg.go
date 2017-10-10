package inputs

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// TypeMsg is a constant for incoming messages
const TypeMsg string = "msg"

type MsgInput struct {
	baseInput
	urn         urns.URN
	text        string
	attachments []flows.Attachment
}

// NewMsgInput creates a new user input based on a message
func NewMsgInput(uuid flows.InputUUID, channel flows.Channel, createdOn time.Time, urn urns.URN, text string, attachments []flows.Attachment) *MsgInput {
	return &MsgInput{
		baseInput:   baseInput{uuid: uuid, channel: channel, createdOn: createdOn},
		urn:         urn,
		text:        text,
		attachments: attachments,
	}
}

// Type returns the type of this event
func (i *MsgInput) Type() string { return TypeMsg }

// Resolve resolves the passed in key to a value, returning an error if the key is unknown
func (i *MsgInput) Resolve(key string) interface{} {
	switch key {

	case "urn":
		return i.urn

	case "text":
		return i.text

	case "attachments":
		return i.attachments
	}
	return i.baseInput.Resolve(key)
}

// Default returns our default value if evaluated in a context, which in this case is the text and attachments combined
func (i *MsgInput) Default() interface{} {
	var parts []string
	if i.text != "" {
		parts = append(parts, i.text)
	}
	for _, attachment := range i.attachments {
		parts = append(parts, attachment.URL())
	}
	return strings.Join(parts, "\n")
}

// String returns our default value if evaluated in a context, our text in our case
func (i *MsgInput) String() string {
	return i.text
}

var _ flows.Input = (*MsgInput)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type msgInputEnvelope struct {
	baseInputEnvelope
	URN         urns.URN           `json:"urn" validate:"urn"`
	Text        string             `json:"text" validate:"required"`
	Attachments []flows.Attachment `json:"attachments,omitempty"`
}

func ReadMsgInput(session flows.Session, envelope *utils.TypedEnvelope) (*MsgInput, error) {
	input := MsgInput{}
	i := msgInputEnvelope{}
	err := json.Unmarshal(envelope.Data, &i)
	if err != nil {
		return nil, err
	}

	err = utils.Validate(i)
	if err != nil {
		return nil, err
	}

	// lookup the channel
	var channel flows.Channel
	if i.ChannelUUID != "" {
		channel, err = session.Assets().GetChannel(i.ChannelUUID)
		if err != nil {
			return nil, err
		}
	}

	input.baseInput.uuid = i.UUID
	input.baseInput.channel = channel
	input.baseInput.createdOn = i.CreatedOn
	input.urn = i.URN
	input.text = i.Text
	input.attachments = i.Attachments
	return &input, nil
}

func (i *MsgInput) MarshalJSON() ([]byte, error) {
	var envelope msgInputEnvelope

	if i.Channel() != nil {
		envelope.baseInputEnvelope.ChannelUUID = i.Channel().UUID()
	}
	envelope.baseInputEnvelope.UUID = i.UUID()
	envelope.baseInputEnvelope.CreatedOn = i.CreatedOn()
	envelope.URN = i.urn
	envelope.Text = i.text
	envelope.Attachments = i.attachments

	return json.Marshal(envelope)
}
