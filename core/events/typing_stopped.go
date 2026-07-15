package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
)

func init() {
	registerType(TypeTypingStopped, func() Event { return &TypingStopped{} })
}

// TypeTypingStopped is the type of our typing stopped event
const TypeTypingStopped string = "typing_stopped"

// TypingStopped events are created when the contact (direction of incoming) or a user (direction of outgoing)
// stops typing. The optional channel, urn and msg_external_id fields identify where the contact last wrote from -
// the source of incoming typing and the destination of outgoing typing.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "typing_stopped",
//	  "created_on": "2019-01-02T15:04:05Z",
//	  "direction": "incoming",
//	  "channel": {"uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf", "name": "Facebook"},
//	  "urn": "tel:+12065551212",
//	  "msg_external_id": "EX12345"
//	}
//
// @event typing_stopped
type TypingStopped struct {
	BaseEvent

	Direction     Direction                `json:"direction" validate:"required,direction"`
	Channel       *assets.ChannelReference `json:"channel,omitempty"`
	URN           urns.URN                 `json:"urn,omitempty" validate:"omitempty,urn"`
	MsgExternalID string                   `json:"msg_external_id,omitempty"`
}

// NewTypingStopped returns a new typing stopped event
func NewTypingStopped(direction Direction, channel *assets.ChannelReference, urn urns.URN, msgExternalID string) *TypingStopped {
	return &TypingStopped{
		BaseEvent:     NewBaseEvent(TypeTypingStopped),
		Direction:     direction,
		Channel:       channel,
		URN:           urn,
		MsgExternalID: msgExternalID,
	}
}

var _ Event = (*TypingStopped)(nil)
