package inspect

import (
	"reflect"
	"sort"
	"strings"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/tools"
	"github.com/nyaruka/goflow/flows"
)

// Templates extracts template values by reading engine tags on a struct
func Templates(s any, localization flows.Localization, include func(i18n.Language, string)) {
	templateValues(reflect.ValueOf(s), localization, include)
}

func templateValues(v reflect.Value, localization flows.Localization, include func(i18n.Language, string)) {
	walk(v,
		func(v reflect.Value) {
			te, ok := v.Interface().(flows.TemplateEnumerator)
			if ok {
				te.EnumerateTemplates(localization, include)
			}
		},
		func(sv reflect.Value, fv reflect.Value, ef *EngineField) {
			if ef.Evaluated {
				extractTemplates(fv, i18n.NilLanguage, include)

				// if this field is also localized, each translation is a template and needs to be included
				if ef.Localized && localization != nil {
					localizable := sv.Interface().(flows.Localizable)

					Translations(localization, localizable.LocalizationUUID(), ef.JSONName, include)
				}
			}
		},
	)
}

func Translations(localization flows.Localization, itemUUID uuids.UUID, property string, include func(i18n.Language, string)) {
	for _, lang := range localization.Languages() {
		for _, v := range localization.GetItemTranslation(lang, itemUUID, property) {
			include(lang, v)
		}
	}
}

// Evaluated tags can be applied to fields of type string, slices of string or map of strings.
// This method extracts template values from any such field.
func extractTemplates(v reflect.Value, lang i18n.Language, include func(i18n.Language, string)) {
	switch typed := v.Interface().(type) {
	case map[string]string:
		for _, i := range typed {
			include(lang, i)
		}
	case []string:
		for _, i := range typed {
			include(lang, i)
		}
	case string:
		include(lang, typed)
	}
}

func TemplatePaths(t reflect.Type) []string {
	paths := []string{}

	walkTypes(t, "", func(path string, ef *EngineField) {
		if ef.Evaluated {
			if ef.Type.Kind() == reflect.Map {
				paths = append(paths, path+".*")
			} else if ef.Type.Kind() == reflect.Slice {
				paths = append(paths, path+"[*]")
			} else {
				paths = append(paths, path)
			}
		}
	})

	sort.Strings(paths)
	return paths
}

// all the paths in the context where contact field references are found
var fieldRefPaths = [][]string{
	{"fields"},
	{"contact", "fields"},
	{"parent", "fields"},
	{"parent", "contact", "fields"},
	{"child", "fields"},
	{"child", "contact", "fields"},
}

// ExtractFromTemplate extracts asset references and parent result references from the given template. Note that
// duplicates are not removed.
func ExtractFromTemplate(template string) ([]assets.Reference, []string) {
	assetRefs := make([]assets.Reference, 0)
	parentRefs := make([]string, 0)

	tools.FindContextRefsInTemplate(template, flows.RunContextTopLevels, func(path []string) {
		if len(path) <= 1 {
			return
		}

		path0 := strings.ToLower(path[0])
		path1 := strings.ToLower(path[1])

		if path0 == "globals" {
			assetRefs = append(assetRefs, assets.NewGlobalReference(path1, ""))
		} else if path0 == "parent" && path1 == "results" && len(path) > 2 {
			parentRefs = append(parentRefs, strings.ToLower(path[2]))
		} else {
			isField, fieldKey := isFieldRefPath(path)
			if isField {
				assetRefs = append(assetRefs, assets.NewFieldReference(fieldKey, ""))
			}
		}
	})
	return assetRefs, parentRefs
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
