package actions

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypePlayMsg, func() flows.Action { return &PlayMsgAction{} })
}

// TypePlayMsg is the type for the play recording action
const TypePlayMsg string = "play_msg"

// PlayMsgAction can be used to communicate with the contact in a voice flow by either reading
// a message with TTS or playing a pre-recorded audio file. If there is an audio file, it takes
// priority and an [event:ivr_play] event is generated. Otherwise the text is used
// and a [event:ivr_say] event is generated.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "play_msg",
//     "audio_url": "http://uploads.temba.io/2353262.m4a",
//     "text": "Hi @contact.name, are you ready to complete today's survey?"
//   }
//
// @action play_msg
type PlayMsgAction struct {
	BaseAction
	voiceAction

	Text     string `json:"text" validate:"required"`
	AudioURL string `json:"audio_url"`
}

// NewPlayMsgAction creates a new play message action
func NewPlayMsgAction(uuid flows.ActionUUID, audioURL string, text string) *PlayMsgAction {
	return &PlayMsgAction{
		BaseAction: NewBaseAction(TypePlayMsg, uuid),
		Text:       text,
		AudioURL:   audioURL,
	}
}

// Validate validates our action is valid and has all the assets it needs
func (a *PlayMsgAction) Validate(assets flows.SessionAssets, context *flows.ValidationContext) error {
	return nil
}

// Execute runs this action
func (a *PlayMsgAction) Execute(run flows.FlowRun, step flows.Step) error {
	// localize and evaluate the message text
	localizedText := run.GetText(utils.UUID(a.UUID()), "text", a.Text)
	evaluatedText, err := run.EvaluateTemplateAsString(localizedText)
	if err != nil {
		a.logError(run, step, err)
	}
	evaluatedText = strings.TrimSpace(evaluatedText)

	// localize the audio URL
	localizedAudioURL := run.GetText(utils.UUID(a.UUID()), "audio_url", a.AudioURL)

	// if we have either an audio URL or backdown text.. tell caller to play this
	if localizedAudioURL != "" {
		a.log(run, step, events.NewIVRPlayEvent(localizedAudioURL, evaluatedText))
	} else if evaluatedText != "" {
		a.log(run, step, events.NewIVRSayEvent(evaluatedText))
	} else {
		a.logError(run, step, fmt.Errorf("need either audio URL or backdown text, skipping"))
	}

	return nil
}
