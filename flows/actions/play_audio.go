package actions

import (
	"strings"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	registerType(TypePlayAudio, func() flows.Action { return &PlayAudioAction{} })
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
	baseAction
	voiceAction

	AudioURL string `json:"audio_url" validate:"required" engine:"localized,evaluated"`
}

// NewPlayAudio creates a new play message action
func NewPlayAudio(uuid flows.ActionUUID, audioURL string) *PlayAudioAction {
	return &PlayAudioAction{
		baseAction: newBaseAction(TypePlayAudio, uuid),
		AudioURL:   audioURL,
	}
}

// Execute runs this action
func (a *PlayAudioAction) Execute(run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	// localize and evaluate audio URL
	localizedAudioURL := run.GetText(uuids.UUID(a.UUID()), "audio_url", a.AudioURL)
	evaluatedAudioURL, err := run.EvaluateTemplate(localizedAudioURL)
	if err != nil {
		logEvent(events.NewError(err))
		return nil
	}

	evaluatedAudioURL = strings.TrimSpace(evaluatedAudioURL)
	if evaluatedAudioURL == "" {
		logEvent(events.NewErrorf("audio URL evaluated to empty, skipping"))
		return nil
	}

	// an IVR flow must have been started with a connection
	connection := run.Session().Trigger().Connection()

	// if we have an audio URL, turn it into a message
	msg := flows.NewIVRMsgOut(connection.URN(), connection.Channel(), "", envs.NilLanguage, evaluatedAudioURL)
	logEvent(events.NewIVRCreated(msg))

	return nil
}
