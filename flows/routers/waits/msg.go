package waits

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/routers/waits/hints"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeMsg, readMsg)
}

// TypeMsg is the type of our message wait
const TypeMsg string = "msg"

// Msg is a wait which waits for an incoming message (i.e. a msg_received event)
type Msg struct {
	baseWait

	// Message waits can indicate to the caller what kind of message the flow is expecting. In the case of flows of type
	// messaging_offline, this should be considered a requirement and the client should only reply with a message containing
	// an attachment of that type. In the case of other flow types this should be considered only a hint to the channel,
	// which may or may not support prompting the contact for media of that type.
	hint flows.Hint
}

// NewMsg creates a new message wait
func NewMsg(timeout *Timeout, hint flows.Hint) *Msg {
	return &Msg{
		baseWait: newBaseWait(TypeMsg, timeout),
		hint:     hint,
	}
}

// Hint returns the hint (optional)
func (w *Msg) Hint() flows.Hint { return w.hint }

// AllowedFlowTypes returns the flow types which this wait is allowed to occur in
func (w *Msg) AllowedFlowTypes() []flows.FlowType {
	return []flows.FlowType{flows.FlowTypeMessaging, flows.FlowTypeMessagingOffline, flows.FlowTypeVoice}
}

// Begin beings waiting at this wait
func (w *Msg) Begin(run flows.Run, log flows.EventCallback) bool {
	// if we have a msg trigger and we're the first thing to happen... then we skip ourselves
	triggerHasMsg := run.Session().Trigger().Type() == triggers.TypeMsg

	if triggerHasMsg && len(run.Session().Runs()) == 1 && len(run.Path()) == 1 {
		return false
	}

	var timeoutSeconds *int
	if w.timeout != nil {
		seconds := w.timeout.Seconds()
		timeoutSeconds = &seconds
	}

	log(events.NewMsgWait(timeoutSeconds, w.expiresOn(run), w.hint))

	return true
}

// Accept returns whether this wait accepts the given resume
func (w *Msg) Accepts(resume flows.Resume) bool {
	switch resume.Type() {
	case resumes.TypeMsg, resumes.TypeWaitExpiration, resumes.TypeRunExpiration:
		return true
	case resumes.TypeWaitTimeout:
		return w.timeout != nil
	}
	return false
}

var _ flows.Wait = (*Msg)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type msgEnvelope struct {
	baseEnvelope

	Hint json.RawMessage `json:"hint,omitempty"`
}

func readMsg(data json.RawMessage) (flows.Wait, error) {
	e := &msgEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	w := &Msg{}

	var err error
	if e.Hint != nil {
		if w.hint, err = hints.Read(e.Hint); err != nil {
			return nil, fmt.Errorf("unable to read hint: %w", err)
		}
	}

	return w, w.unmarshal(&e.baseEnvelope)
}

// MarshalJSON marshals this wait into JSON
func (w *Msg) MarshalJSON() ([]byte, error) {
	e := &msgEnvelope{}

	if err := w.marshal(&e.baseEnvelope); err != nil {
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
