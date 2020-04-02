package inspect

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/nyaruka/goflow/flows"
)

// EngineField is a struct field which is part of the flow spec (i.e. included in JSON) and optionally has a engine tag
type EngineField struct {
	Type      reflect.Type
	JSONName  string
	Localized bool
	Evaluated bool
	Getter    func(reflect.Value) reflect.Value
}

// extracts all engine fields from the given type
func extractEngineFields(t reflect.Type, rt reflect.Type) []*EngineField {
	fields := make([]*EngineField, 0)
	extractEngineFieldsFromType(t, rt, nil, func(f *EngineField) {
		fields = append(fields, f)
	})
	return fields
}

func extractEngineFieldsFromType(ct reflect.Type, rt reflect.Type, loc []int, include func(*EngineField)) {
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)

		var index []int
		index = append(index, loc...)
		index = append(index, f.Index...)

		// if this is an embedded base struct, inspect its fields too
		if f.Anonymous {
			extractEngineFieldsFromType(ct, f.Type, index, include)
			continue
		}

		jsonName := jsonNameTag(f)
		if jsonName == "" {
			continue
		}

		localized, evaluated := parseEngineTag(ct, f)

		include(&EngineField{
			Type:      f.Type,
			JSONName:  jsonName,
			Localized: localized,
			Evaluated: evaluated,
			Getter:    func(v reflect.Value) reflect.Value { return v.FieldByIndex(index) },
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
func parseEngineTag(st reflect.Type, f reflect.StructField) (localized bool, evaluated bool) {
	t := f.Type
	tagVals := strings.Split(f.Tag.Get("engine"), ",")
	localized = false
	evaluated = false

	var l *flows.Localizable

	for _, v := range tagVals {
		if v == "localized" {
			localized = true

			// if a field has localized, the container struct must implement Localizable
			if !st.Implements(reflect.TypeOf(l).Elem()) {
				panic(fmt.Sprintf("engine:localized tag found on field whose container %v doesn't implement Localizable", st))
			}

			// check field is string or slice of strings - the only things that can be localized
			if !(t.Kind() == reflect.String || (t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.String)) {
				panic(fmt.Sprintf("engine:localized tag found on unsupported type %v", t))
			}
		} else if v == "evaluated" {
			evaluated = true

			// check field is string, slice of strings or map of strings - the only things that can be evaluated
			if !(t.Kind() == reflect.String || (t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.String) || (t.Kind() == reflect.Map && t.Elem().Kind() == reflect.String)) {
				panic(fmt.Sprintf("engine:evaluated tag found on unsupported type %v", t))
			}
		}
	}

	return localized, evaluated
}

func walk(v reflect.Value, visitStruct func(reflect.Value), visitField func(reflect.Value, reflect.Value, *EngineField)) {
	// get the real underlying value
	rv := derefValue(v)

	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			walk(rv.Index(i), visitStruct, visitField)
		}
	} else if rv.Kind() == reflect.Struct {
		if visitStruct != nil {
			visitStruct(v)
		}

		fields := extractEngineFields(v.Type(), rv.Type())

		for _, ef := range fields {
			fv := ef.Getter(rv)

			if visitField != nil {
				visitField(v, fv, ef)
			}

			walk(fv, visitStruct, visitField)
		}
	}
}

// gets the actual value if we've been given an interface or pointer
func derefValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func walkTypes(t reflect.Type, path string, visitField func(string, *EngineField)) {
	// get the real underlying type
	rt := derefType(t)

	if rt.Kind() == reflect.Slice {
		walkTypes(rt.Elem(), path+"[*]", visitField)
	} else if rt.Kind() == reflect.Struct {
		fields := extractEngineFields(t, rt)

		for _, ef := range fields {
			fp := path + "." + ef.JSONName
			if visitField != nil {
				visitField(fp, ef)
			}

			walkTypes(ef.Type, fp, visitField)
		}
	}
}

// gets the actual type if we've been given an interface or pointer type
func derefType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Interface || t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
