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
	registerType(TypeTicket, readTicket)
}

// TypeTicket is the type of our ticket modifier
const TypeTicket string = "ticket"

// Ticket opens a ticket for the contact
type Ticket struct {
	baseModifier

	topic    *flows.Topic
	assignee *flows.User
	note     string
}

// NewTicket creates a new ticket modifier
func NewTicket(topic *flows.Topic, assignee *flows.User, note string) *Ticket {
	return &Ticket{
		baseModifier: newBaseModifier(TypeTicket),
		topic:        topic,
		assignee:     assignee,
		note:         note,
	}
}

// Apply applies this modification to the given contact
func (m *Ticket) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) bool {
	// if there's already an open ticket, nothing to do
	if contact.Ticket() != nil {
		return false
	}

	ticket := flows.OpenTicket(m.topic, m.assignee)
	log(events.NewTicketOpened(ticket, m.note))

	contact.SetTicket(ticket)
	return true
}

var _ flows.Modifier = (*Ticket)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type ticketEnvelope struct {
	utils.TypedEnvelope

	Topic    *assets.TopicReference `json:"topic" validate:"required"`
	Assignee *assets.UserReference  `json:"assignee"`
	Note     string                 `json:"note"`
}

func readTicket(assets flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &ticketEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	topic := assets.Topics().Get(e.Topic.UUID)
	if topic == nil {
		missing(e.Topic, nil)
		return nil, ErrNoModifier // can't proceed without a topic
	}

	var assignee *flows.User
	if e.Assignee != nil {
		assignee = assets.Users().Get(e.Assignee.UUID)
		if assignee == nil {
			missing(e.Assignee, nil)
		}
	}

	return NewTicket(topic, assignee, e.Note), nil
}

func (m *Ticket) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&ticketEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Topic:         m.topic.Reference(),
		Assignee:      m.assignee.Reference(),
		Note:          m.note,
	})
}
