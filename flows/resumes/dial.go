package resumes

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
	registerType(TypeDial, readDial)
}

// TypeDial is the type for dial resumes
const TypeDial string = "dial"

// Dial is used when a session is resumed after a number was dialed.
//
//	{
//	  "type": "dial",
//	  "resumed_on": "2021-01-20T12:18:30Z",
//	  "event": {
//	    "type": "dial_ended",
//	    "created_on": "2019-01-02T15:04:05Z",
//	    "dial": {
//	      "status": "answered",
//	      "duration": 10
//	    }
//	  }
//	}
//
// @resume dial
type Dial struct {
	baseResume

	event *events.DialEnded
}

// NewDial creates a new dial resume
func NewDial(event *events.DialEnded) *Dial {
	return &Dial{
		baseResume: newBaseResume(TypeDial),
		event:      event,
	}
}

// Event returns the event this resume is based on
func (r *Dial) Event() flows.Event { return r.event }

// Context for dial resumes additionally exposes the dial object
func (r *Dial) Context(env envs.Environment) map[string]types.XValue {
	c := r.context()
	c.dial = flows.Context(env, r.event.Dial)
	return c.asMap()
}

var _ flows.Resume = (*Dial)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type dialEnvelope struct {
	baseEnvelope

	Event *events.DialEnded `json:"event" validate:"required"`
}

func readDial(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Resume, error) {
	e := &dialEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	r := &Dial{event: e.Event}

	if err := r.unmarshal(sa, &e.baseEnvelope, missing); err != nil {
		return nil, err
	}

	return r, nil
}

// MarshalJSON marshals this resume into JSON
func (r *Dial) MarshalJSON() ([]byte, error) {
	e := &dialEnvelope{Event: r.event}

	if err := r.marshal(&e.baseEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
