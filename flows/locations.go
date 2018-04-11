package flows

import (
	"encoding/json"
	"strings"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// LocationID is the unique identifier for each location, e.g. an OSM ID
type LocationID string

// LocationLevel is a numeric level, e.g. 0 = country, 1 = state
type LocationLevel int

// Location represents a single Location
type Location struct {
	id       LocationID
	level    LocationLevel
	name     string
	aliases  []string
	parent   *Location
	children []*Location
}

// NewLocation creates a new location object
func NewLocation(id LocationID, level LocationLevel, name string) *Location {
	return &Location{id: id, level: level, name: name}
}

// ID gets the id of this location
func (b *Location) ID() LocationID { return b.id }

// Level gets the level of this location
func (b *Location) Level() LocationLevel { return b.level }

// Name gets the name of this location
func (b *Location) Name() string { return b.name }

// Aliases gets the aliases of this location
func (b *Location) Aliases() []string { return b.aliases }

// Parent gets the parent of this location
func (b *Location) Parent() *Location { return b.parent }

// Children gets the children of this location
func (b *Location) Children() []*Location { return b.children }

// Reduce is called when this object needs to be reduced to a primitive
func (b *Location) Reduce() types.XPrimitive { return types.NewXString(b.name) }

func (b *Location) ToXJSON() types.XString { return types.NewXString("TODO") }

var _ types.XValue = (*Location)(nil)

type locationVisitor func(Location *Location)

func (b *Location) visit(visitor locationVisitor) {
	visitor(b)
	for _, child := range b.children {
		child.visit(visitor)
	}
}

// for each level, we maintain some maps for faster lookups
type levelLookup struct {
	byID   map[LocationID]*Location
	byName map[string][]*Location
}

func (l *levelLookup) setIDLookup(id LocationID, location *Location) {
	l.byID[id] = location
}

func (l *levelLookup) addNameLookup(name string, location *Location) {
	name = strings.ToLower(name)
	l.byName[name] = append(l.byName[name], location)
}

// LocationHierarchy is a hierarical tree of locations
type LocationHierarchy struct {
	root         *Location
	levelLookups []*levelLookup
}

// NewLocationHierarchy cretes a new location hierarchy
func NewLocationHierarchy(root *Location, numLevels int) *LocationHierarchy {
	s := &LocationHierarchy{
		root:         root,
		levelLookups: make([]*levelLookup, numLevels),
	}

	for l := 0; l < numLevels; l++ {
		s.levelLookups[l] = &levelLookup{
			byID:   make(map[LocationID]*Location),
			byName: make(map[string][]*Location),
		}
	}

	root.visit(func(Location *Location) { s.addLookups(Location) })
	return s
}

func (s *LocationHierarchy) addLookups(location *Location) {
	lookups := s.levelLookups[int(location.level)]
	lookups.setIDLookup(location.id, location)
	lookups.addNameLookup(location.name, location)

	// include any aliases as names too
	for _, alias := range location.aliases {
		lookups.addNameLookup(alias, location)
	}
}

// Root gets the root location of this hierarchy
func (s *LocationHierarchy) Root() *Location {
	return s.root
}

// FindByID looks for a location in the hierarchy with the given level and ID
func (s *LocationHierarchy) FindByID(id LocationID, level LocationLevel) *Location {
	if int(level) < len(s.levelLookups) {
		return s.levelLookups[int(level)].byID[id]
	}
	return nil
}

// FindByName looks for all locations in the hierarchy with the given level and name or alias
func (s *LocationHierarchy) FindByName(name string, level LocationLevel, parent *Location) []*Location {
	if int(level) < len(s.levelLookups) {
		matches, found := s.levelLookups[int(level)].byName[strings.ToLower(name)]
		if found {
			// if a parent is specified, filter the matches by it
			if parent != nil {
				withParent := make([]*Location, 0)
				for m := range matches {
					if matches[m].parent == parent {
						withParent = append(withParent, matches[m])
					}
				}
				return withParent
			}

			return matches
		}
	}
	return []*Location{}
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type locationEnvelope struct {
	ID       LocationID          `json:"id"`
	Name     string              `json:"name" validate:"required"`
	Aliases  []string            `json:"aliases,omitempty"`
	Children []*locationEnvelope `json:"children,omitempty"`
}

func locationFromEnvelope(envelope *locationEnvelope, currentLevel LocationLevel, parent *Location) *Location {
	location := &Location{
		id:      envelope.ID,
		level:   LocationLevel(currentLevel),
		name:    envelope.Name,
		aliases: envelope.Aliases,
		parent:  parent,
	}

	location.children = make([]*Location, len(envelope.Children))
	for c := range envelope.Children {
		location.children[c] = locationFromEnvelope(envelope.Children[c], currentLevel+1, location)
	}

	return location
}

// ReadLocationHierarchy reads a location hierarchy from the given JSON
func ReadLocationHierarchy(data json.RawMessage) (*LocationHierarchy, error) {
	var le locationEnvelope
	if err := utils.UnmarshalAndValidate(data, &le, "location"); err != nil {
		return nil, err
	}

	root := locationFromEnvelope(&le, LocationLevel(0), nil)

	return NewLocationHierarchy(root, 4), nil
}
