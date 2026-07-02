package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
)

// CallUUID is the UUID of a call
type CallUUID uuids.UUID

// NewCallUUID generates a new UUID for a call
func NewCallUUID() CallUUID { return CallUUID(uuids.NewV7()) }

// CallEnvelope is the serialized form of a call
type CallEnvelope struct {
	UUID    CallUUID                 `json:"uuid"    validate:"required,uuid"`
	Channel *assets.ChannelReference `json:"channel" validate:"required"`
	URN     urns.URN                 `json:"urn"     validate:"required,urn"`
}
