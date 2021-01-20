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
	registerType(TypeRedirectEnd, readRedirectEndResume)
}

// TypeRedirectEnd is the type for redirect end resumes
const TypeRedirectEnd string = "redirect_end"

// RedirectEndResume is used when a session is resumed because a redirect has ended.
//
//   {
//     "type": "redirect_end",
//     "resumed_on": "2021-01-20T12:18:30Z",
//     "response": "answered"
//   }
//
// @resume redirect_end
type RedirectEndResume struct {
	baseResume
	response flows.RedirectResponse
}

// NewRedirectEnd creates a new redirect end resume
func NewRedirectEnd(env envs.Environment, contact *flows.Contact, response flows.RedirectResponse) *RedirectEndResume {
	return &RedirectEndResume{
		baseResume: newBaseResume(TypeRedirectEnd, env, contact),
		response:   response,
	}
}

// Apply applies our state changes and saves any events to the run
func (r *RedirectEndResume) Apply(run flows.FlowRun, logEvent flows.EventCallback) {
	logEvent(events.NewRedirectEnded(r.response))

	r.baseResume.Apply(run, logEvent)
}

var _ flows.Resume = (*RedirectEndResume)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type redirectEndResumeEnvelope struct {
	baseResumeEnvelope
	Response flows.RedirectResponse `json:"response" validate:"required"`
}

func readRedirectEndResume(sessionAssets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Resume, error) {
	e := &redirectEndResumeEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	r := &RedirectEndResume{
		response: e.Response,
	}

	if err := r.unmarshal(sessionAssets, &e.baseResumeEnvelope, missing); err != nil {
		return nil, err
	}

	return r, nil
}

// MarshalJSON marshals this resume into JSON
func (r *RedirectEndResume) MarshalJSON() ([]byte, error) {
	e := &redirectEndResumeEnvelope{
		Response: r.response,
	}

	if err := r.marshal(&e.baseResumeEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
