package flows

import (
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/tools"
)

// RunContextTopLevels are the allowed top-level variables for expression evaluations
var RunContextTopLevels = []string{
	"run",
	"child",
	"parent",
	"contact",
	"input",
	"results",
	"trigger",
	"legacy_extra",
}

var fieldRefPaths = [][]string{
	{"contact", "fields"},
	{"parent", "contact", "fields"},
	{"child", "contact", "fields"},
}

// ExtractFieldReferences extracts fields references from the given template
func ExtractFieldReferences(template string) []*assets.FieldReference {
	fieldRefs := make([]*assets.FieldReference, 0)
	tools.FindContextRefsInTemplate(template, RunContextTopLevels, func(path []string) {
		isField, fieldKey := isFieldRefPath(path)
		if isField {
			fieldRefs = append(fieldRefs, assets.NewFieldReference(fieldKey, ""))
		}
	})
	return fieldRefs
}

func isFieldRefPath(path []string) (bool, string) {
	for _, possible := range fieldRefPaths {
		if len(path) == len(possible)+1 {
			matches := true
			for i := range possible {
				if strings.ToLower(path[i]) != possible[i] {
					matches = false
					break
				}
			}
			if matches {
				return true, strings.ToLower(path[len(possible)])
			}
		}
	}
	return false, ""
}

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
