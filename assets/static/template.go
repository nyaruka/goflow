package static

import (
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/goflow/assets"
)

// Template is a JSON serializable implementation of a template asset
type Template struct {
	UUID_         assets.TemplateUUID    `json:"uuid"         validate:"required,uuid"`
	Name_         string                 `json:"name"`
	Translations_ []*TemplateTranslation `json:"translations"`
}

// NewTemplate creates a new template
func NewTemplate(uuid assets.TemplateUUID, name string, translations []*TemplateTranslation) *Template {
	return &Template{
		UUID_:         uuid,
		Name_:         name,
		Translations_: translations,
	}
}

// UUID returns the UUID of this template
func (t *Template) UUID() assets.TemplateUUID { return t.UUID_ }

// Name returns the name of this template
func (t *Template) Name() string { return t.Name_ }

// Translations returns the translations for this template
func (t *Template) Translations() []assets.TemplateTranslation {
	trs := make([]assets.TemplateTranslation, len(t.Translations_))
	for i := range t.Translations_ {
		trs[i] = t.Translations_[i]
	}
	return trs
}

// TemplateTranslation represents a single template translation
type TemplateTranslation struct {
	Channel_       *assets.ChannelReference `json:"channel"         validate:"required"`
	Content_       string                   `json:"content"         validate:"required"`
	Locale_        i18n.Locale              `json:"locale"          validate:"required"`
	Namespace_     string                   `json:"namespace"`
	VariableCount_ int                      `json:"variable_count"`
}

// NewTemplateTranslation creates a new template translation
func NewTemplateTranslation(channel *assets.ChannelReference, locale i18n.Locale, content string, variableCount int, namespace string) *TemplateTranslation {
	return &TemplateTranslation{
		Channel_:       channel,
		Content_:       content,
		Namespace_:     namespace,
		Locale_:        locale,
		VariableCount_: variableCount,
	}
}

// Content returns the translated content for this template
func (t *TemplateTranslation) Content() string { return t.Content_ }

// Namespace returns the namespace for this template
func (t *TemplateTranslation) Namespace() string { return t.Namespace_ }

// Language returns the locale this translation is in
func (t *TemplateTranslation) Locale() i18n.Locale { return t.Locale_ }

// VariableCount returns the number of variables in this template
func (t *TemplateTranslation) VariableCount() int { return t.VariableCount_ }

// Channel returns the channel this template translation is for
func (t *TemplateTranslation) Channel() *assets.ChannelReference { return t.Channel_ }
