package resumes

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeRunExpiration, readRunExpirationResume)
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
//       "created_on": "2018-01-01T12:00:00.000000Z",
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

// NewRunExpiration creates a new run expired resume with the passed in values
func NewRunExpiration(env envs.Environment, contact *flows.Contact) *RunExpirationResume {
	return &RunExpirationResume{
		baseResume: newBaseResume(TypeRunExpiration, env, contact),
	}
}

// Apply applies our state changes and saves any events to the run
func (r *RunExpirationResume) Apply(run flows.FlowRun, logEvent flows.EventCallback) {
	run.Exit(flows.RunStatusExpired)

	logEvent(events.NewRunExpired(run))

	r.baseResume.Apply(run, logEvent)
}

var _ flows.Resume = (*RunExpirationResume)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readRunExpirationResume(sessionAssets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Resume, error) {
	e := &baseResumeEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	r := &RunExpirationResume{}

	if err := r.unmarshal(sessionAssets, e, missing); err != nil {
		return nil, err
	}

	return r, nil
}

// MarshalJSON marshals this resume into JSON
func (r *RunExpirationResume) MarshalJSON() ([]byte, error) {
	e := &baseResumeEnvelope{}

	if err := r.marshal(e); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
