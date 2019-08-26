package hints

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeLocation, func() flows.Hint { return &LocationHint{} })
}

// TypeLocation is the type of our location hint
const TypeLocation string = "location"

// LocationHint requests a message with a location attachment, i.e. a geo:<lat>,<long>
type LocationHint struct {
	baseHint
}

// NewLocationHint creates a new location hint
func NewLocationHint() *LocationHint {
	return &LocationHint{
		baseHint: newBaseHint(TypeLocation),
	}
}
