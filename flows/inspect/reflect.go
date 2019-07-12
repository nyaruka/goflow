package inspect

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/nyaruka/goflow/flows"
)

// TemplateValuesByTags extracts template values by reading engine tags on a struct
func TemplateValuesByTags(s interface{}, include flows.TemplateIncluder) {
	templateValues(s, reflect.ValueOf(s), include)
}

func templateValues(s interface{}, v reflect.Value, include flows.TemplateIncluder) {
	// get the actual struct if we've been given an interface or pointer
	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for _, ef := range extractEngineFields(v.Type()) {
		fv := v.FieldByIndex(ef.index)

		if ef.evaluated {
			// extract from single strings, slices of strings or maps of string values
			if fv.Type().Kind() == reflect.Map {
				for _, v := range fv.Interface().(map[string]string) {
					include.String(v)
				}
			} else if fv.Type().Kind() == reflect.Slice {
				for _, v := range fv.Interface().([]string) {
					include.String(v)
				}
			} else if fv.Type().Kind() == reflect.String {
				include.String(fv.String())
			} else {
				panic(fmt.Sprintf("engine:evaluated found on field %T.%s which not a supported type (%s)", s, ef.jsonName, fv.Type()))
			}

			if ef.localized {
				localizable, isLocalizable := s.(flows.Localizable)
				if !isLocalizable {
					panic(fmt.Sprintf("engine:localized found on %T which doesn't implement Localizable", s))
				}

				include.Translations(localizable, ef.jsonName)
			}
		}
	}
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
		ft := t.Field(i)

		var index []int
		index = append(index, loc...)
		index = append(index, ft.Index...)

		// if this is an embedded base struct, inspect its fields too
		if ft.Anonymous {

			extractEngineFieldsFromType(ft.Type, index, include)
			continue
		}

		jsonName := jsonNameTag(ft)
		if jsonName == "" {
			continue
		}

		evaluated, localized := parseEngineTag(ft)

		include(&engineField{
			jsonName:  jsonName,
			evaluated: evaluated,
			localized: localized,
			index:     index,
		})
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
	return evaluated, localized
}
