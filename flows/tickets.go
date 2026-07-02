package flows

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/events"
	"github.com/nyaruka/goflow/excellent/types"
)

// Ticket is a ticket in a ticketing system
type Ticket struct {
	uuid     events.TicketUUID
	status   events.TicketStatus
	topic    *Topic
	assignee *User
}

// NewTicket creates a new ticket
func NewTicket(uuid events.TicketUUID, status events.TicketStatus, topic *Topic, assignee *User) *Ticket {
	return &Ticket{
		uuid:     uuid,
		status:   status,
		topic:    topic,
		assignee: assignee,
	}
}

// OpenTicket creates a new ticket. Used by ticketing services to open a new ticket.
func OpenTicket(topic *Topic, assignee *User) *Ticket {
	return NewTicket(events.NewTicketUUID(), events.TicketStatusOpen, topic, assignee)
}

func (t *Ticket) UUID() events.TicketUUID { return t.uuid }

func (t *Ticket) Status() events.TicketStatus          { return t.status }
func (t *Ticket) SetStatus(status events.TicketStatus) { t.status = status }

func (t *Ticket) Topic() *Topic         { return t.topic }
func (t *Ticket) SetTopic(topic *Topic) { t.topic = topic }

func (t *Ticket) Assignee() *User        { return t.assignee }
func (t *Ticket) SetAssignee(user *User) { t.assignee = user }

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

type TicketList struct {
	all []*Ticket
}

func NewTicketList(tickets []*Ticket) *TicketList {
	return &TicketList{all: tickets}
}

func (l *TicketList) All() []*Ticket {
	return l.all
}

func (l *TicketList) Add(t *Ticket) {
	l.all = append(l.all, t)
}

func (l *TicketList) Find(uuid events.TicketUUID) *Ticket {
	for _, t := range l.all {
		if t.uuid == uuid {
			return t
		}
	}
	return nil
}

func (l *TicketList) LastOpen() *Ticket {
	for i := len(l.all) - 1; i >= 0; i-- {
		if l.all[i].status == events.TicketStatusOpen {
			return l.all[i]
		}
	}
	return nil
}

func (l *TicketList) Open() *TicketList {
	open := make([]*Ticket, 0, len(l.all))
	for _, t := range l.all {
		if t.status == events.TicketStatusOpen {
			open = append(open, t)
		}
	}
	return NewTicketList(open)
}

func (l *TicketList) Count() int {
	return len(l.all)
}

// ToXValue returns a representation of this object for use in expressions
func (l *TicketList) ToXValue(env envs.Environment) types.XValue {
	array := make([]types.XValue, len(l.all))
	for i, val := range l.all {
		array[i] = Context(env, val)
	}
	return types.NewXArray(array...)
}

func (l *TicketList) Marshal() []*events.TicketEnvelope {
	envelopes := make([]*events.TicketEnvelope, len(l.all))
	for i, t := range l.all {
		envelopes[i] = t.Marshal()
	}
	return envelopes
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadTicket reads a ticket from the passed in envelope. If the topic or assigned user can't
// be found in the assets, we report the missing asset and return ticket without those.
func ReadTicket(sa SessionAssets, e *events.TicketEnvelope, missing assets.MissingCallback) *Ticket {
	if e.Status == "" {
		e.Status = events.TicketStatusOpen
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

	return &Ticket{uuid: e.UUID, status: e.Status, topic: topic, assignee: assignee}
}

// Marshal marshals a ticket into an envelope.
func (t *Ticket) Marshal() *events.TicketEnvelope {
	var topicRef *assets.TopicReference
	if t.topic != nil {
		topicRef = t.topic.Reference()
	}

	var assigneeRef *assets.UserReference
	if t.assignee != nil {
		assigneeRef = t.assignee.Reference()
	}

	return &events.TicketEnvelope{
		UUID:     t.uuid,
		Status:   t.status,
		Topic:    topicRef,
		Assignee: assigneeRef,
	}
}
