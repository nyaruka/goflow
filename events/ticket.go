package events

import (
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
)

// TicketUUID is the UUID of a ticket
type TicketUUID uuids.UUID

// NewTicketUUID generates a new UUID for a ticket
func NewTicketUUID() TicketUUID { return TicketUUID(uuids.NewV7()) }

type TicketStatus string

const (
	// TicketStatusOpen is the status of an open ticket
	TicketStatusOpen TicketStatus = "open"
	// TicketStatusClosed is the status of a closed ticket
	TicketStatusClosed TicketStatus = "closed"
)

// TicketEnvelope is the serialized form of a ticket
type TicketEnvelope struct {
	UUID     TicketUUID             `json:"uuid"                   validate:"required,uuid"`
	Status   TicketStatus           `json:"status"` // TODO validate:"required,eq=open|eq=closed"`
	Topic    *assets.TopicReference `json:"topic"                  validate:"omitempty"`
	Assignee *assets.UserReference  `json:"assignee,omitempty"     validate:"omitempty"`
}
