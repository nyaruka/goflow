package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeCallMissed, func() flows.Event { return &CallMissed{} })
}

// TypeCallMissed is the type of our call missed event
const TypeCallMissed string = "call_missed"

// CallMissed events are missed when an incoming call is missed.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "call_missed",
//	  "created_on": "2019-01-02T15:04:05Z",
//	  "channel": {"uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf", "name": "Android"}
//	}
//
// @event call_missed
type CallMissed struct {
	BaseEvent

	Channel *assets.ChannelReference `json:"channel"`
}

// NewCallMissed returns a new call missed event
func NewCallMissed(channel *assets.ChannelReference) *CallMissed {
	return &CallMissed{
		BaseEvent: NewBaseEvent(TypeCallMissed),
		Channel:   channel,
	}
}

var _ flows.Event = (*CallMissed)(nil)
