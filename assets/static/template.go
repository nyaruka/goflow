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
	Channel_    *assets.ChannelReference `json:"channel"         validate:"required"`
	Locale_     i18n.Locale              `json:"locale"          validate:"required"`
	Status_     assets.TemplateStatus    `json:"status"          validate:"required"`
	Namespace_  string                   `json:"namespace"`
	Components_ []*TemplateComponent     `json:"components"`
}

// NewTemplateTranslation creates a new template translation
func NewTemplateTranslation(channel *assets.ChannelReference, locale i18n.Locale, status assets.TemplateStatus, namespace string, components []*TemplateComponent) *TemplateTranslation {
	return &TemplateTranslation{
		Channel_:    channel,
		Namespace_:  namespace,
		Status_:     status,
		Locale_:     locale,
		Components_: components,
	}
}

// Components returns the components structure for this template
func (t *TemplateTranslation) Components() []assets.TemplateComponent {
	tcs := make([]assets.TemplateComponent, len(t.Components_))
	for k, tc := range t.Components_ {
		tcs[k] = tc
	}
	return tcs
}

// Namespace returns the namespace for this template
func (t *TemplateTranslation) Namespace() string { return t.Namespace_ }

// Status returns the status for this translation
func (t *TemplateTranslation) Status() assets.TemplateStatus { return t.Status_ }

// Language returns the locale this translation is in
func (t *TemplateTranslation) Locale() i18n.Locale { return t.Locale_ }

// Channel returns the channel this template translation is for
func (t *TemplateTranslation) Channel() *assets.ChannelReference { return t.Channel_ }

type TemplateComponent struct {
	Type_    string           `json:"type"`
	Name_    string           `json:"name"`
	Content_ string           `json:"content"`
	Display_ string           `json:"display"`
	Params_  []*TemplateParam `json:"params"`
}

// Type returns the type for this template component
func (t *TemplateComponent) Type() string { return t.Type_ }

// Name returns the name for this template component
func (t *TemplateComponent) Name() string { return t.Name_ }

// Content returns the content for this template component
func (t *TemplateComponent) Content() string { return t.Content_ }

// Display returns the display for this template component
func (t *TemplateComponent) Display() string { return t.Display_ }

// Params returns the params for this template component
func (t *TemplateComponent) Params() []assets.TemplateParam {
	tps := make([]assets.TemplateParam, len(t.Params_))
	for i := range t.Params_ {
		tps[i] = t.Params_[i]
	}
	return tps
}

// NewTemplateComponent creates a new template param
func NewTemplateComponent(type_, name, content, display string, params []*TemplateParam) *TemplateComponent {
	return &TemplateComponent{Type_: type_, Name_: name, Content_: content, Display_: display, Params_: params}
}

// TemplateParam represents a single parameter for a template translation
type TemplateParam struct {
	Type_ string `json:"type"`
	Name_ string `json:"name"`
}

// Type returns the type for this parameter
func (t *TemplateParam) Type() string { return t.Type_ }

// Name returns the name for this parameter
func (t *TemplateParam) Name() string { return t.Name_ }

// NewTemplateParam creates a new template param
func NewTemplateParam(paramType string, paramName string) *TemplateParam {
	return &TemplateParam{Type_: paramType, Name_: paramName}
}
