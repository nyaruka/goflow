package envs

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/nyaruka/goflow/utils"
)

// LocationLevel is a numeric level, e.g. 0 = country, 1 = state
type LocationLevel int

// LocationPath is a location described by a path Country > State ...
type LocationPath string

// LocationResolver is used to resolve locations from names or hierarchical paths
type LocationResolver interface {
	FindLocations(string, LocationLevel, *Location) []*Location
	FindLocationsFuzzy(string, LocationLevel, *Location) []*Location
	LookupLocation(LocationPath) *Location
}

const (
	LocationPathSeparator = ">"
)

var spaceRegex = regexp.MustCompile(`\s+`)

// IsPossibleLocationPath returns whether the given string could be a location path
func IsPossibleLocationPath(str string) bool {
	return strings.Contains(str, LocationPathSeparator)
}

func NewLocationPath(parts ...string) LocationPath {
	return LocationPath(strings.Join(parts, " "+LocationPathSeparator+" "))
}

func (p LocationPath) join(name string) LocationPath {
	return NewLocationPath(string(p), name)
}

// Name returns the name of the location referenced
func (p LocationPath) Name() string {
	parts := strings.Split(string(p), LocationPathSeparator)
	return strings.TrimSpace(parts[len(parts)-1])
}

// Normalize normalizes this location path
func (p LocationPath) Normalize() LocationPath {
	// trim any period at end
	normalized := strings.TrimRight(string(p), ".")

	// normalize casing and spacing between location parts
	parts := strings.Split(normalized, LocationPathSeparator)
	for i, part := range parts {
		part = spaceRegex.ReplaceAllString(strings.TrimSpace(part), " ")
		part = strings.Title(strings.ToLower(part))
		parts[i] = part
	}

	return NewLocationPath(parts...)
}

// Location represents a single Location
type Location struct {
	level    LocationLevel
	name     string
	path     LocationPath
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
func (l *Location) Path() LocationPath { return l.path }

// Aliases gets the aliases of this location
func (l *Location) Aliases() []string { return l.aliases }

// Parent gets the parent of this location
func (l *Location) Parent() *Location { return l.parent }

// Children gets the children of this location
func (l *Location) Children() []*Location { return l.children }

func (l *Location) String() string { return string(l.path) }

// utility for traversing the location hierarchy
type locationVisitor func(Location *Location)

func (l *Location) visit(visitor locationVisitor) {
	visitor(l)
	for _, child := range l.children {
		child.visit(visitor)
	}
}

type locationPathLookup map[LocationPath]*Location

func (p locationPathLookup) addLookup(path LocationPath, location *Location) {
	p[path.Normalize()] = location
}

func (p locationPathLookup) lookup(path LocationPath) *Location { return p[path.Normalize()] }

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
	h := &LocationHierarchy{}
	h.initializeFromRoot(root, numLevels)
	return h
}

// NewLocationHierarchy cretes a new location hierarchy
func (h *LocationHierarchy) initializeFromRoot(root *Location, numLevels int) {
	h.root = root
	h.levelLookups = make([]locationNameLookup, numLevels)
	h.pathLookup = make(locationPathLookup)

	for i := 0; i < numLevels; i++ {
		h.levelLookups[i] = make(locationNameLookup)
	}

	// traverse the hierarchy to setup paths and lookups
	root.visit(func(location *Location) {
		if location.parent != nil {
			location.path = location.parent.path.join(location.name)
		} else {
			location.path = LocationPath(location.name)
		}

		h.pathLookup.addLookup(location.path, location)
		h.addNameLookups(location)
	})
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
	if level == 0 || IsPossibleLocationPath(name) {
		match := h.pathLookup.lookup(LocationPath(name))
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
				for i := range matches {
					if matches[i].parent == parent {
						withParent = append(withParent, matches[i])
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
func (h *LocationHierarchy) FindByPath(path LocationPath) *Location {
	return h.pathLookup.lookup(path)
}

func (h *LocationHierarchy) UnmarshalJSON(data []byte) error {
	var le locationEnvelope
	if err := utils.UnmarshalAndValidate(data, &le); err != nil {
		return err
	}

	root := locationFromEnvelope(&le, LocationLevel(0), nil)
	h.initializeFromRoot(root, 4)
	return nil
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
	for i := range envelope.Children {
		location.children[i] = locationFromEnvelope(envelope.Children[i], currentLevel+1, location)
	}

	return location
}

// ReadLocationHierarchy reads a location hierarchy from the given JSON
func ReadLocationHierarchy(data json.RawMessage) (*LocationHierarchy, error) {
	var le locationEnvelope
	if err := utils.UnmarshalAndValidate(data, &le); err != nil {
		return nil, err
	}

	root := locationFromEnvelope(&le, LocationLevel(0), nil)

	return NewLocationHierarchy(root, 4), nil
}
