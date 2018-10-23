package waits

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/triggers"
)

func init() {
	RegisterType(TypeMsg, func() flows.Wait { return &MsgWait{} })
}

// TypeMsg is the type of our message wait
const TypeMsg string = "msg"

// MediaType describes a type of media we're waiting for. See note on `media_hint` field below.
type MediaType string

// the different media types
const (
	MediaTypeImage    MediaType = "image"
	MediaTypeAudio    MediaType = "audio"
	MediaTypeVideo    MediaType = "video"
	MediaTypeLocation MediaType = "gps"
)

// MsgWait is a wait which waits for an incoming message (i.e. a msg_received event)
type MsgWait struct {
	baseWait

	// Waits can indicate to the caller what type of media the flow is expecting. In the case of flows of type
	// messaging_offline, this should be considered a requirement and the client should only reply with a message
	// containing an attachment of that type. In the case of other flow types this should be considered only a
	// hint to the channel, which may or may not support prompting the contact for media of that type.
	MediaHint_ MediaType `json:"media_hint,omitempty"`
}

// NewMsgWait creates a new message wait
func NewMsgWait(timeout *int, mediaHint MediaType) *MsgWait {
	return &MsgWait{
		baseWait:   newBaseWait(TypeMsg, timeout),
		MediaHint_: mediaHint,
	}
}

// Begin beings waiting at this wait
func (w *MsgWait) Begin(run flows.FlowRun, step flows.Step) bool {
	if !w.baseWait.Begin(run) {
		return false
	}

	// if we have a msg trigger and we're the first thing to happen... then we skip ourselves
	triggerHasMsg := run.Session().Trigger().Type() == triggers.TypeMsg

	if triggerHasMsg && len(run.Session().Runs()) == 1 && len(run.Path()) == 1 {
		return false
	}

	run.LogEvent(step, events.NewMsgWait(w.TimeoutOn_))
	return true
}

// End ends this wait or returns an error
func (w *MsgWait) End(resume flows.Resume) error {
	if resume.Type() == resumes.TypeMsg {
		return nil
	}

	return w.baseWait.End(resume)
}

var _ flows.Wait = (*MsgWait)(nil)
