package flows

import (
	"slices"
	"text/template"
	"time"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/goflow/envs"
)

type sessionEnvironment struct {
	envs.Environment

	session Session
}

// NewSessionEnvironment creates a new environment from a session's base environment that merges some properties with
// those from the contact.
func NewSessionEnvironment(s Session) envs.Environment {
	return &sessionEnvironment{
		Environment: s.Environment(),
		session:     s,
	}
}

func (e *sessionEnvironment) Timezone() *time.Location {
	contact := e.session.Contact()

	// if we have a contact and they have a timezone that overrides the base enviroment's timezone
	if contact != nil && contact.Timezone() != nil {
		return contact.Timezone()
	}
	return e.Environment.Timezone()
}

func (e *sessionEnvironment) DefaultLanguage() i18n.Language {
	contact := e.session.Contact()

	// if we have a contact and they have a language and it's an allowed language that overrides the base environment's languuage
	if contact != nil && contact.Language() != i18n.NilLanguage && slices.Contains(e.AllowedLanguages(), contact.Language()) {
		return contact.Language()
	}
	return e.Environment.DefaultLanguage()
}

func (e *sessionEnvironment) DefaultCountry() i18n.Country {
	contact := e.session.Contact()

	// if we have a contact and they have a preferred channel with a country that overrides the base environment's country
	if contact != nil {
		cc := contact.Country()
		if cc != i18n.NilCountry {
			return cc
		}
	}
	return e.Environment.DefaultCountry()
}

func (e *sessionEnvironment) DefaultLocale() i18n.Locale {
	return i18n.NewLocale(e.DefaultLanguage(), e.DefaultCountry())
}

func (e *sessionEnvironment) LocationResolver() envs.LocationResolver {
	return e.session.Assets().Locations()
}

// LLMPrompt overrides the base environment to fetch LLM prompts from engine options
func (e *sessionEnvironment) LLMPrompt(name string) *template.Template {
	return e.session.Engine().Options().LLMPrompts[name]
}
