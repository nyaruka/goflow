package modifiers

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeTicketOpen, readTicketOpen)
	registerType("ticket", readTicketOpen) // deprecated
}

// TypeTicketOpen is the type of our ticket modifier
const TypeTicketOpen string = "ticket_open"

// TicketOpen opens a ticket for the contact
type TicketOpen struct {
	baseModifier

	topic    *flows.Topic
	assignee *flows.User
	note     string
}

// NewTicketOpen creates a new ticket open modifier
func NewTicketOpen(topic *flows.Topic, assignee *flows.User, note string) *TicketOpen {
	return &TicketOpen{
		baseModifier: newBaseModifier(TypeTicketOpen),
		topic:        topic,
		assignee:     assignee,
		note:         note,
	}
}

// Apply applies this modification to the given contact
func (m *TicketOpen) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) bool {
	// if there's already an open ticket, nothing to do
	if contact.Tickets().OpenCount() > 0 {
		return false
	}

	ticket := flows.OpenTicket(m.topic, m.assignee)
	log(events.NewTicketOpened(ticket, m.note))

	contact.Tickets().Add(ticket)
	return true
}

var _ flows.Modifier = (*TicketOpen)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type ticketOpenEnvelope struct {
	utils.TypedEnvelope

	Topic    *assets.TopicReference `json:"topic" validate:"required"`
	Assignee *assets.UserReference  `json:"assignee"`
	Note     string                 `json:"note"`
}

func readTicketOpen(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &ticketOpenEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	topic := sa.Topics().Get(e.Topic.UUID)
	if topic == nil {
		missing(e.Topic, nil)
		return nil, ErrNoModifier // can't proceed without a topic
	}

	var assignee *flows.User
	if e.Assignee != nil {
		assignee = sa.Users().Get(e.Assignee.UUID)
		if assignee == nil {
			missing(e.Assignee, nil)
		}
	}

	return NewTicketOpen(topic, assignee, e.Note), nil
}

func (m *TicketOpen) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&ticketOpenEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Topic:         m.topic.Reference(),
		Assignee:      m.assignee.Reference(),
		Note:          m.note,
	})
}
