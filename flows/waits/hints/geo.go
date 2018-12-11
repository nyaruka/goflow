package hints

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeGeo, func() flows.Hint { return &GeoHint{} })
}

// TypeGeo is the type of our geo hint
const TypeGeo string = "Geo"

// GeoHint requests a message with an geo attachment
type GeoHint struct {
	baseHint
}

// NewGeoHint creates a new geo hint
func NewGeoHint() *GeoHint {
	return &GeoHint{
		baseHint: newBaseHint(TypeGeo),
	}
}
