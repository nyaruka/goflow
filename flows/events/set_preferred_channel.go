package events

import "github.com/nyaruka/goflow/flows"

// TypePreferredChannel is the type of our set preferred channel event
const TypePreferredChannel string = "set_preferred_channel"

// PreferredChannelEvent events are created when a contact's preferred channel is changed.
//
// ```
//   {
//     "type": "set_preferred_channel",
//     "created_on": "2006-01-02T15:04:05Z",
//     "channel_uuid": "67a3ac69-e5e0-4ef0-8423-eddf71a71472",
//     "channel_name": "Twilio"
//   }
// ```
//
// @event set_preferred_channel
type PreferredChannelEvent struct {
	BaseEvent
	ChannelUUID flows.ChannelUUID `json:"channel_uuid" validate:"required"`
	ChannelName string            `json:"channel_name"`
}

// NewPreferredChannel returns a new preferred channel event
func NewPreferredChannel(channelUUID flows.ChannelUUID, channelName string) *PreferredChannelEvent {
	return &PreferredChannelEvent{
		BaseEvent:   NewBaseEvent(),
		ChannelUUID: channelUUID,
		ChannelName: channelName,
	}
}

// Type returns the type of this event
func (e *PreferredChannelEvent) Type() string { return TypePreferredChannel }

// Apply applies this event to the given run
func (e *PreferredChannelEvent) Apply(run flows.FlowRun) {
	channel, err := run.Session().Assets().GetChannel(e.ChannelUUID)
	if err != nil {
		return
	}

	run.Contact().SetChannel(channel)
}
