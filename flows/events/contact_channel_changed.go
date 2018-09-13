package events

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeContactChannelChanged, func() flows.Event { return &ContactChannelChangedEvent{} })
}

// TypeContactChannelChanged is the type of our set preferred channel event
const TypeContactChannelChanged string = "contact_channel_changed"

// ContactChannelChangedEvent events are created when a contact's preferred channel is changed.
//
//   {
//     "type": "contact_channel_changed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "channel": {"uuid": "67a3ac69-e5e0-4ef0-8423-eddf71a71472", "name": "Twilio"}
//   }
//
// @event contact_channel_changed
type ContactChannelChangedEvent struct {
	BaseEvent
	callerOrEngineEvent

	Channel *flows.ChannelReference `json:"channel" validate:"required"`
}

// NewContactChannelChangedEvent returns a new preferred channel event
func NewContactChannelChangedEvent(channel *flows.ChannelReference) *ContactChannelChangedEvent {
	return &ContactChannelChangedEvent{
		BaseEvent: NewBaseEvent(),
		Channel:   channel,
	}
}

// Type returns the type of this event
func (e *ContactChannelChangedEvent) Type() string { return TypeContactChannelChanged }

// Validate validates our event is valid and has all the assets it needs
func (e *ContactChannelChangedEvent) Validate(assets flows.SessionAssets) error {
	_, err := assets.Channels().Get(e.Channel.UUID)
	return err
}

// Apply applies this event to the given run
func (e *ContactChannelChangedEvent) Apply(run flows.FlowRun) error {
	if run.Contact() == nil {
		return fmt.Errorf("can't apply event in session without a contact")
	}

	channel, err := run.Session().Assets().Channels().Get(e.Channel.UUID)
	if err != nil {
		return err
	}

	run.Contact().UpdatePreferredChannel(channel)
	return nil
}
