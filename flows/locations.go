package flows

import (
	"strings"

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

// IsPossibleLocationPath returns whether the given string could be a location path
func IsPossibleLocationPath(str string) bool {
	return strings.Contains(str, utils.LocationPathSeparator)
}

// Name returns the name of the location referenced
func (p LocationPath) Name() string {
	parts := strings.Split(string(p), utils.LocationPathSeparator)
	return strings.TrimSpace(parts[len(parts)-1])
}

func (p LocationPath) String() string {
	return string(p)
}

// Repr returns the representation of this type
func (p LocationPath) Repr() string { return "location" }

// Reduce returns the primitive version of this type
func (p LocationPath) Reduce() types.XPrimitive {
	return types.NewXText(string(p))
}

// ToXJSON is called when this type is passed to @(json(...))
func (p LocationPath) ToXJSON() types.XText {
	return p.Reduce().ToXJSON()
}

var _ types.XValue = LocationPath("")
