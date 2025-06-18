package triggers

import (
	"encoding/json"
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeOptIn, readOptInTrigger)
}

// TypeOptIn is the type for sessions triggered by optin/optout events
const TypeOptIn string = "optin"

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
	event flows.Event // optin_started or optin_stopped
	optIn *flows.OptIn
}

// Context for optin triggers includes the optin
func (t *OptInTrigger) Context(env envs.Environment) map[string]types.XValue {
	c := t.context()
	c.optIn = flows.Context(env, t.optIn)
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
func (b *Builder) OptIn(optIn *flows.OptIn, event flows.Event) *OptInBuilder {
	if event.Type() != events.TypeOptInStarted && event.Type() != events.TypeOptInStopped {
		panic("optin trigger event must be of type optin_started or optin_stopped")
	}

	return &OptInBuilder{
		t: &OptInTrigger{
			baseTrigger: newBaseTrigger(TypeOptIn, b.environment, b.flow, b.contact, nil, false, nil),
			event:       event,
			optIn:       optIn,
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

type optInTriggerEnvelope struct {
	baseTriggerEnvelope
	Event json.RawMessage `json:"event" validate:"required"`
}

func readOptInTrigger(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &optInTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	// TODO remove this once all triggers are using real events
	evtType, err := jsonparser.GetString(e.Event, "type")
	if err != nil {
		return nil, fmt.Errorf("error reading type from optin trigger event: %w", err)
	}
	if evtType == "started" {
		e.Event, _ = jsonparser.Set(e.Event, []byte(`"optin_started"`), "type")
		e.Event, _ = jsonparser.Set(e.Event, jsonx.MustMarshal(e.TriggeredOn), "created_on")
	} else if evtType == "stopped" {
		e.Event, _ = jsonparser.Set(e.Event, []byte(`"optin_stopped"`), "type")
		e.Event, _ = jsonparser.Set(e.Event, jsonx.MustMarshal(e.TriggeredOn), "created_on")
	}

	event, err := events.ReadEvent(e.Event)
	if err != nil {
		return nil, fmt.Errorf("error reading optin trigger event: %w", err)
	}

	var optInRef *assets.OptInReference

	switch typed := event.(type) {
	case *events.OptInStartedEvent:
		optInRef = typed.OptIn
	case *events.OptInStoppedEvent:
		optInRef = typed.OptIn
	default:
		panic("optin trigger event must be of type optin_started or optin_stopped")
	}

	optIn := sa.OptIns().Get(optInRef.UUID)
	if optIn == nil {
		missing(optInRef, nil)
	}

	t := &OptInTrigger{
		event: event,
		optIn: optIn,
	}

	if err := t.unmarshal(sa, &e.baseTriggerEnvelope, missing); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *OptInTrigger) MarshalJSON() ([]byte, error) {
	me, err := json.Marshal(t.event)
	if err != nil {
		return nil, fmt.Errorf("error marshaling optin trigger event: %w", err)
	}

	e := &optInTriggerEnvelope{
		Event: me,
	}

	if err := t.marshal(&e.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
