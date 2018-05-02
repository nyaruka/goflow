package flows

import (
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// location levels which can be field types
const (
	LocationLevelState    = utils.LocationLevel(1)
	LocationLevelDistrict = utils.LocationLevel(2)
	LocationLevelWard     = utils.LocationLevel(3)
)

// LocationPath is a location described by a path Country > State ...
type LocationPath string

func (p LocationPath) String() string {
	return string(p)
}

// Reduce returns the primitive version of this type
func (p LocationPath) Reduce() types.XPrimitive {
	return types.NewXText(string(p))
}

// ToXJSON is called when this type is passed to @(json(...))
func (p LocationPath) ToXJSON() types.XText {
	return p.Reduce().ToXJSON()
}

var _ types.XValue = LocationPath("")
