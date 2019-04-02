package flows

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// location levels which can be field types
const (
	LocationLevelState    = utils.LocationLevel(1)
	LocationLevelDistrict = utils.LocationLevel(2)
	LocationLevelWard     = utils.LocationLevel(3)
)

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
