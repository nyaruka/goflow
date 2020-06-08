package flows

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
)

// Template represents messaging templates used by channels types such as WhatsApp
type Template struct {
	assets.Template
}

// NewTemplate returns a new template objects based on the passed in asset
func NewTemplate(t assets.Template) *Template {
	return &Template{Template: t}
}

// Asset returns the underlying asset
func (t *Template) Asset() assets.Template { return t.Template }

// Reference returns the reference for this template
func (t *Template) Reference() *assets.TemplateReference {
	if t == nil {
		return nil
	}
	return assets.NewTemplateReference(t.UUID(), t.Name())
}

// FindTranslation finds the matching translation for the passed in channel and languages (in priority order)
func (t *Template) FindTranslation(channel assets.ChannelUUID, locales []envs.Locale) *TemplateTranslation {
	// first iterate through and find all translations that are for this channel
	candidatesByLocale := make(map[envs.Locale]*TemplateTranslation)
	candidatesByLang := make(map[envs.Language]*TemplateTranslation)
	for _, tr := range t.Template.Translations() {
		if tr.Channel().UUID == channel {
			tt := NewTemplateTranslation(tr)
			candidatesByLocale[tt.Locale()] = tt
			candidatesByLang[tt.Language()] = tt
		}
	}

	// first look for exact locale match
	for _, locale := range locales {
		tt := candidatesByLocale[locale]
		if tt != nil {
			return tt
		}
	}

	// if that fails look for language match
	for _, locale := range locales {
		tt := candidatesByLang[locale.Language]
		if tt != nil {
			return tt
		}
	}

	return nil
}

// TemplateTranslation represents a single translation for a template
type TemplateTranslation struct {
	assets.TemplateTranslation
}

// NewTemplateTranslation returns a new TemplateTranslation for the passed in asset
func NewTemplateTranslation(t assets.TemplateTranslation) *TemplateTranslation {
	return &TemplateTranslation{TemplateTranslation: t}
}

// Locale returns the locale
func (t *TemplateTranslation) Locale() envs.Locale { return envs.NewLocale(t.Language(), t.Country()) }

// Asset returns the underlying asset
func (t *TemplateTranslation) Asset() assets.TemplateTranslation { return t.TemplateTranslation }

var templateRegex = regexp.MustCompile(`({{\d+}})`)

// Substitute substitutes the passed in variables in our template
func (t *TemplateTranslation) Substitute(vars []string) string {
	s := string(t.Content())
	for i, v := range vars {
		s = strings.ReplaceAll(s, fmt.Sprintf("{{%d}}", i+1), v)
	}

	// replace any remaining unmatched items
	s = templateRegex.ReplaceAllString(s, "")

	return s
}

// TemplateAssets is our type for all the templates in an environment
type TemplateAssets struct {
	templates []*Template
	byUUID    map[assets.TemplateUUID]*Template
}

// NewTemplateAssets creates a new template list
func NewTemplateAssets(ts []assets.Template) *TemplateAssets {
	templates := make([]*Template, len(ts))
	byUUID := make(map[assets.TemplateUUID]*Template)
	for i, t := range ts {
		template := NewTemplate(t)
		templates[i] = template
		byUUID[t.UUID()] = template
	}

	return &TemplateAssets{
		templates: templates,
		byUUID:    byUUID,
	}
}

// Get returns the template with the passed in UUID if any
func (a *TemplateAssets) Get(uuid assets.TemplateUUID) *Template {
	return a.byUUID[uuid]
}

// FindTranslation looks through our list of templates to find the template matching the passed in uuid
// If no template or translation is found then empty string is returned
func (a *TemplateAssets) FindTranslation(uuid assets.TemplateUUID, channel *assets.ChannelReference, locales []envs.Locale) *TemplateTranslation {
	// no channel, can't match to a template
	if channel == nil {
		return nil
	}

	template := a.byUUID[uuid]

	// not found, no template
	if template == nil {
		return nil
	}

	// look through our translations looking for a match by both channel and translation
	return template.FindTranslation(channel.UUID, locales)
}
