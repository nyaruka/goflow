package events

import (
	"fmt"

	"github.com/nyaruka/goflow/utils"
)

// Hint tells the caller what type of input the flow is expecting
type Hint interface {
	utils.Typed
}

var registeredHintTypes = map[string](func() Hint){}

// registers a new type of hint
func registerHintType(name string, initFunc func() Hint) {
	registeredHintTypes[name] = initFunc
}

// the base of all hint types
type baseHint struct {
	Type_ string `json:"type" validate:"required"`
}

func newBaseHint(typeName string) baseHint {
	return baseHint{Type_: typeName}
}

// Type returns the type of this hint
func (h *baseHint) Type() string { return h.Type_ }

// ReadHint reads a hint from the given JSON
func ReadHint(data []byte) (Hint, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := registeredHintTypes[typeName]
	if f == nil {
		return nil, fmt.Errorf("unknown type: '%s'", typeName)
	}

	hint := f()
	return hint, utils.UnmarshalAndValidate(data, hint)
}

func init() {
	registerHintType(HintTypeAudio, func() Hint { return &AudioHint{} })
	registerHintType(HintTypeDigits, func() Hint { return &DigitsHint{} })
	registerHintType(HintTypeImage, func() Hint { return &ImageHint{} })
	registerHintType(HintTypeLocation, func() Hint { return &LocationHint{} })
	registerHintType(HintTypeVideo, func() Hint { return &VideoHint{} })
}

// HintTypeAudio is the type of our audio hint
const HintTypeAudio string = "audio"

// AudioHint requests a message with an audio attachment
type AudioHint struct {
	baseHint
}

// NewAudioHint creates a new audio hint
func NewAudioHint() *AudioHint {
	return &AudioHint{baseHint: newBaseHint(HintTypeAudio)}
}

// HintTypeDigits is the type of our digits hint
const HintTypeDigits string = "digits"

// DigitsHint requests a message containing one or more digits
type DigitsHint struct {
	baseHint

	Count        *int   `json:"count,omitempty"`
	TerminatedBy string `json:"terminated_by,omitempty"`
}

// NewFixedDigitsHint creates a new digits hint for a fixed count of digits
func NewFixedDigitsHint(count int) *DigitsHint {
	return &DigitsHint{
		baseHint: newBaseHint(HintTypeDigits),
		Count:    &count,
	}
}

// NewTerminatedDigitsHint creates a new digits hint for a sequence of digits terminated by the given key
func NewTerminatedDigitsHint(terminatedBy string) *DigitsHint {
	return &DigitsHint{
		baseHint:     newBaseHint(HintTypeDigits),
		TerminatedBy: terminatedBy,
	}
}

// HintTypeImage is the type of our image hint
const HintTypeImage string = "image"

// ImageHint requests a message with an image attachment
type ImageHint struct {
	baseHint
}

// NewImageHint creates a new image hint
func NewImageHint() *ImageHint {
	return &ImageHint{baseHint: newBaseHint(HintTypeImage)}
}

// HintTypeLocation is the type of our location hint
const HintTypeLocation string = "location"

// LocationHint requests a message with a location attachment, i.e. a geo:<lat>,<long>
type LocationHint struct {
	baseHint
}

// NewLocationHint creates a new location hint
func NewLocationHint() *LocationHint {
	return &LocationHint{baseHint: newBaseHint(HintTypeLocation)}
}

// HintTypeVideo is the type of our video hint
const HintTypeVideo string = "video"

// VideoHint requests a message with an video attachment
type VideoHint struct {
	baseHint
}

// NewVideoHint creates a new video hint
func NewVideoHint() *VideoHint {
	return &VideoHint{baseHint: newBaseHint(HintTypeVideo)}
}
