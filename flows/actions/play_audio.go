package actions

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypePlayAudio, func() flows.Action { return &PlayAudioAction{} })
}

// TypePlayAudio is the type for the play audio action
const TypePlayAudio string = "play_audio"

// PlayAudioAction can be used to play an audio recording in a voice flow. It will generate an
// [event:ivr_created] event if there is a valid audio URL. This will contain a message which
// the caller should handle as an IVR play command using the audio attachment.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "play_audio",
//     "audio_url": "http://uploads.temba.io/2353262.m4a"
//   }
//
// @action play_audio
type PlayAudioAction struct {
	BaseAction
	voiceAction

	AudioURL string `json:"audio_url" validate:"required"`
}

// NewPlayAudioAction creates a new play message action
func NewPlayAudioAction(uuid flows.ActionUUID, audioURL string) *PlayAudioAction {
	return &PlayAudioAction{
		BaseAction: NewBaseAction(TypePlayAudio, uuid),
		AudioURL:   audioURL,
	}
}

// Execute runs this action
func (a *PlayAudioAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	// localize and evaluate audio URL
	localizedAudioURL := run.GetText(utils.UUID(a.UUID()), "audio_url", a.AudioURL)
	evaluatedAudioURL, err := run.EvaluateTemplate(localizedAudioURL)
	if err != nil {
		logEvent(events.NewErrorEvent(err))
		return nil
	}

	evaluatedAudioURL = strings.TrimSpace(evaluatedAudioURL)
	if evaluatedAudioURL == "" {
		logEvent(events.NewErrorEventf("audio URL evaluated to empty, skipping"))
		return nil
	}

	// an IVR flow must have been started with a connection
	connection := run.Session().Trigger().Connection()

	// if we have an audio URL, turn it into a message
	attachments := []utils.Attachment{utils.Attachment(fmt.Sprintf("audio:%s", evaluatedAudioURL))}
	msg := flows.NewMsgOut(connection.URN(), connection.Channel(), "", attachments, nil, nil)
	logEvent(events.NewIVRCreatedEvent(msg))

	return nil
}

// Inspect inspects this object and any children
func (a *PlayAudioAction) Inspect(inspect func(flows.Inspectable)) {
	inspect(a)
}

// EnumerateTemplates enumerates all expressions on this object and its children
func (a *PlayAudioAction) EnumerateTemplates(localization flows.Localization, include func(string)) {
	include(a.AudioURL)
	flows.EnumerateTemplateTranslations(localization, a, "audio_url", include)
}

// RewriteTemplates rewrites all templates on this object and its children
func (a *PlayAudioAction) RewriteTemplates(localization flows.Localization, rewrite func(string) string) {
	a.AudioURL = rewrite(a.AudioURL)
	flows.RewriteTemplateTranslations(localization, a, "audio_url", rewrite)
}
