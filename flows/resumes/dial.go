package resumes

import (
	"encoding/json"

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
//   {
//     "type": "dial",
//     "resumed_on": "2021-01-20T12:18:30Z",
//     "dial": {
//       "status": "answered",
//       "duration": 15
//     }
//   }
//
// @resume dial
type DialResume struct {
	baseResume

	dial *flows.Dial
}

// NewDial creates a new dial resume
func NewDial(env envs.Environment, contact *flows.Contact, dial *flows.Dial) *DialResume {
	return &DialResume{
		baseResume: newBaseResume(TypeDial, env, contact),
		dial:       dial,
	}
}

// Apply applies our state changes and saves any events to the run
func (r *DialResume) Apply(run flows.FlowRun, logEvent flows.EventCallback) {
	logEvent(events.NewDialEnded(r.dial))

	r.baseResume.Apply(run, logEvent)
}

// Context for dial resumes additionally exposes the dial object
func (r *DialResume) Context(env envs.Environment) map[string]types.XValue {
	c := r.context()
	c.dial = flows.Context(env, r.dial)
	return c.asMap()
}

var _ flows.Resume = (*DialResume)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type dialResumeEnvelope struct {
	baseResumeEnvelope

	Dial *flows.Dial `json:"dial" validate:"required,dive"`
}

func readDialResume(sessionAssets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Resume, error) {
	e := &dialResumeEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	r := &DialResume{dial: e.Dial}

	if err := r.unmarshal(sessionAssets, &e.baseResumeEnvelope, missing); err != nil {
		return nil, err
	}

	return r, nil
}

// MarshalJSON marshals this resume into JSON
func (r *DialResume) MarshalJSON() ([]byte, error) {
	e := &dialResumeEnvelope{Dial: r.dial}

	if err := r.marshal(&e.baseResumeEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
