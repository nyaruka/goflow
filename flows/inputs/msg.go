package inputs

import (
	"encoding/json"
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

// ToXValue returns a representation of this object for use in expressions
func (i *MsgInput) ToXValue(env utils.Environment) types.XValue {
	attachments := types.NewXArray()
	for _, attachment := range i.attachments {
		attachments.Append(types.NewXText(string(attachment)))
	}

	return types.NewXDict(map[string]types.XValue{
		"type":        types.NewXText(i.type_),
		"uuid":        types.NewXText(string(i.uuid)),
		"created_on":  types.NewXDateTime(i.createdOn),
		"channel":     types.ToXValue(env, i.channel),
		"urn":         types.ToXValue(env, i.urn),
		"text":        types.NewXText(i.text),
		"attachments": attachments,
	})
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
