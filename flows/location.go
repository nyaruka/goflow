package flows

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
)

// location levels which can be field types
const (
	LocationLevelState    = envs.LocationLevel(1)
	LocationLevelDistrict = envs.LocationLevel(2)
	LocationLevelWard     = envs.LocationLevel(3)
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
