package hints

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeVideo, func() flows.Hint { return &Video{} })
}

// TypeVideo is the type of our video hint
const TypeVideo string = "video"

// Video requests a message with an video attachment
type Video struct {
	baseHint
}

// NewVideo creates a new video hint
func NewVideo() *Video {
	return &Video{
		baseHint: newBaseHint(TypeVideo),
	}
}
