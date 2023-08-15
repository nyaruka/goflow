package flows

import (
	"regexp"
	"strings"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"golang.org/x/exp/slices"
)

type environment struct {
	envs.Environment

	locationResolver envs.LocationResolver
}

// NewEnvironment creates a new environment
func NewEnvironment(base envs.Environment, la *LocationAssets) envs.Environment {
	var locationResolver envs.LocationResolver

	hierarchies := la.Hierarchies()
	if len(hierarchies) > 0 {
		locationResolver = &assetLocationResolver{hierarchies[0]}
	}

	return &environment{base, locationResolver}
}

func (e *environment) LocationResolver() envs.LocationResolver {
	return e.locationResolver
}

type assetLocationResolver struct {
	locations assets.LocationHierarchy
}

// FindLocations returns locations with the matching name (case-insensitive), level and parent (optional)
func (r *assetLocationResolver) FindLocations(name string, level envs.LocationLevel, parent *envs.Location) []*envs.Location {
	return r.locations.FindByName(name, level, parent)
}

// FindLocationsFuzzy returns matching locations like FindLocations but attempts the following strategies
// to find locations:
//  1. Exact match
//  2. Match with punctuation removed
//  3. Split input into words and try to match each word
//  4. Try to match pairs of words
func (r *assetLocationResolver) FindLocationsFuzzy(text string, level envs.LocationLevel, parent *envs.Location) []*envs.Location {
	// try matching name exactly
	if locations := r.FindLocations(text, level, parent); len(locations) > 0 {
		return locations
	}

	// try with punctuation removed
	stripped := strings.TrimSpace(regexp.MustCompile(`[\s\p{P}]+`).ReplaceAllString(text, ""))
	if locations := r.FindLocations(stripped, level, parent); len(locations) > 0 {
		return locations
	}

	// try on each tokenized word
	re := regexp.MustCompile(`[\p{L}\d]+(-[\p{L}\d]+)*`)
	words := re.FindAllString(text, -1)
	for _, word := range words {
		if locations := r.FindLocations(word, level, parent); len(locations) > 0 {
			return locations
		}
	}

	// try with each pair of words
	for i := 0; i < len(words)-1; i++ {
		wordPair := strings.Join(words[i:i+2], " ")
		if locations := r.FindLocations(wordPair, level, parent); len(locations) > 0 {
			return locations
		}
	}

	return []*envs.Location{}
}

func (r *assetLocationResolver) LookupLocation(path envs.LocationPath) *envs.Location {
	return r.locations.FindByPath(path)
}

// an extended environment which takes some values from a contact if there is one and if the have those values.
type mergedEnvironment struct {
	envs.Environment

	session Session
}

// NewMergedEnvironment creates a new merged environment from a session's base environment and its contact
func NewMergedEnvironment(s Session) envs.Environment {
	return &mergedEnvironment{
		NewEnvironment(s.Environment(), s.Assets().Locations()),
		s,
	}
}

func (e *mergedEnvironment) Timezone() *time.Location {
	contact := e.session.Contact()

	// if we have a contact and they have a timezone that overrides the base enviroment's timezone
	if contact != nil && contact.Timezone() != nil {
		return contact.Timezone()
	}
	return e.Environment.Timezone()
}

func (e *mergedEnvironment) DefaultLanguage() envs.Language {
	contact := e.session.Contact()

	// if we have a contact and they have a language and it's an allowed language that overrides the base environment's languuage
	if contact != nil && contact.Language() != envs.NilLanguage && slices.Contains(e.AllowedLanguages(), contact.Language()) {
		return contact.Language()
	}
	return e.Environment.DefaultLanguage()
}

func (e *mergedEnvironment) DefaultCountry() envs.Country {
	contact := e.session.Contact()

	// if we have a contact and they have a preferred channel with a country that overrides the base environment's country
	if contact != nil {
		cc := contact.Country()
		if cc != envs.NilCountry {
			return cc
		}
	}
	return e.Environment.DefaultCountry()
}

func (e *mergedEnvironment) DefaultLocale() envs.Locale {
	return envs.NewLocale(e.DefaultLanguage(), e.DefaultCountry())
}
