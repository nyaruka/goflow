package runs

import (
	"time"

	"github.com/nyaruka/goflow/utils"
)

// a run specific environment which allows values to be overridden by the contact
type runEnvironment struct {
	utils.Environment
	run *flowRun

	cachedLanguages utils.LanguageList
}

// creates a run environment based on the given run
func newRunEnvironment(base utils.Environment, run *flowRun) *runEnvironment {
	env := &runEnvironment{base, run, nil}
	env.refreshLanguagesCache()
	return env
}

func (e *runEnvironment) Timezone() *time.Location {
	contact := e.run.contact

	// if run has a contact with a timezone, that overrides the enviroment's timezone
	if contact != nil && contact.Timezone() != nil {
		return contact.Timezone()
	}
	return e.run.Session().Environment().Timezone()
}

func (e *runEnvironment) Languages() utils.LanguageList {
	// if contact language has changed, rebuild our cached language list
	if e.run.Contact() != nil && e.cachedLanguages[0] != e.run.Contact().Language() {
		e.refreshLanguagesCache()
	}

	return e.cachedLanguages
}

func (e *runEnvironment) Locations() (*utils.LocationHierarchy, error) {
	sessionAssets := e.run.Session().Assets()
	if sessionAssets.HasLocations() {
		return sessionAssets.GetLocationHierarchy()
	}

	return nil, nil
}

func (e *runEnvironment) refreshLanguagesCache() {
	contact := e.run.contact
	var languages utils.LanguageList

	// if contact has a language, it takes priority
	if contact != nil && contact.Language() != utils.NilLanguage {
		languages = append(languages, contact.Language())
	}

	// next we include any environment languages
	languages = append(languages, e.run.Session().Environment().Languages()...)

	// finally we include the flow native language
	languages = append(languages, e.run.flow.Language())

	e.cachedLanguages = languages.RemoveDuplicates()
}
