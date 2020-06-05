package triggers

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/jsonx"
)

func init() {
	registerType(TypeChannel, readChannelTrigger)
}

// TypeChannel is the type for sessions triggered by channel events
const TypeChannel string = "channel"

// ChannelEventType is the type of event that occurred on the channel
type ChannelEventType string

// different channel event types
const (
	ChannelEventTypeNewConversation ChannelEventType = "new_conversation"
	ChannelEventTypeIncomingCall    ChannelEventType = "incoming_call"
)

// ChannelEvent describes the specific event on the channel that triggered the session
type ChannelEvent struct {
	Type    ChannelEventType         `json:"type" validate:"required"`
	Channel *assets.ChannelReference `json:"channel" validate:"required,dive"`
}

// NewChannelEvent creates a new channel event
func NewChannelEvent(typeName ChannelEventType, channel *assets.ChannelReference) *ChannelEvent {
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

// NewChannel creates a new channel trigger with the passed in values
func NewChannel(env envs.Environment, flow *assets.FlowReference, contact *flows.Contact, event *ChannelEvent, params *types.XObject) *ChannelTrigger {
	if params == nil {
		params = types.XObjectEmpty
	}

	return &ChannelTrigger{
		baseTrigger: newBaseTrigger(TypeChannel, env, flow, contact, nil, false, params),
		event:       event,
	}
}

// NewIncomingCall creates a new channel trigger with the passed in values
func NewIncomingCall(env envs.Environment, flow *assets.FlowReference, contact *flows.Contact, urn urns.URN, channel *assets.ChannelReference) *ChannelTrigger {
	event := NewChannelEvent(ChannelEventTypeIncomingCall, channel)
	connection := flows.NewConnection(channel, urn)

	return &ChannelTrigger{
		baseTrigger: newBaseTrigger(TypeChannel, env, flow, contact, connection, false, types.XObjectEmpty),
		event:       event,
	}
}

var _ flows.Trigger = (*ChannelTrigger)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type channelTriggerEnvelope struct {
	baseTriggerEnvelope
	Event *ChannelEvent `json:"event" validate:"required,dive"`
}

func readChannelTrigger(sessionAssets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &channelTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	t := &ChannelTrigger{
		event: e.Event,
	}

	if err := t.unmarshal(sessionAssets, &e.baseTriggerEnvelope, missing); err != nil {
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

	return jsonx.Marshal(e)
}
