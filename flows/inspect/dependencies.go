package inspect

import (
	"reflect"

	"github.com/nyaruka/goflow/assets"
)

// Dependencies extracts asset dependencies
func Dependencies(s interface{}, include func(assets.Reference)) {
	dependencies(reflect.ValueOf(s), include)
}

func dependencies(v reflect.Value, include func(assets.Reference)) {
	walkFields(v, func(ef *engineField, fv reflect.Value) {
		extractDependencies(fv, include)
	})
}

func extractDependencies(v reflect.Value, include func(assets.Reference)) {
	if v.Kind() == reflect.Slice {
		// field is a slice of asset references
		for i := 0; i < v.Len(); i++ {
			extractDependencies(v.Index(i), include)
		}
	} else {
		// field is a single asset reference
		asRef, isRef := v.Interface().(assets.Reference)
		if isRef && asRef != nil && !asRef.Variable() {
			include(asRef)
		}
	}
}
