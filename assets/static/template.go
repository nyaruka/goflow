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
	Channel_    *assets.ChannelReference `json:"channel"      validate:"required"`
	Locale_     i18n.Locale              `json:"locale"       validate:"required"`
	Components_ []*TemplateComponent     `json:"components"`
	Variables_  []*TemplateVariable      `json:"variables"`
}

// NewTemplateTranslation creates a new template translation
func NewTemplateTranslation(channel *assets.ChannelReference, locale i18n.Locale, components []*TemplateComponent, variables []*TemplateVariable) *TemplateTranslation {
	return &TemplateTranslation{
		Channel_:    channel,
		Locale_:     locale,
		Components_: components,
		Variables_:  variables,
	}
}

// Components returns the components for this template translation
func (t *TemplateTranslation) Components() []assets.TemplateComponent {
	cs := make([]assets.TemplateComponent, len(t.Components_))
	for k, tc := range t.Components_ {
		cs[k] = tc
	}
	return cs
}

// Variables returns the variables for this template translation
func (t *TemplateTranslation) Variables() []assets.TemplateVariable {
	vs := make([]assets.TemplateVariable, len(t.Variables_))
	for i := range t.Variables_ {
		vs[i] = t.Variables_[i]
	}
	return vs
}

// Locale returns the locale this translation is in
func (t *TemplateTranslation) Locale() i18n.Locale { return t.Locale_ }

// Channel returns the channel this template translation is for
func (t *TemplateTranslation) Channel() *assets.ChannelReference { return t.Channel_ }

type TemplateComponent struct {
	Name_      string         `json:"name"`
	Type_      string         `json:"type"`
	Content_   string         `json:"content"`
	Display_   string         `json:"display"`
	Variables_ map[string]int `json:"variables"`
}

// Name returns the name for this template component
func (t *TemplateComponent) Name() string { return t.Name_ }

// Type returns the type for this template component
func (t *TemplateComponent) Type() string { return t.Type_ }

// Content returns the content for this template component
func (t *TemplateComponent) Content() string { return t.Content_ }

// Display returns the display for this template component
func (t *TemplateComponent) Display() string { return t.Display_ }

// Variables returns the variable mapping for this template component
func (t *TemplateComponent) Variables() map[string]int { return t.Variables_ }

// NewTemplateComponent creates a new template param
func NewTemplateComponent(name, type_, content, display string, variables map[string]int) *TemplateComponent {
	return &TemplateComponent{Type_: type_, Name_: name, Content_: content, Display_: display, Variables_: variables}
}

// TemplateVariable represents a single variable for a template translation
type TemplateVariable struct {
	Type_ string `json:"type"`
}

// Type returns the type for this parameter
func (t *TemplateVariable) Type() string { return t.Type_ }

// NewTemplateVariable creates a new template variable
func NewTemplateVariable(paramType string) *TemplateVariable {
	return &TemplateVariable{Type_: paramType}
}
