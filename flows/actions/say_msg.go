package actions

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeSayMsg, func() flows.Action { return &SayMsgAction{} })
}

// TypeSayMsg is the type for the say message action
const TypeSayMsg string = "say_msg"

// SayMsgAction can be used to communicate with the contact in a voice flow by either reading
// a message with TTS or playing a pre-recorded audio file. It will generate an [event:ivr_created]
// event if there is a valid audio URL or backdown text. This will contain a message which
// the caller should handle as an IVR play command if it has an audio attachment, or otherwise
// an IVR say command using the message text.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "say_msg",
//     "audio_url": "http://uploads.temba.io/2353262.m4a",
//     "text": "Hi @contact.name, are you ready to complete today's survey?"
//   }
//
// @action say_msg
type SayMsgAction struct {
	BaseAction
	voiceAction

	Text     string `json:"text" validate:"required"`
	AudioURL string `json:"audio_url,omitempty"`
}

// NewSayMsgAction creates a new say message action
func NewSayMsgAction(uuid flows.ActionUUID, text string, audioURL string) *SayMsgAction {
	return &SayMsgAction{
		BaseAction: NewBaseAction(TypeSayMsg, uuid),
		Text:       text,
		AudioURL:   audioURL,
	}
}

// Validate validates our action is valid and has all the assets it needs
func (a *SayMsgAction) Validate(assets flows.SessionAssets, context *flows.ValidationContext) error {
	return nil
}

// Execute runs this action
func (a *SayMsgAction) Execute(run flows.FlowRun, step flows.Step, logModifier func(flows.Modifier), logEvent func(flows.Event)) error {
	// localize and evaluate the message text
	localizedText := run.GetText(utils.UUID(a.UUID()), "text", a.Text)
	evaluatedText, err := run.EvaluateTemplateAsString(localizedText)
	if err != nil {
		logEvent(events.NewErrorEvent(err))
	}
	evaluatedText = strings.TrimSpace(evaluatedText)

	// localize the audio URL
	localizedAudioURL := run.GetText(utils.UUID(a.UUID()), "audio_url", a.AudioURL)

	// if we have neither an audio URL or backdown text, skip
	if evaluatedText == "" && localizedAudioURL == "" {
		logEvent(events.NewErrorEventf("need either audio URL or backdown text, skipping"))
		return nil
	}

	var attachments []flows.Attachment
	if localizedAudioURL != "" {
		attachments = []flows.Attachment{flows.Attachment(fmt.Sprintf("audio:%s", localizedAudioURL))}
	}

	// an IVR flow must have been started with a connection
	connection := run.Session().Trigger().Connection()

	msg := flows.NewMsgOut(connection.URN(), connection.Channel(), evaluatedText, attachments, nil)
	logEvent(events.NewIVRCreatedEvent(msg))

	return nil
}
