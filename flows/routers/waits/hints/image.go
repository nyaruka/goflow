package hints

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeImage, func() flows.Hint { return &Image{} })
}

// TypeImage is the type of our image hint
const TypeImage string = "image"

// Image requests a message with an image attachment
type Image struct {
	baseHint
}

// NewImage creates a new image hint
func NewImage() *Image {
	return &Image{
		baseHint: newBaseHint(TypeImage),
	}
}
