package inputs

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeMsg, readMsgInput)
}

// TypeMsg is a constant for incoming messages
const TypeMsg string = "msg"

// MsgInput is a message which can be used as input
type MsgInput struct {
	baseInput
	urn         *flows.ContactURN
	text        string
	attachments []utils.Attachment
}

// NewMsgInput creates a new user input based on a message
func NewMsgInput(assets flows.SessionAssets, msg *flows.MsgIn, createdOn time.Time) (*MsgInput, error) {
	// load the channel
	var channel *flows.Channel
	if msg.Channel() != nil {
		channel = assets.Channels().Get(msg.Channel().UUID)
	}

	return &MsgInput{
		baseInput:   newBaseInput(TypeMsg, flows.InputUUID(msg.UUID()), channel, createdOn),
		urn:         flows.NewContactURN(msg.URN(), nil),
		text:        msg.Text(),
		attachments: msg.Attachments(),
	}, nil
}

// Context returns the properties available in expressions
func (i *MsgInput) Context(env utils.Environment) map[string]types.XValue {
	attachments := make([]types.XValue, len(i.attachments))

	for i, attachment := range i.attachments {
		attachments[i] = types.NewXText(string(attachment))
	}

	var urn types.XValue
	if i.urn != nil {
		urn = i.urn.ToXValue(env)
	}

	return map[string]types.XValue{
		"__default__": types.NewXText(i.format()),
		"type":        types.NewXText(i.type_),
		"uuid":        types.NewXText(string(i.uuid)),
		"created_on":  types.NewXDateTime(i.createdOn),
		"channel":     flows.Context(env, i.channel),
		"urn":         urn,
		"text":        types.NewXText(i.text),
		"attachments": types.NewXArray(attachments...),
	}
}

func (i *MsgInput) format() string {
	var parts []string
	if i.text != "" {
		parts = append(parts, i.text)
	}
	for _, attachment := range i.attachments {
		parts = append(parts, attachment.URL())
	}
	return strings.Join(parts, "\n")
}

var _ flows.Input = (*MsgInput)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type msgInputEnvelope struct {
	baseInputEnvelope
	URN         urns.URN           `json:"urn" validate:"omitempty,urn"`
	Text        string             `json:"text"`
	Attachments []utils.Attachment `json:"attachments,omitempty"`
}

func readMsgInput(sessionAssets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Input, error) {
	e := &msgInputEnvelope{}
	err := utils.UnmarshalAndValidate(data, e)
	if err != nil {
		return nil, err
	}

	i := &MsgInput{
		urn:         flows.NewContactURN(e.URN, nil),
		text:        e.Text,
		attachments: e.Attachments,
	}

	if err := i.unmarshal(sessionAssets, &e.baseInputEnvelope, missing); err != nil {
		return nil, err
	}

	return i, nil
}

// MarshalJSON marshals this msg input into JSON
func (i *MsgInput) MarshalJSON() ([]byte, error) {
	e := &msgInputEnvelope{
		URN:         i.urn.URN(),
		Text:        i.text,
		Attachments: i.attachments,
	}

	i.marshal(&e.baseInputEnvelope)

	return json.Marshal(e)
}
