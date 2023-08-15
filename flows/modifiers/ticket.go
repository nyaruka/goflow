package modifiers

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeTicket, readTicketModifier)
}

// TypeTicket is the type of our ticket modifier
const TypeTicket string = "ticket"

// TicketModifier opens a ticket for the contact
type TicketModifier struct {
	baseModifier

	ticketer *flows.Ticketer
	topic    *flows.Topic
	body     string
	assignee *flows.User
}

// NewTicket creates a new ticket modifier
func NewTicket(ticketer *flows.Ticketer, topic *flows.Topic, body string, assignee *flows.User) *TicketModifier {
	return &TicketModifier{
		baseModifier: newBaseModifier(TypeTicket),
		ticketer:     ticketer,
		topic:        topic,
		body:         body,
		assignee:     assignee,
	}
}

// Apply applies this modification to the given contact
func (m *TicketModifier) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) bool {
	// if there's already an open ticket, nothing to do
	if contact.Ticket() != nil {
		return false
	}

	httpLogger := &flows.HTTPLogger{}

	// try to get a ticket service for this ticketer
	svc, err := eng.Services().Ticket(m.ticketer)
	if err != nil {
		log(events.NewError(err))
		return false
	}

	ticket, err := svc.Open(env, contact, m.topic, m.body, m.assignee, httpLogger.Log)
	if err != nil {
		log(events.NewError(err))
	}
	if len(httpLogger.Logs) > 0 {
		log(events.NewTicketerCalled(m.ticketer.Reference(), httpLogger.Logs))
	}
	if ticket != nil {
		log(events.NewTicketOpened(ticket))

		contact.SetTicket(ticket)
		return true
	}

	return false
}

func (m *TicketModifier) Ticketer() *flows.Ticketer {
	return m.ticketer
}

var _ flows.Modifier = (*TicketModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type ticketModifierEnvelope struct {
	utils.TypedEnvelope

	Ticketer *assets.TicketerReference `json:"ticketer" validate:"required"`
	Topic    *assets.TopicReference    `json:"topic" validate:"required"`
	Body     string                    `json:"body"`
	Assignee *assets.UserReference     `json:"assignee"`
}

func readTicketModifier(assets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &ticketModifierEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	ticketer := assets.Ticketers().Get(e.Ticketer.UUID)
	if ticketer == nil {
		missing(e.Ticketer, nil)
		return nil, ErrNoModifier // can't proceed without a ticketer
	}

	topic := assets.Topics().Get(e.Topic.UUID)
	if topic == nil {
		missing(e.Topic, nil)
		return nil, ErrNoModifier // can't proceed without a topic
	}

	var assignee *flows.User
	if e.Assignee != nil {
		assignee = assets.Users().Get(e.Assignee.Email)
		if assignee == nil {
			missing(e.Assignee, nil)
		}
	}

	return NewTicket(ticketer, topic, e.Body, assignee), nil
}

func (m *TicketModifier) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&ticketModifierEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Ticketer:      m.ticketer.Reference(),
		Topic:         m.topic.Reference(),
		Body:          m.body,
		Assignee:      m.assignee.Reference(),
	})
}
