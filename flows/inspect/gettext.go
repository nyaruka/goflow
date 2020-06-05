package inspect

import (
	"reflect"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils/uuids"
)

// LocalizableText extracts localizable text by reading engine tags on a struct
func LocalizableText(s interface{}, include func(uuids.UUID, string, []string, func([]string))) {
	localizableText(reflect.ValueOf(s), include)
}

func localizableText(v reflect.Value, include func(uuids.UUID, string, []string, func([]string))) {
	walk(v, nil, func(sv reflect.Value, fv reflect.Value, ef *EngineField) {
		if ef.Localized {
			localizable := sv.Interface().(flows.Localizable)

			r, w := extractLocalizableText(fv)

			include(localizable.LocalizationUUID(), ef.JSONName, r(), w)
		}
	})
}

// returns read and write functions for the given flow text value
func extractLocalizableText(v reflect.Value) (func() []string, func([]string)) {
	switch typed := v.Interface().(type) {
	case []string:
		r := func() []string {
			return typed
		}
		w := func(n []string) {
			*v.Addr().Interface().(*[]string) = n
		}
		return r, w
	case string:
		r := func() []string {
			return []string{typed}
		}
		w := func(n []string) {
			v.SetString(n[0])
		}
		return r, w
	}
	panic("localized tags can only be applied to strings and slices of strings")
}
