package inspect

import (
	"reflect"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

// DependencyContainer allows flow objects to declare other dependencies
type DependencyContainer interface {
	Dependencies(flows.Localization, func(assets.Reference))
}

// Dependencies extracts dependencies
func Dependencies(s interface{}, localization flows.Localization, include func(assets.Reference)) {
	dependencies(reflect.ValueOf(s), localization, include)
}

func dependencies(v reflect.Value, localization flows.Localization, include func(assets.Reference)) {
	walk(
		v,
		func(sv reflect.Value) {
			// anything which is a DependencyContainer can explicitly provide dependencies
			asDepCon, isDepCon := sv.Interface().(DependencyContainer)
			if isDepCon {
				asDepCon.Dependencies(localization, include)
			}
		},
		func(sv reflect.Value, fv reflect.Value, ef *EngineField) {
			// extract any asset.Reference fields automatically as dependencies
			extractAssetReferences(fv, include)
		},
	)
}

func extractAssetReferences(v reflect.Value, include func(assets.Reference)) {
	if v.Kind() == reflect.Slice {
		// field is a slice of asset references
		for i := 0; i < v.Len(); i++ {
			extractAssetReferences(v.Index(i), include)
		}
	} else {
		// field is a single asset reference
		asRef, isRef := v.Interface().(assets.Reference)
		if isRef && asRef != nil && !asRef.Variable() {
			include(asRef)
		}
	}
}
