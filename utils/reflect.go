package utils

import (
	"reflect"
)

// FieldCallback is a callback for a visit to any struct field
type FieldCallback func(v reflect.Value, tag reflect.StructTag)

// VisitFields visits all struct fields recursively in the given value
func VisitFields(s interface{}, visitor FieldCallback) {
	val := reflect.ValueOf(s)
	traverse(val, visitor)
}

func traverse(v reflect.Value, visitor FieldCallback) {
	v, kind := extractType(v)

	switch kind {
	case reflect.Struct:
		for f := 0; f < v.Type().NumField(); f++ {
			fld := v.Type().Field(f)
			val := v.FieldByIndex(fld.Index)

			visitor(val, fld.Tag)

			traverse(val, visitor)
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			traverse(v.Index(i), visitor)
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			traverse(v.MapIndex(k), visitor)
		}
	case reflect.Ptr, reflect.Interface, reflect.Invalid:
		// error?
	}
}

func extractType(current reflect.Value) (reflect.Value, reflect.Kind) {
BEGIN:
	switch current.Kind() {
	case reflect.Ptr:
		if current.IsNil() {
			return current, reflect.Ptr
		}
		current = current.Elem()
		goto BEGIN

	case reflect.Interface:
		if current.IsNil() {
			return current, reflect.Interface
		}
		current = current.Elem()
		goto BEGIN

	case reflect.Invalid:
		return current, reflect.Invalid
	default:
		return current, current.Kind()
	}
}
