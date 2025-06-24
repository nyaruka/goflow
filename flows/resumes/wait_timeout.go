package resumes

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeWaitTimeout, readWaitTimeoutResume)
}

// TypeWaitTimeout is the type for resuming a session when a wait has timed out
const TypeWaitTimeout string = "wait_timeout"

// WaitTimeoutResume is used when a session is resumed because a wait has timed out
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
type WaitTimeoutResume struct {
	baseResume

	event *events.WaitTimedOutEvent
}

// NewWaitTimeout creates a new timeout resume with the passed in values
func NewWaitTimeout(event *events.WaitTimedOutEvent) *WaitTimeoutResume {
	return &WaitTimeoutResume{
		baseResume: newBaseResume(TypeWaitTimeout),
		event:      event,
	}
}

// Event returns the event this resume is based on
func (r *WaitTimeoutResume) Event() flows.Event { return r.event }

// Apply applies our state changes and saves any events to the run
func (r *WaitTimeoutResume) Apply(run flows.Run, logEvent flows.EventCallback) {
	r.baseResume.Apply(run, logEvent)
}

var _ flows.Resume = (*WaitTimeoutResume)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type waitTimeoutResumeEnvelope struct {
	baseResumeEnvelope

	Event *events.WaitTimedOutEvent `json:"event"` // TODO make required
}

func readWaitTimeoutResume(sessionAssets flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Resume, error) {
	e := &waitTimeoutResumeEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	r := &WaitTimeoutResume{event: e.Event}

	if err := r.unmarshal(sessionAssets, &e.baseResumeEnvelope, missing); err != nil {
		return nil, err
	}

	return r, nil
}

// MarshalJSON marshals this resume into JSON
func (r *WaitTimeoutResume) MarshalJSON() ([]byte, error) {
	e := &waitTimeoutResumeEnvelope{Event: r.event}

	if err := r.marshal(&e.baseResumeEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
