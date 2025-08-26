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

func (b *Builder) OptInStarted(event *events.OptInStarted, optIn *flows.OptIn) *OptInBuilder {
	return &OptInBuilder{
		t: &OptIn{
			baseTrigger: newBaseTrigger(TypeOptIn, event, b.flow, false, nil),
			optIn:       optIn,
		},
	}
}

func (b *Builder) OptInStopped(event *events.OptInStopped, optIn *flows.OptIn) *OptInBuilder {
	return &OptInBuilder{
		t: &OptIn{
			baseTrigger: newBaseTrigger(TypeOptIn, event, b.flow, false, nil),
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

func readOptIn(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &baseEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	t := &OptIn{}

	if err := t.unmarshal(sa, e, missing); err != nil {
		return nil, err
	}

	var optInRef *assets.OptInReference

	switch typed := t.event.(type) {
	case *events.OptInStarted:
		optInRef = typed.OptIn
	case *events.OptInStopped:
		optInRef = typed.OptIn
	default:
		panic("optin trigger event must be of type optin_started or optin_stopped")
	}

	t.optIn = sa.OptIns().Get(optInRef.UUID)
	if t.optIn == nil {
		missing(optInRef, nil)
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *OptIn) MarshalJSON() ([]byte, error) {
	e := &baseEnvelope{}

	if err := t.marshal(e); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
