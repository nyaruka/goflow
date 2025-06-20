package resumes

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
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
//	  "resumed_on": "2000-01-01T00:00:00.000000000-00:00"
//	}
//
// @resume wait_timeout
type WaitTimeoutResume struct {
	baseResume
}

// NewWaitTimeout creates a new timeout resume with the passed in values
func NewWaitTimeout(env envs.Environment) *WaitTimeoutResume {
	return &WaitTimeoutResume{
		baseResume: newBaseResume(TypeWaitTimeout, env),
	}
}

// Apply applies our state changes and saves any events to the run
func (r *WaitTimeoutResume) Apply(run flows.Run, logEvent flows.EventCallback) {
	logEvent(events.NewWaitTimedOut())

	r.baseResume.Apply(run, logEvent)
}

var _ flows.Resume = (*WaitTimeoutResume)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readWaitTimeoutResume(sessionAssets flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Resume, error) {
	e := &baseResumeEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	r := &WaitTimeoutResume{}

	if err := r.unmarshal(sessionAssets, e, missing); err != nil {
		return nil, err
	}

	return r, nil
}

// MarshalJSON marshals this resume into JSON
func (r *WaitTimeoutResume) MarshalJSON() ([]byte, error) {
	e := &baseResumeEnvelope{}

	if err := r.marshal(e); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
