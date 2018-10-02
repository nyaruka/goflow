package resumes

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeWaitTimeout, ReadWaitTimeoutResume)
}

// TypeWaitTimeout is the type for resuming a session when a wait has timed out
const TypeWaitTimeout string = "wait_timeout"

// WaitTimeoutResume is used when a session is resumed because a wait has timed out
//
//   {
//     "type": "wait_timeout",
//     "contact": {
//       "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
//       "name": "Bob",
//       "language": "fra",
//       "fields": {"gender": {"text": "Male"}},
//       "groups": []
//     },
//     "resumed_on": "2000-01-01T00:00:00.000000000-00:00"
//   }
//
// @resume wait_timeout
type WaitTimeoutResume struct {
	baseResume
}

// NewWaitTimeoutResume creates a new timeout resume with the passed in values
func NewWaitTimeoutResume(env utils.Environment, contact *flows.Contact) *WaitTimeoutResume {
	return &WaitTimeoutResume{
		baseResume: newBaseResume(env, contact),
	}
}

// Type returns the type of this resume
func (r *WaitTimeoutResume) Type() string { return TypeWaitTimeout }

// Apply applies our state changes and saves any events to the run
func (r *WaitTimeoutResume) Apply(run flows.FlowRun, step flows.Step) error {
	// clear the last input on the run
	run.SetInput(nil)
	run.LogEvent(step, events.NewWaitTimedOutEvent())

	return r.baseResume.Apply(run, step)
}

var _ flows.Resume = (*WaitTimeoutResume)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadWaitTimeoutResume reads a timeout resume
func ReadWaitTimeoutResume(session flows.Session, data json.RawMessage) (flows.Resume, error) {
	resume := &WaitTimeoutResume{}
	e := baseResumeEnvelope{}
	if err := utils.UnmarshalAndValidate(data, &e); err != nil {
		return nil, err
	}

	if err := unmarshalBaseResume(session, &resume.baseResume, &e); err != nil {
		return nil, err
	}

	return resume, nil
}
