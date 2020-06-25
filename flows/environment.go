package flows

import (
	"regexp"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
)

type environment struct {
	envs.Environment

	locations assets.LocationHierarchy
}

// NewEnvironment creates a new environment
func NewEnvironment(base envs.Environment, locations assets.LocationHierarchy) Environment {
	return &environment{base, locations}
}

// HasLocations returns whether this environment has location support
func (e *environment) HasLocations() bool {
	return e.locations != nil
}

// FindLocations returns locations with the matching name (case-insensitive), level and parent (optional)
func (e *environment) FindLocations(name string, level utils.LocationLevel, parent *utils.Location) []*utils.Location {
	return e.locations.FindByName(name, level, parent)
}

// FindLocationsFuzzy returns matching locations like FindLocations but attempts the following strategies
// to find locations:
//   1. Exact match
//   2. Match with punctuation removed
//   3. Split input into words and try to match each word
//   4. Try to match pairs of words
func (e *environment) FindLocationsFuzzy(text string, level utils.LocationLevel, parent *utils.Location) []*utils.Location {
	// try matching name exactly
	if locations := e.FindLocations(text, level, parent); len(locations) > 0 {
		return locations
	}

	// try with punctuation removed
	stripped := strings.TrimSpace(regexp.MustCompile(`\W+`).ReplaceAllString(text, ""))
	if locations := e.FindLocations(stripped, level, parent); len(locations) > 0 {
		return locations
	}

	// try on each tokenized word
	words := regexp.MustCompile(`\W+`).Split(text, -1)
	for _, word := range words {
		if locations := e.FindLocations(word, level, parent); len(locations) > 0 {
			return locations
		}
	}

	// try with each pair of words
	for i := 0; i < len(words)-1; i++ {
		wordPair := strings.Join(words[i:i+2], " ")
		if locations := e.FindLocations(wordPair, level, parent); len(locations) > 0 {
			return locations
		}
	}

	return []*utils.Location{}
}

func (e *environment) LookupLocation(path utils.LocationPath) *utils.Location {
	return e.locations.FindByPath(path)
}

var _ Environment = (*environment)(nil)
