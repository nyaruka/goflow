package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
)

func init() {
	registerType(TypeTypingStarted, func() Event { return &TypingStarted{} })
}

// TypeTypingStarted is the type of our typing started event
const TypeTypingStarted string = "typing_started"

// TypingStarted events are created when the contact (direction of incoming) or a user (direction of outgoing)
// starts typing. The optional channel, urn and msg_external_id fields identify where the contact last wrote from -
// the source of incoming typing and the destination of outgoing typing.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "typing_started",
//	  "created_on": "2019-01-02T15:04:05Z",
//	  "direction": "incoming",
//	  "channel": {"uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf", "name": "Facebook"},
//	  "urn": "tel:+12065551212",
//	  "msg_external_id": "EX12345"
//	}
//
// @event typing_started
type TypingStarted struct {
	BaseEvent

	Direction     Direction                `json:"direction" validate:"required,direction"`
	Channel       *assets.ChannelReference `json:"channel,omitempty"`
	URN           urns.URN                 `json:"urn,omitempty" validate:"omitempty,urn"`
	MsgExternalID string                   `json:"msg_external_id,omitempty"`
}

// NewTypingStarted returns a new typing started event
func NewTypingStarted(direction Direction, channel *assets.ChannelReference, urn urns.URN, msgExternalID string) *TypingStarted {
	return &TypingStarted{
		BaseEvent:     NewBaseEvent(TypeTypingStarted),
		Direction:     direction,
		Channel:       channel,
		URN:           urn,
		MsgExternalID: msgExternalID,
	}
}

var _ Event = (*TypingStarted)(nil)
