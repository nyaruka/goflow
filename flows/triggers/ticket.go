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
//	    "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	    "type": "ticket_closed",
//	    "created_on": "2006-01-02T15:04:05Z",
//	    "ticket_uuid": "019905d4-5f7b-71b8-bcb8-6a68de2d91d2"
//	  },
//	  "ticket": {
//	    "uuid": "019905d4-5f7b-71b8-bcb8-6a68de2d91d2",
//	    "topic": {"uuid": "472a7a73-96cb-4736-b567-056d987cc5b4", "name": "Weather"}
//	  },
//	  "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//	}
//
// @trigger ticket
type Ticket struct {
	baseTrigger

	ticket *flows.Ticket
}

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

func (b *Builder) TicketClosed(event *events.TicketClosed, ticket *flows.Ticket) *TicketBuilder {
	return &TicketBuilder{
		t: &Ticket{
			baseTrigger: newBaseTrigger(TypeTicket, event, b.flow, false, nil),
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

	Ticket *flows.TicketEnvelope `json:"ticket"`
}

func readTicket(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &ticketEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	t := &Ticket{}

	if e.Ticket != nil {
		t.ticket = e.Ticket.Unmarshal(sa, missing)
	}

	if err := t.unmarshal(sa, &e.baseEnvelope, missing); err != nil {
		return nil, err
	}

	if event, ok := t.event.(*events.TicketClosed); ok {
		// TODO remove once there are no more triggers passing the ticket by event
		if t.ticket == nil {
			t.ticket = event.Ticket.Unmarshal(sa, missing)
		}
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *Ticket) MarshalJSON() ([]byte, error) {
	e := &ticketEnvelope{
		Ticket: t.ticket.Marshal(),
	}

	if err := t.marshal(&e.baseEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
