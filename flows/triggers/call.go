package triggers

import (
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
}

var _ flows.Trigger = (*Call)(nil)

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// CallBuilder is a builder for call type triggers
type CallBuilder struct {
	t *Call
}

func (b *Builder) CallReceived(event *events.CallReceived) *CallBuilder {
	return &CallBuilder{
		t: &Call{
			baseTrigger: newBaseTrigger(TypeCall, event, b.flow, false, nil),
		},
	}
}

func (b *Builder) CallMissed(event *events.CallMissed) *CallBuilder {
	return &CallBuilder{
		t: &Call{
			baseTrigger: newBaseTrigger(TypeCall, event, b.flow, false, nil),
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

func readCall(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &baseEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	t := &Call{}
	if err := t.unmarshal(sa, e, missing); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *Call) MarshalJSON() ([]byte, error) {
	e := &baseEnvelope{}

	if err := t.marshal(e); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
