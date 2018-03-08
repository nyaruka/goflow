package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
)

// TypeContactChannelChanged is the type of our set preferred channel event
const TypeContactChannelChanged string = "contact_channel_changed"

// ContactChannelChangedEvent events are created when a contact's preferred channel is changed.
//
// ```
//   {
//     "type": "contact_channel_changed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "channel": {"uuid": "67a3ac69-e5e0-4ef0-8423-eddf71a71472", "name": "Twilio"}
//   }
// ```
//
// @event contact_channel_changed
type ContactChannelChangedEvent struct {
	BaseEvent
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

// Apply applies this event to the given run
func (e *ContactChannelChangedEvent) Apply(run flows.FlowRun) error {
	channel, err := run.Session().Assets().GetChannel(e.Channel.UUID)
	if err != nil {
		return err
	}

	contactURNs := run.Contact().URNs()

	// change any tel URNs this contact has to use this channel (other URN types may be channel specific)
	if channel.SupportsScheme(urns.TelScheme) {
		for _, urn := range contactURNs {
			if urn.URN.Scheme() == urns.TelScheme {
				urn.SetChannel(channel)
			}
		}
	}

	// if our scheme isn't the highest priority
	if len(contactURNs) > 0 && !channel.SupportsScheme(contactURNs[0].Scheme()) {

		// find the highest priority supported by this channel
		var newPreferredURN *flows.ContactURN
		for _, urn := range contactURNs {
			if channel.SupportsScheme(urn.Scheme()) {
				newPreferredURN = urn
			}
		}

		// update the highest URN of the right scheme to be highest priority
		if newPreferredURN != nil {
			newURNs := make(flows.URNList, 1)
			newURNs[0] = newPreferredURN

			for _, urn := range contactURNs {
				if urn.URN != newPreferredURN.URN {
					newURNs = append(newURNs, urn)
				}
			}
		}
	}

	return nil
}
