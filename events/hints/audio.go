package hints

func init() {
	registerType(TypeAudio, func() Hint { return &Audio{} })
}

// TypeAudio is the type of our audio hint
const TypeAudio string = "audio"

// Audio requests a message with an audio attachment
type Audio struct {
	baseHint
}

// NewAudio creates a new audio hint
func NewAudio() *Audio {
	return &Audio{
		baseHint: newBaseHint(TypeAudio),
	}
}
