package flows

import (
	"encoding/json"
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

// LocationHierarchySet defines the unordered set of all location hierarchies for a session
type LocationHierarchySet struct {
	hierarchies []*utils.LocationHierarchy
}

// NewLocationHierarchySet creates a new location hierarchy set from the given list of hierarchies
func NewLocationHierarchySet(hierarchies []*utils.LocationHierarchy) *LocationHierarchySet {
	return &LocationHierarchySet{hierarchies: hierarchies}
}

// All returns all hierarchies in this location hierarchy set
func (s *LocationHierarchySet) All() []*utils.LocationHierarchy {
	return s.hierarchies
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadLocationHierarchySet reads a location hierarchy set from the given JSON
func ReadLocationHierarchySet(data json.RawMessage) (*LocationHierarchySet, error) {
	items, err := utils.UnmarshalArray(data)
	if err != nil {
		return nil, err
	}

	hierarchies := make([]*utils.LocationHierarchy, len(items))
	for d := range items {
		if hierarchies[d], err = utils.ReadLocationHierarchy(items[d]); err != nil {
			return nil, err
		}
	}

	return NewLocationHierarchySet(hierarchies), nil
}
