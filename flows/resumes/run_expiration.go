package resumes

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeRunExpiration, ReadRunExpirationResume)
}

// TypeRunExpiration is the type for resuming a session when a run has expired
const TypeRunExpiration string = "run_expiration"

// RunExpirationResume is used when a session is resumed because the waiting run has expired
//
//   {
//     "type": "run_expiration",
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
// @resume run_expiration
type RunExpirationResume struct {
	baseResume
}

// NewRunExpirationResume creates a new run expired resume with the passed in values
func NewRunExpirationResume(env utils.Environment, contact *flows.Contact) *RunExpirationResume {
	return &RunExpirationResume{
		baseResume: newBaseResume(env, contact),
	}
}

// Type returns the type of this resume
func (r *RunExpirationResume) Type() string { return TypeRunExpiration }

// Apply applies our state changes and saves any events to the run
func (r *RunExpirationResume) Apply(run flows.FlowRun, step flows.Step) error {
	run.Exit(flows.RunStatusExpired)
	run.AddEvent(step, nil, events.NewRunExpiredEvent(run))

	return r.baseResume.Apply(run, step)
}

var _ flows.Resume = (*RunExpirationResume)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadRunExpirationResume reads a run expired resume
func ReadRunExpirationResume(session flows.Session, data json.RawMessage) (flows.Resume, error) {
	resume := &RunExpirationResume{}
	e := baseResumeEnvelope{}
	if err := utils.UnmarshalAndValidate(data, &e); err != nil {
		return nil, err
	}

	if err := unmarshalBaseResume(session, &resume.baseResume, &e); err != nil {
		return nil, err
	}

	return resume, nil
}
