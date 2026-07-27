package inputs

import (
	"encoding/json"
	"strings"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/core/events"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeMsg, readMsg)
}

// TypeMsg is a constant for incoming messages
const TypeMsg string = "msg"

// Msg is a message which can be used as input
type Msg struct {
	baseInput

	urn         *core.URN
	text        string
	attachments []utils.Attachment
	externalID  string
	payload     json.RawMessage
}

// NewMsg creates a new user input based on a message event
func NewMsg(sa flows.SessionAssets, evt *events.MsgReceived) *Msg {
	// load the channel
	var channel *core.Channel
	if evt.Msg.Channel() != nil {
		channel = sa.Channels().Get(evt.Msg.Channel().UUID)
	}

	return &Msg{
		baseInput:   newBaseInput(TypeMsg, core.InputUUID(evt.UUID()), channel, evt.CreatedOn()),
		urn:         core.NewURN(evt.Msg.URN().Scheme(), evt.Msg.URN().Path(), "", nil),
		text:        evt.Msg.Text(),
		attachments: evt.Msg.Attachments(),
		externalID:  evt.Msg.ExternalID(),
		payload:     evt.Msg.Payload(),
	}
}

// Context returns the properties available in expressions
//
//	__default__:text -> the text and attachments
//	uuid:text -> the UUID of the input
//	created_on:datetime -> the creation date of the input
//	channel:channel -> the channel that the input was received on
//	urn:text -> the contact URN that the input was received on
//	text:text -> the text part of the input
//	attachments:[]text -> any attachments on the input
//	external_id:text -> the external ID of the input
//	payload:any -> any structured data on the input, e.g. a form submission
//
// @context input
func (i *Msg) Context(env envs.Environment) map[string]types.XValue {
	attachments := make([]types.XValue, len(i.attachments))

	for i, attachment := range i.attachments {
		attachments[i] = types.NewXText(string(attachment))
	}

	var urn types.XValue
	if i.urn != nil {
		urn = i.urn.ToXValue(env)
	}

	payload := types.JSONToXValue(i.payload)
	if payload == nil {
		payload = types.XObjectEmpty
	}

	return map[string]types.XValue{
		"__default__": types.NewXText(i.format()),
		"type":        types.NewXText(i.type_),
		"uuid":        types.NewXText(string(i.uuid)),
		"created_on":  types.NewXDateTime(i.createdOn),
		"channel":     core.Context(env, i.channel),
		"urn":         urn,
		"text":        types.NewXText(i.text),
		"attachments": types.NewXArray(attachments...),
		"external_id": types.NewXText(i.externalID),
		"payload":     payload,
	}
}

func (i *Msg) format() string {
	var parts []string
	if i.text != "" {
		parts = append(parts, i.text)
	}
	for _, attachment := range i.attachments {
		parts = append(parts, attachment.URL())
	}
	return strings.Join(parts, "\n")
}

var _ flows.Input = (*Msg)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type msgEnvelope struct {
	baseEnvelope

	URN         urns.URN           `json:"urn" validate:"omitempty,urn"`
	Text        string             `json:"text"`
	Attachments []utils.Attachment `json:"attachments,omitempty"`
	ExternalID  string             `json:"external_id,omitempty"`
	Payload     json.RawMessage    `json:"payload,omitempty"`
}

func readMsg(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Input, error) {
	e := &msgEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	i := &Msg{
		urn:         core.NewURN(e.URN.Scheme(), e.URN.Path(), "", nil),
		text:        e.Text,
		attachments: e.Attachments,
		externalID:  e.ExternalID,
		payload:     e.Payload,
	}

	if err := i.unmarshal(sa, &e.baseEnvelope, missing); err != nil {
		return nil, err
	}

	return i, nil
}

// MarshalJSON marshals this msg input into JSON
func (i *Msg) MarshalJSON() ([]byte, error) {
	e := &msgEnvelope{
		URN:         i.urn.Identity(),
		Text:        i.text,
		Attachments: i.attachments,
		ExternalID:  i.externalID,
		Payload:     i.payload,
	}

	i.marshal(&e.baseEnvelope)

	return jsonx.Marshal(e)
}
