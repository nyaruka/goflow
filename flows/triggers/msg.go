package triggers

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/inputs"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeMsg, ReadMsgTrigger)
}

// TypeMsg is the type for message triggered sessions
const TypeMsg string = "msg"

// MsgTrigger is used when a session was triggered by a message being recieved by the caller
//
//   {
//     "type": "msg",
//     "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
//     "contact": {
//       "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
//       "name": "Bob"
//     },
//     "msg": {
//       "uuid": "2d611e17-fb22-457f-b802-b8f7ec5cda5b",
//       "channel": {"uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf", "name": "Twilio"},
//       "urn": "tel:+12065551212",
//       "text": "hi there",
//       "attachments": ["https://s3.amazon.com/mybucket/attachment.jpg"]
//     },
//     "keyword_match": {
//       "type": "first_word",
//       "keyword": "start"
//     },
//     "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//   }
//
// @trigger msg
type MsgTrigger struct {
	baseTrigger
	msg   *flows.MsgIn
	match *KeywordMatch
}

// KeywordMatchType describes how the message matched a keyword
type KeywordMatchType string

// the different types of keyword match
const (
	KeywordMatchTypeFirstWord KeywordMatchType = "first_word"
	KeywordMatchTypeOnlyWord  KeywordMatchType = "only_word"
)

// KeywordMatch describes why the message triggered a session
type KeywordMatch struct {
	Type    KeywordMatchType `json:"type" validate:"required"`
	Keyword string           `json:"keyword" validate:"required"`
}

// NewMsgTrigger creates a new message trigger
func NewMsgTrigger(env utils.Environment, contact *flows.Contact, flow *assets.FlowReference, params types.XValue, msg *flows.MsgIn, match *KeywordMatch, triggeredOn time.Time) flows.Trigger {
	return &MsgTrigger{
		baseTrigger: baseTrigger{environment: env, contact: contact, flow: flow, triggeredOn: triggeredOn},
		msg:         msg,
		match:       match,
	}
}

// InitializeRun performs additional initialization when we visit our first node
func (t *MsgTrigger) InitializeRun(run flows.FlowRun, step flows.Step) error {
	// update the run's input
	input, err := inputs.NewMsgInput(run.Session().Assets(), t.msg, t.triggeredOn)
	if err != nil {
		return err
	}

	run.SetInput(input)
	run.LogEvent(step, events.NewMsgReceivedEvent(t.msg))
	return nil
}

// Type returns the type of this trigger
func (t *MsgTrigger) Type() string { return TypeMsg }

// Resolve resolves the given key when this trigger is referenced in an expression
func (t *MsgTrigger) Resolve(env utils.Environment, key string) types.XValue {
	switch key {
	case "type":
		return types.NewXText(TypeMsg)
	}

	return t.baseTrigger.Resolve(env, key)
}

// ToXJSON is called when this type is passed to @(json(...))
func (t *MsgTrigger) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, t, "type", "params").ToXJSON(env)
}

var _ flows.Trigger = (*MsgTrigger)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type msgTriggerEnvelope struct {
	baseTriggerEnvelope
	Msg   *flows.MsgIn  `json:"msg" validate:"required,dive"`
	Match *KeywordMatch `json:"keyword_match,omitempty" validate:"omitempty,dive"`
}

// ReadMsgTrigger reads a message trigger
func ReadMsgTrigger(session flows.Session, data json.RawMessage) (flows.Trigger, error) {
	trigger := &MsgTrigger{}
	e := msgTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, &e); err != nil {
		return nil, err
	}

	if err := unmarshalBaseTrigger(session, &trigger.baseTrigger, &e.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	trigger.msg = e.Msg
	trigger.match = e.Match

	return trigger, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *MsgTrigger) MarshalJSON() ([]byte, error) {
	var envelope msgTriggerEnvelope

	if err := marshalBaseTrigger(&t.baseTrigger, &envelope.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	envelope.Msg = t.msg
	envelope.Match = t.match

	return json.Marshal(envelope)
}
