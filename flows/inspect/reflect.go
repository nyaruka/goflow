package inspect

import (
	"fmt"
	"reflect"
	"strings"
)

// a struct field which is part of the flow spec (i.e. included in JSON) and optionally has a engine tag
type engineField struct {
	jsonName  string
	localized bool
	evaluated bool
	index     []int
}

// extracts all engine fields from the given type
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

		localized, evaluated := parseEngineTag(f)

		include(&engineField{
			jsonName:  jsonName,
			localized: localized,
			evaluated: evaluated,
			index:     index,
		})
	}
}

func walkFields(v reflect.Value, visit func(*engineField, reflect.Value)) {
	v = derefValue(v)

	if v.Type().Kind() != reflect.Struct {
		return
	}

	for _, ef := range extractEngineFields(v.Type()) {
		fv := v.FieldByIndex(ef.index)

		visit(ef, fv)

		fv = derefValue(fv)

		if fv.Kind() == reflect.Struct {
			walkFields(fv, visit)
		} else if fv.Kind() == reflect.Slice {
			for i := 0; i < fv.Len(); i++ {
				walkFields(fv.Index(i), visit)
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

// gets the JSON name of the given field
func jsonNameTag(f reflect.StructField) string {
	tagVals := strings.Split(f.Tag.Get("json"), ",")
	if len(tagVals) > 0 {
		return tagVals[0]
	}

	return ""
}

// parses the engine tag on a field if it exists
func parseEngineTag(f reflect.StructField) (localized bool, evaluated bool) {
	t := f.Type
	tagVals := strings.Split(f.Tag.Get("engine"), ",")
	localized = false
	evaluated = false

	for _, v := range tagVals {
		if v == "localized" {
			localized = true

			if !(t.Kind() == reflect.String || (t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.String)) {
				panic(fmt.Sprintf("engine:localized tag found on unsupported type %v", t))
			}
		} else if v == "evaluated" {
			evaluated = true

			if !(t.Kind() == reflect.String || (t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.String) || (t.Kind() == reflect.Map && t.Elem().Kind() == reflect.String)) {
				panic(fmt.Sprintf("engine:evaluated tag found on unsupported type %v", t))
			}
		}
	}

	return localized, evaluated
}
