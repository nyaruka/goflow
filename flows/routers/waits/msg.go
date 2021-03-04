package waits

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/routers/waits/hints"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

func init() {
	registerType(TypeMsg, readMsgWait, readActivatedMsgWait)
}

// TypeMsg is the type of our message wait
const TypeMsg string = "msg"

// MsgWait is a wait which waits for an incoming message (i.e. a msg_received event)
type MsgWait struct {
	baseWait

	// Message waits can indicate to the caller what kind of message the flow is expecting. In the case of flows of type
	// messaging_offline, this should be considered a requirement and the client should only reply with a message containing
	// an attachment of that type. In the case of other flow types this should be considered only a hint to the channel,
	// which may or may not support prompting the contact for media of that type.
	hint flows.Hint
}

// NewMsgWait creates a new message wait
func NewMsgWait(timeout *Timeout, hint flows.Hint) *MsgWait {
	return &MsgWait{
		baseWait: newBaseWait(TypeMsg, timeout),
		hint:     hint,
	}
}

// Hint returns the hint (optional)
func (w *MsgWait) Hint() flows.Hint { return w.hint }

// AllowedFlowTypes returns the flow types which this wait is allowed to occur in
func (w *MsgWait) AllowedFlowTypes() []flows.FlowType {
	return []flows.FlowType{flows.FlowTypeMessaging, flows.FlowTypeMessagingOffline, flows.FlowTypeVoice}
}

// Begin beings waiting at this wait
func (w *MsgWait) Begin(run flows.FlowRun, log flows.EventCallback) flows.ActivatedWait {
	var timeoutSeconds *int

	if w.timeout != nil {
		seconds := w.timeout.Seconds()
		timeoutSeconds = &seconds
	}

	// if we have a msg trigger and we're the first thing to happen... then we skip ourselves
	triggerHasMsg := run.Session().Trigger().Type() == triggers.TypeMsg

	if triggerHasMsg && len(run.Session().Runs()) == 1 && len(run.Path()) == 1 {
		return nil
	}

	log(events.NewMsgWait(timeoutSeconds, w.hint))

	return NewActivatedMsgWait(timeoutSeconds, w.hint)
}

// End ends this wait or returns an error
func (w *MsgWait) End(resume flows.Resume) error {
	switch resume.Type() {
	case resumes.TypeMsg, resumes.TypeRunExpiration:
		return nil
	case resumes.TypeWaitTimeout:
		if w.timeout == nil {
			return errors.Errorf("can't end with timeout as wait doesn't have a timeout")
		}
		return nil
	}
	return w.resumeTypeError(resume)
}

var _ flows.Wait = (*MsgWait)(nil)

type ActivatedMsgWait struct {
	baseActivatedWait

	hint flows.Hint
}

func NewActivatedMsgWait(timeoutSeconds *int, hint flows.Hint) *ActivatedMsgWait {
	return &ActivatedMsgWait{
		baseActivatedWait: baseActivatedWait{type_: TypeMsg, timeoutSeconds: timeoutSeconds},
		hint:              hint,
	}
}

// Hint returns the hint (optional)
func (w *ActivatedMsgWait) Hint() flows.Hint { return w.hint }

var _ flows.ActivatedWait = (*ActivatedMsgWait)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type msgWaitEnvelope struct {
	baseWaitEnvelope

	Hint json.RawMessage `json:"hint,omitempty"`
}

func readMsgWait(data json.RawMessage) (flows.Wait, error) {
	e := &msgWaitEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	w := &MsgWait{}

	var err error
	if e.Hint != nil {
		if w.hint, err = hints.ReadHint(e.Hint); err != nil {
			return nil, errors.Wrap(err, "unable to read hint")
		}
	}

	return w, w.unmarshal(&e.baseWaitEnvelope)
}

// MarshalJSON marshals this wait into JSON
func (w *MsgWait) MarshalJSON() ([]byte, error) {
	e := &msgWaitEnvelope{}

	if err := w.marshal(&e.baseWaitEnvelope); err != nil {
		return nil, err
	}

	var err error
	if w.hint != nil {
		if e.Hint, err = jsonx.Marshal(w.hint); err != nil {
			return nil, err
		}
	}

	return jsonx.Marshal(e)
}

type activatedMsgWaitEnvelope struct {
	baseActivatedWaitEnvelope

	Hint json.RawMessage `json:"hint,omitempty"`
}

func readActivatedMsgWait(data json.RawMessage) (flows.ActivatedWait, error) {
	e := &activatedMsgWaitEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	w := &ActivatedMsgWait{}

	var err error
	if e.Hint != nil {
		if w.hint, err = hints.ReadHint(e.Hint); err != nil {
			return nil, errors.Wrap(err, "unable to read hint")
		}
	}

	return w, w.unmarshal(&e.baseActivatedWaitEnvelope)
}

// MarshalJSON marshals this wait into JSON
func (w *ActivatedMsgWait) MarshalJSON() ([]byte, error) {
	e := &activatedMsgWaitEnvelope{}

	if err := w.marshal(&e.baseActivatedWaitEnvelope); err != nil {
		return nil, err
	}

	var err error
	if w.hint != nil {
		if e.Hint, err = jsonx.Marshal(w.hint); err != nil {
			return nil, err
		}
	}

	return jsonx.Marshal(e)
}
