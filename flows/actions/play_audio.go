package actions

import (
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

func init() {
	RegisterType(TypePlayAudio, func() flows.Action { return &PlayAudioAction{} })
}

// TypePlayAudio is the type for the play audio action
const TypePlayAudio string = "play_audio"

// PlayAudioAction can be used to play an audio recording in a voice flow. It will generate
// an [event:ivr_play] event.
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

// Validate validates our action is valid and has all the assets it needs
func (a *PlayAudioAction) Validate(assets flows.SessionAssets, context *flows.ValidationContext) error {
	return nil
}

// Execute runs this action
func (a *PlayAudioAction) Execute(run flows.FlowRun, step flows.Step) error {
	// localize and evaluate audio URL
	localizedAudioURL := run.GetText(utils.UUID(a.UUID()), "audio_url", a.AudioURL)
	evaluatedAudioURL, err := run.EvaluateTemplateAsString(localizedAudioURL)
	if err != nil {
		a.logError(run, step, err)
		return nil
	}

	evaluatedAudioURL = strings.TrimSpace(evaluatedAudioURL)
	if evaluatedAudioURL == "" {
		a.logError(run, step, errors.Errorf("audio URL evaluated to empty, skipping"))
		return nil
	}

	// if we have an audio URL, tell caller to play it
	a.log(run, step, events.NewIVRPlayEvent(evaluatedAudioURL, ""))

	return nil
}
