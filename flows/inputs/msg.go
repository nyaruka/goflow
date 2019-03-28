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

// Resolve resolves the given key when this input is referenced in an expression
func (i *MsgInput) Resolve(env utils.Environment, key string) types.XValue {
	switch strings.ToLower(key) {
	case "urn":
		if i.urn != nil {
			return i.urn.Context(env)
		}
		return nil
	case "text":
		return types.NewXText(i.text)
	case "attachments":
		attachments := types.NewXArray()
		for _, attachment := range i.attachments {
			attachments.Append(types.NewXText(string(attachment)))
		}
		return attachments
	}
	return i.baseInput.Resolve(env, key)
}

// Describe returns a representation of this type for error messages
func (i *MsgInput) Describe() string { return "input" }

// Reduce is called when this object needs to be reduced to a primitive
func (i *MsgInput) Reduce(env utils.Environment) types.XPrimitive {
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
func (i *MsgInput) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, i, "uuid", "created_on", "channel", "type", "urn", "text", "attachments").ToXJSON(env)
}

var _ types.XValue = (*MsgInput)(nil)
var _ types.XResolvable = (*MsgInput)(nil)
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
