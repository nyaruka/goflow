package flows

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils/uuids"
)

// TicketUUID is the UUID of a ticket
type TicketUUID uuids.UUID

// Ticket is a ticket in a ticketing system
type Ticket struct {
	UUID       TicketUUID                `json:"uuid"`
	Ticketer   *assets.TicketerReference `json:"ticketer"`
	Subject    string                    `json:"subject"`
	Body       string                    `json:"body"`
	ExternalID string                    `json:"external_id,omitempty"`
}

// NewTicket creates a new ticket
func NewTicket(uuid TicketUUID, ticketer *assets.TicketerReference, subject, body, externalID string) *Ticket {
	return &Ticket{
		UUID:       uuid,
		Ticketer:   ticketer,
		Subject:    subject,
		Body:       body,
		ExternalID: externalID,
	}
}

// Ticketer represents a ticket issuing system.
type Ticketer struct {
	assets.Ticketer
}

// NewTicketer returns a new classifier object from the given classifier asset
func NewTicketer(asset assets.Ticketer) *Ticketer {
	return &Ticketer{Ticketer: asset}
}

// Asset returns the underlying asset
func (t *Ticketer) Asset() assets.Ticketer { return t.Ticketer }

// Reference returns a reference to this classifier
func (t *Ticketer) Reference() *assets.TicketerReference {
	return assets.NewTicketerReference(t.UUID(), t.Name())
}

// TicketerAssets provides access to all ticketer assets
type TicketerAssets struct {
	byUUID map[assets.TicketerUUID]*Ticketer
}

// NewTicketerAssets creates a new set of ticketer assets
func NewTicketerAssets(ticketers []assets.Ticketer) *TicketerAssets {
	s := &TicketerAssets{
		byUUID: make(map[assets.TicketerUUID]*Ticketer, len(ticketers)),
	}
	for _, asset := range ticketers {
		s.byUUID[asset.UUID()] = NewTicketer(asset)
	}
	return s
}

// Get returns the ticketer with the given UUID
func (s *TicketerAssets) Get(uuid assets.TicketerUUID) *Ticketer {
	return s.byUUID[uuid]
}
