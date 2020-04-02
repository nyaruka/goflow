package inspect

import (
	"reflect"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils/uuids"
)

// LocalizedText extracts localizable text by reading engine tags on a struct
func LocalizedText(s interface{}, include func(uuids.UUID, string, []string)) {
	localizedText(reflect.ValueOf(s), include)
}

func localizedText(v reflect.Value, include func(uuids.UUID, string, []string)) {
	walk(v, nil, func(sv reflect.Value, fv reflect.Value, ef *EngineField) {
		if ef.Localized {
			localizable := sv.Interface().(flows.Localizable)

			include(localizable.LocalizationUUID(), ef.JSONName, extractLocalizedText(fv))
		}
	})
}

// Localized tags can be applied to fields of type string or slices of string
func extractLocalizedText(v reflect.Value) []string {
	switch typed := v.Interface().(type) {
	case []string:
		return typed
	case string:
		return []string{typed}
	}
	panic("localized tags can only be applied to strings and slices of strings")
}
