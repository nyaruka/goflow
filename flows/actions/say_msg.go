package actions

import (
	"context"
	"strings"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	registerType(TypeSayMsg, func() flows.Action { return &SayMsg{} })
}

// TypeSayMsg is the type for the say message action
const TypeSayMsg string = "say_msg"

// SayMsg can be used to communicate with the contact in a voice flow by either reading
// a message with TTS or playing a pre-recorded audio file. It will generate an [event:ivr_created]
// event if there is a valid audio URL or backdown text. This will contain a message which
// the caller should handle as an IVR play command if it has an audio attachment, or otherwise
// an IVR say command using the message text.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "say_msg",
//	  "audio_url": "http://uploads.temba.io/2353262.m4a",
//	  "text": "Hi @contact.name, are you ready to complete today's survey?"
//	}
//
// @action say_msg
type SayMsg struct {
	baseAction
	voiceAction

	Text     string `json:"text"                 validate:"required,max=10000" engine:"localized,evaluated"`
	AudioURL string `json:"audio_url,omitempty"`
}

// NewSayMsg creates a new say message action
func NewSayMsg(uuid flows.ActionUUID, text string, audioURL string) *SayMsg {
	return &SayMsg{
		baseAction: newBaseAction(TypeSayMsg, uuid),
		Text:       text,
		AudioURL:   audioURL,
	}
}

// Execute runs this action
func (a *SayMsg) Execute(ctx context.Context, run flows.Run, step flows.Step, log flows.EventLogger) error {
	// localize and evaluate the message text
	localizedText, textLang := run.GetText(uuids.UUID(a.UUID()), "text", a.Text)
	evaluatedText, _ := run.EvaluateTemplate(localizedText, log)
	evaluatedText = strings.TrimSpace(evaluatedText)

	// localize the audio URL
	localizedAudioURL, _ := run.GetText(uuids.UUID(a.UUID()), "audio_url", a.AudioURL)

	// if we have neither an audio URL or backdown text, skip
	if evaluatedText == "" && localizedAudioURL == "" {
		log(events.NewError("Can't say message without audio URL or backdown text", ""))
		return nil
	}

	// an IVR flow must have been started with a call
	call := run.Session().Call()

	msg := flows.NewIVRMsgOut(call.URN(), call.Channel().Reference(), evaluatedText, localizedAudioURL, currentLocale(run, textLang))
	log(events.NewIVRCreated(msg))

	return nil
}
