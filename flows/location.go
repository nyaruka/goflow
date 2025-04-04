package flows

import (
	"regexp"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
)

// location levels which can be field types
const (
	LocationLevelState    = envs.LocationLevel(1)
	LocationLevelDistrict = envs.LocationLevel(2)
	LocationLevelWard     = envs.LocationLevel(3)
)

// LocationAssets provides access to location assets and implements envs.LocationResolver
type LocationAssets struct {
	hierarchies []assets.LocationHierarchy
}

// NewLocationAssets creates a new set of location assets
func NewLocationAssets(hierarchies []assets.LocationHierarchy) *LocationAssets {
	return &LocationAssets{hierarchies: hierarchies}
}

// FindLocations returns locations with the matching name (case-insensitive), level and parent (optional)
func (s *LocationAssets) FindLocations(env envs.Environment, name string, level envs.LocationLevel, parent *envs.Location) []*envs.Location {
	if len(s.hierarchies) > 0 {
		return s.hierarchies[0].FindByName(env, name, level, parent)
	}
	return nil
}

// FindLocationsFuzzy returns matching locations like FindLocations but attempts the following strategies
// to find locations:
//  1. Exact match
//  2. Match with punctuation removed
//  3. Split input into words and try to match each word
//  4. Try to match pairs of words
func (s *LocationAssets) FindLocationsFuzzy(env envs.Environment, text string, level envs.LocationLevel, parent *envs.Location) []*envs.Location {
	// try matching name exactly
	if locations := s.FindLocations(env, text, level, parent); len(locations) > 0 {
		return locations
	}

	// try with punctuation removed
	stripped := strings.TrimSpace(regexp.MustCompile(`[\s\p{P}]+`).ReplaceAllString(text, ""))
	if locations := s.FindLocations(env, stripped, level, parent); len(locations) > 0 {
		return locations
	}

	// try on each tokenized word
	re := regexp.MustCompile(`[\p{L}\d]+(-[\p{L}\d]+)*`)
	words := re.FindAllString(text, -1)
	for _, word := range words {
		if locations := s.FindLocations(env, word, level, parent); len(locations) > 0 {
			return locations
		}
	}

	// try with each pair of words
	for i := 0; i < len(words)-1; i++ {
		wordPair := strings.Join(words[i:i+2], " ")
		if locations := s.FindLocations(env, wordPair, level, parent); len(locations) > 0 {
			return locations
		}
	}

	return nil
}

func (s *LocationAssets) LookupLocation(path envs.LocationPath) *envs.Location {
	if len(s.hierarchies) > 0 {
		return s.hierarchies[0].FindByPath(path)
	}
	return nil
}

var _ envs.LocationResolver = (*LocationAssets)(nil)
