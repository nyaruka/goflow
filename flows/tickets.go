package flows

import (
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// TicketUUID is the UUID of a ticket
type TicketUUID uuids.UUID

type TicketStatus string

const (
	// TicketStatusOpen is the status of an open ticket
	TicketStatusOpen TicketStatus = "open"
	// TicketStatusClosed is the status of a closed ticket
	TicketStatusClosed TicketStatus = "closed"
)

// NewTicketUUID generates a new UUID for a ticket
func NewTicketUUID() TicketUUID { return TicketUUID(uuids.NewV7()) }

// Ticket is a ticket in a ticketing system
type Ticket struct {
	uuid           TicketUUID
	status         TicketStatus
	topic          *Topic
	assignee       *User
	lastActivityOn time.Time
}

// NewTicket creates a new ticket
func NewTicket(uuid TicketUUID, status TicketStatus, topic *Topic, assignee *User, lastActivityOn time.Time) *Ticket {
	return &Ticket{
		uuid:           uuid,
		status:         status,
		topic:          topic,
		assignee:       assignee,
		lastActivityOn: lastActivityOn,
	}
}

// OpenTicket creates a new ticket. Used by ticketing services to open a new ticket.
func OpenTicket(topic *Topic, assignee *User) *Ticket {
	return NewTicket(NewTicketUUID(), TicketStatusOpen, topic, assignee, dates.Now())
}

func (t *Ticket) UUID() TicketUUID { return t.uuid }

func (t *Ticket) Status() TicketStatus          { return t.status }
func (t *Ticket) SetStatus(status TicketStatus) { t.status = status }

func (t *Ticket) Topic() *Topic         { return t.topic }
func (t *Ticket) SetTopic(topic *Topic) { t.topic = topic }

func (t *Ticket) Assignee() *User        { return t.assignee }
func (t *Ticket) SetAssignee(user *User) { t.assignee = user }

func (t *Ticket) LastActivityOn() time.Time { return t.lastActivityOn }
func (t *Ticket) RecordActivity()           { t.lastActivityOn = dates.Now() }

// Context returns the properties available in expressions
//
//	uuid:text -> the UUID of the ticket
//	topic:any -> the topic of the ticket
//	assignee:any -> the assignee of the ticket
//
// @context ticket
func (t *Ticket) Context(env envs.Environment) map[string]types.XValue {
	return map[string]types.XValue{
		"uuid":     types.NewXText(string(t.uuid)),
		"topic":    Context(env, t.topic),
		"assignee": Context(env, t.assignee),
	}
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type TicketEnvelope struct {
	UUID           TicketUUID             `json:"uuid"                   validate:"required,uuid"`
	Status         TicketStatus           `json:"status"`
	Topic          *assets.TopicReference `json:"topic"                  validate:"omitempty"`
	Assignee       *assets.UserReference  `json:"assignee,omitempty"     validate:"omitempty"`
	LastActivityOn time.Time              `json:"last_activity_on"`
}

// Unmarshal unmarshals a ticket from the passed in envelope. If the topic or assigned user can't
// be found in the assets, we report the missing asset and return ticket without those.
func (e *TicketEnvelope) Unmarshal(sa SessionAssets, missing assets.MissingCallback) *Ticket {
	if e.Status == "" {
		e.Status = TicketStatusOpen
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
		assignee = sa.Users().Get(e.Assignee.UUID)
		if assignee == nil {
			missing(e.Assignee, nil)
		}
	}

	return &Ticket{uuid: e.UUID, status: e.Status, topic: topic, assignee: assignee, lastActivityOn: e.LastActivityOn}
}

// Marshal marshals a ticket into an envelope.
func (t *Ticket) Marshal() *TicketEnvelope {
	var topicRef *assets.TopicReference
	if t.topic != nil {
		topicRef = t.topic.Reference()
	}

	var assigneeRef *assets.UserReference
	if t.assignee != nil {
		assigneeRef = t.assignee.Reference()
	}

	return &TicketEnvelope{
		UUID:           t.uuid,
		Status:         t.status,
		Topic:          topicRef,
		Assignee:       assigneeRef,
		LastActivityOn: t.lastActivityOn,
	}
}
