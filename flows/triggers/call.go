package triggers

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeCall, readCall)
}

// TypeCall is the type for call triggered sessions
const TypeCall string = "call"

// Call is used to trigger a session for a new incoming or missed call
//
//	{
//	  "type": "call",
//	  "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
//	  "event": {
//	    "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	    "type": "call_received",
//	    "created_on": "2006-01-02T15:04:05Z",
//	    "call": {
//	      "uuid": "0198ce92-ff2f-7b07-b158-b21ab168ebba",
//	      "channel": {"uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf", "name": "Twilio"},
//	      "urn": "tel:+12065551212"
//	    }
//	  },
//	  "triggered_on": "2000-01-01T00:00:00.000000000-00:00"
//	}
//
// @trigger call
type Call struct {
	baseTrigger

	event flows.Event // call_received or call_missed
}

func (t *Call) Event() flows.Event { return t.event }

var _ flows.Trigger = (*Call)(nil)

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// CallBuilder is a builder for call type triggers
type CallBuilder struct {
	t *Call
}

// Call returns a call trigger builder
func (b *Builder) Call(e flows.Event) *CallBuilder {
	if e.Type() != events.TypeCallReceived && e.Type() != events.TypeCallMissed {
		panic("call trigger event must be of type call_received or call_missed")
	}

	return &CallBuilder{
		t: &Call{
			baseTrigger: newBaseTrigger(TypeCall, b.flow, false, nil),
			event:       e,
		},
	}
}

// Build builds the trigger
func (b *CallBuilder) Build() *Call {
	return b.t
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type callEnvelope struct {
	baseEnvelope

	Event json.RawMessage `json:"event" validate:"required"`
}

func readCall(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &callEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	event, err := events.Read(e.Event)
	if err != nil {
		return nil, fmt.Errorf("error reading call trigger event: %w", err)
	}

	t := &Call{event: event}

	if err := t.unmarshal(sa, &e.baseEnvelope, missing); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *Call) MarshalJSON() ([]byte, error) {
	me, err := json.Marshal(t.event)
	if err != nil {
		return nil, fmt.Errorf("error marshaling optin trigger event: %w", err)
	}

	e := &callEnvelope{
		Event: me,
	}

	if err := t.marshal(&e.baseEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
