package inspect

import (
	"reflect"
	"strings"

	"github.com/nyaruka/goflow/flows"
)

// TemplateValues extracts template values by reading engine tags on a struct
func TemplateValues(s flows.Localizable, localization flows.Localization, include func(string)) {
	templateValues(reflect.ValueOf(s), s, localization, include)
}

func templateValues(v reflect.Value, l flows.Localizable, localization flows.Localization, include func(string)) {
	v = derefValue(v)

	if v.Type().Kind() != reflect.Struct {
		return
	}

	for _, ef := range extractEngineFields(v.Type()) {
		//fmt.Printf("== %v.%s\n", v, ef.jsonName)

		fv := v.FieldByIndex(ef.index)

		if ef.evaluated {
			extractTemplatesFromField(fv, include)

			// if this field is also localized, each translation is a template and needs to be included
			if ef.localized && l != nil {
				for _, lang := range localization.Languages() {
					translations := localization.GetTranslations(lang)
					for _, v := range translations.GetTextArray(l.LocalizationUUID(), ef.jsonName) {
						include(v)
					}
				}
			}
		}

		fv = derefValue(fv)

		if fv.Kind() == reflect.Struct {
			templateValues(fv, nil, localization, include)
		} else if fv.Kind() == reflect.Slice {
			for i := 0; i < fv.Len(); i++ {
				templateValues(fv.Index(i), nil, localization, include)
			}
		}
	}
}

// gets the actual value if we've been given an interface or pointer
func derefValue(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}

type engineField struct {
	jsonName  string
	evaluated bool
	localized bool
	index     []int
}

func extractEngineFields(t reflect.Type) []*engineField {
	fields := make([]*engineField, 0)
	extractEngineFieldsFromType(t, nil, func(f *engineField) {
		fields = append(fields, f)
	})
	return fields
}

func extractEngineFieldsFromType(t reflect.Type, loc []int, include func(*engineField)) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		var index []int
		index = append(index, loc...)
		index = append(index, f.Index...)

		// if this is an embedded base struct, inspect its fields too
		if f.Anonymous {
			extractEngineFieldsFromType(f.Type, index, include)
			continue
		}

		jsonName := jsonNameTag(f)
		if jsonName == "" {
			continue
		}

		evaluated, localized := parseEngineTag(f)

		include(&engineField{
			jsonName:  jsonName,
			evaluated: evaluated,
			localized: localized,
			index:     index,
		})
	}
}

// Evaluated tags can be applied to fields of type string, slices of string or map of strings.
// This method extracts templates from any such field.
func extractTemplatesFromField(v reflect.Value, include func(string)) {
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

// gets the JSON name of the given field
func jsonNameTag(f reflect.StructField) string {
	tagVals := strings.Split(f.Tag.Get("json"), ",")
	if len(tagVals) > 0 {
		return tagVals[0]
	}
	return ""
}

// parses the engine tag on a field if it exists
func parseEngineTag(f reflect.StructField) (evaluated bool, localized bool) {
	tagVals := strings.Split(f.Tag.Get("engine"), ",")
	evaluated = false
	localized = false
	for _, v := range tagVals {
		if v == "evaluated" {
			evaluated = true
		} else if v == "localized" {
			localized = true
		}
	}

	// TODO check if tags are legal for type of f

	return evaluated, localized
}
