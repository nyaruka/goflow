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

var fieldRefPaths = [][]string{
	{"fields"},
	{"contact", "fields"},
	{"parent", "fields"},
	{"parent", "contact", "fields"},
	{"child", "fields"},
	{"child", "contact", "fields"},
}

// ExtractFieldReferences extracts fields references from the given template
func ExtractFieldReferences(template string) []*assets.FieldReference {
	fieldRefs := make([]*assets.FieldReference, 0)
	tools.FindContextRefsInTemplate(template, flows.RunContextTopLevels, func(path []string) {
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
