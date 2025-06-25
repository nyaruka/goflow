package triggers

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeTicket, readTicket)
}

// TypeTicket is the type for sessions triggered by ticket events
const TypeTicket string = "ticket"

// Ticket is used when a session was triggered by a ticket event (for now only closed events).
//
//	{
//	  "type": "ticket",
//	  "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
//	  "event": {
//	      "type": "ticket_closed",
//	      "created_on": "2006-01-02T15:04:05Z",
//	      "ticket": {
//	          "uuid": "58e9b092-fe42-4173-876c-ff45a14a24fe",
//	          "topic": {"uuid": "472a7a73-96cb-4736-b567-056d987cc5b4", "name": "Weather"}
//	      }
//	  },
//	  "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//	}
//
// @trigger ticket
type Ticket struct {
	baseTrigger

	event  *events.TicketClosed
	ticket *flows.Ticket
}

func (t *Ticket) Event() flows.Event { return t.event }

// Context for ticket triggers includes the ticket
func (t *Ticket) Context(env envs.Environment) map[string]types.XValue {
	c := t.context()
	c.ticket = flows.Context(env, t.ticket)
	return c.asMap()
}

var _ flows.Trigger = (*Ticket)(nil)

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// TicketBuilder is a builder for ticket type triggers
type TicketBuilder struct {
	t *Ticket
}

// Ticket returns a ticket trigger builder
func (b *Builder) Ticket(ticket *flows.Ticket, event *events.TicketClosed) *TicketBuilder {
	return &TicketBuilder{
		t: &Ticket{
			baseTrigger: newBaseTrigger(TypeTicket, b.flow, false, nil),
			event:       event,
			ticket:      ticket,
		},
	}
}

// Build builds the trigger
func (b *TicketBuilder) Build() *Ticket {
	return b.t
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type ticketEnvelope struct {
	baseEnvelope

	Event *events.TicketClosed `json:"event" validate:"required"`
}

func readTicket(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &ticketEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	// older sessions will have events that aren't really events so fix 'em
	if e.Event.Type() == "closed" {
		e.Event.Type_ = events.TypeTicketClosed
		e.Event.CreatedOn_ = e.TriggeredOn // ensure we have a created on time
	}

	t := &Ticket{
		event:  e.Event,
		ticket: e.Event.Ticket.Unmarshal(sa, missing),
	}

	if err := t.unmarshal(sa, &e.baseEnvelope, missing); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *Ticket) MarshalJSON() ([]byte, error) {
	e := &ticketEnvelope{
		Event: t.event,
	}

	if err := t.marshal(&e.baseEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
