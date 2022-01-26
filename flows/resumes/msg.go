package resumes

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/inputs"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeMsg, readMsgResume)
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
//       "created_on": "2018-01-01T12:00:00.000000Z",
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
	msg *flows.MsgIn
}

// NewMsg creates a new message resume with the passed in values
func NewMsg(env envs.Environment, contact *flows.Contact, msg *flows.MsgIn) *MsgResume {
	return &MsgResume{
		baseResume: newBaseResume(TypeMsg, env, contact),
		msg:        msg,
	}
}

// Msg returns the msg this resume is based on
func (r *MsgResume) Msg() *flows.MsgIn { return r.msg }

// Apply applies our state changes and saves any events to the run
func (r *MsgResume) Apply(run flows.Run, logEvent flows.EventCallback) {
	// do base changes (contact, environment)
	r.baseResume.Apply(run, logEvent)

	// update our input
	input := inputs.NewMsg(run.Session().Assets(), r.msg, r.ResumedOn())

	run.Session().SetInput(input)

	logEvent(events.NewMsgReceived(r.msg))
}

var _ flows.Resume = (*MsgResume)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type msgResumeEnvelope struct {
	baseResumeEnvelope
	Msg *flows.MsgIn `json:"msg" validate:"required,dive"`
}

func readMsgResume(sessionAssets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Resume, error) {
	e := &msgResumeEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	r := &MsgResume{
		msg: e.Msg,
	}

	if err := r.unmarshal(sessionAssets, &e.baseResumeEnvelope, missing); err != nil {
		return nil, err
	}

	return r, nil
}

// MarshalJSON marshals this resume into JSON
func (r *MsgResume) MarshalJSON() ([]byte, error) {
	e := &msgResumeEnvelope{
		Msg: r.msg,
	}

	if err := r.marshal(&e.baseResumeEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
