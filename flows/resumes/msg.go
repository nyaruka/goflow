package resumes

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/inputs"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeMsg, readMsg)
}

// TypeMsg is the type for resuming a session with a message
const TypeMsg string = "msg"

// Msg is used when a session is resumed with a new message from the contact
//
//	{
//	  "type": "msg",
//	  "resumed_on": "2000-01-01T00:00:00.000000000-00:00",
//	  "event": {
//	    "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	    "type": "msg_received",
//	    "created_on": "2006-01-02T15:04:05Z",
//	    "msg": {
//	      "channel": {"uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf", "name": "Twilio"},
//	      "urn": "tel:+12065551212",
//	      "text": "hi there",
//	      "attachments": ["https://s3.amazon.com/mybucket/attachment.jpg"]
//	    }
//	  }
//	}
//
// @resume msg
type Msg struct {
	baseResume

	event *events.MsgReceived
}

// NewMsg creates a new message resume with the passed in values
func NewMsg(event *events.MsgReceived) *Msg {
	return &Msg{
		baseResume: newBaseResume(TypeMsg),
		event:      event,
	}
}

// Event returns the event this resume is based on
func (r *Msg) Event() flows.Event { return r.event }

func (r *Msg) Input(sa flows.SessionAssets) flows.Input {
	return inputs.NewMsg(sa, r.event)
}

var _ flows.Resume = (*Msg)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type msgEnvelope struct {
	baseEnvelope

	Event *events.MsgReceived `json:"event" validate:"required"`
}

func readMsg(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Resume, error) {
	e := &msgEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	r := &Msg{event: e.Event}

	if err := r.unmarshal(sa, &e.baseEnvelope, missing); err != nil {
		return nil, err
	}

	return r, nil
}

// MarshalJSON marshals this resume into JSON
func (r *Msg) MarshalJSON() ([]byte, error) {
	e := &msgEnvelope{
		Event: r.event,
	}

	if err := r.marshal(&e.baseEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
