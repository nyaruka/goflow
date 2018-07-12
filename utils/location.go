package utils

import (
	"encoding/json"
	"strings"
)

// LocationLevel is a numeric level, e.g. 0 = country, 1 = state
type LocationLevel int

const (
	LocationPathSeparator       = ">"
	LocationPaddedPathSeparator = " > "
)

// Location represents a single Location
type Location struct {
	level    LocationLevel
	name     string
	path     string
	aliases  []string
	parent   *Location
	children []*Location
}

// NewLocation creates a new location object
func NewLocation(level LocationLevel, name string) *Location {
	return &Location{level: level, name: name}
}

// Level gets the level of this location
func (l *Location) Level() LocationLevel { return l.level }

// Name gets the name of this location
func (l *Location) Name() string { return l.name }

// Path gets the full path of this location
func (l *Location) Path() string { return l.path }

// Aliases gets the aliases of this location
func (l *Location) Aliases() []string { return l.aliases }

// Parent gets the parent of this location
func (l *Location) Parent() *Location { return l.parent }

// Children gets the children of this location
func (l *Location) Children() []*Location { return l.children }

func (l *Location) String() string { return l.path }

// utility for traversing the location hierarchy
type locationVisitor func(Location *Location)

func (l *Location) visit(visitor locationVisitor) {
	visitor(l)
	for _, child := range l.children {
		child.visit(visitor)
	}
}

type locationPathLookup map[string]*Location

func (p locationPathLookup) addLookup(path string, location *Location) {
	p[strings.ToLower(path)] = location
}

func (p locationPathLookup) lookup(path string) *Location { return p[strings.ToLower(path)] }

// location names aren't always unique in a given level - i.e. you can have two wards with the same name, but different parents
type locationNameLookup map[string][]*Location

func (n locationNameLookup) addLookup(name string, location *Location) {
	name = strings.ToLower(name)
	n[name] = append(n[name], location)
}

func (n locationNameLookup) lookup(name string) []*Location { return n[strings.ToLower(name)] }

// LocationHierarchy is a hierarical tree of locations
type LocationHierarchy struct {
	root *Location

	// for faster lookups
	levelLookups []locationNameLookup
	pathLookup   locationPathLookup
}

// NewLocationHierarchy cretes a new location hierarchy
func NewLocationHierarchy(root *Location, numLevels int) *LocationHierarchy {
	h := &LocationHierarchy{
		root:         root,
		levelLookups: make([]locationNameLookup, numLevels),
		pathLookup:   make(locationPathLookup),
	}

	for l := 0; l < numLevels; l++ {
		h.levelLookups[l] = make(locationNameLookup)
	}

	// traverse the hierarchy to setup paths and lookups
	root.visit(func(location *Location) {
		if location.parent != nil {
			location.path = strings.Join([]string{location.parent.path, location.name}, LocationPaddedPathSeparator)
		} else {
			location.path = location.name
		}

		h.pathLookup.addLookup(location.path, location)
		h.addNameLookups(location)
	})
	return h
}

func (h *LocationHierarchy) addNameLookups(location *Location) {
	lookups := h.levelLookups[int(location.level)]
	lookups.addLookup(location.name, location)

	// include any aliases as names too
	for _, alias := range location.aliases {
		lookups.addLookup(alias, location)
	}
}

// Root gets the root location of this hierarchy (typically a country)
func (h *LocationHierarchy) Root() *Location {
	return h.root
}

// FindByName looks for all locations in the hierarchy with the given level and name or alias
func (h *LocationHierarchy) FindByName(name string, level LocationLevel, parent *Location) []*Location {

	// try it as a path first if it looks possible
	if level == 0 || strings.Contains(name, LocationPathSeparator) {
		match := h.pathLookup.lookup(name)
		if match != nil {
			return []*Location{match}
		}
	}

	if int(level) < len(h.levelLookups) {
		matches := h.levelLookups[int(level)].lookup(name)
		if matches != nil {
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

// FindByPath looks for a location in the hierarchy with the given path
func (h *LocationHierarchy) FindByPath(path string) *Location {
	return h.pathLookup.lookup(strings.ToLower(path))
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type locationEnvelope struct {
	Name     string              `json:"name" validate:"required"`
	Aliases  []string            `json:"aliases,omitempty"`
	Children []*locationEnvelope `json:"children,omitempty"`
}

func locationFromEnvelope(envelope *locationEnvelope, currentLevel LocationLevel, parent *Location) *Location {
	location := &Location{
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
	if err := UnmarshalAndValidate(data, &le); err != nil {
		return nil, err
	}

	root := locationFromEnvelope(&le, LocationLevel(0), nil)

	return NewLocationHierarchy(root, 4), nil
}
