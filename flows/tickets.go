package flows

import (
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// TicketUUID is the UUID of a ticket
type TicketUUID uuids.UUID

type baseTicket struct {
	UUID       TicketUUID `json:"uuid"`
	Subject    string     `json:"subject"`
	Body       string     `json:"body"`
	ExternalID string     `json:"external_id,omitempty"`
}

// Ticket is a ticket in a ticketing system
type Ticket struct {
	Ticketer *Ticketer
	baseTicket
}

// TicketReference is a ticket with a reference to the ticketer
type TicketReference struct {
	Ticketer *assets.TicketerReference `json:"ticketer"`
	baseTicket
}

// NewTicket creates a new ticket
func NewTicket(ticketer *Ticketer, subject, body, externalID string) *Ticket {
	return newTicket(TicketUUID(uuids.New()), ticketer, subject, body, externalID)
}

// creates a new ticket
func newTicket(uuid TicketUUID, ticketer *Ticketer, subject, body, externalID string) *Ticket {
	return &Ticket{
		baseTicket: baseTicket{
			UUID:       uuid,
			Subject:    subject,
			Body:       body,
			ExternalID: externalID,
		},
		Ticketer: ticketer,
	}
}

// Reference converts this ticket to a ticket reference suitable for marshaling
func (t *Ticket) Reference() *TicketReference {
	return &TicketReference{
		baseTicket: t.baseTicket,
		Ticketer:   t.Ticketer.Reference(),
	}
}

// ToXValue returns a representation of this object for use in expressions
//
//   uuid:text -> the UUID of the ticket
//   subject:text -> the subject of the ticket
//   body:text -> the body of the ticket
//
// @context ticket
func (t *Ticket) ToXValue(env envs.Environment) types.XValue {
	return types.NewXObject(map[string]types.XValue{
		"uuid":    types.NewXText(string(t.UUID)),
		"subject": types.NewXText(t.Subject),
		"body":    types.NewXText(t.Body),
	})
}

// TicketList defines a contact's list of tickets
type TicketList struct {
	tickets []*Ticket
}

// NewTicketList creates a new ticket list
func NewTicketList(sa SessionAssets, refs []*TicketReference, missing assets.MissingCallback) *TicketList {
	tickets := make([]*Ticket, 0, len(refs))

	for _, ref := range refs {
		ticketer := sa.Ticketers().Get(ref.Ticketer.UUID)
		if ticketer == nil {
			missing(ref.Ticketer, nil)
		} else {
			tickets = append(tickets, newTicket(ref.UUID, ticketer, ref.Subject, ref.Body, ref.ExternalID))
		}
	}
	return &TicketList{tickets: tickets}
}

// returns a clone of this ticket list
func (l *TicketList) clone() *TicketList {
	tickets := make([]*Ticket, len(l.tickets))
	copy(tickets, l.tickets)
	return &TicketList{tickets: tickets}
}

// returns this ticket list as a slice of ticket references
func (l *TicketList) references() []*TicketReference {
	refs := make([]*TicketReference, len(l.tickets))
	for i, ticket := range l.tickets {
		refs[i] = ticket.Reference()
	}
	return refs
}

// Adds adds the given ticket to this ticket list
func (l *TicketList) Add(ticket *Ticket) {
	l.tickets = append(l.tickets, ticket)
}

// Count returns the number of tickets
func (l *TicketList) Count() int {
	return len(l.tickets)
}

// ToXValue returns a representation of this object for use in expressions
func (l TicketList) ToXValue(env envs.Environment) types.XValue {
	array := make([]types.XValue, len(l.tickets))
	for i, ticket := range l.tickets {
		array[i] = ticket.ToXValue(env)
	}
	return types.NewXArray(array...)
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
