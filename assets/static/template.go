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
	Channel_       *assets.ChannelReference   `json:"channel"         validate:"required"`
	Content_       string                     `json:"content"         validate:"required"`
	Locale_        i18n.Locale                `json:"locale"          validate:"required"`
	Namespace_     string                     `json:"namespace"`
	VariableCount_ int                        `json:"variable_count"`
	Params_        map[string][]TemplateParam `json:"params"`
}

// NewTemplateTranslation creates a new template translation
func NewTemplateTranslation(channel *assets.ChannelReference, locale i18n.Locale, content string, variableCount int, namespace string, params map[string][]TemplateParam) *TemplateTranslation {
	return &TemplateTranslation{
		Channel_:       channel,
		Content_:       content,
		Namespace_:     namespace,
		Locale_:        locale,
		VariableCount_: variableCount,
		Params_:        params,
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

// Params returns the params for this template translation
func (t *TemplateTranslation) Params() map[string][]assets.TemplateParam {
	prs := make(map[string][]assets.TemplateParam, len(t.Params_))
	for k, v := range t.Params_ {
		compParams := make([]assets.TemplateParam, len(v))
		for i, pr := range v {
			compParams[i] = (assets.TemplateParam)(&pr)
		}
		prs[k] = compParams
	}
	return prs
}

// TemplateParam represents a single parameter for a template translation
type TemplateParam struct {
	Type_ string `json:"type"`
}

// Type returns the type for this parameter
func (t *TemplateParam) Type() string { return t.Type_ }

// NewTemplateParam creates a new template param
func NewTemplateParam(paramType string) TemplateParam {
	return TemplateParam{Type_: paramType}
}
