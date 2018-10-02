package resumes

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/inputs"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeMsg, ReadMsgResume)
}

// TypeMsg is the type for resuming a session with a message
const TypeMsg string = "msg"

// MsgResume is used when a session is resumed with a new message from the contact
//
//   {
//     "type": "msg",
//     "contact": {
//       "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
//       "name": "Bob",
//       "language": "fra",
//       "fields": {"gender": {"text": "Male"}},
//       "groups": []
//     },
//     "msg": {
//       "uuid": "2d611e17-fb22-457f-b802-b8f7ec5cda5b",
//       "channel": {"uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf", "name": "Twilio"},
//       "urn": "tel:+12065551212",
//       "text": "hi there",
//       "attachments": ["https://s3.amazon.com/mybucket/attachment.jpg"]
//     },
//     "resumed_on": "2000-01-01T00:00:00.000000000-00:00"
//   }
//
// @resume msg
type MsgResume struct {
	baseResume
	Msg *flows.MsgIn
}

// NewMsgResume creates a new message resume with the passed in values
func NewMsgResume(env utils.Environment, contact *flows.Contact, msg *flows.MsgIn) *MsgResume {
	return &MsgResume{
		baseResume: newBaseResume(env, contact),
		Msg:        msg,
	}
}

// Type returns the type of this resume
func (r *MsgResume) Type() string { return TypeMsg }

// Apply applies our state changes and saves any events to the run
func (r *MsgResume) Apply(run flows.FlowRun, step flows.Step) error {
	var channel *flows.Channel
	var err error
	if r.Msg.Channel() != nil {
		channel, err = run.Session().Assets().Channels().Get(r.Msg.Channel().UUID)
		if err != nil {
			return err
		}
	}

	// update this run's input
	input := inputs.NewMsgInput(flows.InputUUID(r.Msg.UUID()), channel, r.ResumedOn(), r.Msg.URN(), r.Msg.Text(), r.Msg.Attachments())
	run.SetInput(input)
	run.ResetExpiration(nil)
	run.AddEvent(step, events.NewMsgReceivedEvent(r.Msg))

	return r.baseResume.Apply(run, step)
}

var _ flows.Resume = (*MsgResume)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type msgResumeEnvelope struct {
	baseResumeEnvelope
	Msg *flows.MsgIn `json:"msg" validate:"required,dive"`
}

// ReadMsgResume reads a message resume
func ReadMsgResume(session flows.Session, data json.RawMessage) (flows.Resume, error) {
	resume := &MsgResume{}
	e := msgResumeEnvelope{}
	if err := utils.UnmarshalAndValidate(data, &e); err != nil {
		return nil, err
	}

	if err := unmarshalBaseResume(session, &resume.baseResume, &e.baseResumeEnvelope); err != nil {
		return nil, err
	}

	resume.Msg = e.Msg

	return resume, nil
}
