package triggers

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeChannel, readChannel)
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
	Channel *assets.ChannelReference `json:"channel" validate:"required"`
}

// Channel is used when a session was triggered by a channel event
//
//	{
//	  "type": "channel",
//	  "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
//	  "event": {
//	      "type": "new_conversation",
//	      "channel": {"uuid": "58e9b092-fe42-4173-876c-ff45a14a24fe", "name": "Facebook"}
//	  },
//	  "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//	}
//
// @trigger channel
type Channel struct {
	baseTrigger
	event *ChannelEvent
}

var _ flows.Trigger = (*Channel)(nil)

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// ChannelBuilder is a builder for channel type triggers
type ChannelBuilder struct {
	t *Channel
}

// Channel returns a channel trigger builder
func (b *Builder) Channel(channel *assets.ChannelReference, eventType ChannelEventType) *ChannelBuilder {
	return &ChannelBuilder{
		t: &Channel{
			baseTrigger: newBaseTrigger(TypeChannel, b.flow, false, nil),
			event:       &ChannelEvent{Type: eventType, Channel: channel},
		},
	}
}

// WithParams sets the params for the trigger
func (b *ChannelBuilder) WithParams(params *types.XObject) *ChannelBuilder {
	b.t.params = params
	return b
}

// Build builds the trigger
func (b *ChannelBuilder) Build() *Channel {
	return b.t
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type channelEnvelope struct {
	baseEnvelope

	Event *ChannelEvent `json:"event" validate:"required"`
}

func readChannel(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &channelEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	t := &Channel{
		event: e.Event,
	}

	if err := t.unmarshal(sa, &e.baseEnvelope, missing); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *Channel) MarshalJSON() ([]byte, error) {
	e := &channelEnvelope{
		Event: t.event,
	}

	if err := t.marshal(&e.baseEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
