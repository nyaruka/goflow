package flows

import (
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/goflow/assets"
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
func (t *Template) FindTranslation(channel *Channel, locales []i18n.Locale) *TemplateTranslation {
	// find all translations for this channel
	candidates := make(map[string]*TemplateTranslation)
	candidateLocales := make([]string, 0, 5)
	for _, tr := range t.Template.Translations() {
		if tr.Channel().UUID == channel.UUID() {
			candidates[string(tr.Locale())] = NewTemplateTranslation(tr)
			candidateLocales = append(candidateLocales, string(tr.Locale()))
		}
	}
	if len(candidates) == 0 {
		return nil
	}

	match := i18n.NewBCP47Matcher(candidateLocales...).ForLocales(locales...)
	return candidates[match]
}

// TemplateTranslation represents a single translation for a template
type TemplateTranslation struct {
	assets.TemplateTranslation
}

// NewTemplateTranslation returns a new TemplateTranslation for the passed in asset
func NewTemplateTranslation(t assets.TemplateTranslation) *TemplateTranslation {
	return &TemplateTranslation{TemplateTranslation: t}
}

// Asset returns the underlying asset
func (t *TemplateTranslation) Asset() assets.TemplateTranslation { return t.TemplateTranslation }

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
