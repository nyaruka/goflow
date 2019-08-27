package hints

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeImage, func() flows.Hint { return &ImageHint{} })
}

// TypeImage is the type of our image hint
const TypeImage string = "image"

// ImageHint requests a message with an image attachment
type ImageHint struct {
	baseHint
}

// NewImageHint creates a new image hint
func NewImageHint() *ImageHint {
	return &ImageHint{
		baseHint: newBaseHint(TypeImage),
	}
}
