package triggers

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
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
	ChannelEventTypeIncomingCall    ChannelEventType = "incoming_call"
	ChannelEventTypeMissedCall      ChannelEventType = "missed_call"
	ChannelEventTypeNewConversation ChannelEventType = "new_conversation"
	ChannelEventTypeReferral        ChannelEventType = "referral"
)

// ChannelEvent describes the specific event on the channel that triggered the session
type ChannelEvent struct {
	Type    ChannelEventType         `json:"type" validate:"required"`
	Channel *assets.ChannelReference `json:"channel" validate:"required,dive"`
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

var _ flows.Trigger = (*ChannelTrigger)(nil)

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// ChannelBuilder is a builder for channel type triggers
type ChannelBuilder struct {
	t *ChannelTrigger
}

// Channel returns a channel trigger builder
func (b *Builder) Channel(channel *assets.ChannelReference, eventType ChannelEventType) *ChannelBuilder {
	return &ChannelBuilder{
		t: &ChannelTrigger{
			baseTrigger: newBaseTrigger(TypeChannel, b.environment, b.flow, b.contact, nil, false, nil),
			event:       &ChannelEvent{Type: eventType, Channel: channel},
		},
	}
}

// WithConnection sets the channel connection for the trigger
func (b *ChannelBuilder) WithConnection(urn urns.URN) *ChannelBuilder {
	b.t.connection = flows.NewConnection(b.t.event.Channel, urn)
	return b
}

// WithParams sets the params for the trigger
func (b *ChannelBuilder) WithParams(params *types.XObject) *ChannelBuilder {
	b.t.params = params
	return b
}

// Build builds the trigger
func (b *ChannelBuilder) Build() *ChannelTrigger {
	return b.t
}

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
