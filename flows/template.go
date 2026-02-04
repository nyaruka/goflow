package flows

import (
	"fmt"
	"strings"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/stringsx"
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

// Templating generates a templating object for the passed in translation and variables
func (t *Template) Templating(tt *TemplateTranslation, vars []string) *MsgTemplating {
	// cross-reference with asset to get variable types and pad out any missing variables
	variables := make([]*TemplatingVariable, len(tt.Variables()))
	for i, v := range tt.Variables() {
		value := ""
		if i < len(vars) {
			value = vars[i]
		}
		variables[i] = &TemplatingVariable{Type: v.Type(), Value: value}
	}

	// create a list of components that have variables
	components := make([]*TemplatingComponent, 0, len(tt.Components()))
	for _, comp := range tt.Components() {
		if len(comp.Variables()) > 0 {
			components = append(components, &TemplatingComponent{
				Type:      comp.Type(),
				Name:      comp.Name(),
				Variables: comp.Variables(),
			})
		}
	}

	return NewMsgTemplating(t.Reference(), components, variables)
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

// Preview returns message content which will act as a preview of a message sent with this template
func (t *TemplateTranslation) Preview(vars []*TemplatingVariable) *MsgContent {
	var text []string
	var attachments []utils.Attachment
	var quickReplies []QuickReply

	for _, comp := range t.Components() {
		content := comp.Content()
		for key, index := range comp.Variables() {
			variable := vars[index]

			if variable.Type == "text" {
				content = strings.ReplaceAll(content, fmt.Sprintf("{{%s}}", key), variable.Value)
			} else if (variable.Type == "image" || variable.Type == "video" || variable.Type == "document") && utils.IsValidAttachment(variable.Value) {
				attachments = append(attachments, utils.Attachment(variable.Value))
			}
		}

		if content != "" {
			if comp.Type() == "header/text" || comp.Type() == "body/text" || comp.Type() == "footer/text" {
				text = append(text, content)
			} else if strings.HasPrefix(comp.Type(), "button/") {
				quickReplies = append(quickReplies, QuickReply{Text: stringsx.TruncateEllipsis(content, MaxQuickReplyTextLength)})
			}
		}
	}

	return &MsgContent{Text: strings.Join(text, "\n\n"), Attachments: attachments, QuickReplies: quickReplies}
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
