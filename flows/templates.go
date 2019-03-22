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

// checks whether the given context path is a reference to a contact field
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
func EnumerateTemplateArray(templates []string, include func(string)) {
	for _, template := range templates {
		include(template)
	}
}

// RewriteTemplateArray rewrites each template in the array
func RewriteTemplateArray(templates []string, rewrite func(string) string) {
	for t := range templates {
		templates[t] = rewrite(templates[t])
	}
}

func EnumerateTemplateTranslations(localization Localization, localizable Localizable, key string, include func(string)) {
	for _, lang := range localization.Languages() {
		translations := localization.GetTranslations(lang)
		for _, tpl := range translations.GetTextArray(localizable.LocalizationUUID(), key) {
			include(tpl)
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

// wrapper for an asset reference to make it inspectable
type inspectableReference struct {
	ref assets.Reference
}

// InspectReference inspects the given asset reference if it's non-nil
func InspectReference(ref assets.Reference, inspect func(Inspectable)) {
	if ref != nil {
		inspectableReference{ref: ref}.Inspect(inspect)
	}
}

// Inspect inspects this object and any children
func (r inspectableReference) Inspect(inspect func(Inspectable)) {
	inspect(r)
}

// EnumerateTemplates enumerates all expressions on this object and its children
func (r inspectableReference) EnumerateTemplates(localization Localization, include func(string)) {
	if r.ref != nil && r.ref.Variable() {
		switch typed := r.ref.(type) {
		case *assets.GroupReference:
			include(typed.NameMatch)
		case *assets.LabelReference:
			include(typed.NameMatch)
		}
	}
}

// RewriteTemplates rewrites all templates on this object and its children
func (r inspectableReference) RewriteTemplates(localization Localization, rewrite func(string) string) {
	if r.ref != nil && r.ref.Variable() {
		switch typed := r.ref.(type) {
		case *assets.GroupReference:
			typed.NameMatch = rewrite(typed.NameMatch)
		case *assets.LabelReference:
			typed.NameMatch = rewrite(typed.NameMatch)
		}
	}
}

// EnumerateDependencies enumerates all dependencies on this object and its children
func (r inspectableReference) EnumerateDependencies(localization Localization, include func(assets.Reference)) {
	if r.ref != nil && !r.ref.Variable() {
		include(r.ref)
	}
}

// EnumerateResults enumerates all potential results on this object
// Asset references can't contain results.
func (r inspectableReference) EnumerateResults(include func(*ResultSpec)) {}
