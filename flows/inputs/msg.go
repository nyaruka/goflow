package inputs

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// TypeMsg is a constant for incoming messages
const TypeMsg string = "msg"

// MsgInput is a message which can be used as input
type MsgInput struct {
	baseInput
	urn         *flows.ContactURN
	text        string
	attachments flows.AttachmentList
}

// NewMsgInput creates a new user input based on a message
func NewMsgInput(uuid flows.InputUUID, channel flows.Channel, createdOn time.Time, urn urns.URN, text string, attachments []flows.Attachment) *MsgInput {
	return &MsgInput{
		baseInput:   baseInput{uuid: uuid, channel: channel, createdOn: createdOn},
		urn:         flows.NewContactURN(urn, nil),
		text:        text,
		attachments: attachments,
	}
}

// Type returns the type of this event
func (i *MsgInput) Type() string { return TypeMsg }

// Resolve resolves the given key when this input is referenced in an expression
func (i *MsgInput) Resolve(key string) types.XValue {
	switch key {
	case "type":
		return types.NewXText(TypeMsg)
	case "urn":
		return i.urn
	case "text":
		return types.NewXText(i.text)
	case "attachments":
		return i.attachments
	}
	return i.baseInput.Resolve(key)
}

// Describe returns a representation of this type for error messages
func (i *MsgInput) Describe() string { return "input" }

// Reduce is called when this object needs to be reduced to a primitive
func (i *MsgInput) Reduce() types.XPrimitive {
	var parts []string
	if i.text != "" {
		parts = append(parts, i.text)
	}
	for _, attachment := range i.attachments {
		parts = append(parts, attachment.URL())
	}
	return types.NewXText(strings.Join(parts, "\n"))
}

// ToXJSON is called when this type is passed to @(json(...))
func (i *MsgInput) ToXJSON() types.XText {
	return types.ResolveKeys(i, "uuid", "created_on", "channel", "type", "urn", "text", "attachments").ToXJSON()
}

var _ types.XValue = (*MsgInput)(nil)
var _ types.XResolvable = (*MsgInput)(nil)
var _ flows.Input = (*MsgInput)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type msgInputEnvelope struct {
	baseInputEnvelope
	URN         urns.URN             `json:"urn" validate:"omitempty,urn"`
	Text        string               `json:"text" validate:"required"`
	Attachments flows.AttachmentList `json:"attachments,omitempty"`
}

// ReadMsgInput reads a message input from the given JSON
func ReadMsgInput(session flows.Session, data json.RawMessage) (*MsgInput, error) {
	input := MsgInput{}
	i := msgInputEnvelope{}
	err := utils.UnmarshalAndValidate(data, &i, "")
	if err != nil {
		return nil, err
	}

	// lookup the channel
	var channel flows.Channel
	if i.Channel != nil {
		channel, err = session.Assets().GetChannel(i.Channel.UUID)
		if err != nil {
			return nil, err
		}
	}

	input.baseInput.uuid = i.UUID
	input.baseInput.channel = channel
	input.baseInput.createdOn = i.CreatedOn
	input.urn = flows.NewContactURN(i.URN, nil)
	input.text = i.Text
	input.attachments = i.Attachments
	return &input, nil
}

// MarshalJSON marshals this msg input into JSON
func (i *MsgInput) MarshalJSON() ([]byte, error) {
	var envelope msgInputEnvelope

	if i.Channel() != nil {
		envelope.baseInputEnvelope.Channel = i.Channel().Reference()
	}
	envelope.baseInputEnvelope.UUID = i.UUID()
	envelope.baseInputEnvelope.CreatedOn = i.CreatedOn()
	envelope.URN = i.urn.URN
	envelope.Text = i.text
	envelope.Attachments = i.attachments

	return json.Marshal(envelope)
}
