package resumes

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeWaitTimeout, readWaitTimeout)
}

// TypeWaitTimeout is the type for resuming a session when a wait has timed out
const TypeWaitTimeout string = "wait_timeout"

// WaitTimeout is used when a session is resumed because a wait has timed out
//
//	{
//	  "type": "wait_timeout",
//	  "resumed_on": "2000-01-01T00:00:00.000000000-00:00",
//	  "event": {
//	    "type": "wait_timed_out",
//	    "created_on": "2006-01-02T15:04:05Z",
//	    "run_uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a"
//	  }
//	}
//
// @resume wait_timeout
type WaitTimeout struct {
	baseResume

	event *events.WaitTimedOut
}

// NewWaitTimeout creates a new timeout resume with the passed in values
func NewWaitTimeout(event *events.WaitTimedOut) *WaitTimeout {
	return &WaitTimeout{
		baseResume: newBaseResume(TypeWaitTimeout),
		event:      event,
	}
}

// Event returns the event this resume is based on
func (r *WaitTimeout) Event() flows.Event { return r.event }

// Apply applies our state changes and saves any events to the run
func (r *WaitTimeout) Apply(run flows.Run, logEvent flows.EventCallback) {
	r.baseResume.Apply(run, logEvent)
}

var _ flows.Resume = (*WaitTimeout)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type waitTimeoutEnvelope struct {
	baseEnvelope

	Event *events.WaitTimedOut `json:"event"` // TODO make required
}

func readWaitTimeout(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Resume, error) {
	e := &waitTimeoutEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	r := &WaitTimeout{event: e.Event}

	if err := r.unmarshal(sa, &e.baseEnvelope, missing); err != nil {
		return nil, err
	}

	return r, nil
}

// MarshalJSON marshals this resume into JSON
func (r *WaitTimeout) MarshalJSON() ([]byte, error) {
	e := &waitTimeoutEnvelope{Event: r.event}

	if err := r.marshal(&e.baseEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
