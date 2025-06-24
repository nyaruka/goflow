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
	registerType(TypeDial, readDialResume)
}

// TypeDial is the type for dial resumes
const TypeDial string = "dial"

// DialResume is used when a session is resumed after a number was dialed.
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
type DialResume struct {
	baseResume

	event *events.DialEndedEvent
}

// NewDial creates a new dial resume
func NewDial(event *events.DialEndedEvent) *DialResume {
	return &DialResume{
		baseResume: newBaseResume(TypeDial),
		event:      event,
	}
}

// Event returns the event this resume is based on
func (r *DialResume) Event() flows.Event { return r.event }

// Context for dial resumes additionally exposes the dial object
func (r *DialResume) Context(env envs.Environment) map[string]types.XValue {
	c := r.context()
	c.dial = flows.Context(env, r.event.Dial)
	return c.asMap()
}

var _ flows.Resume = (*DialResume)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type dialResumeEnvelope struct {
	baseResumeEnvelope

	Event *events.DialEndedEvent `json:"event"`          // TODO make required
	Dial  *flows.Dial            `json:"dial,omitempty"` // used by older sessions
}

func readDialResume(sessionAssets flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Resume, error) {
	e := &dialResumeEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	r := &DialResume{event: e.Event}

	// older resumes will have dial instead of event so convert that into an event
	if e.Dial != nil {
		r.event = &events.DialEndedEvent{
			BaseEvent: events.BaseEvent{Type_: events.TypeDialEnded, CreatedOn_: e.baseResumeEnvelope.ResumedOn},
			Dial:      e.Dial,
		}
	}

	if err := r.unmarshal(sessionAssets, &e.baseResumeEnvelope, missing); err != nil {
		return nil, err
	}

	return r, nil
}

// MarshalJSON marshals this resume into JSON
func (r *DialResume) MarshalJSON() ([]byte, error) {
	e := &dialResumeEnvelope{Event: r.event}

	if err := r.marshal(&e.baseResumeEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
