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
//     "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//   }
//
// @trigger msg
type MsgTrigger struct {
	baseTrigger
	Msg *flows.MsgIn
}

// NewMsgTrigger creates a new message trigger
func NewMsgTrigger(env utils.Environment, contact *flows.Contact, flow *assets.FlowReference, params types.XValue, msg *flows.MsgIn, triggeredOn time.Time) flows.Trigger {
	return &MsgTrigger{
		baseTrigger: baseTrigger{environment: env, contact: contact, flow: flow, triggeredOn: triggeredOn},
		Msg:         msg,
	}
}

// InitializeRun performs additional initialization when we visit our first node
func (t *MsgTrigger) InitializeRun(run flows.FlowRun) error {
	var channel *flows.Channel
	var err error
	if t.Msg.Channel() != nil {
		channel, err = run.Session().Assets().Channels().Get(t.Msg.Channel().UUID)
		if err != nil {
			return err
		}
	}

	// TODO this method is basically the same as MsgResume.Apply

	// update this run's input
	input := inputs.NewMsgInput(flows.InputUUID(t.Msg.UUID()), channel, t.TriggeredOn(), t.Msg.URN(), t.Msg.Text(), t.Msg.Attachments())
	run.SetInput(input)
	run.ResetExpiration(nil)
	run.AddEvent(nil, nil, events.NewMsgReceivedEvent(t.Msg))
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
	Msg *flows.MsgIn `json:"msg" validate:"required,dive"`
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

	trigger.Msg = e.Msg

	return trigger, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *MsgTrigger) MarshalJSON() ([]byte, error) {
	var envelope msgTriggerEnvelope

	if err := marshalBaseTrigger(&t.baseTrigger, &envelope.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	envelope.Msg = t.Msg

	return json.Marshal(envelope)
}
