package inputs

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeMsg, readMsgInput)
}

// TypeMsg is a constant for incoming messages
const TypeMsg string = "msg"

// MsgInput is a message which can be used as input
type MsgInput struct {
	baseInput

	urn         *flows.ContactURN
	text        string
	attachments []utils.Attachment
	externalID  string
}

// NewMsg creates a new user input based on a message
func NewMsg(assets flows.SessionAssets, msg *flows.MsgIn, createdOn time.Time) *MsgInput {
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
		externalID:  msg.ExternalID(),
	}
}

// Context returns the properties available in expressions
//
//   __default__:text -> the text and attachments
//   uuid:text -> the UUID of the input
//   created_on:datetime -> the creation date of the input
//   channel:channel -> the channel that the input was received on
//   urn:text -> the contact URN that the input was received on
//   text:text -> the text part of the input
//   attachments:[]text -> any attachments on the input
//   external_id:text -> the external ID of the input
//
// @context input
func (i *MsgInput) Context(env envs.Environment) map[string]types.XValue {
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
		"external_id": types.NewXText(i.externalID),
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
	ExternalID  string             `json:"external_id,omitempty"`
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
		externalID:  e.ExternalID,
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
		ExternalID:  i.externalID,
	}

	i.marshal(&e.baseInputEnvelope)

	return jsonx.Marshal(e)
}
