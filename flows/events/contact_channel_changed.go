package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeContactChannelChanged, func() flows.Event { return &ContactChannelChangedEvent{} })
}

// TypeContactChannelChanged is the type of our set preferred channel event
const TypeContactChannelChanged string = "contact_channel_changed"

// ContactChannelChangedEvent events are created when the preferred channel of the contact has been changed.
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

	Channel *assets.ChannelReference `json:"channel" validate:"required"`
}

// NewContactChannelChangedEvent returns a new preferred channel event
func NewContactChannelChangedEvent(channel *assets.ChannelReference) *ContactChannelChangedEvent {
	return &ContactChannelChangedEvent{
		BaseEvent: NewBaseEvent(),
		Channel:   channel,
	}
}

// Type returns the type of this event
func (e *ContactChannelChangedEvent) Type() string { return TypeContactChannelChanged }
