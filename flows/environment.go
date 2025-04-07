package flows

import (
	"slices"
	"text/template"
	"time"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/goflow/envs"
)

type assetsEnvironment struct {
	envs.Environment

	sa SessionAssets
}

// NewAssetsEnvironment creates a new environment that has access to assets which can be used for resolving locations.
func NewAssetsEnvironment(base envs.Environment, sa SessionAssets) envs.Environment {
	return &assetsEnvironment{Environment: base, sa: sa}
}

func (e *assetsEnvironment) LocationResolver() envs.LocationResolver {
	return e.sa.Locations()
}

type sessionEnvironment struct {
	envs.Environment

	session Session
}

// NewSessionEnvironment creates a new environment from a session's base environment that merges some properties with
// those from the contact.
func NewSessionEnvironment(s Session) envs.Environment {
	return &sessionEnvironment{
		Environment: NewAssetsEnvironment(s.Environment(), s.Assets()),
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

// LLMPrompt overrides the base environment to fetch LLM prompts from engine options
func (e *sessionEnvironment) LLMPrompt(name string) *template.Template {
	return e.session.Engine().Options().LLMPrompts[name]
}
