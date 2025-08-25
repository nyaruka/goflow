package triggers

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeOptIn, readOptIn)
}

// TypeOptIn is the type for sessions triggered by optin/optout events
const TypeOptIn string = "optin"

// OptIn is used when a session was triggered by an optin or optout.
//
//	{
//	  "type": "optin",
//	  "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
//	  "event": {
//	    "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	    "type": "optin_started",
//	    "created_on": "2006-01-02T15:04:05Z",
//	    "optin": {
//	      "uuid": "248be71d-78e9-4d71-a6c4-9981d369e5cb",
//	      "name": "Joke Of The Day"
//	    },
//	    "channel": {
//	      "uuid": "4bb288a0-7fca-4da1-abe8-59a593aff648",
//	      "name": "Facebook"
//	    }
//	  },
//	  "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//	}
//
// @trigger optin
type OptIn struct {
	baseTrigger

	event flows.Event // optin_started or optin_stopped
	optIn *flows.OptIn
}

// Context for optin triggers includes the optin
func (t *OptIn) Context(env envs.Environment) map[string]types.XValue {
	c := t.context()
	c.optIn = flows.Context(env, t.optIn)
	return c.asMap()
}

var _ flows.Trigger = (*OptIn)(nil)

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// OptInBuilder is a builder for optin type triggers
type OptInBuilder struct {
	t *OptIn
}

// OptIn returns a optin trigger builder
func (b *Builder) OptIn(optIn *flows.OptIn, event flows.Event) *OptInBuilder {
	if event.Type() != events.TypeOptInStarted && event.Type() != events.TypeOptInStopped {
		panic("optin trigger event must be of type optin_started or optin_stopped")
	}

	return &OptInBuilder{
		t: &OptIn{
			baseTrigger: newBaseTrigger(TypeOptIn, b.flow, false, nil),
			event:       event,
			optIn:       optIn,
		},
	}
}

// Build builds the trigger
func (b *OptInBuilder) Build() *OptIn {
	return b.t
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type optInEnvelope struct {
	baseEnvelope

	Event json.RawMessage `json:"event" validate:"required"`
}

func readOptIn(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &optInEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	event, err := events.Read(e.Event)
	if err != nil {
		return nil, fmt.Errorf("error reading optin trigger event: %w", err)
	}

	var optInRef *assets.OptInReference

	switch typed := event.(type) {
	case *events.OptInStarted:
		optInRef = typed.OptIn
	case *events.OptInStopped:
		optInRef = typed.OptIn
	default:
		panic("optin trigger event must be of type optin_started or optin_stopped")
	}

	optIn := sa.OptIns().Get(optInRef.UUID)
	if optIn == nil {
		missing(optInRef, nil)
	}

	t := &OptIn{
		event: event,
		optIn: optIn,
	}

	if err := t.unmarshal(sa, &e.baseEnvelope, missing); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *OptIn) MarshalJSON() ([]byte, error) {
	me, err := json.Marshal(t.event)
	if err != nil {
		return nil, fmt.Errorf("error marshaling optin trigger event: %w", err)
	}

	e := &optInEnvelope{
		Event: me,
	}

	if err := t.marshal(&e.baseEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
