package types

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils/jsonx"
)

// Template is a JSON serializable implementation of a template asset
type Template struct {
	t struct {
		UUID         assets.TemplateUUID    `json:"uuid"         validate:"required,uuid"`
		Name         string                 `json:"name"`
		Translations []*TemplateTranslation `json:"translations"`
	}
}

// NewTemplate creates a new template
func NewTemplate(uuid assets.TemplateUUID, name string, translations []*TemplateTranslation) *Template {
	t := &Template{}
	t.t.UUID = uuid
	t.t.Name = name
	t.t.Translations = translations
	return t
}

// UUID returns the UUID of this template
func (t *Template) UUID() assets.TemplateUUID { return t.t.UUID }

// Name returns the name of this template
func (t *Template) Name() string { return t.t.Name }

// Translations returns the translations for this template
func (t *Template) Translations() []assets.TemplateTranslation {
	trs := make([]assets.TemplateTranslation, len(t.t.Translations))
	for i := range t.t.Translations {
		trs[i] = t.t.Translations[i]
	}
	return trs
}

// UnmarshalJSON is our unmarshaller for json data
func (t *Template) UnmarshalJSON(data []byte) error { return jsonx.Unmarshal(data, &t.t) }

// MarshalJSON is our marshaller for json data
func (t *Template) MarshalJSON() ([]byte, error) { return jsonx.Marshal(t.t) }

// TemplateTranslation represents a single template translation
type TemplateTranslation struct {
	t struct {
		Channel       assets.ChannelReference `json:"channel"         validate:"required"`
		Content       string                  `json:"content"         validate:"required"`
		Language      envs.Language           `json:"language"        validate:"required"`
		Country       envs.Country            `json:"country,omitempty"`
		VariableCount int                     `json:"variable_count"`
	}
}

// NewTemplateTranslation creates a new template translation
func NewTemplateTranslation(channel assets.ChannelReference, language envs.Language, country envs.Country, content string, variableCount int) *TemplateTranslation {
	t := &TemplateTranslation{}
	t.t.Channel = channel
	t.t.Content = content
	t.t.Language = language
	t.t.Country = country
	t.t.VariableCount = variableCount
	return t
}

// Content returns the translated content for this template
func (t *TemplateTranslation) Content() string { return t.t.Content }

// Language returns the language this translation is in
func (t *TemplateTranslation) Language() envs.Language { return t.t.Language }

// Country returns the country this translation is for if any
func (t *TemplateTranslation) Country() envs.Country { return t.t.Country }

// VariableCount returns the number of variables in this template
func (t *TemplateTranslation) VariableCount() int { return t.t.VariableCount }

// Channel returns the channel this template translation is for
func (t *TemplateTranslation) Channel() assets.ChannelReference { return t.t.Channel }

// UnmarshalJSON is our unmarshaller for json data
func (t *TemplateTranslation) UnmarshalJSON(data []byte) error { return jsonx.Unmarshal(data, &t.t) }

// MarshalJSON is our marshaller for json data
func (t *TemplateTranslation) MarshalJSON() ([]byte, error) { return jsonx.Marshal(t.t) }
