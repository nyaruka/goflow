package flows

import (
	"github.com/nyaruka/goflow/assets"
)

// EnumerateTemplateArray enumerates each template in the array
func EnumerateTemplateArray(templates []string, callback func(string)) {
	for _, template := range templates {
		callback(template)
	}
}

// RewriteTemplateArray rewrites each template in the array
func RewriteTemplateArray(templates []string, rewrite func(string) string) {
	for t := range templates {
		templates[t] = rewrite(templates[t])
	}
}

func EnumerateTemplateTranslations(localization Localization, localizable Localizable, key string, callback func(string)) {
	for _, lang := range localization.Languages() {
		translations := localization.GetTranslations(lang)
		for _, tpl := range translations.GetTextArray(localizable.LocalizationUUID(), key) {
			callback(tpl)
		}
	}
}

func RewriteTemplateTranslations(localization Localization, localizable Localizable, key string, rewrite func(string) string) {
	for _, lang := range localization.Languages() {
		translations := localization.GetTranslations(lang)

		templates := translations.GetTextArray(localizable.LocalizationUUID(), key)
		rewritten := make([]string, len(templates))
		for t := range templates {
			rewritten[t] = rewrite(templates[t])
		}
		translations.SetTextArray(localizable.LocalizationUUID(), key, rewritten)
	}
}

func EnumerateTemplatesInGroupReferences(groups []*assets.GroupReference, callback func(string)) {
	for _, group := range groups {
		if group.NameMatch != "" {
			callback(group.NameMatch)
		}
	}
}

func RewriteTemplatesInGroupReferences(groups []*assets.GroupReference, rewrite func(string) string) {
	for _, group := range groups {
		if group.NameMatch != "" {
			group.NameMatch = rewrite(group.NameMatch)
		}
	}
}

func EnumerateTemplatesInLabelReferences(labels []*assets.LabelReference, callback func(string)) {
	for _, label := range labels {
		if label.NameMatch != "" {
			callback(label.NameMatch)
		}
	}
}

func RewriteTemplatesInLabelReferences(labels []*assets.LabelReference, rewrite func(string) string) {
	for _, label := range labels {
		if label.NameMatch != "" {
			label.NameMatch = rewrite(label.NameMatch)
		}
	}
}
