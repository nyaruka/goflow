package hints

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeAudio, func() flows.Hint { return &AudioHint{} })
}

// TypeAudio is the type of our audio hint
const TypeAudio string = "audio"

// AudioHint requests a message with an audio attachment
type AudioHint struct {
	baseHint
}

// NewAudioHint creates a new audio hint
func NewAudioHint() *AudioHint {
	return &AudioHint{
		baseHint: newBaseHint(TypeAudio),
	}
}
