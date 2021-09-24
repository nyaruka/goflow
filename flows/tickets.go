package flows

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// TicketUUID is the UUID of a ticket
type TicketUUID uuids.UUID

// Ticket is a ticket in a ticketing system
type Ticket struct {
	uuid       TicketUUID
	ticketer   *Ticketer
	topic      *Topic
	body       string
	externalID string
	assignee   *User
}

// NewTicket creates a new ticket
func NewTicket(uuid TicketUUID, ticketer *Ticketer, topic *Topic, body, externalID string, assignee *User) *Ticket {
	return &Ticket{
		uuid:       uuid,
		ticketer:   ticketer,
		topic:      topic,
		body:       body,
		externalID: externalID,
		assignee:   assignee,
	}
}

// OpenTicket creates a new ticket. Used by ticketing services to open a new ticket.
func OpenTicket(ticketer *Ticketer, topic *Topic, body string, assignee *User) *Ticket {
	return NewTicket(TicketUUID(uuids.New()), ticketer, topic, body, "", assignee)
}

func (t *Ticket) UUID() TicketUUID        { return t.uuid }
func (t *Ticket) Ticketer() *Ticketer     { return t.ticketer }
func (t *Ticket) Topic() *Topic           { return t.topic }
func (t *Ticket) Body() string            { return t.body }
func (t *Ticket) ExternalID() string      { return t.externalID }
func (t *Ticket) SetExternalID(id string) { t.externalID = id }
func (t *Ticket) Assignee() *User         { return t.assignee }

// Context returns the properties available in expressions
//
//   uuid:text -> the UUID of the ticket
//   subject:text -> the subject of the ticket
//   body:text -> the body of the ticket
//
// @context ticket
func (t *Ticket) Context(env envs.Environment) map[string]types.XValue {
	return map[string]types.XValue{
		"uuid":     types.NewXText(string(t.uuid)),
		"topic":    Context(env, t.topic),
		"body":     types.NewXText(t.body),
		"assignee": Context(env, t.assignee),
	}
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type ticketEnvelope struct {
	UUID       TicketUUID                `json:"uuid"                   validate:"required,uuid4"`
	Ticketer   *assets.TicketerReference `json:"ticketer"               validate:"omitempty,dive"`
	Topic      *assets.TopicReference    `json:"topic"                  validate:"omitempty,dive"`
	Body       string                    `json:"body"`
	ExternalID string                    `json:"external_id,omitempty"`
	Assignee   *assets.UserReference     `json:"assignee,omitempty"     validate:"omitempty,dive"`
}

// ReadTicket decodes a contact from the passed in JSON. If the ticketer or assigned user can't
// be found in the assets, we report the missing asset and return ticket without those.
func ReadTicket(sa SessionAssets, data []byte, missing assets.MissingCallback) (*Ticket, error) {
	e := &ticketEnvelope{}

	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	ticketer := sa.Ticketers().Get(e.Ticketer.UUID)
	if ticketer == nil {
		missing(e.Ticketer, nil)
	}

	var topic *Topic
	if e.Topic != nil {
		topic = sa.Topics().Get(e.Topic.UUID)
		if topic == nil {
			missing(e.Topic, nil)
		}
	}

	var assignee *User
	if e.Assignee != nil {
		assignee = sa.Users().Get(e.Assignee.Email)
		if assignee == nil {
			missing(e.Assignee, nil)
		}
	}

	return &Ticket{
		uuid:       e.UUID,
		ticketer:   ticketer,
		topic:      topic,
		body:       e.Body,
		externalID: e.ExternalID,
		assignee:   assignee,
	}, nil
}

// MarshalJSON marshals this ticket into JSON
func (t *Ticket) MarshalJSON() ([]byte, error) {
	var ticketerRef *assets.TicketerReference
	if t.ticketer != nil {
		ticketerRef = t.ticketer.Reference()
	}

	var topicRef *assets.TopicReference
	if t.topic != nil {
		topicRef = t.topic.Reference()
	}

	var assigneeRef *assets.UserReference
	if t.assignee != nil {
		assigneeRef = t.assignee.Reference()
	}

	return jsonx.Marshal(&ticketEnvelope{
		UUID:       t.uuid,
		Ticketer:   ticketerRef,
		Topic:      topicRef,
		Body:       t.body,
		ExternalID: t.externalID,
		Assignee:   assigneeRef,
	})
}

// TicketList defines a contact's list of tickets
type TicketList struct {
	tickets []*Ticket
}

// NewTicketList creates a new ticket list
func NewTicketList(tickets []*Ticket) *TicketList {
	return &TicketList{tickets: tickets}
}

// returns a clone of this ticket list
func (l *TicketList) clone() *TicketList {
	tickets := make([]*Ticket, len(l.tickets))
	copy(tickets, l.tickets)
	return &TicketList{tickets: tickets}
}

// Adds adds the given ticket to this ticket list
func (l *TicketList) Add(ticket *Ticket) {
	l.tickets = append(l.tickets, ticket)
}

// All returns all tickets in this ticket list
func (l *TicketList) All() []*Ticket {
	return l.tickets
}

// Count returns the number of tickets
func (l *TicketList) Count() int {
	return len(l.tickets)
}

// ToXValue returns a representation of this object for use in expressions
func (l TicketList) ToXValue(env envs.Environment) types.XValue {
	array := make([]types.XValue, len(l.tickets))
	for i, ticket := range l.tickets {
		array[i] = Context(env, ticket)
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
	if t == nil {
		return nil
	}
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
