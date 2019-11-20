package inspect

import (
	"reflect"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/tools"
	"github.com/nyaruka/goflow/flows"
)

// Templates extracts template values by reading engine tags on a struct
func Templates(s interface{}, localization flows.Localization, include func(string)) {
	templateValues(reflect.ValueOf(s), localization, include)
}

func templateValues(v reflect.Value, localization flows.Localization, include func(string)) {
	walk(v, nil, func(sv reflect.Value, fv reflect.Value, ef *EngineField) {
		if ef.Evaluated {
			extractTemplates(fv, include)

			// if this field is also localized, each translation is a template and needs to be included
			if ef.Localized && localization != nil {
				localizable := sv.Interface().(flows.Localizable)

				for _, lang := range localization.Languages() {
					translations := localization.GetTranslations(lang)
					for _, v := range translations.GetTextArray(localizable.LocalizationUUID(), ef.JSONName) {
						include(v)
					}
				}
			}
		}
	})
}

// Evaluated tags can be applied to fields of type string, slices of string or map of strings.
// This method extracts template values from any such field.
func extractTemplates(v reflect.Value, include func(string)) {
	switch typed := v.Interface().(type) {
	case map[string]string:
		for _, i := range typed {
			include(i)
		}
	case []string:
		for _, i := range typed {
			include(i)
		}
	case string:
		include(v.String())
	}
}

func TemplatePaths(t reflect.Type, base string, include func(string)) {
	walkTypes(t, base, func(path string, ef *EngineField) {
		if ef.Evaluated {
			if ef.Type.Kind() == reflect.Map || ef.Type.Kind() == reflect.Slice {
				include(path + "[*]")
			} else {
				include(path)
			}
		}
	})
}

var fieldRefPaths = [][]string{
	{"fields"},
	{"contact", "fields"},
	{"parent", "fields"},
	{"parent", "contact", "fields"},
	{"child", "fields"},
	{"child", "contact", "fields"},
}

// ExtractFromTemplate extracts asset references from the given template
func ExtractFromTemplate(template string) []assets.Reference {
	refs := make([]assets.Reference, 0)
	tools.FindContextRefsInTemplate(template, flows.RunContextTopLevels, func(path []string) {
		if len(path) == 2 && path[0] == "globals" {
			refs = append(refs, assets.NewGlobalReference(strings.ToLower(path[1]), ""))
		} else {
			isField, fieldKey := isFieldRefPath(path)
			if isField {
				refs = append(refs, assets.NewFieldReference(fieldKey, ""))
			}
		}
	})
	return refs
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
