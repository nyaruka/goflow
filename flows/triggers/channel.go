package triggers

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeChannel, ReadChannelTrigger)
}

// TypeChannel is the type for sessions triggered by channel events
const TypeChannel string = "channel"

// ChannelEvent describes the specific event on the channel that triggered the session
type ChannelEvent struct {
	Type    string                   `json:"type" validate:"required"`
	Channel *assets.ChannelReference `json:"channel" validate:"required,dive"`
}

// NewChannelEvent creates a new channel event
func NewChannelEvent(typeName string, channel *assets.ChannelReference) *ChannelEvent {
	return &ChannelEvent{Type: typeName, Channel: channel}
}

// ChannelTrigger is used when a session was triggered by a channel event
//
//   {
//     "type": "channel",
//     "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
//     "contact": {
//       "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
//       "name": "Bob",
//       "created_on": "2018-01-01T12:00:00.000000Z"
//     },
//     "event": {
//         "type": "new_conversation",
//         "channel": {"uuid": "58e9b092-fe42-4173-876c-ff45a14a24fe", "name": "Facebook"}
//     },
//     "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//   }
//
// @trigger channel
type ChannelTrigger struct {
	baseTrigger
	event *ChannelEvent
}

// NewChannelTrigger creates a new channel trigger with the passed in values
func NewChannelTrigger(env utils.Environment, flow *assets.FlowReference, contact *flows.Contact, event *ChannelEvent, params types.XValue, triggeredOn time.Time) *ChannelTrigger {
	return &ChannelTrigger{
		baseTrigger: newBaseTrigger(TypeChannel, env, flow, contact, params, triggeredOn),
		event:       event,
	}
}

// ToXJSON is called when this type is passed to @(json(...))
func (t *ChannelTrigger) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, t, "type", "params").ToXJSON(env)
}

var _ flows.Trigger = (*ChannelTrigger)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type channelTriggerEnvelope struct {
	baseTriggerEnvelope
	Event *ChannelEvent `json:"event" validate:"required,dive"`
}

// ReadChannelTrigger reads a channel trigger
func ReadChannelTrigger(session flows.Session, data json.RawMessage) (flows.Trigger, error) {
	e := &channelTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	t := &ChannelTrigger{
		event: e.Event,
	}

	if err := t.unmarshal(session, &e.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *ChannelTrigger) MarshalJSON() ([]byte, error) {
	e := &channelTriggerEnvelope{
		Event: t.event,
	}

	if err := t.marshal(&e.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	return json.Marshal(e)
}
