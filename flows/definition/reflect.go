package definition

import (
	"reflect"
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// ExtractTemplates extracts a list of evaluate-able templates from the given flow
func ExtractTemplates(flow flows.Flow) []string {
	templates := make([]string, 0)
	visitor := func(p reflect.Value, v reflect.Value, tag reflect.StructTag) {
		localize, evaluate := parseEngineTag(tag)
		if evaluate {
			templates = append(templates, extractTemplates(v, localize)...)

			if localize {
				// find the parent struct's UUID and this field's JSON name
				uuid := findStructUUID(p)
				key := ""
				jsonTag := strings.Split(tag.Get("json"), ",")
				if len(jsonTag) > 0 {
					key = jsonTag[0]
				}

				//fmt.Printf("field %+v localized at uuid=%s key=%s\n", v, uuid, key)

				for _, lang := range flow.Localization().Languages() {
					translations := flow.Localization().GetTranslations(lang)
					localizedTemplates := translations.GetTextArray(uuid, key)
					templates = append(templates, localizedTemplates...)
				}
			}
		}
	}

	// Flow and node structs use unexported fields so we have to unpack those manually - also
	// actions and routers are the only things that contain templates or localizable fields.
	for _, node := range flow.Nodes() {
		utils.VisitFields(node.Actions(), visitor)
		utils.VisitFields(node.Router(), visitor)
	}

	return templates
}

func parseEngineTag(tag reflect.StructTag) (localize bool, evaluate bool) {
	for _, tagVal := range strings.Split(tag.Get("engine"), ",") {
		if tagVal == "localize" {
			localize = true
		} else if tagVal == "evaluate" {
			evaluate = true
		}
	}
	return
}

func findStructUUID(v reflect.Value) utils.UUID {
	for f := 0; f < v.Type().NumField(); f++ {
		fld := v.Type().Field(f)

		// look on an encapsulatd field for the UUID
		if fld.Anonymous {
			uuid := findStructUUID(v.FieldByIndex(fld.Index))
			if uuid != "" {
				return uuid
			}
		}

		jsonTag := strings.Split(fld.Tag.Get("json"), ",")
		if len(jsonTag) > 0 && jsonTag[0] == "uuid" {
			uuidValue := v.FieldByIndex(fld.Index)
			return utils.UUID(uuidValue.String())
		}
	}
	return utils.UUID("")
}

func extractTemplates(v reflect.Value, localized bool) []string {
	i := v.Interface()
	switch typed := i.(type) {
	case []string:
		return typed
	case map[string]string:
		return []string{}
	case string:
		return []string{typed}
	default:
		panic("engine tag can only be applied to fields to type string, []string or map[string]string")
	}
}
