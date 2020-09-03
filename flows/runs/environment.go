package runs

import (
	"time"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
)

// an extended environment which takes some values from a contact if there is one and if the have those values.
type runEnvironment struct {
	envs.Environment

	run *flowRun
}

// creates a run environment based on the given run
func newRunEnvironment(base envs.Environment, run *flowRun) envs.Environment {
	return &runEnvironment{
		flows.NewEnvironment(base, run.Session().Assets().Locations()),
		run,
	}
}

func (e *runEnvironment) Timezone() *time.Location {
	contact := e.run.Contact()

	// if we have a contact and they have a timezone that overrides the base enviroment's timezone
	if contact != nil && contact.Timezone() != nil {
		return contact.Timezone()
	}
	return e.Environment.Timezone()
}

func (e *runEnvironment) DefaultLanguage() envs.Language {
	contact := e.run.Contact()

	// if we have a contact and they have a language and it's an allowed language that overrides the base environment's languuage
	if contact != nil && contact.Language() != envs.NilLanguage && isAllowedLanguage(e, contact.Language()) {
		return contact.Language()
	}
	return e.Environment.DefaultLanguage()
}

func (e *runEnvironment) DefaultCountry() envs.Country {
	contact := e.run.Contact()

	// if we have a contact and they have a preferred channel with a country that overrides the base environment's country
	if contact != nil {
		cc := contact.Country()
		if cc != envs.NilCountry {
			return cc
		}
	}
	return e.Environment.DefaultCountry()
}

func (e *runEnvironment) DefaultLocale() envs.Locale {
	return envs.NewLocale(e.DefaultLanguage(), e.DefaultCountry())
}

func isAllowedLanguage(e envs.Environment, language envs.Language) bool {
	for _, l := range e.AllowedLanguages() {
		if language == l {
			return true
		}
	}
	return false
}
