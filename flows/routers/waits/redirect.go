package waits

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeRedirect, readRedirectWait, readActivatedRedirectWait)
}

// TypeRedirect is the type of our redirect wait
const TypeRedirect string = "redirect"

// RedirectWait is a wait which waits for a redirect to happen or fail to happen
type RedirectWait struct {
	baseWait
}

// NewRedirectWait creates a new redirect wait
func NewRedirectWait() *RedirectWait {
	return &RedirectWait{
		baseWait: newBaseWait(TypeRedirect, nil),
	}
}

// AllowedFlowTypes returns the flow types which this wait is allowed to occur in
func (w *RedirectWait) AllowedFlowTypes() []flows.FlowType {
	return []flows.FlowType{flows.FlowTypeMessaging, flows.FlowTypeMessagingOffline, flows.FlowTypeVoice}
}

// Begin beings waiting at this wait
func (w *RedirectWait) Begin(run flows.FlowRun, log flows.EventCallback) flows.ActivatedWait {
	log(events.NewRedirectWait())

	return NewActivatedRedirectWait()
}

// End ends this wait or returns an error
func (w *RedirectWait) End(resume flows.Resume) error {
	if resume.Type() == resumes.TypeRedirectEnd {
		return nil
	}
	return w.resumeTypeError(resume)
}

var _ flows.Wait = (*RedirectWait)(nil)

type ActivatedRedirectWait struct {
	baseActivatedWait
}

func NewActivatedRedirectWait() *ActivatedRedirectWait {
	return &ActivatedRedirectWait{
		baseActivatedWait: baseActivatedWait{type_: TypeRedirect},
	}
}

var _ flows.ActivatedWait = (*ActivatedRedirectWait)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readRedirectWait(data json.RawMessage) (flows.Wait, error) {
	e := &baseWaitEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	w := &RedirectWait{}

	return w, w.unmarshal(e)
}

// MarshalJSON marshals this wait into JSON
func (w *RedirectWait) MarshalJSON() ([]byte, error) {
	e := &baseWaitEnvelope{}

	if err := w.marshal(e); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}

func readActivatedRedirectWait(data json.RawMessage) (flows.ActivatedWait, error) {
	e := &baseActivatedWaitEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	w := &ActivatedRedirectWait{}

	return w, w.unmarshal(e)
}

// MarshalJSON marshals this wait into JSON
func (w *ActivatedRedirectWait) MarshalJSON() ([]byte, error) {
	e := &baseActivatedWaitEnvelope{}

	if err := w.marshal(e); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
