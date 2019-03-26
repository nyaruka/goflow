package types

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
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
func (t *Template) Translations() []*TemplateTranslation { return t.t.Translations }

// UnmarshalJSON is our unmarshaller for json data
func (t *Template) UnmarshalJSON(data []byte) error { return json.Unmarshal(data, &t.t) }

// MarshalJSON is our marshaller for json data
func (t *Template) MarshalJSON() ([]byte, error) { return json.Marshal(t.t) }

// TemplateTranslation represents a single template translation
type TemplateTranslation struct {
	t struct {
		Channel       assets.ChannelReference `json:"channel"         validate:"required"`
		Content       string                  `json:"content"         validate:"required"`
		Language      utils.Language          `json:"language"        validate:"required"`
		VariableCount int                     `json:"variable_count"`
	}
}

// NewTemplateTranslation creates a new template translation
func NewTemplateTranslation(channel assets.ChannelReference, language utils.Language, content string, variableCount int) *TemplateTranslation {
	t := &TemplateTranslation{}
	t.t.Channel = channel
	t.t.Content = content
	t.t.Language = language
	t.t.VariableCount = variableCount
	return t
}

// Content returns the translated content for this template
func (t *TemplateTranslation) Content() string { return t.t.Content }

// Language returns the language this translation is in
func (t *TemplateTranslation) Language() utils.Language { return t.t.Language }

// VariableCount returns the number of variables in this template
func (t *TemplateTranslation) VariableCount() int { return t.t.VariableCount }

// Channel returns the channel this template translation is for
func (t *TemplateTranslation) Channel() assets.ChannelReference { return t.t.Channel }

// UnmarshalJSON is our unmarshaller for json data
func (t *TemplateTranslation) UnmarshalJSON(data []byte) error { return json.Unmarshal(data, &t.t) }

// MarshalJSON is our marshaller for json data
func (t *TemplateTranslation) MarshalJSON() ([]byte, error) { return json.Marshal(t.t) }
