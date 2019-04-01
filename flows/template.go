package flows

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
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

// FindTranslation finds the matching translation for the passed in channel and languages (in priority order)
func (t *Template) FindTranslation(channel assets.ChannelUUID, langs []utils.Language) *TemplateTranslation {
	// first iterate through and find all translations that are for this channel
	matches := make(map[utils.Language]assets.TemplateTranslation)
	for _, tr := range t.Template.Translations() {
		if tr.Channel().UUID == channel {
			matches[tr.Language()] = tr
		}
	}

	// now find the first that matches our language
	for _, lang := range langs {
		tr := matches[lang]
		if tr != nil {
			return NewTemplateTranslation(tr)
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

// FindTranslation looks through our list of templates to find the template matching the passed in uuid
// If no template or translation is found then empty string is returned
func (l *TemplateAssets) FindTranslation(uuid assets.TemplateUUID, channel *assets.ChannelReference, langs []utils.Language) *TemplateTranslation {
	// no channel, can't match to a template
	if channel == nil {
		return nil
	}

	template := l.byUUID[uuid]

	// not found, no template
	if template == nil {
		return nil
	}

	// look through our translations looking for a match by both channel and translation
	return template.FindTranslation(channel.UUID, langs)
}
