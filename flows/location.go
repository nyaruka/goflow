package flows

import (
	"strings"

	"github.com/nyaruka/goflow/assets"
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

// Describe returns a representation of this type for error messages
func (p LocationPath) Describe() string { return "location" }

// Reduce returns the primitive version of this type
func (p LocationPath) Reduce(env utils.Environment) types.XPrimitive {
	return types.NewXText(string(p))
}

// ToXJSON is called when this type is passed to @(json(...))
func (p LocationPath) ToXJSON(env utils.Environment) types.XText {
	return p.Reduce(env).ToXJSON(env)
}

var _ types.XValue = LocationPath("")

// LocationAssets provides access to location assets
type LocationAssets struct {
	hierarchies []assets.LocationHierarchy
}

// NewLocationAssets creates a new set of location assets
func NewLocationAssets(hierarchies []assets.LocationHierarchy) *LocationAssets {
	return &LocationAssets{hierarchies: hierarchies}
}

// Hierarchies returns all hierarchies
func (s *LocationAssets) Hierarchies() []assets.LocationHierarchy {
	return s.hierarchies
}
