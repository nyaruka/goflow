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
//     "channel": {"uuid": "67a3ac69-e5e0-4ef0-8423-eddf71a71472", "name": "Twilio"}
//   }
// ```
//
// @event set_preferred_channel
type PreferredChannelEvent struct {
	BaseEvent
	Channel *flows.ChannelReference `json:"channel" validate:"required"`
}

// NewPreferredChannel returns a new preferred channel event
func NewPreferredChannel(channel *flows.ChannelReference) *PreferredChannelEvent {
	return &PreferredChannelEvent{
		BaseEvent: NewBaseEvent(),
		Channel:   channel,
	}
}

// Type returns the type of this event
func (e *PreferredChannelEvent) Type() string { return TypePreferredChannel }

// Apply applies this event to the given run
func (e *PreferredChannelEvent) Apply(run flows.FlowRun) error {
	channel, err := run.Session().Assets().GetChannel(e.Channel.UUID)
	if err != nil {
		return err
	}

	run.Contact().SetChannel(channel)
	return nil
}
