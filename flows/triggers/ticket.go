package triggers

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
	"github.com/pkg/errors"
)

func init() {
	registerType(TypeTicket, readTicketTrigger)
}

// TypeTicket is the type for sessions triggered by ticket events
const TypeTicket string = "ticket"

// TicketEventType is the type of event that occurred on the ticket
type TicketEventType string

// different ticket event types
const (
	TicketEventTypeClosed TicketEventType = "closed"
)

// TicketEvent describes the specific event on the ticket that triggered the session
type TicketEvent struct {
	type_  TicketEventType
	ticket *flows.Ticket
}

// TicketTrigger is used when a session was triggered by a ticket event
//
//   {
//     "type": "ticket",
//     "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
//     "contact": {
//       "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
//       "name": "Bob",
//       "created_on": "2018-01-01T12:00:00.000000Z"
//     },
//     "event": {
//         "type": "closed",
//         "ticket": {
//             "uuid": "58e9b092-fe42-4173-876c-ff45a14a24fe",
//             "ticketer": {"uuid": "19dc6346-9623-4fe4-be80-538d493ecdf5", "name": "Support Tickets"},
//             "topic": {"uuid": "472a7a73-96cb-4736-b567-056d987cc5b4", "name": "Weather"},
//             "body": "Where are my shoes?"
//         }
//     },
//     "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//   }
//
// @trigger ticket
type TicketTrigger struct {
	baseTrigger
	event *TicketEvent
}

// Context for ticket triggers includes the ticket
func (t *TicketTrigger) Context(env envs.Environment) map[string]types.XValue {
	c := t.context()
	c.ticket = flows.Context(env, t.event.ticket)
	return c.asMap()
}

var _ flows.Trigger = (*TicketTrigger)(nil)

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// TicketBuilder is a builder for ticket type triggers
type TicketBuilder struct {
	t *TicketTrigger
}

// Ticket returns a ticket trigger builder
func (b *Builder) Ticket(ticket *flows.Ticket, eventType TicketEventType) *TicketBuilder {
	return &TicketBuilder{
		t: &TicketTrigger{
			baseTrigger: newBaseTrigger(TypeTicket, b.environment, b.flow, b.contact, nil, false, nil),
			event:       &TicketEvent{type_: eventType, ticket: ticket},
		},
	}
}

// Build builds the trigger
func (b *TicketBuilder) Build() *TicketTrigger {
	return b.t
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type ticketEventEnvelope struct {
	Type   TicketEventType `json:"type"   validate:"required"`
	Ticket json.RawMessage `json:"ticket" validate:"required"`
}

type ticketTriggerEnvelope struct {
	baseTriggerEnvelope
	Event ticketEventEnvelope `json:"event" validate:"required,dive"`
}

func readTicketTrigger(sa flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &ticketTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	t := &TicketTrigger{
		event: &TicketEvent{
			type_: e.Event.Type,
		},
	}

	var err error
	t.event.ticket, err = flows.ReadTicket(sa, e.Event.Ticket, missing)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read ticket")
	}

	if err := t.unmarshal(sa, &e.baseTriggerEnvelope, missing); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *TicketTrigger) MarshalJSON() ([]byte, error) {
	ticket, err := jsonx.Marshal(t.event.ticket)
	if err != nil {
		return nil, err
	}

	e := &ticketTriggerEnvelope{
		Event: ticketEventEnvelope{
			Type:   t.event.type_,
			Ticket: ticket,
		},
	}

	if err := t.marshal(&e.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
