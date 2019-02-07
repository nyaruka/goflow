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
	visitor := func(v reflect.Value, tag reflect.StructTag) {
		engineTag := strings.Split(tag.Get("engine"), ",")
		for _, tagVal := range engineTag {
			if tagVal == "evaluate" {
				templates = append(templates, extractTemplates(v)...)
			}
		}
	}

	// flow and node structs use unexported fields so we have to unpack those manually
	for _, node := range flow.Nodes() {
		utils.VisitFields(node.Actions(), visitor)
		utils.VisitFields(node.Router(), visitor)
	}

	return templates
}

func extractTemplates(v reflect.Value) []string {
	i := v.Interface()
	switch typed := i.(type) {
	case []string:
		return typed
	case map[string]string:
		return []string{}
	case string:
		return []string{typed}
	default:
		panic("engine:\"evaluate\" tag can only be applied to fields to type string, []string or map[string]string")
	}
}
