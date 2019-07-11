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
	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		fv := v.Field(i)
		ft := t.Field(i)

		fmt.Printf("inspecting field %s %s\n", ft.Name, ft.Type.Kind())

		if ft.Anonymous {
			templateValues(s, fv, include)
		}

		evaluated, localized := parseEngineTag(ft)
		jsonName := jsonNameTag(ft)

		if evaluated {
			if ft.Type.Kind() == reflect.Map {
				for _, v := range fv.Interface().(map[string]string) {
					include.String(v)
				}
			} else if ft.Type.Kind() == reflect.Slice {
				for _, v := range fv.Interface().([]string) {
					include.String(v)
				}
			} else {
				include.String(fv.String())
			}

			if localized {
				include.Translations(s.(flows.Localizable), jsonName)
			}
		}
	}
}

func jsonNameTag(f reflect.StructField) string {
	tagVals := strings.Split(f.Tag.Get("json"), ",")
	if len(tagVals) > 0 {
		return tagVals[0]
	}
	return ""
}

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
