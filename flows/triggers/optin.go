package triggers

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeOptIn, readOptInTrigger)
}

// TypeOptIn is the type for sessions triggered by optin/optout events
const TypeOptIn string = "optin"

// OptInEventType is the type of event that occurred on the optin
type OptInEventType string

// different optin event types
const (
	OptInEventTypeStarted OptInEventType = "started"
	OptInEventTypeStopped OptInEventType = "stopped"
)

// OptInEvent describes the specific event on the ticket that triggered the session
type OptInEvent struct {
	type_ OptInEventType
	optIn *flows.OptIn
}

// OptInTrigger is used when a session was triggered by an optin or optout.
//
//	{
//	  "type": "optin",
//	  "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
//	  "contact": {
//	    "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
//	    "name": "Bob",
//	    "created_on": "2018-01-01T12:00:00.000000Z"
//	  },
//	  "event": {
//	      "type": "started",
//	      "optin": {
//	          "uuid": "248be71d-78e9-4d71-a6c4-9981d369e5cb",
//	          "name": "Joke Of The Day"
//	      }
//	  },
//	  "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//	}
//
// @trigger optin
type OptInTrigger struct {
	baseTrigger
	event *OptInEvent
}

// Context for optin triggers includes the optin
func (t *OptInTrigger) Context(env envs.Environment) map[string]types.XValue {
	c := t.context()
	c.optIn = flows.Context(env, t.event.optIn)
	return c.asMap()
}

var _ flows.Trigger = (*OptInTrigger)(nil)

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// OptInBuilder is a builder for optin type triggers
type OptInBuilder struct {
	t *OptInTrigger
}

// OptIn returns a optin trigger builder
func (b *Builder) OptIn(optIn *flows.OptIn, eventType OptInEventType) *OptInBuilder {
	return &OptInBuilder{
		t: &OptInTrigger{
			baseTrigger: newBaseTrigger(TypeOptIn, b.environment, b.flow, b.contact, nil, false, nil),
			event:       &OptInEvent{type_: eventType, optIn: optIn},
		},
	}
}

// Build builds the trigger
func (b *OptInBuilder) Build() *OptInTrigger {
	return b.t
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type optInEventEnvelope struct {
	Type  OptInEventType         `json:"type"  validate:"required"`
	OptIn *assets.OptInReference `json:"optin" validate:"required"`
}

type optInTriggerEnvelope struct {
	baseTriggerEnvelope
	Event *optInEventEnvelope `json:"event" validate:"required"`
}

func readOptInTrigger(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &optInTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	t := &OptInTrigger{
		event: &OptInEvent{
			type_: e.Event.Type,
		},
	}

	t.event.optIn = sa.OptIns().Get(e.Event.OptIn.UUID)
	if t.event.optIn == nil {
		missing(e.Event.OptIn, nil)
	}

	if err := t.unmarshal(sa, &e.baseTriggerEnvelope, missing); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *OptInTrigger) MarshalJSON() ([]byte, error) {
	e := &optInTriggerEnvelope{
		Event: &optInEventEnvelope{
			Type:  t.event.type_,
			OptIn: t.event.optIn.Reference(),
		},
	}

	if err := t.marshal(&e.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
